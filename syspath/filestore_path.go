package syspath

import (
	"fmt"
	"regexp"
)

const (
	filestoreRegexStr     = `^file-stores\/([^\/]+)\.json$`
	filestoreFileRegexStr = `^file-stores-files\/([^\/]+)\/(.+)$`
)

var (
	filestoreRegex     *regexp.Regexp
	filestoreFileRegex *regexp.Regexp
)

func init() {
	filestoreRegex = regexp.MustCompile(filestoreRegexStr)
	filestoreFileRegex = regexp.MustCompile(filestoreFileRegexStr)
}

func IsFilestorePath(path string) bool {
	return IsFilestoreMetaPath(path) || IsFilestoreFilePath(path)
}

func IsFilestoreMetaPath(path string) bool {
	return topLevelDirectoryIs(path, "file-stores")
}

func IsFilestoreFilePath(path string) bool {
	return topLevelDirectoryIs(path, "file-stores-files")
}

func GetFilestoreNameFromPath(path string) (string, error) {
	matches := filestoreRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a filestore path", path)
	}

	return matches[1], nil
}

type FilestoreFilePath struct {
	FilestoreName string
	RelativePath  string
}

func ParseFileStorePath(path string) (*FilestoreFilePath, error) {
	matches := filestoreFileRegex.FindStringSubmatch(path)
	if len(matches) != 3 {
		return nil, fmt.Errorf("path %q is not a filestore file path", path)
	}

	return &FilestoreFilePath{
		FilestoreName: matches[1],
		RelativePath:  matches[2],
	}, nil
}
