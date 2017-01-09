package cblib

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

func setupAddrs(paddr string, maddr string) {
	cb.CB_ADDR = paddr
	preIdx := strings.Index(paddr, "://")
	baseAddress := paddr[preIdx+3:]
	postIdx := strings.Index(baseAddress, ":")
	if postIdx != -1 {
		baseAddress = baseAddress[:postIdx]
	}
	// cb.CB_MSG_ADDR = baseAddress + ":1883"
	cb.CB_MSG_ADDR = maddr
}

func convertPermissionsNames(perms map[string]interface{}) map[string]interface{} {
	rval := map[string]interface{}{}
	for key, val := range perms {
		switch key {
		case "CodeServices":
			rval["services"] = val
		case "Collections":
			rval["collections"] = val
		case "DevicesList":
			rval["devices"] = val
		case "MsgHistory":
			rval["messagehistory"] = val
		case "Portals":
			rval["portals"] = val
		case "Push":
			rval["push"] = val
		case "Topics":
			rval["topics"] = val
		case "UsersList":
			rval["users"] = val
		default:
			rval[key] = "unknown"
		}
	}
	return rval
}
