package cblib

import (
	//"fmt"
	cb "github.com/clearblade/Go-SDK"
	"math/rand"
	"reflect"
	"strings"
)

type compare func(sliceOfSystemResources *[]interface{}, i int, j int) bool

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

// Bubble sort, compare by map key
func sortByMapKey(arrayPointer *[]interface{}, sortKey string) {
	if arrayPointer == nil {
		return
	}
	array := *arrayPointer
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(array)-1; i++ {
			needToSwap := compareWithKey(sortKey, arrayPointer, i+1, i)
			if needToSwap {
				swap(arrayPointer, i, i+1)
				swapped = true
			}
		}
	}
}

// Bubble sort, compare by function
func sortByFunction(arrayPointer *[]interface{}, compareFn compare) {
	if arrayPointer == nil {
		return
	}
	array := *arrayPointer
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(array)-1; i++ {
			needToSwap := compareFn(arrayPointer, i+1, i)
			if needToSwap {
				swap(arrayPointer, i, i+1)
				swapped = true
			}
		}
	}
}

func swap(array *[]interface{}, i, j int) {
	tmp := (*array)[j]
	(*array)[j] = (*array)[i]
	(*array)[i] = tmp
}

func isString(input interface{}) bool {
	return input != nil && reflect.TypeOf(input).Name() == "string"
}

func compareWithKey(sortKey string, sliceOfCodeServices *[]interface{}, i, j int) bool {
	slice := *sliceOfCodeServices

	map1, castSuccess1 := slice[i].(map[string]interface{})
	map2, castSuccess2 := slice[j].(map[string]interface{})

	if !castSuccess1 || !castSuccess2 {
		return false
	}

	name1 := map1[sortKey]
	name2 := map2[sortKey]

	if !isString(name1) || !isString(name2) {
		return false
	}

	return name1.(string) < name2.(string)
}

func randSeq(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createFilePath(args ...string) string {
	return strings.Join(args, "/")
}
