package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type edgesSection struct {
	run *cb.SystemUploadDryRun
}

func newEdgesSection(run *cb.SystemUploadDryRun) *edgesSection {
	return &edgesSection{run: run}
}

func (l *edgesSection) Title() string {
	return "EDGES"
}

func (l *edgesSection) HasChanges() bool {
	return len(l.run.EdgesToCreate)+len(l.run.EdgesToUpdate)+len(l.run.EdgeColumnsToAdd)+len(l.run.EdgeColumnsToDelete) > 0
}

func (l *edgesSection) String() string {
	sb := strings.Builder{}

	for _, edge := range l.run.EdgesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", edge))
	}

	for _, edge := range l.run.EdgesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", edge))
	}

	if len(l.run.EdgeColumnsToAdd) > 0 {
		sb.WriteString("Schema Columns to Add: ")
		writeList(&sb, l.run.EdgeColumnsToAdd)
	}

	if len(l.run.EdgeColumnsToDelete) > 0 {
		sb.WriteString("Schema Columns to Delete: ")
		writeList(&sb, l.run.EdgeColumnsToDelete)
	}

	return sb.String()
}
