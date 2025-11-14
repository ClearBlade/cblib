package syspath

import (
	"fmt"
	"regexp"
)

const (
	servicePathRegexStr = `^code\/services\/([^\/]+)\/([^\/]+)\.(?:js|json|js.map)$`
	libraryPathRegexStr = `^code\/libraries\/([^\/]+)\/([^\/]+)\.(?:js|json)$`
)

var (
	servicePathRegex *regexp.Regexp
	libraryPathRegex *regexp.Regexp
)

func init() {
	servicePathRegex = regexp.MustCompile(servicePathRegexStr)
	libraryPathRegex = regexp.MustCompile(libraryPathRegexStr)
}

func IsCodePath(path string) bool {
	return topLevelDirectoryIs(path, "code")
}

func GetServiceNameFromPath(path string) (string, error) {
	matches := servicePathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 3 {
		return "", fmt.Errorf("path %q is not a service path", path)
	}

	if matches[1] != matches[2] {
		return "", fmt.Errorf("service name %q does not match directory name %q", matches[2], matches[1])
	}

	return matches[1], nil
}

func GetLibraryNameFromPath(path string) (string, error) {
	matches := libraryPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 3 {
		return "", fmt.Errorf("path %q is not a library path", path)
	}

	if matches[1] != matches[2] {
		return "", fmt.Errorf("library name %q does not match directory name %q", matches[2], matches[1])
	}

	return matches[1], nil
}
