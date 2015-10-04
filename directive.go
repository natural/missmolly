package missmolly

import (
	"net/http"

	"github.com/gorilla/mux"
)

//
//
var DirectiveRegistry = map[string]Directive{
//	"init": InitDirective{},
//	"http": HttpDirective{},
}

type InitDirective struct {
}

// A DirectiveHandler can interpret a config directive.
//
type DirectiveHandler interface {
	ProcessDirective(*ServerContext, []byte) (bool, error)
}

//
//
type ServerContext struct {
	Confg      *Config
	RootMux    *mux.Router
	Directives Directives
}

//
//
type Directives interface {
	Routes() ServerRoutes
}

type Directive interface {
	http.Handler
	Filter(*http.Request, http.ResponseWriter) (bool, error)
}
