package directive

import (
	"github.com/inconshreveable/log15"
	"github.com/natural/missmolly/mm"
)

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
	// mb check cert and key files exist and are readable
	c.AddEndpoint(host, cf, kf)
	log15.Info("directive:http.items", "h", h)
	return false, nil
}
