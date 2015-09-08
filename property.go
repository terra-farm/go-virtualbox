package virtualbox

import (
	"bufio"
	_ "fmt"
	"strings"
)

func SetProperty(key, val string) error {
	return vbm("setproperty", key, val)
}

func guestPropertyGet(uuid, key string) (string, error) {
	val, err := vbmOut("guestproperty", "get", uuid, key)
	if err == nil {
		res := reVMProperty.FindStringSubmatch(val)
		if res != nil {
			val = res[1]
		} else {
			val = ""
		}
	}

	return val, err
}

func guestPropertyWait(uuid, key string, timeout int) (string, error) {
	val, err := vbmOut("guestproperty", "wait", uuid, key, "--timeout", string(timeout))
	if err == nil {
		res := reVMProperty.FindStringSubmatch(val)
		if res != nil {
			val = res[1]
		} else {
			val = ""
		}
	}
	return val, err
}

func guestPropertySet(uuid, key, val string) error {
	return vbm("guestproperty", "set", uuid, key, val)
}

func guestPropertyDel(uuid, key string) error {
	return vbm("guestproperty", "delete", uuid, key)
}

func guestPropertyEnumerate(uuid string) (map[string]string, error) {
	var vals map[string]string
	out, err := vbmOut("guestproperty", "enumerate", uuid)
	//val, err := vbmOut("guestproperty", "wait", uuid, key, "--timeout", string(timeout))
	if err == nil {
		vals = make(map[string]string)
		s := bufio.NewScanner(strings.NewReader(out))
		for s.Scan() {
			res := reVMPropertyEnumerate.FindStringSubmatch(s.Text())
			if res == nil {
				continue
			}
			/*if true {
				for v, k := range res {
					fmt.Printf("Index: %d %s\n", v, k)
				}
			}*/

			if len(res) != 4 {
				continue
			}
			vals[res[1]] = res[2]
		}
	}
	return vals, err
}
