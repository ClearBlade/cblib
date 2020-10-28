package cblib

import (
	"testing"

	rt "github.com/clearblade/cblib/resourcetree"
	"github.com/stretchr/testify/assert"
)

func TestGetDiffForIndexesSucceeds(t *testing.T) {

	tests := []struct {
		a       []*rt.Index
		b       []*rt.Index
		added   []*rt.Index
		removed []*rt.Index
	}{

		// they are the same, no diff
		{
			[]*rt.Index{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			[]*rt.Index{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			[]*rt.Index{},
			[]*rt.Index{},
		},

		// added one item to local
		{
			[]*rt.Index{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			[]*rt.Index{{Name: "a"}, {Name: "b"}},
			[]*rt.Index{{Name: "c"}},
			[]*rt.Index{},
		},

		// removed one item from local
		{
			[]*rt.Index{{Name: "a"}, {Name: "b"}},
			[]*rt.Index{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			[]*rt.Index{},
			[]*rt.Index{{Name: "c"}},
		},

		// change index type for the same column
		{
			[]*rt.Index{{Name: "a", IndexType: rt.IndexNonUnique}},
			[]*rt.Index{{Name: "a", IndexType: rt.IndexUnique}},
			[]*rt.Index{{Name: "a", IndexType: rt.IndexNonUnique}},
			[]*rt.Index{{Name: "a", IndexType: rt.IndexUnique}},
		},
	}

	for _, tt := range tests {
		added, removed := DiffIndexesFull(tt.a, tt.b)
		assert.Equal(t, tt.added, added)
		assert.Equal(t, tt.removed, removed)
	}

}
