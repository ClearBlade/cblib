package syspath

import (
	"fmt"
	"regexp"
)

const (
	adaptorRegexStr         = `^adapters\/([^\/]+)\/([^\/]+)\.json$`
	adaptorFileMetaRegexStr = `^adapters\/([^\/]+)\/files\/([^\/]+)\/([^\/]+)\.json$`
	adaptorFileDataRegexStr = `^adapters\/([^\/]+)\/files\/([^\/]+)\/([^\/]+)$`
)

var (
	adaptorRegex         *regexp.Regexp
	adaptorFileMetaRegex *regexp.Regexp
	adaptorFileDataRegex *regexp.Regexp
)

func init() {
	adaptorRegex = regexp.MustCompile(adaptorRegexStr)
	adaptorFileMetaRegex = regexp.MustCompile(adaptorFileMetaRegexStr)
	adaptorFileDataRegex = regexp.MustCompile(adaptorFileDataRegexStr)
}

func IsAdaptorPath(path string) bool {
	return topLevelDirectoryIs(path, "adapters")
}

func GetAdaptorNameFromPath(path string) (string, error) {
	matches := adaptorRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 3 {
		return "", fmt.Errorf("path %q is not an adaptor path", path)
	}

	if matches[1] != matches[2] {
		return "", fmt.Errorf("adaptor name %q does not match directory name %q", matches[2], matches[1])
	}

	return matches[1], nil
}

func GetAdaptorFileMetaNameFromPath(path string) (string, string, error) {
	matches := adaptorFileMetaRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("path %q is not an adaptor path", path)
	}

	if matches[2] != matches[3] {
		return "", "", fmt.Errorf("adaptor meta file name %q does not match directory name %q", matches[3], matches[2])
	}

	return matches[1], matches[2], nil
}

func GetAdaptorFileDataNameFromPath(path string) (string, string, error) {
	matches := adaptorFileDataRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("path %q is not an adaptor path", path)
	}

	if matches[2] != matches[3] {
		return "", "", fmt.Errorf("adaptor data file name %q does not match directory name %q", matches[3], matches[2])
	}

	return matches[1], matches[2], nil
}
