package dryRun

import (
	cb "github.com/clearblade/Go-SDK"
)

func newUsersSection(run *cb.SystemUploadDryRun) dryRunSection {
	return &schemaSection{
		title:           "USERS",
		creates:         run.UsersToCreate,
		updates:         run.UsersToUpdate,
		columnsToAdd:    run.UserColumnsToAdd,
		columnsToDelete: run.UserColumnsToDelete,
	}
}
