package logger

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

var ErrLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)

func Fatal(err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	ErrLog.Output(2, trace)
	os.Exit(1)
}
