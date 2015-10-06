package directive

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
	"github.com/robertkrimen/otto"
)

//
type InitDirective struct{}

func (d *InitDirective) Process(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
	s, ok := items["init"].(string)
	if !ok {
		return false, errors.New("init value not string")
	}

	sz := strings.TrimSpace(s)
	if bs, err := ioutil.ReadFile(sz); err == nil {
		s = string(bs)
	} else {
		if !strings.HasSuffix(err.Error(), "no such file or directory") {
			log.Info("directive.init.readfile", "error", err)
		}
	}
	c.OnInit(func(vm *otto.Otto) error {
		v, err := vm.Run(s)
		log.Info("vm.init", "source-len", len(s), "error", err, "value", v)
		return err
	})
	return false, nil
}
