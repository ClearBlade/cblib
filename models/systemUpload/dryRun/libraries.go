package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type librariesSection struct {
	run *cb.SystemUploadDryRun
}

func newLibrariesSection(run *cb.SystemUploadDryRun) *librariesSection {
	return &librariesSection{run: run}
}

func (l *librariesSection) Title() string {
	return "LIBRARIES"
}

func (l *librariesSection) HasChanges() bool {
	return len(l.run.LibrariesToCreate)+len(l.run.LibrariesToUpdate) > 0
}

func (l *librariesSection) String() string {
	sb := strings.Builder{}

	for _, library := range l.run.LibrariesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", library))
	}

	for _, library := range l.run.LibrariesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", library))
	}

	return sb.String()
}
