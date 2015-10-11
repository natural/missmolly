package config

import (
	"gopkg.in/yaml.v2"

	"github.com/natural/missmolly/log"
)

// Config holds the run-time application configuration.
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
		log.Debug("config.init.unmarshal", "error", err)
		return nil, err
	}
	cfg := &Config{SourceItems: src}
	log.Debug("config.init.success", "bytes", len(cfg.SourceItems))
	return cfg, nil
}
