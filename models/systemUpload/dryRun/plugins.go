package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type pluginsSection struct {
	run *cb.SystemUploadDryRun
}

func newPluginsSection(run *cb.SystemUploadDryRun) *pluginsSection {
	return &pluginsSection{run: run}
}

func (l *pluginsSection) Title() string {
	return "PLUGINS"
}

func (l *pluginsSection) HasChanges() bool {
	return len(l.run.PluginsToCreate)+len(l.run.PluginsToUpdate) > 0
}

func (l *pluginsSection) String() string {
	sb := strings.Builder{}

	for _, plugin := range l.run.PluginsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", plugin))
	}

	for _, plugin := range l.run.PluginsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", plugin))
	}

	return sb.String()
}
