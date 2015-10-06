package main

import (
	"flag"

	"github.com/inconshreveable/log15"
	"github.com/natural/missmolly/log"
	"github.com/natural/missmolly/server"
)

// main program options their defaults
//
var opts = struct {
	config   string
	loglevel log15.Lvl
}{
	"missmolly.conf",
	log15.LvlWarn,
}

// slurp those options
//
func init() {
	flag.StringVar(&opts.config, "config", opts.config, "config file name")
	ln := "debug"
	flag.StringVar(&ln, "loglevel", ln, "log level")
	flag.Parse()

	lv, _ := log15.LvlFromString(ln)
	lh := log15.MultiHandler(log15.LvlFilterHandler(lv, log15.StderrHandler))
	log.Log.SetHandler(lh)

	opts.loglevel = lv
}

// run that server.
//
func main() {
	srv, err := server.NewFromFile(opts.config)
	if err != nil {
		log.Error("main", "error", err)
		log.Fatal(err)
	}
	log.Fatal(srv.Run())
}
