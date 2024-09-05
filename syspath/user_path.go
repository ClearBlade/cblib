package syspath

import (
	"fmt"
	"regexp"
)

const (
	userDataPathRegexStr = `^users\/([^\/]+)\.json$`
	userRolePathRegexStr = `^users\/roles\/([^\/]+)\.json$`
)

var (
	userDataPathRegex *regexp.Regexp
	userRolePathRegex *regexp.Regexp
)

func init() {
	userDataPathRegex = regexp.MustCompile(userDataPathRegexStr)
	userRolePathRegex = regexp.MustCompile(userRolePathRegexStr)
}

func IsUserPath(path string) bool {
	return topLevelDirectoryIs(path, "users")
}

func IsUserSchemaPath(path string) bool {
	return path == "users/schema.json"
}

func GetUserEmailFromDataPath(path string) (string, error) {
	matches := userDataPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a user data path", path)
	}

	return matches[1], nil
}

func GetUserEmailFromRolePath(path string) (string, error) {
	matches := userRolePathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a user role path", path)
	}

	return matches[1], nil
}
