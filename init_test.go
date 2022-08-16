//go:build !go1.13
// +build !go1.13

package virtualbox

import (
	"testing"
)

func init() {
	Debug = LogF
	Debug("Using Verbose Log")
	Debug("testing.Verbose=%v", testing.Verbose())
}
