package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type cachesSection struct {
	run *cb.SystemUploadDryRun
}

func newCachesSection(run *cb.SystemUploadDryRun) *cachesSection {
	return &cachesSection{run: run}
}

func (l *cachesSection) Title() string {
	return "SHARED CACHES"
}

func (l *cachesSection) HasChanges() bool {
	return len(l.run.CachesToCreate)+len(l.run.CachesToUpdate) > 0
}

func (l *cachesSection) String() string {
	sb := strings.Builder{}

	for _, cache := range l.run.CachesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", cache))
	}

	for _, cache := range l.run.CachesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", cache))
	}

	return sb.String()
}
