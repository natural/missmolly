package log

import (
	"os"

	"github.com/inconshreveable/log15"
)

// Set the package global `Log` to the discard handler; any code that imports
// missmolly as a library will not have our logging output.
//
func init() {
	Log.SetHandler(log15.DiscardHandler())
}

var (
	// package-level logger
	Log = log15.New()

	// expose these for ease of use
	Error = Log.Error
	Warn  = Log.Warn
	Info  = Log.Info
	Debug = Log.Debug

	// track the log level as set by Setup
	CurrentLevel = log15.LvlDebug
)

// Like log.Fatal but uses log15 instead.
//
func Fatal(v ...interface{}) {
	Log.Crit("fatal", "msg", v)
	os.Exit(1)
}

// Set or reset the package logger.
//
func Setup(n string) {
	lvl := CurrentLevel
	v, err := log15.LvlFromString(n)
	if err == nil {
		lvl = v
		CurrentLevel = lvl
	}
	Log.SetHandler(log15.LvlFilterHandler(lvl, log15.StderrHandler))
}
