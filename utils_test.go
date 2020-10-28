package cblib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessURLsSucceeds(t *testing.T) {
	tests := []struct {
		platformURL          string
		messagingURL         string
		expectedPlatformURL  string
		expectedMessagingURL string
	}{

		// canonical
		{
			"https://platform.clearblade.com", "platform.clearblade.com",
			"https://platform.clearblade.com", "platform.clearblade.com:1883",
		},

		// platform and messaging override port
		{
			"https://platform.clearblade.com:8080", "platform.clearblade.com:8883",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:8883",
		},

		// platform has trailing slash
		{
			"https://platform.clearblade.com:8080/", "platform.clearblade.com",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:1883",
		},

		// does not specify messaging
		{
			"https://platform.clearblade.com:8080/", "",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:1883",
		},
	}

	for _, tt := range tests {
		platformURL, messagingURL, err := processURLs(tt.platformURL, tt.messagingURL)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.expectedPlatformURL, platformURL)
		assert.Equal(t, tt.expectedMessagingURL, messagingURL)
	}
}

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
	diff := getDiffForColumns(local, backend, []string{})
	if len(diff.remove) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(diff.remove))
	}
	if diff.remove[0].(map[string]interface{})["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, diff.remove[0].(map[string]interface{})["ColumnName"].(string))
	}
	if len(diff.add) != 0 {
		test.Errorf("Expected to add 0 elements but got %d elements", len(diff.add))
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
	diff := getDiffForColumns(local, backend, defaultColumns)
	if len(diff.remove) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(diff.remove))
	}
	if diff.remove[0].(map[string]interface{})["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, diff.remove[0].(map[string]interface{})["ColumnName"].(string))
	}
	if len(diff.add) != 1 {
		test.Errorf("Expected to add 1 element but got %d elements", len(diff.add))
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

	diff := getDiffForColumns(local, backend, DefaultEdgeColumns)

	if len(diff.remove) != 0 {
		t.Errorf("Expected to remove 0 elements but got %d elements", len(diff.remove))
	}

	if len(diff.add) != 0 {
		t.Errorf("Expected to add 0 elements but got %d elements", len(diff.add))
	}

}

func Test_DeletingAllUserRoles(t *testing.T) {
	backend := []string{"Authenticated", "Administrator"}
	local := []interface{}{}

	roleDiff := diffRoles(local, convertStringSliceToInterfaceSlice(backend))

	if len(roleDiff.remove) != 2 {
		t.Errorf("Expected to remove 2 elements but got %d elements", len(roleDiff.remove))
	}

	if len(roleDiff.add) != 0 {
		t.Errorf("Expected to add 0 elements but got %d elements", len(roleDiff.add))
	}

}

func TestFilterSliceSuceeds(t *testing.T) {

	tests := []struct {
		items     []interface{}
		predicate func(interface{}) bool
		expected  []interface{}
	}{
		// filter even
		{
			[]interface{}{1, 2, 3, 4, 5},
			func(item interface{}) bool { return item.(int)%2 == 0 },
			[]interface{}{2, 4},
		},

		// filter odd
		{
			[]interface{}{1, 2, 3, 4, 5},
			func(item interface{}) bool { return item.(int)%2 != 0 },
			[]interface{}{1, 3, 5},
		},
	}

	for _, tt := range tests {
		filtered := FilterSlice(tt.items, tt.predicate)
		assert.Equal(t, tt.expected, filtered)
	}
}

func TestDiffSliceUsingSucceeds(t *testing.T) {

	tests := []struct {
		a    []interface{}
		b    []interface{}
		diff []interface{}
	}{
		{[]interface{}{1, 2, 3}, []interface{}{}, []interface{}{1, 2, 3}},
		{[]interface{}{1, 2, 3}, []interface{}{1}, []interface{}{2, 3}},
		{[]interface{}{1, 2, 3}, []interface{}{1, 2}, []interface{}{3}},
		{[]interface{}{1, 2, 3}, []interface{}{1, 2, 3}, []interface{}{}},
		{[]interface{}{2, 3}, []interface{}{1, 2, 3}, []interface{}{}},
		{[]interface{}{3}, []interface{}{1, 2, 3}, []interface{}{}},
		{[]interface{}{}, []interface{}{1, 2, 3}, []interface{}{}},
	}

	for _, tt := range tests {
		diff := DiffSliceUsing(tt.a, tt.b, func(a, b interface{}) bool { return a == b })
		assert.Equal(t, tt.diff, diff)
	}
}
