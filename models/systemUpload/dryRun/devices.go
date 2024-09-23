package dryRun

import (
	cb "github.com/clearblade/Go-SDK"
)

func newDevicesSection(run *cb.SystemUploadDryRun) dryRunSection {
	return &schemaSection{
		title:           "DEVICES",
		creates:         run.DevicesToCreate,
		updates:         run.DevicesToUpdate,
		columnsToAdd:    run.DeviceColumnsToAdd,
		columnsToDelete: run.DeviceColumnsToDelete,
	}
}
