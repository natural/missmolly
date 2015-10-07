package config

import (
	"gopkg.in/yaml.v2"

	"github.com/natural/missmolly/log"
)

// Struct Config holds the run-time application configuration; the main function
// builds and populates one of these.
//
type Config struct {
	SourceItems []map[string]interface{}
}

//
//
func New(bs []byte) (*Config, error) {
	src := []map[string]interface{}{}
	err := yaml.Unmarshal(bs, &src)
	if err != nil {
		log.Error("config.init.unmarshal", "error", err)
		return nil, err
	}
	cfg := &Config{
		SourceItems: src,
	}
	log.Info("config.init.success", "bytes", len(bs))
	return cfg, nil
}
