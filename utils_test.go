package cblib

import (
	"fmt"
	"testing"
)

func TestBubbleSort_String(test *testing.T) {
	// input, truth to be tested against output
	rubric := [][][]string{
		{{"c", "a", "b"}, {"a", "b", "c"}},
		{{"1beta", "0alpha"}, {"0alpha", "1beta"}},
		{{"ngf", "35g", "][3", "KAH"}, {"35g", "KAH", "][3", "ngf"}}}

	err := executeBubbleSortTests(rubric)
	if err != nil {
		test.Error(err.Error())
	}
}

func TestBubbleSort_Empty(test *testing.T) {
	rubric := [][][]string{
		{{"", "", ""}, {"", "", ""}},
		{{}, {}}}

	err := executeBubbleSortTests(rubric)
	if err != nil {
		test.Error(err.Error())
	}
}

func TestBubbleSort_NilEntry(test *testing.T) {
	sortMe := []interface{}{"", "a", nil}
	truth := []interface{}{nil, "", "a"}
	sortByFunction(&sortMe, compareString)
	for i := 0; i < len(sortMe); i++ {
		if sortMe[i] == nil && truth[i] == nil {
			continue
		}
		if (sortMe[i] == nil && truth[i] != nil) || (sortMe[i] != nil && truth[i] == nil) {
			test.Errorf("Failed to sort input %d in %s", i, sortMe)
			return
		}
		if sortMe[i].(string) != truth[i].(string) {
			test.Errorf("Failed to sort input %d", i)
			return
		}
	}

}

func TestBubbleSort_NilArray(test *testing.T) {
	sortByFunction(nil, compareString)
}

func TestIsString_Normal(t *testing.T) {
	var rubric map[interface{}]bool = map[interface{}]bool{
		"hello": true,
		"":      true,
		nil:     false}
	for input, truth := range rubric {
		var output bool = isString(input)
		if output != truth {
			t.Errorf("Result: %s is a string: %t", input, output)
		}
	}
}

// Helper Test Method

func executeBubbleSortTests(inputOutput [][][]string) error {
	var sortMe []interface{}
	for t := 0; t < len(inputOutput); t++ {
		testIO := inputOutput[t]
		var input []string = testIO[0]
		var output []string = testIO[1]

		sortMe = make([]interface{}, len(input))
		for i := 0; i < len(sortMe); i++ {
			sortMe[i] = input[i]
		}

		sortByFunction(&sortMe, compareString)

		for i := 0; i < len(input); i++ {
			if sortMe[i].(string) != output[i] {
				return fmt.Errorf("Failed to sort input %d on test: %d", i, t)
			}
		}

	}
	return nil
}

// custom Compare function used for bubblesort
func compareString(slice *[]interface{}, i int, j int) bool {
	s1 := (*slice)[i]
	s2 := (*slice)[j]
	if s1 == nil {
		return true
	}
	if s2 == nil {
		return false
	}
	return s1.(string) < s2.(string)
}

func TestFindDiff_NoDefaultColumns(test *testing.T) {
	removeColName := "test2"
	local := []interface{}{
		map[string]interface{}{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
	}
	backend := []interface{}{
		map[string]interface{}{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		map[string]interface{}{
			"ColumnName": removeColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	removeDiff := findDiff(backend, local, columnExists([]string{}))
	if len(removeDiff) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(removeDiff))
	}
	if removeDiff[0].(map[string]interface{})["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, removeDiff[0].(map[string]interface{})["ColumnName"].(string))
	}
	addDiff := findDiff(local, backend, columnExists([]string{}))
	if len(addDiff) != 0 {
		test.Errorf("Expected to add 0 elements but got %d elements", len(addDiff))
	}
}

func TestFindDiff_WithDefaultColumns(test *testing.T) {
	removeColName := "test2"
	addColName := "test3"
	local := []interface{}{
		map[string]interface{}{
			"ColumnName": "user_id",
			"ColumnType": "string",
			"PK":         true,
		},
		map[string]interface{}{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		map[string]interface{}{
			"ColumnName": addColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	backend := []interface{}{
		map[string]interface{}{
			"ColumnName": "user_id",
			"ColumnType": "string",
			"PK":         true,
		},
		map[string]interface{}{
			"ColumnName": "creation_date",
			"ColumnType": "string",
			"PK":         false,
		},
		map[string]interface{}{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		map[string]interface{}{
			"ColumnName": removeColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	defaultColumns := []string{"user_id", "creation_date"}
	removeDiff := findDiff(backend, local, columnExists(defaultColumns))
	if len(removeDiff) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(removeDiff))
	}
	if removeDiff[0].(map[string]interface{})["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, removeDiff[0].(map[string]interface{})["ColumnName"].(string))
	}
	addDiff := findDiff(local, backend, columnExists(defaultColumns))
	if len(addDiff) != 1 {
		test.Errorf("Expected to add 1 element but got %d elements", len(addDiff))
	}
}

func Test_IsDefaultColumn(t *testing.T) {
	noDefaultColumns := isDefaultColumn([]string{}, "test")
	if noDefaultColumns {
		t.Errorf("Should return false when no default columns")
	}

	match := isDefaultColumn([]string{"one", "two"}, "two")
	if !match {
		t.Errorf("Should return as a match")
	}

	noMatch := isDefaultColumn([]string{"one", "two"}, "three")
	if noMatch {
		t.Errorf("Should not return as a match")
	}
}

func Test_DiffEdgeColumnsWithNoCustomColumns(t *testing.T) {
	backend := []interface{}{
		map[string]interface{}{"ColumnName": "edge_key", "ColumnType": "string", "PK": true}, map[string]interface{}{"ColumnName": "novi_system_key", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "system_key", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "system_secret", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "token", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "name", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "description", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "location", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "mac_address", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "public_addr", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "public_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "local_addr", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "local_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_tls_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_ws_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_wss_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_auth_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "broker_ws_auth_port", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "first_talked", "ColumnType": "bigint", "PK": false}, map[string]interface{}{"ColumnName": "last_talked", "ColumnType": "bigint", "PK": false}, map[string]interface{}{"ColumnName": "communication_style", "ColumnType": "int", "PK": false}, map[string]interface{}{"ColumnName": "last_seen_version", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "policy_name", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "resolver_func", "ColumnType": "string", "PK": false}, map[string]interface{}{"ColumnName": "sync_edge_tables", "ColumnType": "string", "PK": false},
	}
	local := []interface{}{}

	removeDiff := findDiff(backend, local, columnExists(DefaultEdgeColumns))
	if len(removeDiff) != 0 {
		t.Errorf("Expected to remove 0 elements but got %d elements", len(removeDiff))
	}

	addDiff := findDiff(local, backend, columnExists(DefaultEdgeColumns))
	if len(addDiff) != 0 {
		t.Errorf("Expected to add 0 elements but got %d elements", len(addDiff))
	}

}
