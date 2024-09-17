package dryRun

import (
	cb "github.com/clearblade/Go-SDK"
)

func getCollectionName(collection *cb.CollectionUpdate) string {
	return collection.Name
}

func newCollectionsSection(run *cb.SystemUploadDryRun) dryRunSection {
	creates := mapList(run.CollectionsToCreate, getCollectionName)
	update := mapList(run.CollectionsToUpdate, getCollectionName)
	return newSimpleSection("COLLECTIONS", creates, update)
}
