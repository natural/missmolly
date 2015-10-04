package config

import (
	"gopkg.in/yaml.v2"

	"github.com/natural/missmolly/log"
)

// Struct Config holds the run-time application configuration; the main function
// builds and populates one of these.
//
type Config struct {
}

//
//
type flexmap map[string]interface{}

//
//
func (f flexmap) keys() []string {
	a := make([]string, len(f))
	i := 0
	for k, _ := range f {
		a[i] = k
		i++
	}
	return a
}

//
//
func New(bs []byte) (*Config, error) {
	decls := []flexmap{}
	err := yaml.Unmarshal(bs, &decls)
	if err != nil {
		return nil, err
	}
	log.Info("config", "init", true, "decls-count",
		len(decls), "decl-keys", decls[0].keys(), "error", err)
	return &Config{}, nil
}
