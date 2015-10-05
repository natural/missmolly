package directive

import (
	"github.com/natural/missmolly/api"
	"github.com/natural/missmolly/log"
)

//
//
func HttpDirective(c api.ServerManipulator, items map[string]interface{}) (bool, error) {
	h := map[string]string{}
	if err := api.Remarshal(items, &h); err != nil {
		log.Error("directive.http.remarshal", "error", err)
		return false, err
	}
	host, cf, kf := h["http"], h["certfile"], h["keyfile"]
	if host == "" {
		host = h["https"]
	}
	// mb check cert and key files exist and are readable
	c.Endpoint(host, cf, kf, host == h["https"])
	//log.Info("directive.http.items", "h", h)
	return false, nil
}
