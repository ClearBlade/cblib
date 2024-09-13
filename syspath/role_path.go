package syspath

import (
	"fmt"
	"regexp"
)

const (
	rolePathRegexStr = `^roles\/[^\/]+\.json$`
)

var (
	rolePathRegex *regexp.Regexp
)

func init() {
	rolePathRegex = regexp.MustCompile(rolePathRegexStr)
}

func IsRolePath(path string) bool {
	return topLevelDirectoryIs(path, "roles")
}

func GetRoleNameFromPath(path string) (string, error) {
	matches := rolePathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a roles path", path)
	}

	return matches[1], nil
}
