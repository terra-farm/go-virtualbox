package virtualbox

import (
	"os"
	"regexp"
	"testing"
)

func TestGetRegexp(t *testing.T) {
	var str = os.Getenv("TEST_STRING")
	if len(str) <= 0 {
		str = "Value: foo"
	}
	var re *regexp.Regexp
	var reStr = os.Getenv("TEST_GETREGEXP")
	if len(reStr) <= 0 {
		re = getRegexp
	} else {
		re = regexp.MustCompile(reStr)
	}
	var match = re.FindStringSubmatch(str)
	t.Log("match:", match)
	if len(match) != 2 {
		t.Fatal("No match")
	}
	if match[0] != str {
		t.Fatal("No global match")
	}
	if match[1] != "foo" {
		t.Fatal("No value match")
	}
}
