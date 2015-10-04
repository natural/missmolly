package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/natural/missmolly/config"
	"github.com/robertkrimen/otto"
)

//
//
func NewFromBytes(bs []byte) (*Server, error) {
	c, err := config.New(bs)
	if err != nil {
		log.Fatal(err)
	}
	return New(c)
}

//
//
func New(c *config.Config) (*Server, error) {
	r := mux.NewRouter()
	o := otto.New()
	s := &Server{
		Config: c,
		Router: r,
		VM:     o,
	}
	return s, nil
}

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
}

func (s *Server) ListenAndServe() error {
	// all hosts in conf ofc
	srv := &http.Server{
		Addr:           ":7373",
		Handler:        s, //app per server, or...
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return srv.ListenAndServe()
}

//
//
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm := s.NewVM(w, r)
	js := `
        var header = response.header()
        response.write('hello from javascript: ' + request.method);
        console.log("header:", header.X, header.Y)
    `
	_, err := vm.Run(js)
	log.Printf("vm error: %v", err)
}

//
//
func (s *Server) NewVM(w http.ResponseWriter, r *http.Request) {
	vm := s.VM.Copy() // faster than channel

	// build the request object
	_, err := vm.Object(fmt.Sprintf("request = {method: '%v'}", r.Method))
	if err != nil {
		log.Fatal(err)
	}

	// build the response object
	vres, err := vm.Object("response = {}")
	if err != nil {
		log.Fatal(err)
	}
	resvtbl := map[string]interface{}{
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
	for k, v := range resvtbl {
		vres.Set(k, v)
	}

}
