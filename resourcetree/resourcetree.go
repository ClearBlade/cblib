// Package resourcetree provides a virtual tree to represent all the resources
// of a system. It has types and functions for adding, removing, and diffing
// resources.
package resourcetree

// Tree represents a virtual tree of resources.
type Tree struct {
	Collections []Collection
}

// TreeReader interface specifies a reader for a virtual tree.
type TreeReader interface {
	CollectionsReader
}

// TreeWriter interface specifies a writer for a virtual tree.
type TreeWriter interface {
	CollectionsWriter
}

// CollectionsReader interface specifies a function for reading collections.
type CollectionsReader interface {
	ReadCollections() ([]Collection, error)
}

// CollectionsWriter interface specifies a function for writing collections.
type CollectionsWriter interface {
	WriteCollections([]Collection) error
}

// NOTE: add other resource interfaces below...
