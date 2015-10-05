package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/robertkrimen/otto"

	"github.com/natural/missmolly/config"
	"github.com/natural/missmolly/directive"
	"github.com/natural/missmolly/log"
)

//
//
func NewFromBytes(bs []byte) (*Server, error) {
	c, err := config.New(bs)
	if err != nil {
		log.Error("server:new", "error", err, "cond", 0)
		return nil, err
	}
	return New(c)
}

//
//
func New(c *config.Config) (*Server, error) {
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
			log.Warn("server:new.directive", "missing", rd)
		}
	}

	// apply the directives in registry order
	for i, dn := range directive.Keys() {
		for j, w := range dss[dn] {
			w.d.Process(s, w.i)
			log.Info("server:new.process", "w", w, "i", i, "j", j)
		}
	}

	log.Info("server:new.directives", "count", len(c.Directives))
	return s, nil
}

//
//
func NewFromFile(fn string) (*Server, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(bs)
}

// Root handler, delegates to an internal mux.
//
type Server struct {
	Config *config.Config
	Router *mux.Router
	VM     *otto.Otto

	Hosts []string

	inits []func(*otto.Otto) error
}

func (s *Server) init() error {
	for _, f := range s.inits {
		if err := f(s.VM); err != nil {
			return err
		}
	}
	return nil
}

// implement the ServerManipulator interface:
func (s *Server) OnInit(f func(*otto.Otto) error) error {
	s.inits = append(s.inits, f)
	return nil
}

func (s *Server) AddHost(host, certFile, keyFile string) {
	s.Hosts = append(s.Hosts, host)
}

func (s *Server) AddHandler(path string, handler http.Handler) {
	s.Router.Handle(path, handler)
}

//
//
func (s *Server) ListenAndServe() error {
	// all hosts in conf ofc
	srv := &http.Server{
		Addr:           ":7373",
		Handler:        s, //app per server, or...
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.init()
	return srv.ListenAndServe()
}

//
//
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm, err := s.NewVM(w, r)
	if err != nil {
		log.Error("server", "error", err, "cond", 1)
		return
	}
	js := `
        response.write('hello from javascript: ' + r.Method);
        console.log("client accepts:", r.Header["Accept"])
    `
	if _, err := vm.Run(js); err != nil {
		log.Error("vm:run", "path", r.URL, "error", err)
	} else {
		log.Info("vm:run", "path", r.URL, "status", "?")
	}

}

//
//
func (s *Server) NewVM(w http.ResponseWriter, r *http.Request) (*otto.Otto, error) {
	// faster than channel, mb need a diff way to queue vm copies.  still slow.
	vm := s.VM.Copy()
	vm.Set("r", r)
	vm.Set("w", w)

	// con, err := vm.Object("console")
	// if err != nil {
	// 	return nil, err
	// }
	// con.Set("log", func(items ...interface{}) {
	// 	log.Info("request", items...)
	// })

	// build some kind of api instead:
	// build the request object
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

	// build the response object
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
