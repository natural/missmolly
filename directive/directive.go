package directive

import (
	"net/http"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"
)

var DirectiveRegistry = map[string]Directive{
	"init": &InitDirective{},
	"http": &HttpDirective{},
}

type DirectiveHandlerFunc func(*ServerContext, []byte) (*ServerContext, error)

// Directive interprets a map of objects, possibly creating http handlers, possibly
// modifying the server, etc.
//
type Directive interface {
	Process(ServerContext, map[string]interface{}) (bool, error)
}

type ServerContext interface {
	RootRouter() *mux.Router
}

//
//
type Directives interface {
	//	Routes() ServerRoutes
}

type DirectiveFilter interface {
	Filter(*http.Request, http.ResponseWriter) (bool, error)
}

// Map an interfacy thing into something else.
//
func Remarshal(in map[string]interface{}, out interface{}) error {
	bs, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(bs, out)
}
