package directive

import (
	"github.com/natural/missmolly/mm"
)

//
//
type LocationDirective struct {
}

//
//
func (d *LocationDirective) Process(c mm.ServerManipulator, items map[string]interface{}) (bool, error) {
	return false, nil
}
