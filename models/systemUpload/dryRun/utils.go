package dryRun

import (
	"fmt"
	"strings"
)

func toStringArray(val interface{}) []string {
	result := make([]string, 0)

	arr := val.([]interface{})
	for _, item := range arr {
		result = append(result, item.(string))
	}

	return result
}

func writeDryRunSection(sb *strings.Builder, title, contents string) {
	sb.WriteString(fmt.Sprintf("-- %s --\n", title))
	sb.WriteString(contents)
	sb.WriteString("\n\n")
}
