package syspath

import (
	"fmt"
	"regexp"
)

const (
	triggerPathRegexStr = `^triggers\/([^\/]+)\.json$`
)

var (
	triggerPathRegex *regexp.Regexp
)

func init() {
	triggerPathRegex = regexp.MustCompile(triggerPathRegexStr)
}

func IsTriggerPath(path string) bool {
	return topLevelDirectoryIs(path, "triggers")
}

func GetTriggerNameFromPath(path string) (string, error) {
	matches := triggerPathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a trigger path", path)
	}

	return matches[1], nil
}
