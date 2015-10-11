package builtin

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
	"github.com/yuin/gopher-lua"
)

//
//
type InitDirective struct{}

//
//
func (d *InitDirective) Name() string {
	return dirinit
}

//
//
func (d *InitDirective) Package() string {
	return dirpkg
}

//
//
func (d *InitDirective) Accept(decl map[string]interface{}) bool {
	_, ok := decl[dirinit]
	return ok
}

//
//
func (d *InitDirective) Apply(c api.Server, items map[string]interface{}) (bool, error) {
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

	c.OnInit(func(L *lua.LState) error {
		L.SetGlobal("motd", lua.LString("hello, world oninit"))
		return nil
	})

	return false, nil
}
