package syspath

import (
	"fmt"
	"regexp"
)

const (
	edgePathRegexStr = `^edges\/([^\/]+)\.json$`
)

var (
	edgePathRegex *regexp.Regexp
)

func init() {
	edgePathRegex = regexp.MustCompile(edgePathRegexStr)
}

func IsEdgePath(path string) bool {
	return topLevelDirectoryIs(path, "edges")
}

func IsEdgeSchemaPath(path string) bool {
	return path == "edges/schema.json"
}

func GetEdgeNameFromPath(path string) (string, error) {
	matches := edgePathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not an edge data path", path)
	}

	return matches[1], nil
}
