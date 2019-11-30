package virtualbox

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var logger = log.New(os.Stderr, "", 0)

func logLn(msg string) {
	// logger.SetPrefix("\t" + time.Now().Format("2006-01-02 15:04:05") + " ")
	logger.SetPrefix("\t  ")
	logger.Print(msg + "\n")
}

// LogF
func LogF(format string, args ...interface{}) {
	Verbose = testing.Verbose()
	if !Verbose {
		return
	}
	logLn(fmt.Sprintf(format, args...))
}
