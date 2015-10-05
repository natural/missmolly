package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/robertkrimen/otto"

	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/config"
	"github.com/natural/missmolly/directive"
	"github.com/natural/missmolly/log"
)

//
//
func NewFromFile(fn string) (api.ServerManipulator, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Error("server.new.file", "error", err)
		return nil, err
	}
	return NewFromBytes(bs)
}

//
//
func NewFromBytes(bs []byte) (api.ServerManipulator, error) {
	c, err := config.New(bs)
	if err != nil {
		log.Error("server.new.bytes", "error", err)
		return nil, err
	}
	return New(c)
}

//
//
type dirsrc struct {
	dir directive.Directive
	src map[string]interface{}
}

//
//
func New(c *config.Config) (api.ServerManipulator, error) {
	s := &server{
		config: c,
		router: mux.NewRouter(),
		vm:     otto.New(),
	}

	// locate the all directives used
	dss := map[string][]dirsrc{}
	for _, rd := range c.RawItems {
		if k, d := directive.Select(rd); d != nil {
			c.Directives = append(c.Directives, d)
			if _, ok := dss[k]; ok {
				dss[k] = append(dss[k], dirsrc{d, rd})
			} else {
				dss[k] = []dirsrc{{d, rd}}
			}
		} else {
			log.Warn("server.new.directive", "missing", rd)
		}
	}

	// apply the directives in registry order
	for _, dn := range directive.Keys() {
		for _, w := range dss[dn] {
			w.dir.Process(s, w.src)
			log.Info("server.new.process", "directive", dn)
		}
	}

	log.Info("server.new.directives", "count", len(c.Directives))
	return s, nil
}

type endpoint struct {
	addr     string
	certfile string
	keyfile  string
	tls      bool
}

// Root handler, delegates to an internal mux.
//
type server struct {
	config *config.Config
	router *mux.Router
	vm     *otto.Otto

	epts   []endpoint
	inits  []func(*otto.Otto) error
	httpds map[string]*http.Server
}

// implement the ServerManipulator interface:
//
func (s *server) OnInit(f func(*otto.Otto) error) {
	s.inits = append(s.inits, f)
}

//
//
func (s *server) Endpoint(host, certfile, keyfile string, tls bool) {
	s.epts = append(s.epts, endpoint{host, certfile, keyfile, tls})
}

func (s *server) HttpServer(host string) *http.Server {
	return s.httpds[host]
}

//
//
func (s *server) Handler(path string, handler http.Handler) {
	s.router.Handle(path, handler)
}

//
//
func (s *server) Run() (err error) {
	if len(s.epts) == 0 {
		return errors.New("no server endpoints; config missing http and/or https?")
	}
	for _, f := range s.inits {
		if err := f(s.vm); err != nil {
			return err
		}
	}
	wg := sync.WaitGroup{}
	s.httpds = map[string]*http.Server{}
	for _, ep := range s.epts {
		srv := &http.Server{
			Addr:           ep.addr,
			Handler:        s,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		s.httpds[ep.addr] = srv
		wg.Add(1)
		if ep.tls {
			go func() {
				defer wg.Done()
				log.Info("server.run.listen-tls", "addr", srv.Addr)
				if e := srv.ListenAndServeTLS(ep.certfile, ep.keyfile); e != nil {
					log.Error("server.run.listen-tls", "error", e)
					err = e
				}
			}()
		} else {
			go func() {
				defer wg.Done()
				log.Info("server.run.listen", "addr", srv.Addr)
				if e := srv.ListenAndServe(); e != nil {
					log.Error("server.run.listen", "error", e)
					err = e
				}
			}()
		}

	}
	wg.Wait()
	return
}

//
//
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm, err := s.NewVM(w, r)
	if err != nil {
		log.Error("server.new.vm", "error", err)
		return
	}
	js := `
        response.write('hello from javascript: ' + r.Method);
    `
	if val, err := vm.Run(js); err == nil {
		log.Info("server.vm.run", "path", r.URL, "value", val)
	} else {
		log.Error("server.vm run", "path", r.URL, "error", err)
	}

}

//
//
func (s *server) NewVM(w http.ResponseWriter, r *http.Request) (*otto.Otto, error) {
	// faster than channel, mb need a diff way to queue vm copies.  still slow.
	// sync.Pool not faster either.
	vm := s.vm.Copy()

	// build some kind of api instead?
	vm.Set("r", r)
	vm.Set("w", w)

	// build the request object:
	vreq, err := vm.Object(fmt.Sprintf("request = {}"))
	if err != nil {
		return nil, err
	}
	vtbl := map[string]interface{}{
		"method": r.Method,
	}
	for k, v := range vtbl {
		vreq.Set(k, v)
	}

	// build the response object:
	vres, err := vm.Object("response = {}")
	if err != nil {
		return nil, err
	}
	vtbl = map[string]interface{}{
		"write": func(v string) {
			w.Write([]byte(v))
		},
		"status": func(i int) {
			w.WriteHeader(i)
		},
		"header": func() http.Header {
			return w.Header()
		},
	}
	for k, v := range vtbl {
		vres.Set(k, v)
	}

	return vm, nil
}
