package maputil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetIfMissingSucceeds(t *testing.T) {
	m := map[string]interface{}{"a": 1}

	assert.False(t, SetIfMissing(m, "a", 1))

	assert.True(t, SetIfMissing(m, "b", 2))

	b, found := LookupKey(m, "b")
	assert.True(t, found)
	assert.NotNil(t, b)
}
