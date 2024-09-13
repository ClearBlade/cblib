package syspath

import (
	"fmt"
	"regexp"
)

const (
	bucketSetRegexStr     = `^bucket-sets\/([^\/]+)\.json$`
	bucketSetFileRegexStr = `^bucket-set-files\/([^\/]+)\/([^\/]+)\/(.+)$`
)

var (
	bucketSetRegex     *regexp.Regexp
	bucketSetFileRegex *regexp.Regexp
)

func init() {
	bucketSetRegex = regexp.MustCompile(bucketSetRegexStr)
	bucketSetFileRegex = regexp.MustCompile(bucketSetFileRegexStr)
}

type FullBucketPath struct {
	BucketName   string
	Box          string
	RelativePath string
}

func IsBucketSetPath(path string) bool {
	return IsBucketSetMetaPath(path) || IsBucketSetFilePath(path)
}

func IsBucketSetMetaPath(path string) bool {
	return topLevelDirectoryIs(path, "bucket-sets")
}

func IsBucketSetFilePath(path string) bool {
	return topLevelDirectoryIs(path, "bucket-set-files")
}

func GetBucketSetNameFromPath(path string) (string, error) {
	matches := bucketSetRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a bucket set path", path)
	}

	return matches[1], nil
}

func ParseBucketPath(path string) (*FullBucketPath, error) {
	matches := bucketSetFileRegex.FindStringSubmatch(path)
	if len(matches) != 4 {
		return nil, fmt.Errorf("path %q is not a bucket set file path", path)
	}

	return &FullBucketPath{
		BucketName:   matches[1],
		Box:          matches[2],
		RelativePath: matches[3],
	}, nil
}
