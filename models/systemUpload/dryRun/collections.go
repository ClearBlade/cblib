package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type collectionsSection struct {
	run *cb.SystemUploadDryRun
}

func newCollectionsSection(run *cb.SystemUploadDryRun) *collectionsSection {
	return &collectionsSection{run: run}
}

func (l *collectionsSection) Title() string {
	return "COLLECTIONS"
}

func (l *collectionsSection) HasChanges() bool {
	return len(l.run.CollectionsToCreate)+len(l.run.CollectionsToUpdate) > 0
}

func (l *collectionsSection) String() string {
	sb := strings.Builder{}

	for _, collection := range l.run.CollectionsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", collection.Name))
	}

	for _, collection := range l.run.CollectionsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", collection.Name))
	}

	return sb.String()
}
