package index

import (
	d "github.com/clearblade/cblib/diff"
	rt "github.com/clearblade/cblib/resourcetree"
)

// IndexDiff implements the `Differ` interface.
type IndexDiff struct {
	A       []*rt.Index
	B       []*rt.Index
	Added   []*rt.Index
	Removed []*rt.Index
}

func (idxdiff *IndexDiff) Prepare() {
	idxdiff.Added = make([]*rt.Index, 0, len(idxdiff.A))
	idxdiff.Removed = make([]*rt.Index, 0, len(idxdiff.B))
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
	idxdiff.Added = append(idxdiff.Added, idxdiff.A[i])
}

func (idxdiff *IndexDiff) Drop(j int) {
	idxdiff.Removed = append(idxdiff.Removed, idxdiff.B[j])
}

// DiffIndexesFull takes two slices of indexes diffs them using *IndexDiff.
func DiffIndexesFull(after, before []*rt.Index) *IndexDiff {
	diff := IndexDiff{after, before, nil, nil}
	d.Diff(&diff)
	return &diff
}
