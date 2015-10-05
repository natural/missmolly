package directive

import (
	"github.com/natural/missmolly/api"
)

//
//
type LocationDirective struct {
}

//
//
func (d *LocationDirective) Process(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
	return false, nil
}
