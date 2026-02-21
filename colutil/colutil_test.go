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

	diff, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

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
	diff, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err != nil {
		test.Fatalf("Unexpected error: %s", err)
	}
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
	diff, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err != nil {
		test.Fatalf("Unexpected error: %s", err)
	}
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

func Test_NormalizeType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"boolean", "bool"},
		{"bool", "bool"},
		{"integer", "int"},
		{"int", "int"},
		{"text", "string"},
		{"string", "string"},
		{"character varying", "string"},
		{"varchar", "string"},
		{"real", "float"},
		{"float", "float"},
		{"double precision", "double"},
		{"double", "double"},
		{"bytea", "blob"},
		{"blob", "blob"},
		{"timestamp without time zone", "timestamp"},
		{"timestamp", "timestamp"},
		{"bigint", "bigint"},
		{"uuid", "uuid"},
	}

	for _, tc := range tests {
		result := normalizeType(tc.input)
		if result != tc.expected {
			t.Errorf("normalizeType(%q) = %q, expected %q", tc.input, result, tc.expected)
		}
	}
}

func Test_IsValidColumnType(t *testing.T) {
	validTypes := []string{
		// App types
		"string", "int", "bigint", "float", "double", "blob", "uuid", "timestamp", "bool", "counter", "autoincrement",
		// PostgreSQL native types
		"int8",
		"bigserial", "serial8",
		"bit", "bit varying", "varbit",
		"boolean",
		"box",
		"bytea",
		"character", "char", "character varying", "varchar",
		"cidr",
		"circle",
		"date",
		"double precision", "float8", "float4",
		"inet",
		"integer", "int4", "int2",
		"interval",
		"json", "jsonb",
		"line", "lseg",
		"macaddr", "macaddr8",
		"money",
		"numeric", "decimal",
		"path",
		"pg_lsn", "pg_snapshot",
		"point", "polygon",
		"real",
		"smallint",
		"smallserial", "serial2",
		"serial", "serial4",
		"text",
		"time", "time without time zone", "time with time zone", "timetz",
		"timestamp without time zone", "timestamp with time zone", "timestamptz",
		"tsquery", "tsvector",
		"txid_snapshot",
		"xml",
	}
	for _, typ := range validTypes {
		if !isValidColumnType(typ) {
			t.Errorf("isValidColumnType(%q) = false, expected true", typ)
		}
	}

	invalidTypes := []string{"str", "bogus", "INT", "Bool", ""}
	for _, typ := range invalidTypes {
		if isValidColumnType(typ) {
			t.Errorf("isValidColumnType(%q) = true, expected false", typ)
		}
	}
}

func Test_ValidateColumnTypes_RejectsInvalid(t *testing.T) {
	columns := []map[string]interface{}{
		{"ColumnName": "good_col", "ColumnType": "string"},
		{"ColumnName": "bad_col", "ColumnType": "str"},
	}
	err := validateColumnTypes(columns)
	if err == nil {
		t.Fatal("Expected error for invalid column type 'str', got nil")
	}
}

func Test_ValidateColumnTypes_AcceptsValid(t *testing.T) {
	columns := []map[string]interface{}{
		{"ColumnName": "col1", "ColumnType": "string"},
		{"ColumnName": "col2", "ColumnType": "bool"},
		{"ColumnName": "col3", "ColumnType": "boolean"},
	}
	err := validateColumnTypes(columns)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
}

func Test_ColumnExists_NormalizesTypes(t *testing.T) {
	// "boolean" (PostgreSQL) should match "bool" (app type)
	colA := map[string]interface{}{"ColumnName": "active", "ColumnType": "boolean"}
	colB := map[string]interface{}{"ColumnName": "active", "ColumnType": "bool"}
	if !columnExists(colA, colB) {
		t.Error("columnExists should treat 'boolean' and 'bool' as the same type")
	}

	// "text" (PostgreSQL) should match "string" (app type)
	colC := map[string]interface{}{"ColumnName": "name", "ColumnType": "text"}
	colD := map[string]interface{}{"ColumnName": "name", "ColumnType": "string"}
	if !columnExists(colC, colD) {
		t.Error("columnExists should treat 'text' and 'string' as the same type")
	}

	// Different names should not match even with same type
	colE := map[string]interface{}{"ColumnName": "col1", "ColumnType": "string"}
	colF := map[string]interface{}{"ColumnName": "col2", "ColumnType": "string"}
	if columnExists(colE, colF) {
		t.Error("columnExists should not match columns with different names")
	}
}

func Test_DiffWithPsqlTypeAliases_NoFalseRemoval(t *testing.T) {
	// Local schema uses PostgreSQL type "boolean", backend has app type "bool".
	// These should be treated as the same column — no removal or addition.
	local := []map[string]interface{}{
		{"ColumnName": "active", "ColumnType": "boolean", "UserDefined": true},
	}
	backend := []map[string]interface{}{
		{"ColumnName": "active", "ColumnType": "bool", "UserDefined": true},
	}
	diff, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if len(diff.Removed) != 0 {
		t.Errorf("Expected 0 removals but got %d — type alias caused false removal", len(diff.Removed))
	}
	if len(diff.Added) != 0 {
		t.Errorf("Expected 0 additions but got %d — type alias caused false addition", len(diff.Added))
	}
}

func Test_TypeModifiers_AcceptedAndNormalized(t *testing.T) {
	// varchar(128) should be valid and match "string" on the backend
	local := []map[string]interface{}{
		{"ColumnName": "name", "ColumnType": "varchar(128)", "UserDefined": true},
		{"ColumnName": "score", "ColumnType": "numeric(10,2)", "UserDefined": true},
		{"ColumnName": "created", "ColumnType": "timestamp(6) without time zone", "UserDefined": true},
	}
	backend := []map[string]interface{}{
		{"ColumnName": "name", "ColumnType": "string", "UserDefined": true},
		{"ColumnName": "score", "ColumnType": "numeric", "UserDefined": true},
		{"ColumnName": "created", "ColumnType": "timestamp", "UserDefined": true},
	}
	diff, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err != nil {
		t.Fatalf("Unexpected error for types with modifiers: %s", err)
	}
	if len(diff.Removed) != 0 {
		t.Errorf("Expected 0 removals but got %d — type modifier caused false removal", len(diff.Removed))
	}
	if len(diff.Added) != 0 {
		t.Errorf("Expected 0 additions but got %d — type modifier caused false addition", len(diff.Added))
	}
}

func Test_DiffWithInvalidType_ReturnsError(t *testing.T) {
	local := []map[string]interface{}{
		{"ColumnName": "bad_col", "ColumnType": "str", "UserDefined": true},
	}
	backend := []map[string]interface{}{
		{"ColumnName": "bad_col", "ColumnType": "string", "UserDefined": true},
	}
	_, err := GetDiffForColumnsWithDynamicListOfDefaultColumns(local, backend)
	if err == nil {
		t.Fatal("Expected error for invalid column type 'str', got nil")
	}
}
