package config

import (
	"gopkg.in/yaml.v2"

	"github.com/natural/missmolly/directive"
	"github.com/natural/missmolly/log"
)

// Struct Config holds the run-time application configuration; the main function
// builds and populates one of these.
//
type Config struct {
	Directives directive.Directives
	RawItems   []map[string]interface{}
}

//
//
func New(bs []byte) (*Config, error) {
	rds := []map[string]interface{}{}
	err := yaml.Unmarshal(bs, &rds)
	if err != nil {
		log.Error("config.init.unmarshal", "error", err)
		return nil, err
	}
	drs := directive.Directives{}
	cfg := &Config{
		Directives: drs,
		RawItems:   rds,
	}
	log.Info("config.init.success", "bytes", len(bs))
	return cfg, nil
}
