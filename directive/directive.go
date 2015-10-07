package directive

import (
	"github.com/natural/missmolly/api"
)

// Package directive registry.
//
var reg = api.DirectiveRegistry{
	&InitDirective{},
	&HttpDirective{},
	&LocationDirective{},
}

// Returns all of the Directives in the package registry.
//
func All() []api.Directive {
	ds := []api.Directive{}
	for _, d := range reg {
		ds = append(ds, d)
	}
	return ds
}

const (
	DIR_PKG = "builtin"

	DIR_HTTP = "http"
	DIR_INIT = "init"
	DIR_LOC  = "location"
)
