package colutil

import (
	"testing"
)

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
	backend := []map[string]interface{}{
		map[string]interface{}{"ColumnName": "edge_key", "ColumnType": "string", "PK": true, "UserDefined": false}, map[string]interface{}{"ColumnName": "novi_system_key", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "system_key", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "system_secret", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "token", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "name", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "description", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "location", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "mac_address", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "public_addr", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "public_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "local_addr", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "local_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_tls_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_ws_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_wss_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_auth_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "broker_ws_auth_port", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "first_talked", "ColumnType": "bigint", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "last_talked", "ColumnType": "bigint", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "communication_style", "ColumnType": "int", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "last_seen_version", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "policy_name", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "resolver_func", "ColumnType": "string", "PK": false, "UserDefined": false}, map[string]interface{}{"ColumnName": "sync_edge_tables", "ColumnType": "string", "PK": false, "UserDefined": false},
	}
	local := []map[string]interface{}{}

	diff := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)

	if len(diff.Removed) != 0 {
		t.Errorf("Expected to remove 0 elements but got %d elements", len(diff.Removed))
	}

	if len(diff.Added) != 0 {
		t.Errorf("Expected to add 0 elements but got %d elements", len(diff.Added))
	}

}

func TestFindDiff_WithDefaultColumns(test *testing.T) {
	removeColName := "test2"
	addColName := "test3"
	local := []map[string]interface{}{
		{
			"ColumnName":  "user_id",
			"ColumnType":  "string",
			"PK":          true,
			"UserDefined": false,
		},
		{
			"ColumnName":  "test",
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
		{
			"ColumnName":  addColName,
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
	}
	backend := []map[string]interface{}{
		{
			"ColumnName":  "user_id",
			"ColumnType":  "string",
			"PK":          true,
			"UserDefined": false,
		},
		{
			"ColumnName":  "creation_date",
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": false,
		},
		{
			"ColumnName":  "test",
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
		{
			"ColumnName":  removeColName,
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
	}
	diff := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if len(diff.Removed) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(diff.Removed))
	}
	if diff.Removed[0]["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, diff.Removed[0]["ColumnName"].(string))
	}
	if len(diff.Added) != 1 {
		test.Errorf("Expected to add 1 element but got %d elements", len(diff.Added))
	}
}

func TestFindDiff_NoDefaultColumns(test *testing.T) {
	removeColName := "test2"
	local := []map[string]interface{}{
		{
			"ColumnName":  "test",
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
	}
	backend := []map[string]interface{}{
		{
			"ColumnName":  "test",
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
		{
			"ColumnName":  removeColName,
			"ColumnType":  "string",
			"PK":          false,
			"UserDefined": true,
		},
	}
	diff := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if len(diff.Removed) != 1 {
		test.Errorf("Expected to remove 1 element but got %d elements", len(diff.Removed))
	}
	if diff.Removed[0]["ColumnName"].(string) != removeColName {
		test.Errorf("Expected column name to be '%s' but got '%s'\n", removeColName, diff.Removed[0]["ColumnName"].(string))
	}
	if len(diff.Added) != 0 {
		test.Errorf("Expected to add 0 elements but got %d elements", len(diff.Added))
	}
}
