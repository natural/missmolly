package directive

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/natural/missmolly/mm"
	"github.com/robertkrimen/otto"
)

//

func InitDirective(c mm.ServerManipulator, items map[string]interface{}) (bool, error) {
	s, ok := items["init"].(string)
	if !ok {
		return false, errors.New("init value not string")
	}

	// map "file:foo" to a file read
	pfx := "file:"
	if strings.HasPrefix(s, pfx) {
		if bs, err := ioutil.ReadFile(s[len(pfx):]); err != nil {
			log15.Error("directive:init.readfile", "error", err)
			return false, err
		} else {
			// mb parse w otto here too
			s = string(bs)
		}
	}

	c.OnInit(func(vm *otto.Otto) error {
		_, err := vm.Run(s)
		log15.Info("vm:init", "source", s, "error", err)
		return err
	})
	return false, nil
}

//
//
func HttpDirective(c mm.ServerManipulator, items map[string]interface{}) (bool, error) {
	h := map[string]string{}
	if err := mm.Remarshal(items, &h); err != nil {
		log15.Error("directive:http.remarshal", "error", err)
		return false, err
	}
	host, cf, kf := h["http"], h["certfile"], h["keyfile"]
	if host == "" {
		host = h["https"]
	}
	c.AddHost(host, cf, kf)
	log15.Info("directive:http.items", "h", h)
	return false, nil
}

//
//
type LocationDirective struct {
}

//
//
func (d *LocationDirective) Process(c mm.ServerManipulator, items map[string]interface{}) (bool, error) {
	return false, nil
}
