package directive

import (
	"github.com/natural/missmolly/mm"
)

// Directive interprets a map of objects, possibly creating http handlers, possibly
// modifying the server, etc.
//
type Directive interface {
	Process(mm.ServerManipulator, map[string]interface{}) (bool, error)
}

type Directives []Directive

// Global directive registry.
//
var reg = Registry{
	{"builtin", "init", FromFunc(InitDirective)},
	{"builtin", "http", FromFunc(HttpDirective)},
	{"builtin", "https", FromFunc(HttpDirective)},
	{"builtin", "location", &LocationDirective{}},
}

//
//
func Keys() []string {
	ks := make([]string, len(reg))
	for i, r := range reg {
		ks[i] = r.Key
	}
	return ks
}

// RegistryEntry is a tuple of (key, directive).
//
type RegistryEntry struct {
	Package   string
	Key       string
	Directive Directive
}

// Registry is a slice of RegistryEntry.
//
type Registry []RegistryEntry

// Matches items by key to a directive.
//
func Select(items map[string]interface{}) (string, Directive) {
	for _, r := range reg {
		if _, ok := items[r.Key]; ok {
			return r.Key, r.Directive
		}
	}
	return "", nil
}

//
type DirectiveFunc func(mm.ServerManipulator, map[string]interface{}) (bool, error)

//
//
func FromFunc(f DirectiveFunc) Directive {
	return dfw{f}
}

//
//
type dfw struct{ df DirectiveFunc }

func (d dfw) Process(srv mm.ServerManipulator, items map[string]interface{}) (bool, error) {
	return d.df(srv, items)
}
