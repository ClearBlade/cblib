package dryRun

import (
	"strings"
)

/**
 * Generic printing for dry run section with a schema.
 * This section has a create, update, schema column add, and schema column delete
 */
type schemaSection struct {
	title           string
	creates         []string
	updates         []string
	columnsToAdd    []string
	columnsToDelete []string
}

func (s *schemaSection) Title() string {
	return s.title
}

func (s *schemaSection) HasChanges() bool {
	return len(s.updates)+len(s.creates)+len(s.columnsToAdd)+len(s.columnsToDelete) > 0
}

func (s *schemaSection) String() string {
	sb := strings.Builder{}

	// Write the creates and updates first
	sb.WriteString(newSimpleSection(s.title, s.creates, s.updates).String())

	// Write the schema changes
	if len(s.columnsToAdd) > 0 {
		sb.WriteString("Schema Columns to Add: ")
		writeList(&sb, s.columnsToAdd)
	}

	if len(s.columnsToDelete) > 0 {
		sb.WriteString("Schema Columns to Delete: ")
		writeList(&sb, s.columnsToDelete)
	}

	return sb.String()
}
