package virtualbox

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", 0)

func logLn(msg string) {
	// logger.SetPrefix("\t" + time.Now().Format("2006-01-02 15:04:05") + " ")
	logger.SetPrefix("\t  ")
	logger.Print(msg + "\n")
}

func logF(format string, args ...interface{}) {
	logLn(fmt.Sprintf(format, args...))
}

func init() {
	// XXX testing.Verbose() always returns false
	// if testing.Verbose() {
	// 	Log = logF
	// }
	Debug = logF
	Debug("Using Verbose Log")
}
