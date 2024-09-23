package dryRun

import (
	cb "github.com/clearblade/Go-SDK"
)

func newEdgesSection(run *cb.SystemUploadDryRun) dryRunSection {
	return &schemaSection{
		title:           "EDGES",
		creates:         run.EdgesToCreate,
		updates:         run.EdgesToUpdate,
		columnsToAdd:    run.EdgeColumnsToAdd,
		columnsToDelete: run.EdgeColumnsToDelete,
	}
}
