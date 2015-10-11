// Package api provides the library interface for missmolly.
//
package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"
)

// Server slurps listen endpoints, lua initializers, and builds
// routes.  It can also serve http via Run().
//
type Server interface {
	OnInit(func(L *lua.LState) error)
	Endpoint(string, string, string, bool)
	Route(string) *mux.Route
	Run() error

	NewLuaState(w http.ResponseWriter, r *http.Request) *lua.LState
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
	Apply(Server, map[string]interface{}) (bool, error)
	Name() string
	Package() string
}

// Type DirectiveRegistry is a slice of Directives.
//
type DirectiveRegistry []Directive
