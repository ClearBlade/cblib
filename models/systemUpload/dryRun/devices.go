package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type devicesSection struct {
	run *cb.SystemUploadDryRun
}

func newDevicesSection(run *cb.SystemUploadDryRun) *devicesSection {
	return &devicesSection{run: run}
}

func (l *devicesSection) Title() string {
	return "DEVICES"
}

func (l *devicesSection) HasChanges() bool {
	return len(l.run.DeviceColumnsToAdd)+len(l.run.DeviceColumnsToDelete)+len(l.run.DevicesToCreate)+len(l.run.DevicesToUpdate) > 0
}

func (l *devicesSection) String() string {
	sb := strings.Builder{}

	for _, device := range l.run.DevicesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", device))
	}

	for _, device := range l.run.DevicesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", device))
	}

	if len(l.run.DeviceColumnsToAdd) > 0 {
		sb.WriteString("Schema Columns to Add:\n")
		for i, device := range l.run.DeviceColumnsToAdd {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(device)
		}
	}

	if len(l.run.DeviceColumnsToDelete) > 0 {
		sb.WriteString("Schema Columns to Delete:\n")
		for i, device := range l.run.DeviceColumnsToDelete {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(device)
		}
	}

	return sb.String()
}
