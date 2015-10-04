package directive

type InitDirective struct {
}

func (d *InitDirective) Process(c ServerContext, items map[string]interface{}) (bool, error) {
	return false, nil
}

type HttpDirective struct {
}

func (d *HttpDirective) Process(c ServerContext, items map[string]interface{}) (bool, error) {
	return false, nil
}

type LocationDirective struct {
}

func (d *LocationDirective) Process(c ServerContext, items map[string]interface{}) (bool, error) {
	return false, nil
}
