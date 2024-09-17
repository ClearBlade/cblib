package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type externalDatabasesSection struct {
	run *cb.SystemUploadDryRun
}

func newExternalDatabasesSection(run *cb.SystemUploadDryRun) *externalDatabasesSection {
	return &externalDatabasesSection{run: run}
}

func (l *externalDatabasesSection) Title() string {
	return "EXTERNAL DATABASES"
}

func (l *externalDatabasesSection) HasChanges() bool {
	return len(l.run.ExternalDbsToCreate)+len(l.run.ExternalDbsToUpdate) > 0
}

func (l *externalDatabasesSection) String() string {
	sb := strings.Builder{}

	for _, db := range l.run.ExternalDbsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", db))
	}

	for _, db := range l.run.ExternalDbsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", db))
	}

	return sb.String()
}
