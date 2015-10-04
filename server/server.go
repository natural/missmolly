package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/natural/missmolly/config"
	"github.com/robertkrimen/otto"
)

//
//
func New(bs []byte) *Server {
	c, err := config.New(bs)
	if err != nil {
		log.Fatal(err)
	}
	return NewFromConfig(c)
}

func NewFromConfig(c *config.Config) *Server {
	r := mux.NewRouter()
	o := otto.New()

	s := &Server{
		Config: c,
		Router: r,
		VM:     o,
	}

	return s
}

func NewFromFile(file string) *Server {
	bs := []byte{}
	return New(bs)
}

// Root handler, delegates to an internal mux.
//
type Server struct {
	Config *config.Config
	Router *mux.Router
	VM     *otto.Otto
}

//
//
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vm := s.VM.Copy() // faster than channel
	_, err := vm.Object(fmt.Sprintf("request = {method: '%v'}", r.Method))
	if err != nil {
		log.Fatal(err)
	}
	vres, err := vm.Object("response = {}")
	if err != nil {
		log.Fatal(err)
	}
	vres.Set("write", func(v string) {
		w.Write([]byte(v))
	})
	vres.Set("status", func(i int) {
		w.WriteHeader(i)
	})
	vres.Set("headers", func() http.Header {
		return w.Header()
	})

	//
	js := "response.write('hello from javascript: ' + request.method)"
	//w.Header()["Content-Type"] = []string{"text/json"}
	_, _ = vm.Run(js)
	//log.Printf("vm value: %v err: %v", v, err)
}

//
//
type ServerRoutes []*ServerRoute

//
//
type ServerRoute struct {
	Path string
}
