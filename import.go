package virtualbox

// ImportOV imports ova or ovf from the given path
func ImportOV(path string) error {
	_, _, err := Manage().run("import", path)
	return err
}
