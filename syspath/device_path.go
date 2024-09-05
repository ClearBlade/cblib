package syspath

import (
	"fmt"
	"regexp"
)

const (
	deviceDataPathRegexStr = `^devices\/([^\/]+)\.json$`
	deviceRolePathRegexStr = `^devices\/roles\/([^\/]+)\.json$`
)

var (
	deviceDataPathRegex *regexp.Regexp
	deviceRolePathRegex *regexp.Regexp
)

func init() {
	deviceDataPathRegex = regexp.MustCompile(deviceDataPathRegexStr)
	deviceRolePathRegex = regexp.MustCompile(deviceRolePathRegexStr)
}

func IsDevicePath(path string) bool {
	return topLevelDirectoryIs(path, "devices")
}

func IsDeviceSchemaPath(path string) bool {
	return path == "devices/schema.json"
}

func GetDeviceNameFromDataPath(path string) (string, error) {
	matches := deviceDataPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a device data path", path)
	}

	return matches[1], nil
}

func GetDeviceNameFromRolePath(path string) (string, error) {
	matches := deviceRolePathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a device role path", path)
	}

	return matches[1], nil
}
