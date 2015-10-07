// Package api provides the library interface for missmolly.
//
package api

import (
	"github.com/gorilla/mux"
	"github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
)

// Directives are given objects that implement ServerManipulator.
//
type ServerManipulator interface {
	OnInit(func(L *lua.LState) error)
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

// Directive interprets a map of objects, possibly creating http
// handlers, possibly modifying the server, etc.
//
type Directive interface {
	Accept(map[string]interface{}) bool
	Name() string
	Package() string
	Process(ServerManipulator, map[string]interface{}) (bool, error)
}

// Type DirectiveRegistry is a slice of Directives.
//
type DirectiveRegistry []Directive
