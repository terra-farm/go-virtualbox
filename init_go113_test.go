// +build go1.13

package virtualbox

import (
	"flag"
	"testing"
)

func init() {
	testing.Init()
	flag.Parse()
	Debug = LogF
	Debug("Using Verbose Log")
	Debug("testing.Verbose=%v", testing.Verbose())
}
