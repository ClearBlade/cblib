package syspath

import (
	"fmt"
	"regexp"
)

const (
	collectionPathRegexStr = `^data\/([^\/]+)\.json$`
)

var (
	collectionPathRegex *regexp.Regexp
)

func init() {
	collectionPathRegex = regexp.MustCompile(collectionPathRegexStr)
}

func IsCollectionPath(path string) bool {
	return topLevelDirectoryIs(path, "data")
}

func GetCollectionNameFromPath(path string) (string, error) {
	matches := collectionPathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a collection path", path)
	}

	return matches[1], nil
}
