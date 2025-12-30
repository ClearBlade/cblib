package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type DryRun struct {
	sections []dryRunSection
	*cb.SystemUploadDryRun
}

func New(run *cb.SystemUploadDryRun) (DryRun, error) {
	return DryRun{
		SystemUploadDryRun: run,
		sections: []dryRunSection{
			newAdaptorsSection(run),
			newBucketSetsSection(run),
			newFileStoresSection(run),
			newCollectionsSection(run),
			newSimpleSection("DEPLOYMENTS", run.DeploymentsToCreate, run.DeploymentsToUpdate),
			newDevicesSection(run),
			newEdgesSection(run),
			newSimpleSection("EXTERNAL DATABASES", run.ExternalDbsToCreate, run.ExternalDbsToUpdate),
			newSimpleSection("LIBRARIES", run.LibrariesToCreate, run.LibrariesToUpdate),
			newMessageHistorySection(run),
			newMessageTypeTriggersSection(run),
			newSimpleSection("PLUGINS", run.PluginsToCreate, run.PluginsToUpdate),
			newSimpleSection("PORTALS", run.PluginsToCreate, run.PortalsToUpdate),
			newSimpleSection("ROLES", run.RolesToCreate, run.RolesToUpdate),
			newSimpleSection("SECRETS", run.SecretsToCreate, run.SecretsToUpdate),
			newSimpleSection("SHARED CACHES", run.CachesToCreate, run.CachesToUpdate),
			newSimpleSection("SERVICES", run.ServicesToCreate, run.ServicesToUpdate),
			newSimpleSection("TIMERS", run.TimersToCreate, run.TimersToUpdate),
			newSimpleSection("TRIGGERS", run.TriggersToCreate, run.TriggersToUpdate),
			newUsersSection(run),
			newSimpleSection("WEBHOOKS", run.WebhooksToCreate, run.WebhooksToUpdate),
		},
	}, nil
}

func (d *DryRun) String() string {
	if len(d.Errors) > 0 {
		return fmt.Sprintf("cannot push: %s", strings.Join(d.Errors, "\n"))
	}

	sb := strings.Builder{}
	if len(d.Warnings) > 0 {
		for _, warning := range d.Warnings {
			sb.WriteString(fmt.Sprintf("Warning: %s\n", warning))
		}
	}

	if !d.HasChanges() {
		return sb.String()
	}

	sb.WriteString("The following changes will be made:\n")
	for _, section := range d.sections {
		if section.HasChanges() {
			writeDryRunSection(&sb, section)
		}
	}

	return sb.String()
}

func (d *DryRun) HasChanges() bool {
	if len(d.Errors) > 0 {
		return false
	}

	for _, section := range d.sections {
		if section.HasChanges() {
			return true
		}
	}

	return false
}

func (d *DryRun) HasErrors() bool {
	return len(d.Errors) > 0
}
