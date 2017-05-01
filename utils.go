package cblib

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
	"strings"
	"fmt"
	"reflect"
)

type compare func(sliceOfCodeServices *[]interface{}, i int, j int) bool 

func setupAddrs(paddr string, maddr string) {
	cb.CB_ADDR = paddr

	preIdx := strings.Index(paddr, "://")
	if maddr == "" {
		if preIdx != -1 {
			maddr = paddr[preIdx+3:]
		} else {
			maddr = paddr
		}
	}
	postIdx := strings.Index(maddr, ":")
	if postIdx != -1 {
		cb.CB_MSG_ADDR = maddr[:postIdx] + ":1883"
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

func bubbleSort(arrayzor *[]interface{}, fn compare) {
	array := *arrayzor
	swapped := true;
	for swapped {
		swapped = false
		for i := 0; i < len(array) - 1; i++ {
			needToSwap := fn(arrayzor, i+1, i)
			if  needToSwap {
				swap(arrayzor, i, i + 1)
				swapped = true
			}
		}
	}
}

func swap(arrayzor *[]interface{}, i, j int) {
	tmp := (*arrayzor)[j]
	(*arrayzor)[j] = (*arrayzor)[i]
	(*arrayzor)[i] = tmp
}

func isString(input interface{}) bool {
	return input != nil && reflect.TypeOf(input).Name() == "string"
}
