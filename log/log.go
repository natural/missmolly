package log

import (
	"log"

	"github.com/inconshreveable/log15"
)

//
//
func init() {
	Log.SetHandler(log15.DiscardHandler())
}

//
//
var (
	Log   = log15.New()
	Error = Log.Error
	Warn  = Log.Warn
	Info  = Log.Info
	Debug = Log.Debug

	LvlCurrent = log15.LvlDebug
)

//
//
func Fatal(err error) {
	log.Fatal(err)
}

//
//
func Setup(n string) {
	lvl := LvlCurrent
	v, err := log15.LvlFromString(n)
	if err == nil {
		lvl = v
	}
	Log.SetHandler(log15.LvlFilterHandler(lvl, log15.StderrHandler))
}
