package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type triggersSection struct {
	run *cb.SystemUploadDryRun
}

func newTriggersSection(run *cb.SystemUploadDryRun) *triggersSection {
	return &triggersSection{run: run}
}

func (l *triggersSection) Title() string {
	return "TRIGGERS"
}

func (l *triggersSection) HasChanges() bool {
	return len(l.run.TriggersToCreate)+len(l.run.TriggersToUpdate) > 0
}

func (l *triggersSection) String() string {
	sb := strings.Builder{}

	for _, trigger := range l.run.TriggersToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", trigger))
	}

	for _, trigger := range l.run.TriggersToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", trigger))
	}

	return sb.String()
}
