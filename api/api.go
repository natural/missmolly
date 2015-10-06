// Package api provides the library interface for missmolly.
//
package api

import (
	"github.com/gorilla/mux"
	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

//
//
type ServerManipulator interface {
	OnInit(func(*otto.Otto) error)
	Endpoint(string, string, string, bool)
	Route(string) *mux.Route
	Run() error
}

// Map an interfacy thing into something else.
//
func Remarshal(in interface{}, out interface{}) error {
	bs, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, out)
}
