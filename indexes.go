package cblib

import (
	"fmt"

	rt "github.com/clearblade/cblib/resourcetree"
)

// IndexDiff implements the `Differ` interface.
type IndexDiff struct {
	A      []*rt.Index
	B      []*rt.Index
	Result []*rt.Index
}

func (idxdiff *IndexDiff) Prepare() {
	idxdiff.Result = make([]*rt.Index, 0, len(idxdiff.A))
}

func (idxdiff *IndexDiff) LenA() int {
	return len(idxdiff.A)
}

func (idxdiff *IndexDiff) LenB() int {
	return len(idxdiff.B)
}

func (idxdiff *IndexDiff) Same(i, j int) bool {
	a := idxdiff.A[i]
	b := idxdiff.B[j]
	return a.Name == b.Name && a.IndexType == b.IndexType
}

func (idxdiff *IndexDiff) Keep(i int) {
	idxdiff.Result = append(idxdiff.Result, idxdiff.A[i])
}

// DiffIndexesFull takes two indexes, and returns two slices. The first slice
// contains the items that are in `left` but not in `right` (added), and the second
// slice returns the items that are in `right` but not in `left` (removed).
func DiffIndexesFull(left, right []*rt.Index) ([]*rt.Index, []*rt.Index) {
	added := IndexDiff{left, right, nil}
	Diff(&added)

	removed := IndexDiff{right, left, nil}
	Diff(&removed)

	return added.Result, removed.Result
}

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
