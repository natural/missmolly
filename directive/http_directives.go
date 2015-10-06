package directive

import (
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
)

type HttpDirective struct {
	Http     string
	Https    string
	CertFile string
	KeyFile  string
}

//
//
func (d *HttpDirective) Process(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
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
