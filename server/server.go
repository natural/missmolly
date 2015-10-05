package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/inconshreveable/log15"
	"github.com/robertkrimen/otto"

	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/config"
	"github.com/natural/missmolly/directive"
	"github.com/natural/missmolly/log"
)

//
//
func NewFromBytes(bs []byte) (api.ServerManipulator, error) {
	c, err := config.New(bs)
	if err != nil {
		log.Error("server.new.config", "error", err)
		return nil, err
	}
	return New(c)
}

//
//
func New(c *config.Config) (api.ServerManipulator, error) {
	s := &Server{
		Config: c,
		Router: mux.NewRouter(),
		VM:     otto.New(),
	}

	type ds struct {
		d directive.Directive
		i map[string]interface{}
	}

	// locate the all directives used
	dss := map[string][]ds{}
	for _, rd := range c.RawItems {
		if k, d := directive.Select(rd); d != nil {
			c.Directives = append(c.Directives, d)
			if _, ok := dss[k]; ok {
				dss[k] = append(dss[k], ds{d, rd})
			} else {
				dss[k] = []ds{{d, rd}}
			}
		} else {
			log.Warn("server.new.directive", "missing", rd)
		}
	}

	// apply the directives in registry order
	for i, dn := range directive.Keys() {
		for j, w := range dss[dn] {
			w.d.Process(s, w.i)
			log.Info("server.new.process", "w", w, "i", i, "j", j)
		}
	}

	log.Info("server.new.directives", "count", len(c.Directives))
	return s, nil
}

//
//
func NewFromFile(fn string) (api.ServerManipulator, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(bs)
}

type endpoint struct {
	addr     string
	certfile string
	keyfile  string
	tls      bool
}

// Root handler, delegates to an internal mux.
//
type Server struct {
	Config *config.Config
	Router *mux.Router
	VM     *otto.Otto

	epts   []endpoint
	inits  []func(*otto.Otto) error
	httpds map[string]*http.Server
}

// implement the ServerManipulator interface:
//
func (s *Server) OnInit(f func(*otto.Otto) error) {
	s.inits = append(s.inits, f)
}

//
//
func (s *Server) Endpoint(host, certfile, keyfile string, tls bool) {
	s.epts = append(s.epts, endpoint{host, certfile, keyfile, tls})
}

func (s *Server) HttpServer(host string) *http.Server {
	return s.httpds[host]
}

//
//
func (s *Server) Handler(path string, handler http.Handler) {
	s.Router.Handle(path, handler)
}

//
//
func (s *Server) Run() (err error) {
	if len(s.epts) == 0 {
		return errors.New("no server endpoints; config missing http and/or https?")
	}
	for _, f := range s.inits {
		if err := f(s.VM); err != nil {
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
				log15.Info("server:run.listen-tls", "addr", srv.Addr)
				if e := srv.ListenAndServeTLS(ep.certfile, ep.keyfile); e != nil {
					log15.Error("server:run.listen-tls", "error", e)
					err = e
				}
			}()
		} else {
			go func() {
				defer wg.Done()
				log15.Info("server:run.listen", "addr", srv.Addr)
				if e := srv.ListenAndServe(); e != nil {
					log15.Error("server:run.listen", "error", e)
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
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm, err := s.NewVM(w, r)
	if err != nil {
		log.Error("server.new.vm", "error", err)
		return
	}
	js := `
        response.write('hello from javascript: ' + r.Method);
    `
	if _, err := vm.Run(js); err != nil {
		log.Error("server.vm run", "path", r.URL, "error", err)
	} else {
		log.Info("server.vm.run", "path", r.URL, "status", "?")
	}

}

//
//
func (s *Server) NewVM(w http.ResponseWriter, r *http.Request) (*otto.Otto, error) {
	// faster than channel, mb need a diff way to queue vm copies.  still slow.
	// sync.Pool not much faster either.
	vm := s.VM.Copy()

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
