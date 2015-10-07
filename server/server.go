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

// NewFromFile builds a ServerManipulator from the given config file.
// File should be YAML data.
//
func NewFromFile(fn string) (api.ServerManipulator, error) {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Error("server.new.file", "error", err)
		return nil, err
	}
	return NewFromBytes(bs)
}

// NewFromBytes builds a ServerManipulator from the given byte slice.
// Contents should be YAML.
//
func NewFromBytes(bs []byte) (api.ServerManipulator, error) {
	c, err := config.New(bs)
	if err != nil {
		log.Error("server.new.bytes", "error", err)
		return nil, err
	}
	return New(c)
}

// New builds a ServerManipulator from the given Config struct.
//
func New(c *config.Config) (api.ServerManipulator, error) {
	r := mux.NewRouter()
	s := &server{config: c, router: r, rootvm: lua.NewState()}

	for _, dir := range directive.All() {
		for i, decl := range c.SourceItems {
			if dir.Accept(decl) {
				dir.Process(s, decl)
				log.Info("server.new",
					"process.directive", dir.Name(), "block", i)
			}
		}
	}

	log.Info("server.new", "error", nil)
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

	rootvm *lua.LState
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
	for _, f := range s.inits {
		if err := f(s.rootvm); err != nil {
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

// This method needs to morph into something that routes the request;
// the stuff bits already here need to move to some kind of ContentHandler.
//
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	L := s.NewVM(w, r)
	defer L.Close()

	err := L.DoString(`response.write(request.user_agent())`)
	if err != nil {
		log.Error("server.vm run", "path", r.URL, "error", err)
	}
}

//
func (s *server) NewVM(w http.ResponseWriter, r *http.Request) *lua.LState {
	// smaller sizes == faster benchmarks
	L := lua.NewState(lua.Options{
		CallStackSize:       120,
		RegistrySize:        128, // smallest per source
		SkipOpenLibs:        true,
		IncludeGoStackTrace: true,
	})

	// copy some globals from the root vm into this one
	for _, g := range []string{"motd", "app", "session"} {
		L.SetGlobal(g, s.rootvm.GetGlobal(g))
	}

	// make the "response" api (incomplete)
	lres := L.NewTable()
	L.SetGlobal("response", lres)
	L.SetFuncs(lres, map[string]lua.LGFunction{
		"write": func(L *lua.LState) int {
			s := L.CheckString(1)
			c, err := w.Write([]byte(s))
			e := ""
			if err != nil {
				e = err.Error()
			}
			L.Push(lua.LString(e))
			L.Push(lua.LNumber(c))
			return 2
		},
	})

	// make the "request" api (incomplete)
	lreq := L.NewTable()
	L.SetGlobal("request", lreq)
	lreq.RawSetString("method", lua.LString(r.Method))
	lreq.RawSetString("remote_addr", lua.LString(r.RemoteAddr))
	lreq.RawSetString("request_uri", lua.LString(r.RequestURI))
	lreq.RawSetString("url", lua.LString(r.URL.String()))
	lreq.RawSetString("proto", lua.LString(r.Proto))
	lreq.RawSetString("proto_major", lua.LNumber(r.ProtoMajor))
	lreq.RawSetString("proto_minor", lua.LNumber(r.ProtoMinor))
	lreq.RawSetString("content_length", lua.LNumber(r.ContentLength))
	lreq.RawSetString("close", lua.LBool(r.Close))
	lreq.RawSetString("host", lua.LString(r.Host))

	L.SetFuncs(lreq, map[string]lua.LGFunction{
		"header": func(L *lua.LState) int {
			// conv header map and return it
			return 0
		},
		"user_agent": func(L *lua.LState) int {
			L.Push(lua.LString(r.UserAgent()))
			return 1
		},
	})
	return L
}
