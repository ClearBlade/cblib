package colutil

import (
	"fmt"
	"regexp"

	"github.com/clearblade/cblib/diff"
	"github.com/clearblade/cblib/listutil"
)

// typeModifierRe matches PostgreSQL type modifiers like (128), (10,2), (6).
var typeModifierRe = regexp.MustCompile(`\s*\([^)]*\)`)

func GetDiffForColumnsWithDynamicListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []map[string]interface{}) (*diff.UnsafeDiff[map[string]interface{}], error) {
	if err := validateColumnTypes(localSchemaInterfaces); err != nil {
		return nil, err
	}
	return listutil.CompareListsAndFilter[map[string]interface{}](localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a map[string]interface{}) bool {
		// if the UserDefined key exists, that means the column exists on the backend.
		userDefined, ok := a["UserDefined"].(bool)
		if ok {
			return userDefined
		}
		// if the UserDefined key does not exist, that means the column does not exist yet on the backend and should be added.
		return true
	}), nil
}

func GetDiffForColumnsWithStaticListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []map[string]interface{}, defaultColumns []string) (*diff.UnsafeDiff[map[string]interface{}], error) {
	if err := validateColumnTypes(localSchemaInterfaces); err != nil {
		return nil, err
	}
	return listutil.CompareListsAndFilter(localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a map[string]interface{}) bool {
		return !isDefaultColumn(defaultColumns, a["ColumnName"].(string))
	}), nil
}

func columnExists(colA, colB map[string]interface{}) bool {
	return colA["ColumnName"].(string) == colB["ColumnName"].(string) &&
		normalizeType(colA["ColumnType"].(string)) == normalizeType(colB["ColumnType"].(string))
}

func isDefaultColumn(defaultColumns []string, colName string) bool {
	for i := 0; i < len(defaultColumns); i++ {
		if defaultColumns[i] == colName {
			return true
		}
	}
	return false
}

// normalizeType maps PostgreSQL native type names to app-level type names.
// This matches the Typed() function in the clearblade server (postgres/dbOps.go).
func normalizeType(t string) string {
	t = typeModifierRe.ReplaceAllString(t, "")
	switch t {
	// PostgreSQL types with ClearBlade app type equivalents
	case "character varying", "varchar", "text":
		return "string"
	case "integer", "int4":
		return "int"
	case "int8":
		return "bigint"
	case "real", "float4":
		return "float"
	case "double precision", "float8":
		return "double"
	case "bytea":
		return "blob"
	case "boolean":
		return "bool"
	case "timestamp without time zone":
		return "timestamp"
	default:
		return t
	}
}

// isValidColumnType checks whether a type string is a valid ClearBlade app type.
// Only the 10 app types (string, int, bool, timestamp, float, bigint, double,
// jsonb, blob, uuid) plus the internal types counter and autoincrement are accepted.
func isValidColumnType(t string) bool {
	switch t {
	case "string", "int", "bool", "timestamp", "float", "bigint", "double",
		"jsonb", "blob", "uuid",
		"counter", "autoincrement":
		return true
	default:
		return false
	}
}

func validateColumnTypes(columns []map[string]interface{}) error {
	for _, col := range columns {
		colType, ok := col["ColumnType"].(string)
		if !ok {
			return fmt.Errorf("column %q is missing ColumnType", col["ColumnName"])
		}
		if !isValidColumnType(colType) {
			return fmt.Errorf("invalid column type %q for column %q", colType, col["ColumnName"])
		}
	}
	return nil
}
