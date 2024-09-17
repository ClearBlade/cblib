package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type deploymentsSection struct {
	run *cb.SystemUploadDryRun
}

func newDeploymentsSection(run *cb.SystemUploadDryRun) *deploymentsSection {
	return &deploymentsSection{run: run}
}

func (l *deploymentsSection) Title() string {
	return "DEPLOYMENTS"
}

func (l *deploymentsSection) HasChanges() bool {
	return len(l.run.DeploymentsToCreate)+len(l.run.DeploymentsToUpdate) > 0
}

func (l *deploymentsSection) String() string {
	sb := strings.Builder{}

	for _, deployment := range l.run.DeploymentsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", deployment))
	}

	for _, deployment := range l.run.DeploymentsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", deployment))
	}

	return sb.String()
}
