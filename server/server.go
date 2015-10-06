package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuin/gopher-lua"

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

	epts   []endpoint
	inits  []func(L *lua.LState) error
	httpds map[string]*http.Server
}

// implement the ServerManipulator interface:
//
func (s *server) OnInit(f func(L *lua.LState) error) {
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
func (s *server) Route(path string) *mux.Route {
	return s.router.Path(path)
}

//
//
func (s *server) Run() (err error) {
	if len(s.epts) == 0 {
		return errors.New("no server endpoints; config missing http and/or https?")
	}
	// for _, f := range s.inits {
	// 	if err := f(s.vm); err != nil {
	// 		return err
	// 	}
	// }
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

type VmResponseWriter struct {
	w http.ResponseWriter
}

func (w *VmResponseWriter) Write(s string) (int, error) {
	return w.w.Write([]byte(s))
}

//
//
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("response", func(L *lua.LState) int {
		api := map[string]lua.LGFunction{
			"write": func(L *lua.LState) int {
				s := L.CheckString(1)
				w.Write([]byte(s))
				L.Push(nil)
				return 1
			},
		}
		t := L.NewTable()
		L.SetFuncs(t, api)
		L.Push(t)
		return 1
	})
	// yeah no
	_ = L.DoString(`response = require "response"
                    response.write("hello, world (lua)")`)
	// if err != nil {
	// 	log.Error("server.vm run", "path", r.URL, "error", err)
	// } else {
	// 	log.Info("server.vm.run", "path", r.URL, "value", "")
	// }
}
