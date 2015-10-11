package builtin

import (
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
)

//
//
type HttpDirective struct {
	Http     string
	Https    string
	CertFile string
	KeyFile  string
}

//
//
func (d *HttpDirective) Name() string {
	return dirhttp
}

//
//
func (d *HttpDirective) Package() string {
	return dirpkg
}

//
//
func (d *HttpDirective) Accept(decl map[string]interface{}) bool {
	_, ok := decl[dirhttp]
	if !ok {
		_, ok = decl[dirhttp+"s"]
	}
	return ok
}

//
//
func (d *HttpDirective) Apply(c api.Server, items map[string]interface{}) (bool, error) {
	if err := api.Remarshal(items, d); err != nil {
		log.Error("directive.http.remarshal", "error", err)
		return false, err
	}
	host, cf, kf := d.Http, d.CertFile, d.KeyFile
	if host == "" {
		host = d.Https
	}
	// mb check cert and key files exist and are readable
	c.Endpoint(host, cf, kf, host == d.Https)
	return false, nil
}
