package syspath

import (
	"fmt"
	"regexp"
)

const (
	timerPathRegexStr = `^timers\/([^\/]+)\.json$`
)

var (
	timerPathRegex *regexp.Regexp
)

func init() {
	timerPathRegex = regexp.MustCompile(timerPathRegexStr)
}

func IsTimerPath(path string) bool {
	return topLevelDirectoryIs(path, "timers")
}

func GetTimerNameFromPath(path string) (string, error) {
	matches := timerPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a timer path", path)
	}

	return matches[1], nil
}
