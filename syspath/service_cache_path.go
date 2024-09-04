package syspath

import (
	"fmt"
	"regexp"
)

const (
	serviceCacheRegexStr = `^shared\-caches\/([^\/]+)\.json$`
)

var (
	serviceCacheRegex *regexp.Regexp
)

func init() {
	serviceCacheRegex = regexp.MustCompile(serviceCacheRegexStr)
}

func IsServiceCachePath(path string) bool {
	return topLevelDirectoryIs(path, "shared-caches")
}

func GetServiceCacheNameFromPath(path string) (string, error) {
	matches := serviceCacheRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a service cache path", path)
	}

	return matches[1], nil
}
