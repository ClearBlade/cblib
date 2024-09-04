package syspath

import (
	"fmt"
	"regexp"
)

const (
	secretPathRegexStr = `^secrets\/([^\/]+)\.json$`
)

var (
	secretPathRegex *regexp.Regexp
)

func init() {
	secretPathRegex = regexp.MustCompile(secretPathRegexStr)
}

func IsSecretPath(path string) bool {
	return topLevelDirectoryIs(path, "secrets")
}

func GetSecretNameFromPath(path string) (string, error) {
	matches := secretPathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a secret path", path)
	}

	return matches[1], nil
}
