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

func New(dryRun *cb.SystemUploadDryRun) (DryRun, error) {
	return DryRun{
		SystemUploadDryRun: dryRun,
		sections: []dryRunSection{
			newAdaptorsSection(dryRun),
			newBucketSetsSection(dryRun),
			newCollectionsSection(dryRun),
			newDevicesSection(dryRun),
			newEdgesSection(dryRun),
			newDeploymentsSection(dryRun),
			newLibrariesSection(dryRun),
			newServicesSection(dryRun),
		},
	}, nil
}

func (d *DryRun) String() string {
	if len(d.Errors) > 0 {
		return fmt.Sprintf("cannot push: %s", strings.Join(d.Errors, "\n"))
	}

	if !d.HasChanges() {
		return ""
	}

	sb := strings.Builder{}
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
