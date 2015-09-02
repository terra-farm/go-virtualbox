package virtualbox

func ModifyMacAddress(uuid string, solt int, val string) error {
	return vbm("modifyvm", uuid, "--macaddress", string(solt), val)
}
