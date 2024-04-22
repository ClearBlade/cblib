package colutil

import (
	"github.com/clearblade/cblib/diff"
	"github.com/clearblade/cblib/listutil"
)

func GetDiffForColumnsWithDynamicListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []interface{}) *diff.UnsafeDiff {
	return listutil.CompareListsAndFilter(localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a interface{}) bool {
		return a.(map[string]interface{})["UserDefined"].(bool)
	})
}

func GetDiffForColumnsWithStaticListOfDefaultColumns(localSchemaInterfaces, backendSchemaInterfaces []interface{}, defaultColumns []string) *diff.UnsafeDiff {
	return listutil.CompareListsAndFilter(localSchemaInterfaces, backendSchemaInterfaces, columnExists, func(a interface{}) bool {
		return !isDefaultColumn(defaultColumns, a.(map[string]interface{})["ColumnName"].(string))
	})
}

func columnExists(colA interface{}, colB interface{}) bool {
	if colA.(map[string]interface{})["ColumnName"].(string) == colB.(map[string]interface{})["ColumnName"].(string) && colA.(map[string]interface{})["ColumnType"].(string) == colB.(map[string]interface{})["ColumnType"].(string) {
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
