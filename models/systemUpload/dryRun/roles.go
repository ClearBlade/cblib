package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type rolesSection struct {
	run *cb.SystemUploadDryRun
}

func newRolesSection(run *cb.SystemUploadDryRun) *rolesSection {
	return &rolesSection{run: run}
}

func (l *rolesSection) Title() string {
	return "ROLES"
}

func (l *rolesSection) HasChanges() bool {
	return len(l.run.RolesToCreate)+len(l.run.RolesToUpdate) > 0
}

func (l *rolesSection) String() string {
	sb := strings.Builder{}

	for _, role := range l.run.RolesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", role))
	}

	for _, role := range l.run.RolesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", role))
	}

	return sb.String()
}
