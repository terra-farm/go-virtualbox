package virtualbox

import "strconv"

func modifyMacAddress(uuid string, solt int, val string) error {
	return vbm("modifyvm", uuid, "--macaddress"+strconv.Itoa(solt), val)
}
