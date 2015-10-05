// Package mm provides the api for using  missmolly as a library.
//
package mm

import (
	"net/http"

	"github.com/robertkrimen/otto"
	"gopkg.in/yaml.v2"
)

//
//
type ServerManipulator interface {
	OnInit(func(*otto.Otto) error) error
	AddHost(string, string, string)
	AddHandler(string, http.Handler)
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
