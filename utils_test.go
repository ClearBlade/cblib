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
	local := []map[string]interface{}{
		{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
	}
	backend := []map[string]interface{}{
		{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		{
			"ColumnName": removeColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	removeDiff := findDiff(backend, local, []string{})
	if len(removeDiff) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(removeDiff))
	}
	if removeDiff[0]["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, removeDiff[0]["ColumnName"].(string))
	}
	addDiff := findDiff(local, backend, []string{})
	if len(addDiff) != 0 {
		test.Errorf("Expected to add 0 elements but got %d elements", len(addDiff))
	}
}

func TestFindDiff_WithDefaultColumns(test *testing.T) {
	removeColName := "test2"
	addColName := "test3"
	local := []map[string]interface{}{
		{
			"ColumnName": "user_id",
			"ColumnType": "string",
			"PK":         true,
		},
		{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		{
			"ColumnName": addColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	backend := []map[string]interface{}{
		{
			"ColumnName": "user_id",
			"ColumnType": "string",
			"PK":         true,
		},
		{
			"ColumnName": "creation_date",
			"ColumnType": "string",
			"PK":         false,
		},
		{
			"ColumnName": "test",
			"ColumnType": "string",
			"PK":         false,
		},
		{
			"ColumnName": removeColName,
			"ColumnType": "string",
			"PK":         false,
		},
	}
	defaultColumns := []string{"user_id", "creation_date"}
	removeDiff := findDiff(backend, local, defaultColumns)
	if len(removeDiff) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(removeDiff))
	}
	if removeDiff[0]["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, removeDiff[0]["ColumnName"].(string))
	}
	addDiff := findDiff(local, backend, defaultColumns)
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
