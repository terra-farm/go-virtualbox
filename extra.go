package virtualbox

// SetExtra sets extra data. Name could be "global"|<uuid>|<vmname>
func SetExtra(name, key, val string) error {
	return vbm("setextradata", name, key, val)
}

// DelExtraData deletes extra data. Name could be "global"|<uuid>|<vmname>
func DelExtra(name, key string) error {
	return vbm("setextradata", name, key)
}

// DelExtraData deletes extra data. Name could be "global"|<uuid>|<vmname>
func GetExtra(name, key string) (string, error) {
	val, err := vbmOut("getextradata", name, key)

	if err == nil {
		res := reVMProperty.FindStringSubmatch(val)
		if res != nil {
			val = res[1]
		}
	}
	return val, err
}
