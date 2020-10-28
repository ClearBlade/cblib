package cblib

// Differ defines an interface for a collection that can be diffed against
// another collection, it works with indexes to let the caller handle type
// since Go doesn't have generics. Similar to sort.Interface.
// see: https://golang.org/pkg/sort/#Interface
type Differ interface {
	Prepare()
	LenA() int
	LenB() int
	Same(i, j int) bool
	Keep(i int)
}

// Diff computes the difference between A and B in the given `Differ`.
func Diff(diff Differ) {

	diff.Prepare()

	for idx := 0; idx < diff.LenA(); idx++ {

		found := false

		for jdx := 0; jdx < diff.LenB(); jdx++ {
			if diff.Same(idx, jdx) {
				found = true
			}
		}

		if !found {
			diff.Keep(idx)
		}
	}
}

type IntDiff struct {
	A      []int
	B      []int
	Result []int
}

func (idiff *IntDiff) Prepare() {
	idiff.Result = make([]int, 0, len(idiff.A))
}

func (idiff *IntDiff) LenA() int {
	return len(idiff.A)
}

func (idiff *IntDiff) LenB() int {
	return len(idiff.B)
}

func (idiff *IntDiff) Same(i, j int) bool {
	return idiff.A[i] == idiff.B[j]
}

func (idiff *IntDiff) Keep(i int) {
	if idiff.Result == nil {
		idiff.Result = make([]int, 0, len(idiff.A))
	}
	idiff.Result = append(idiff.Result, idiff.A[i])
}
