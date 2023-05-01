package libraries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostorderLibraries(t *testing.T) {
	libraries := []Library{
		{name: "1", dependencies: []string{"2"}},
		{name: "2", dependencies: []string{"4", "5"}},
		{name: "3", dependencies: []string{"6", "7"}},
		{name: "4", dependencies: []string{"8"}},
		{name: "5", dependencies: []string{}},
		{name: "6", dependencies: []string{"9", "10"}},
		{name: "7", dependencies: []string{}},
		{name: "8", dependencies: []string{}},
		{name: "9", dependencies: []string{}},
		{name: "10", dependencies: []string{}},
	}

	orderedLibraries := PostorderLibraries(libraries)
	orderedLibraryNames := make([]string, 0)
	for _, library := range orderedLibraries {
		orderedLibraryNames = append(orderedLibraryNames, library.GetName())
	}

	assert.Equal(t, []string{"8", "4", "5", "2", "1", "9", "10", "6", "7", "3"}, orderedLibraryNames)

	libraries = []Library{
		{name: "1", dependencies: []string{"2"}},
		{name: "2", dependencies: []string{"3"}},
		{name: "3", dependencies: []string{}},
		{name: "4", dependencies: []string{"3", "5"}},
		{name: "5", dependencies: []string{}},
	}

	orderedLibraries = PostorderLibraries(libraries)
	orderedLibraryNames = make([]string, 0)
	for _, library := range orderedLibraries {
		orderedLibraryNames = append(orderedLibraryNames, library.GetName())
	}

	assert.Equal(t, []string{"3", "2", "1", "5", "4"}, orderedLibraryNames)
}
