package log

import (
	"log"

	"github.com/inconshreveable/log15"
)

var (
	Log   = log15.New()
	Error = Log.Error
	Warn  = Log.Warn
	Info  = Log.Info
	Debug = Log.Debug
)

//
//
func Fatal(err error) {
	log.Fatal(err)
}

//
//
func init() {
	Log.SetHandler(log15.StderrHandler)
}
