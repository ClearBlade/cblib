package colutil

import (
	"fmt"

	"github.com/clearblade/cblib/diff"
	"github.com/clearblade/cblib/listutil"
)

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
	switch t {
	case "character varying", "varchar", "text":
		return "string"
	case "integer":
		return "int"
	case "real":
		return "float"
	case "bytea":
		return "blob"
	case "boolean":
		return "bool"
	case "double precision":
		return "double"
	case "timestamp without time zone":
		return "timestamp"
	default:
		return t
	}
}

// isValidColumnType checks whether a type string is a recognized app type or PostgreSQL type.
// This matches the PsqlType() function in the clearblade server (postgres/dbOps.go).
func isValidColumnType(t string) bool {
	switch t {
	// App types
	case "string", "int", "bigint", "float", "double", "blob", "uuid", "timestamp", "bool", "counter", "autoincrement":
		return true
	// PostgreSQL native types
	case "text", "integer", "real", "boolean", "double precision", "bytea", "bigserial",
		"character varying", "varchar", "timestamp without time zone":
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
