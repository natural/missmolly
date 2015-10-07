package directive

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/natural/missmolly/api"
)

//
//
type LocationDirective struct {
	Location   string
	Content    string
	Match      Match
	Middleware []string

	Auth   string
	Nested []LocationDirective
}

//
//
func (d *LocationDirective) Name() string {
	return DIR_LOC
}

//
//
func (d *LocationDirective) Package() string {
	return DIR_PKG
}

//
//
func (d *LocationDirective) Accept(decl map[string]interface{}) bool {
	_, ok := decl[DIR_LOC]
	return ok
}

//
//
func (d *LocationDirective) Process(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
	if err := api.Remarshal(items, d); err != nil {
		return false, err
	}
	if d.Location == "" {
		return false, errors.New("location directive missing path")
	}
	r := c.Route(strings.TrimSpace(d.Location))
	d.Match.ApplyTo(r)
	return false, nil
}

//
//
type Match struct {
	Prefix  string
	Hosts   []string
	Methods []string
	Schemes []string
	Queries []string
	Headers []map[string]string
	Custom  []string
}

//
//
func (m Match) ApplyTo(r *mux.Route) *mux.Route {
	if m.Prefix != "" {
		r = r.PathPrefix(m.Prefix)
	}
	if m.Hosts != nil && len(m.Hosts) > 0 {
		for _, host := range m.Hosts {
			r = r.Host(host)
		}
	}
	if m.Methods != nil && len(m.Methods) > 0 {
		r = r.Methods(m.Methods...)
	}
	// schemes, queries, headers, custom, etc

	r = r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK!"))
	})
	return r
}
