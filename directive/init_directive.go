package directive

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/inconshreveable/log15"
	"github.com/natural/missmolly/api"
	"github.com/robertkrimen/otto"
)

//

func InitDirective(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
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
