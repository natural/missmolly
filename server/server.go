package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/yuin/gopher-lua"

	"github.com/natural/glox"
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/config"
	"github.com/natural/missmolly/directive"
	"github.com/natural/missmolly/log"
)

// NewFromFile builds a Server from the given config file.
// File should be YAML data.
//
func NewFromFile(fn string) (api.Server, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(bs)
}

// NewFromBytes builds a Server from the given byte slice.
// Contents should be YAML.
//
func NewFromBytes(bs []byte) (api.Server, error) {
	c, err := config.New(bs)
	if err != nil {
		return nil, err
	}
	return New(c)
}

// New builds a Server from the given Config struct.
//
func New(c *config.Config) (api.Server, error) {
	r := mux.NewRouter()
	s := &Server{
		Config:    c,
		Router:    r,
		RootState: lua.NewState(),
	}

	for _, dir := range directive.All() {
		for i, decl := range c.SourceItems {
			if dir.Accept(decl) {
				dir.Apply(s, decl)
				log.Debug("directive.process", "name", dir.Name(), "block", i)
			}
		}
	}
	log.Debug("server.new", "server", s)
	return s, nil
}

//
//
type Endpoint struct {
	Addr     string
	CertFile string
	KeyFile  string
	TLS      bool
}

// Root handler, delegates to an internal mux.
//
type Server struct {
	Config *config.Config
	Router *mux.Router

	Endpoints   []Endpoint
	HttpServers map[string]*http.Server
	InitFuncs   []func(L *lua.LState) error
	RootState   *lua.LState
}

// implement the Server interface:
//
func (s *Server) OnInit(f func(L *lua.LState) error) {
	s.InitFuncs = append(s.InitFuncs, f)
}

//
//
func (s *Server) Endpoint(host, certfile, keyfile string, tls bool) {
	s.Endpoints = append(s.Endpoints, Endpoint{host, certfile, keyfile, tls})
}

func (s *Server) HttpServer(host string) *http.Server {
	return s.HttpServers[host]
}

//
//
func (s *Server) Route(path string) *mux.Route {
	return s.Router.Path(path)
}

//
//
func (s *Server) Run() (err error) {
	if len(s.Endpoints) == 0 {
		return errors.New("no server endpoints; config missing http and/or https?")
	}
	L := s.RootState
	if L == nil {
		return errors.New("no root lua state")
	}
	for _, init := range s.InitFuncs {
		if err := init(L); err != nil {
			return err
		}
	}
	wg := sync.WaitGroup{}
	s.HttpServers = map[string]*http.Server{}
	for _, ep := range s.Endpoints {
		srv := &http.Server{
			Addr:           ep.Addr,
			Handler:        s.Router,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		s.HttpServers[ep.Addr] = srv
		wg.Add(1)
		if ep.TLS {
			go func() {
				defer wg.Done()
				log.Debug("server.run", "server", srv)
				log.Info("server.run.listen-tls", "addr", srv.Addr)
				if e := srv.ListenAndServeTLS(ep.CertFile, ep.KeyFile); e != nil {
					log.Error("server.run.listen-tls", "error", e)
					err = e
				}
			}()
		} else {
			go func() {
				defer wg.Done()
				log.Debug("server.run", "server", srv)
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

// This method needs to morph into something that routes the request;
// the stuff bits already here need to move to some kind of ContentHandler.
//
// func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	s.Router.ServeHTTP(w, r)
// 	//L := s.NewLuaState(w, r)
// 	//defer L.Close()
// 	//err := L.DoFile("test.lua")
// 	//err := L.DoString("response.write('hello, world')")
// 	//if err != nil {
// 	//	log.Error("server.request.lua", "path", r.URL, "error", err)
// 	//}
// }

//
//
func (s *Server) NewLuaState(w http.ResponseWriter, r *http.Request) *lua.LState {
	// smaller sizes == faster benchmarks
	L := lua.NewState(lua.Options{
		CallStackSize:       120,
		RegistrySize:        128, // smallest per source
		SkipOpenLibs:        true,
		IncludeGoStackTrace: true,
	})
	// benchmark killer... make optional?
	//L.OpenLibs()

	glox.LCopyGlobal(s.RootState, L, "motd", "app", "session")
	L.SetGlobal("request", glox.LHttpRequest(L, r))
	L.SetGlobal("response", glox.LHttpResponseWriter(L, w))
	return L
}
