package syspath

import (
	"fmt"
	"regexp"
)

const (
	externalDbRegexStr = `^external-databases\/([^\/]+)\.json$`
)

var (
	externalDbRegex *regexp.Regexp
)

func init() {
	externalDbRegex = regexp.MustCompile(externalDbRegexStr)
}

func IsExternalDbPath(path string) bool {
	return topLevelDirectoryIs(path, "external-databases")
}

func GetExternalDbNameFromPath(path string) (string, error) {
	matches := externalDbRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not an external db path", path)
	}

	return matches[1], nil
}
