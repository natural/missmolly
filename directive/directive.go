package directive

import (
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/directive/builtin"
)

// Package directive registry.
//
var reg = api.DirectiveRegistry{
	&builtin.InitDirective{},
	&builtin.HttpDirective{},
	&builtin.LocationDirective{},
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

// Mosdef will give you an error later.
//
func Register(d api.Directive) error {
	reg = append(reg, d)
	return nil
}
