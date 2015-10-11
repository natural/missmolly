package main

import (
	"flag"

	"github.com/natural/missmolly/log"
	"github.com/natural/missmolly/server"
)

// name some options and defaults.
//
var opts = struct {
	config   string
	loglevel string
}{
	"missmolly.conf",
	log.CurrentLevel.String(),
}

// slurp those options.
//
func init() {
	flag.StringVar(&opts.loglevel, "loglevel", opts.loglevel, "log level")
	flag.StringVar(&opts.config, "config", opts.config, "config file name")
	flag.Parse()
	log.Setup(opts.loglevel)
}

// run that server.
//
func main() {
	if srv, err := server.NewFromFile(opts.config); err == nil {
		log.Fatal(srv.Run())
	} else {
		log.Fatal(err)
	}
}
