package cblib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntDiffSucceeds(t *testing.T) {

	tests := []struct {
		a       []int
		b       []int
		added   []int
		removed []int
	}{
		{[]int{1, 2, 3}, []int{}, []int{1, 2, 3}, []int{}},
		{[]int{1, 2, 3}, []int{1}, []int{2, 3}, []int{}},
		{[]int{1, 2, 3}, []int{1, 2}, []int{3}, []int{}},
		{[]int{1, 2, 3}, []int{1, 2, 3}, []int{}, []int{}},
		{[]int{2, 3}, []int{1, 2, 3}, []int{}, []int{1}},
		{[]int{3}, []int{1, 2, 3}, []int{}, []int{1, 2}},
		{[]int{}, []int{1, 2, 3}, []int{}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		intDiff := IntDiff{tt.a, tt.b, nil, nil}
		Diff(&intDiff)
		assert.Equal(t, tt.added, intDiff.Added)
		assert.Equal(t, tt.removed, intDiff.Removed)
	}
}
