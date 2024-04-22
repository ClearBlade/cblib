package colutil

import (
	"github.com/clearblade/cblib/diff"
	"github.com/clearblade/cblib/listutil"
)

func GetDiffForColumnsWithDynamicListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []map[string]interface{}) *diff.UnsafeDiff[map[string]interface{}] {
	return listutil.CompareListsAndFilter[map[string]interface{}](localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a map[string]interface{}) bool {
		return a["UserDefined"].(bool)
	})
}

func GetDiffForColumnsWithStaticListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []map[string]interface{}, defaultColumns []string) *diff.UnsafeDiff[map[string]interface{}] {
	return listutil.CompareListsAndFilter(localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a map[string]interface{}) bool {
		return !isDefaultColumn(defaultColumns, a["ColumnName"].(string))
	})
}

func columnExists(colA, colB map[string]interface{}) bool {
	if colA["ColumnName"].(string) == colB["ColumnName"].(string) && colA["ColumnType"].(string) == colB["ColumnType"].(string) {
		return true
	}
	return false
}

func isDefaultColumn(defaultColumns []string, colName string) bool {
	for i := 0; i < len(defaultColumns); i++ {
		if defaultColumns[i] == colName {
			return true
		}
	}
	return false
}
