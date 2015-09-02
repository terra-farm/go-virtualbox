package virtualbox

func SetProperty(key, val string) error {
	return vbm("setproperty", key, val)
}

func GuestPropertyGet(uuid, key string) (string, error) {
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

func GuestPropertyWait(uuid, key string) (string, error) {
	val, err := vbmOut("guestproperty", "wait", uuid, key)
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

func GuestPropertySet(uuid, key, val string) error {
	return vbm("guestproperty", "set", uuid, key, val)
}

func GuestPropertyDel(uuid, key string) error {
	return vbm("guestproperty", "delete", uuid, key)
}

func GuestPropertyEnumerate(uuid, key string) (string, error) {
	return vbmOut("guestproperty", "enumerate", uuid, key)
}
