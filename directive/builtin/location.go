package builtin

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
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
	return dirloc
}

//
//
func (d *LocationDirective) Package() string {
	return dirpkg
}

//
//
func (d *LocationDirective) Accept(decl map[string]interface{}) bool {
	_, ok := decl[dirloc]
	return ok
}

//
//
func (d *LocationDirective) Apply(c api.Server, items map[string]interface{}) (bool, error) {
	if err := api.Remarshal(items, d); err != nil {
		return false, err
	}
	if d.Location == "" {
		return false, errors.New("location directive missing path")
	}
	r := c.Route(strings.TrimSpace(d.Location))
	r = d.Match.SetRoute(r)
	r = d.SetContentHandler(c, r)
	return false, nil
}

func (d *LocationDirective) SetContentHandler(c api.Server, r *mux.Route) *mux.Route {
	//
	r = r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		L := c.NewLuaState(w, r)
		defer L.Close()

		// if file, read once and set string
		err := L.DoFile("test.lua")

		// if string, set string

		// if value in handler map, set handler

		// if values in handler map, set handlers

		if err != nil {
			log.Error("server.request.lua", "path", r.URL, "error", err)
		}

	})
	return r
}

//
//
type Match struct {
	Prefix  string
	Hosts   []string
	Methods []string
	Schemes []string
	Queries []map[string]string
	Headers []map[string]string
	Custom  []string
}

//
//
func (m Match) SetRoute(r *mux.Route) *mux.Route {
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
	if m.Schemes != nil && len(m.Schemes) > 0 {
		for _, scheme := range m.Schemes {
			r = r.Schemes(scheme)
		}
	}
	if m.Queries != nil && len(m.Queries) > 0 {
		for _, m := range m.Queries {
			for k, v := range m {
				r = r.Queries(k, v)
			}
		}
	}
	if m.Headers != nil && len(m.Headers) > 0 {
		for _, m := range m.Headers {
			for k, v := range m {
				r = r.Headers(k, v)
			}
		}
	}
	if m.Custom != nil && len(m.Custom) > 0 {
		// lookup custom function by name
		// func(r *http.Request, rm *RouteMatch) bool
	}
	return r
}
