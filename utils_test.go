package cblib

import (
	"fmt"
	"testing"

	"github.com/clearblade/cblib/models/roles"
	"github.com/stretchr/testify/assert"
)

func TestProcessURLsSucceeds(t *testing.T) {
	tests := []struct {
		platformURL          string
		messagingURL         string
		expectedPlatformURL  string
		expectedMessagingURL string
	}{

		// canonical
		{
			"https://platform.clearblade.com", "platform.clearblade.com",
			"https://platform.clearblade.com", "platform.clearblade.com:1883",
		},

		// platform and messaging override port
		{
			"https://platform.clearblade.com:8080", "platform.clearblade.com:8883",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:8883",
		},

		// platform has trailing slash
		{
			"https://platform.clearblade.com:8080/", "platform.clearblade.com",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:1883",
		},

		// does not specify messaging
		{
			"https://platform.clearblade.com:8080/", "",
			"https://platform.clearblade.com:8080", "platform.clearblade.com:1883",
		},
	}

	for _, tt := range tests {
		platformURL, messagingURL, err := processURLs(tt.platformURL, tt.messagingURL)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.expectedPlatformURL, platformURL)
		assert.Equal(t, tt.expectedMessagingURL, messagingURL)
	}
}

func TestBubbleSort_String(test *testing.T) {
	// input, truth to be tested against output
	rubric := [][][]string{
		{{"c", "a", "b"}, {"a", "b", "c"}},
		{{"1beta", "0alpha"}, {"0alpha", "1beta"}},
		{{"ngf", "35g", "][3", "KAH"}, {"35g", "KAH", "][3", "ngf"}}}

	err := executeBubbleSortTests(rubric)
	if err != nil {
		test.Error(err.Error())
	}
}

func TestBubbleSort_Empty(test *testing.T) {
	rubric := [][][]string{
		{{"", "", ""}, {"", "", ""}},
		{{}, {}}}

	err := executeBubbleSortTests(rubric)
	if err != nil {
		test.Error(err.Error())
	}
}

func TestBubbleSort_NilEntry(test *testing.T) {
	sortMe := []interface{}{"", "a", nil}
	truth := []interface{}{nil, "", "a"}
	sortByFunction(&sortMe, compareString)
	for i := 0; i < len(sortMe); i++ {
		if sortMe[i] == nil && truth[i] == nil {
			continue
		}
		if (sortMe[i] == nil && truth[i] != nil) || (sortMe[i] != nil && truth[i] == nil) {
			test.Errorf("Failed to sort input %d in %s", i, sortMe)
			return
		}
		if sortMe[i].(string) != truth[i].(string) {
			test.Errorf("Failed to sort input %d", i)
			return
		}
	}

}

func TestBubbleSort_NilArray(test *testing.T) {
	sortByFunction(nil, compareString)
}

func TestIsString_Normal(t *testing.T) {
	var rubric map[interface{}]bool = map[interface{}]bool{
		"hello": true,
		"":      true,
		nil:     false}
	for input, truth := range rubric {
		var output bool = isString(input)
		if output != truth {
			t.Errorf("Result: %s is a string: %t", input, output)
		}
	}
}

// Helper Test Method

func executeBubbleSortTests(inputOutput [][][]string) error {
	var sortMe []interface{}
	for t := 0; t < len(inputOutput); t++ {
		testIO := inputOutput[t]
		var input []string = testIO[0]
		var output []string = testIO[1]

		sortMe = make([]interface{}, len(input))
		for i := 0; i < len(sortMe); i++ {
			sortMe[i] = input[i]
		}

		sortByFunction(&sortMe, compareString)

		for i := 0; i < len(input); i++ {
			if sortMe[i].(string) != output[i] {
				return fmt.Errorf("Failed to sort input %d on test: %d", i, t)
			}
		}

	}
	return nil
}

// custom Compare function used for bubblesort
func compareString(slice *[]interface{}, i int, j int) bool {
	s1 := (*slice)[i]
	s2 := (*slice)[j]
	if s1 == nil {
		return true
	}
	if s2 == nil {
		return false
	}
	return s1.(string) < s2.(string)
}

func Test_DeletingAllUserRoles(t *testing.T) {
	backend := []string{"Authenticated", "Administrator"}
	local := []interface{}{}

	roleDiff := roles.DiffRoles(local, convertStringSliceToInterfaceSlice(backend))

	if len(roleDiff.Removed) != 2 {
		t.Errorf("Expected to remove 2 elements but got %d elements", len(roleDiff.Removed))
	}

	if len(roleDiff.Added) != 0 {
		t.Errorf("Expected to add 0 elements but got %d elements", len(roleDiff.Added))
	}

}
