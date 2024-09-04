package syspath

import "regexp"

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
	matches := rolePathRegex.FindStringSubmatch(path)
	return len(matches) == 1
}
