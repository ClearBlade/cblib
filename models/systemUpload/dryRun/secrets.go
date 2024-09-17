package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type secretsSection struct {
	run *cb.SystemUploadDryRun
}

func newSecretsSection(run *cb.SystemUploadDryRun) *secretsSection {
	return &secretsSection{run: run}
}

func (l *secretsSection) Title() string {
	return "SECRETS"
}

func (l *secretsSection) HasChanges() bool {
	return len(l.run.SecretsToCreate)+len(l.run.SecretsToUpdate) > 0
}

func (l *secretsSection) String() string {
	sb := strings.Builder{}

	for _, secret := range l.run.SecretsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", secret))
	}

	for _, secret := range l.run.SecretsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", secret))
	}

	return sb.String()
}
