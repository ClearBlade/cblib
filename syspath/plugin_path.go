package syspath

import (
	"fmt"
	"regexp"
)

const (
	pluginPathRegexStr = `^plugins\/([^\/]+)\.json$`
)

var (
	pluginPathRegex *regexp.Regexp
)

func init() {
	pluginPathRegex = regexp.MustCompile(pluginPathRegexStr)
}

func IsPluginPath(path string) bool {
	return topLevelDirectoryIs(path, "plugins")
}

func GetPluginNameFromPath(path string) (string, error) {
	matches := pluginPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a plugin path", path)
	}

	return matches[1], nil
}
