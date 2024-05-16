package cblib

import (
	"fmt"

	rt "github.com/clearblade/cblib/resourcetree"
)

// index utilities

func handleIndex(index *rt.Index, onUnique func() error, onNonunique func() error) error {
	switch index.IndexType {
	case rt.IndexUnique:
		if onUnique != nil {
			return onUnique()
		}
	case rt.IndexNonUnique:
		if onNonunique != nil {
			return onNonunique()
		}
	default:
		return fmt.Errorf("unknown index type: %s", index.IndexType)
	}
	return nil
}

func doDropIndex(index *rt.Index, onUnique func() error, onNonunique func() error) error {
	err := handleIndex(index, onUnique, onNonunique)
	if err != nil {
		return fmt.Errorf("unable to drop index: %s", err)
	}
	return nil
}

func doCreateIndex(index *rt.Index, onUnique func() error, onNonunique func() error) error {
	err := handleIndex(index, onUnique, onNonunique)
	if err != nil {
		return fmt.Errorf("unable to create index: %s", err)
	}
	return nil
}
