package remote

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRemoteName(t *testing.T) {
	tests := []struct {
		name    string
		isValid bool
	}{
		{"foo", true},
		{"bar", true},
		{"foo-bar", true},
		{"foo_bar", true},
		{"foo0", true},
		{"foo-0", true},
		{"-foo", false},
		{"0foo", false},
		{"foo!", false},
		{"!foo", false},
		{"foo bar", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("(%s)", tt.name), func(t *testing.T) {
			if tt.isValid {
				assert.Nil(t, validateRemoteName(tt.name))
			} else {
				assert.NotNil(t, validateRemoteName(tt.name))
			}
		})
	}
}
