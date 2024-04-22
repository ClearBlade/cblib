package listutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterSliceSuceeds(t *testing.T) {

	tests := []struct {
		items     []interface{}
		predicate func(interface{}) bool
		expected  []interface{}
	}{
		// filter even
		{
			[]interface{}{1, 2, 3, 4, 5},
			func(item interface{}) bool { return item.(int)%2 == 0 },
			[]interface{}{2, 4},
		},

		// filter odd
		{
			[]interface{}{1, 2, 3, 4, 5},
			func(item interface{}) bool { return item.(int)%2 != 0 },
			[]interface{}{1, 3, 5},
		},
	}

	for _, tt := range tests {
		filtered := FilterSlice(tt.items, tt.predicate)
		assert.Equal(t, tt.expected, filtered)
	}
}
