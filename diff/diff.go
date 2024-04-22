package diff

// Differ defines an interface for a collection that can be diffed against
// another collection, it works with indices to let the caller handle the
// type since Go doesn't have generics. Similar to sort.Interface.
// see: https://golang.org/pkg/sort/#Interface
type Differ interface {
	LenA() int
	LenB() int
	Same(i, j int) bool
	Keep(i int)
}

// DifferPrepare defines an interface for an optional prepare operation that
// implementors might want to use for allocating memory, etc.
type DifferPrepare interface {
	Prepare()
}

// DifferDrop defines an interface for an optional drop operation that captures
// dropped values.
type DifferDrop interface {
	Drop(j int)
}

// Diff computes the difference between A (after) and B (before) in the given
// `Differ`. Values that are in A but not in B will be sent to `Keep` (added),
// while values that are in B but not in A will be sent to `Drop` (removed).
func Diff(diff Differ) {

	diffPrepare, ok := diff.(DifferPrepare)
	if ok {
		diffPrepare.Prepare()
	}

	intersectA := make(map[int]struct{}, diff.LenA())
	intersectB := make(map[int]struct{}, diff.LenB())

	for idx := 0; idx < diff.LenA(); idx++ {
		for jdx := 0; jdx < diff.LenB(); jdx++ {
			if diff.Same(idx, jdx) {
				intersectA[idx] = struct{}{}
				intersectB[jdx] = struct{}{}
			}
		}
	}

	for idx := 0; idx < diff.LenA(); idx++ {
		if _, ok := intersectA[idx]; !ok {
			diff.Keep(idx)
		}
	}

	diffDrop, ok := diff.(DifferDrop)
	if !ok {
		return
	}

	for jdx := 0; jdx < diff.LenB(); jdx++ {
		if _, ok := intersectB[jdx]; !ok {
			diffDrop.Drop(jdx)
		}
	}
}

// IntDiff implements the `Differ` interface for a slice of integers.
type IntDiff struct {
	After   []int
	Before  []int
	Added   []int
	Removed []int
}

func (idiff *IntDiff) Prepare() {
	idiff.Added = make([]int, 0, len(idiff.After))
	idiff.Removed = make([]int, 0, len(idiff.Before))
}

func (idiff *IntDiff) LenA() int {
	return len(idiff.After)
}

func (idiff *IntDiff) LenB() int {
	return len(idiff.Before)
}

func (idiff *IntDiff) Same(i, j int) bool {
	return idiff.After[i] == idiff.Before[j]
}

func (idiff *IntDiff) Keep(i int) {
	idiff.Added = append(idiff.Added, idiff.After[i])
}

func (idiff *IntDiff) Drop(j int) {
	idiff.Removed = append(idiff.Removed, idiff.Before[j])
}

// StringDiff implements the `Differ` interface for a slice of strings.
type StringDiff struct {
	After   []string
	Before  []string
	Added   []string
	Removed []string
}

func (sdiff *StringDiff) Prepare() {
	sdiff.Added = make([]string, 0, len(sdiff.After))
	sdiff.Removed = make([]string, 0, len(sdiff.Before))
}

func (sdiff *StringDiff) LenA() int {
	return len(sdiff.After)
}

func (sdiff *StringDiff) LenB() int {
	return len(sdiff.Before)
}

func (sdiff *StringDiff) Same(i, j int) bool {
	return sdiff.After[i] == sdiff.Before[j]
}

func (sdiff *StringDiff) Keep(i int) {
	sdiff.Added = append(sdiff.Added, sdiff.After[i])
}

func (sdiff *StringDiff) Drop(j int) {
	sdiff.Removed = append(sdiff.Removed, sdiff.Before[j])
}

// UnsafeDiff implements the `Differ` interface for an unsafe slice of interfaces.
type UnsafeDiff struct {
	After   []interface{}
	Before  []interface{}
	Added   []interface{}
	Removed []interface{}
	Compare func(interface{}, interface{}) bool
}

func (udiff *UnsafeDiff) Prepare() {
	udiff.Added = make([]interface{}, 0, len(udiff.After))
	udiff.Removed = make([]interface{}, 0, len(udiff.Before))
}

func (udiff *UnsafeDiff) LenA() int {
	return len(udiff.After)
}

func (udiff *UnsafeDiff) LenB() int {
	return len(udiff.Before)
}

func (udiff *UnsafeDiff) Same(i, j int) bool {
	return udiff.Compare(udiff.After[i], udiff.Before[j])
}

func (udiff *UnsafeDiff) Keep(i int) {
	udiff.Added = append(udiff.Added, udiff.After[i])
}

func (udiff *UnsafeDiff) Drop(j int) {
	udiff.Removed = append(udiff.Removed, udiff.Before[j])
}
