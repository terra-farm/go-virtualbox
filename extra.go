package virtualbox

// SetExtra sets extra data. Name could be "global"|<uuid>|<vmname>
func SetExtra(name, key, val string) error {
	_, _, err := Manage().run("setextradata", name, key, val)
	return err
}

// DelExtra deletes extra data. Name could be "global"|<uuid>|<vmname>
func DelExtra(name, key string) error {
	_, _, err := Manage().run("setextradata", name, key)
	return err
}
