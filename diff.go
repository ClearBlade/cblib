package cblib

import (
	"fmt"
	"reflect"
	"strings"
)

type Stack struct {
	name      string
	stringRep string
	stack     []string
}

var (
	names             *Stack
	ignores           map[string][]string
	uniqueKeys        map[string]string
	suppressErrors    []int
	printedErrorCount int
	runtimeStack      []byte
)

func init() {
	runtimeStack = make([]byte, 1000000)
	printedErrorCount = 0
	suppressErrors = []int{0}
	names = NewStack("names")
	ignores = map[string][]string{
		"system.json":           []string{"platformURL"},
		"system.json:data":      []string{"appID", "collectionID"},
		"system.json:libraries": []string{"version"},
		"system.json:services":  []string{"current_version"},
		"users.json":            []string{"user_id", "creation_date"},
	}
	uniqueKeys = map[string]string{
		"system.json:data":        "name",
		"system.json:data:schema": "ColumnName",
		"system.json:libraries":   "name",
		"system.json:services":    "name",
		"system.json:timers":      "name",
		"system.json:triggers":    "name",
		"system.json:users":       "ColumnName",
		"users.json":              "email",
	}
}

func printErr(strFmt string, args ...interface{}) {
	if showErrors() {
		printedErrorCount++
		newArgs := append([]interface{}{names.stringRep}, args...)
		fmt.Printf("In %s: "+strFmt, newArgs...)
	}
}

func NewStack(name string) *Stack {
	return &Stack{
		name:      name,
		stringRep: "",
		stack:     make([]string, 0),
	}
}

func (s *Stack) push(item string) {
	s.stack = append(s.stack, item)
	s.stringRep = strings.Join(s.stack, ":")
}

func (s *Stack) top() (string, error) {
	rval := ""
	if len(s.stack) == 0 {
		return rval, fmt.Errorf("Attempt to get top of stack for empty stack %s", s.name)
	}
	return s.stack[len(s.stack)-1], nil
}

func (s *Stack) pop() (string, error) {
	rval := ""
	if len(s.stack) == 0 {
		return rval, fmt.Errorf("Attempt to pop empty stack %s", s.name)
	}
	rval, s.stack = s.stack[len(s.stack)-1], s.stack[:len(s.stack)-1]
	s.stringRep = strings.Join(s.stack, ":")
	return rval, nil
}

func diffSystemDotJSON(a, b map[string]interface{}) int {
	names.push("system.json")
	defer names.pop()
	diffMap(a, b)
	fmt.Printf("%d Total Errors\n", printedErrorCount)
	return printedErrorCount
}

func diffUsersDotJSON(a, b []interface{}) int {
	names.push("users.json")
	defer names.pop()
	return diffSlice(a, b)
}

func diffUnknownTypes(key string, a, b interface{}) int {
	if !sameTypes(a, b) {
		return 1
	}
	if outerType(a) == "map" {
		if key != "" {
			names.push(key)
			defer names.pop()
		}
		return diffMap(a.(map[string]interface{}), b.(map[string]interface{}))
	} else if outerType(a) == "slice" {
		if key != "" {
			names.push(key)
			defer names.pop()
		}
		return diffSlice(a.([]interface{}), b.([]interface{}))
	} else if a == b {
		return 0
	}
	printErr("Found differing values: %v != %v\n", a, b)
	return 1
}

func diffMap(a, b map[string]interface{}) int {
	totalErrors := 0
	checkedKeys := map[string]bool{}
	for aKey, aVal := range a {
		checkedKeys[aKey] = true
		if shouldIgnore(aKey) {
			continue
		}
		if bVal, ok := b[aKey]; ok {
			totalErrors += diffUnknownTypes(aKey, aVal, bVal)
		} else {
			totalErrors++
			printErr("Item %s in first map missing in second map\n", aKey)
		}
	}
	for bKey, _ := range b {
		_, ok := checkedKeys[bKey]
		if shouldIgnore(bKey) || ok {
			continue
		}
		if _, ok := a[bKey]; !ok {
			printErr("Item %s in second map missing in first map\n", bKey)
			totalErrors++
		}
	}
	return totalErrors
}

func diffSlice(a, b []interface{}) int {
	if len(a) > 0 {
		if reflect.TypeOf(a[0]).String() == "map[string]interface {}" {
			pushErrorContext()
			defer popErrorContext()
		}
	}
	// Assumption
	totalErrors := 0
	if !sameTypes(a, b) {
		return 1
	}
	if len(a) != len(b) {
		printErr("Slices are of different length: %d != %d\n", len(a), len(b))
		totalErrors++
	}
	totalErrors += diffTwoSlices(a, b)
	return totalErrors
}

func getUniqueKeyInfo(valSlice []interface{}) (string, bool) {
	if len(valSlice) == 0 {
		return "", false
	}
	oneVal := valSlice[0]
	if outerType(oneVal) != "map" {
		return "", false
	}
	if uniqueKey, haveOne := uniqueKeys[names.stringRep]; haveOne {
		return uniqueKey, true
	}
	return "", false
}

func findMatchInOtherSlice(b []interface{}, uniqueKey string, uniqueVal interface{}) map[string]interface{} {
	for _, bValIF := range b {
		bVal := bValIF.(map[string]interface{})
		if valForKeyIF, ok := bVal[uniqueKey]; ok {
			if uniqueVal == valForKeyIF {
				return bVal
			}
		}
	}
	return nil
}

func valInSlice(val interface{}, slice []interface{}) bool {
	for _, sliceVal := range slice {
		if sliceVal == val {
			return true
		}
	}
	return false
}

func diffKeyedSlices(a, b []interface{}, uniqueKey string) int {
	myErrors := 0
	seenKeyVals := []interface{}{}
	for _, aValIF := range a {
		aVal := aValIF.(map[string]interface{})
		if valForKey, ok := aVal[uniqueKey]; ok {
			seenKeyVals = append(seenKeyVals, valForKey)
			bVal := findMatchInOtherSlice(b, uniqueKey, valForKey)
			if bVal == nil {
				myErrors++
				printErr("Item %s:%v not found in other system\n", uniqueKey, valForKey)
			}
			pushErrorContext()
			defer popErrorContext()
			myErrors += diffMap(aVal, bVal)
		} else {
			printErr("Item supposedly with uniqueKey doesn't have one: %#v\n", aVal)
			return -1
		}
	}
	//  Now, we're just finding entries in b that aren't in a
	for _, bValIF := range b {
		bVal := bValIF.(map[string]interface{})
		if valForKey, ok := bVal[uniqueKey]; ok {
			if !valInSlice(valForKey, seenKeyVals) {
				printErr("Key %s with value %v not found in first system\n",
					uniqueKey, valForKey)
				myErrors++
			}
		} else {
			printErr("Item supposedly with uniqueKey doesn't have one: %#v\n", bVal)
			return -1
		}
	}
	return myErrors
}

func diffTwoSlices(a, b []interface{}) int {
	uniqueKey, useUniqueKey := getUniqueKeyInfo(a)
	if useUniqueKey {
		if myErrors := diffKeyedSlices(a, b, uniqueKey); myErrors != -1 {
			return myErrors
		}
	}
	return diffUnkeyedSlices(a, b)
}

func diffUnkeyedSlices(a, b []interface{}) int {
	totalErrors := 0
	printsBefore := printedErrorCount
	for _, aVal := range a {
		found := false
		for _, bVal := range b {
			blockErrors()
			errCount := diffUnknownTypes("", aVal, bVal)
			unblockErrors()
			if errCount == 0 {
				found = true
				break
			}
		}
		if !found {
			totalErrors++
			if printsBefore == printedErrorCount {
				printErr("Could not find item %#v in other slice\n", aVal)
			}
		}
	}

	printsBefore = printedErrorCount
	for _, bVal := range b {
		found := false
		for _, aVal := range a {
			blockErrors()
			errCount := diffUnknownTypes("", bVal, aVal)
			unblockErrors()
			if errCount == 0 {
				found = true
				break
			}
		}
		if !found {
			totalErrors++
			if printsBefore == printedErrorCount {
				printErr("Could not find item %#v in other slice\n", bVal)
			}
		}
	}
	return totalErrors
}

func shouldIgnore(key string) bool {
	if keyList, ok := ignores[names.stringRep]; ok {
		for _, ignoreKey := range keyList {
			if ignoreKey == key {
				return true
			}
		}
	}
	return false
}

func sameTypes(a, b interface{}) bool {
	typeA := reflect.TypeOf(a).String()
	typeB := reflect.TypeOf(b).String()
	rval := typeA == typeB
	if !rval {
		printErr("Encountered two different types: %s != %s\n", typeA, typeB)
	}
	return rval
}

func outerType(a interface{}) string {
	return reflect.ValueOf(a).Kind().String()
}

func showErrors() bool {
	return suppressErrors[len(suppressErrors)-1] == 0
}

func pushErrorContext() {
	suppressErrors = append(suppressErrors, 0)
}

func popErrorContext() {
	suppressErrors = suppressErrors[:len(suppressErrors)-1]
}

func blockErrors() {
	suppressErrors[len(suppressErrors)-1] = suppressErrors[len(suppressErrors)-1] + 1
}

func unblockErrors() {
	suppressErrors[len(suppressErrors)-1] = suppressErrors[len(suppressErrors)-1] - 1
}
