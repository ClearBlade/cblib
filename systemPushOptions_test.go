package cblib

import (
	"testing"
)

func TestEmptyOptions(t *testing.T) {
	SetRootDir(".")
	opts := systemPushOptions{}
	regex := opts.GetFileRegex()
	if regex.String() != "" {
		t.Fatalf("Expected empty regex, got %s", regex.String())
	}
}

func TestOnlyCode(t *testing.T) {
	SetRootDir(".")
	opts := systemPushOptions{
		AllServices:  true,
		AllLibraries: true,
	}

	regex := opts.GetFileRegex()
	t.Fatalf(regex.String())
	if regex.String() != "" {
		t.Fatalf("Expected empty regex, got %s", regex.String())
	}
}
