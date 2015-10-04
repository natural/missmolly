package missmolly

// Struct Config holds the run-time application configuration; the main function
// builds and populates one of these.
//
type Config struct {
}

func NewConfig(bs []byte) (*Config, error) {
	return &Config{}, nil
}
