package cblib

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

func setupAddrs(paddr string, maddr string) {
	cb.CB_ADDR = paddr

	postIdx := strings.Index(maddr, ":")
	if postIdx != -1 {
		cb.CB_MSG_ADDR = maddr
	} else {
		cb.CB_MSG_ADDR = maddr + ":1883" 
	}
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
