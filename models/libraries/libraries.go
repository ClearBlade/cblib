package libraries

import (
	"strings"
)

type Library struct {
	rawLibrary   map[string]interface{}
	name         string
	dependencies []string
}

func NewLibraryFromMap(rawLibrary map[string]interface{}) Library {
	dependencies := strings.Split(rawLibrary["dependencies"].(string), ",")
	return Library{
		rawLibrary:   rawLibrary,
		name:         rawLibrary["name"].(string),
		dependencies: dependencies,
	}
}

func (l *Library) GetName() string {
	return l.name
}

func (l *Library) GetDependencies() []string {
	return l.dependencies
}

func (l *Library) GetMap() map[string]interface{} {
	return l.rawLibrary
}

// returns a list of Library in postorder
func PostorderLibraries(libs []Library) []Library {
	// Create a map to keep track of visited libraries
	visited := make(map[string]bool)

	// Create a slice to hold the result libraries
	result := []Library{}

	// Helper function to traverse dependencies recursively
	var traverse func(lib Library)

	traverse = func(lib Library) {
		// Check if the current library has already been visited
		if visited[lib.GetName()] {
			return
		}

		// Mark the current library as visited
		visited[lib.GetName()] = true

		// Traverse the dependencies of the current library
		for _, dep := range lib.GetDependencies() {
			// Find the dependent library in the slice of libraries
			for _, l := range libs {
				if l.GetName() == dep {
					// Recursively traverse the dependencies of the dependent library
					traverse(l)
				}
			}
		}

		// Add the current library to the result slice
		result = append(result, lib)
	}

	// Traverse the dependencies of each library in the input slice
	for _, lib := range libs {
		traverse(lib)
	}

	return result
}
