package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type portalsSection struct {
	run *cb.SystemUploadDryRun
}

func newPortalsSection(run *cb.SystemUploadDryRun) *portalsSection {
	return &portalsSection{run: run}
}

func (l *portalsSection) Title() string {
	return "PORTALS"
}

func (l *portalsSection) HasChanges() bool {
	return len(l.run.PortalsToCreate)+len(l.run.PortalsToUpdate) > 0
}

func (l *portalsSection) String() string {
	sb := strings.Builder{}

	for _, portal := range l.run.PortalsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", portal))
	}

	for _, portal := range l.run.PortalsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", portal))
	}

	return sb.String()
}
