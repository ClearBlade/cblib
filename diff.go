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

// IntDiff implements the `Differ` interface for a slice of integers.
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
	idiff.Result = append(idiff.Result, idiff.A[i])
}

// StringDiff implements the `Differ` interface for a slice of strings.
type StringDiff struct {
	A      []string
	B      []string
	Result []string
}

func (sdiff *StringDiff) Prepare() {
	sdiff.Result = make([]string, 0, len(sdiff.A))
}

func (sdiff *StringDiff) LenA() int {
	return len(sdiff.A)
}

func (sdiff *StringDiff) LenB() int {
	return len(sdiff.B)
}

func (sdiff *StringDiff) Same(i, j int) bool {
	return sdiff.A[i] == sdiff.B[j]
}

func (sdiff *StringDiff) Keep(i int) {
	sdiff.Result = append(sdiff.Result, sdiff.A[i])
}

// UnsafeDiff implements the `Differ` interface for an unsafe slice of interfaces.
type UnsafeDiff struct {
	A       []interface{}
	B       []interface{}
	Result  []interface{}
	Compare func(interface{}, interface{}) bool
}

func (udiff *UnsafeDiff) Prepare() {
	udiff.Result = make([]interface{}, 0, len(udiff.A))
}

func (udiff *UnsafeDiff) LenA() int {
	return len(udiff.A)
}

func (udiff *UnsafeDiff) LenB() int {
	return len(udiff.B)
}

func (udiff *UnsafeDiff) Same(i, j int) bool {
	return udiff.Compare(udiff.A[i], udiff.B[j])
}

func (udiff *UnsafeDiff) Keep(i int) {
	udiff.Result = append(udiff.Result, udiff.A[i])
}
