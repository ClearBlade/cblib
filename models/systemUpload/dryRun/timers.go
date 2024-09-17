package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type timersSection struct {
	run *cb.SystemUploadDryRun
}

func newTimersSection(run *cb.SystemUploadDryRun) *timersSection {
	return &timersSection{run: run}
}

func (l *timersSection) Title() string {
	return "TIMERS"
}

func (l *timersSection) HasChanges() bool {
	return len(l.run.TimersToCreate)+len(l.run.TimersToUpdate) > 0
}

func (l *timersSection) String() string {
	sb := strings.Builder{}

	for _, timer := range l.run.TimersToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", timer))
	}

	for _, timer := range l.run.TimersToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", timer))
	}

	return sb.String()
}
