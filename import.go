package virtualbox

//ImportOV imports ova or ovf from the given path
func ImportOV(path string) error {
	return vbm("import", path)
}
