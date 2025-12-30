package syspath

import (
	pth "path"
	"strings"
)

func topLevelDirectoryIs(path, dirName string) bool {
	components := strings.Split(path, "/")
	return components[0] == dirName
}

func getFileExtension(path string) string {
	components := strings.Split(path, ".")
	if len(components) <= 1 {
		return ""
	}

	return strings.ToLower(components[len(components)-1])
}

func IsJsFile(path string) bool {
	return getFileExtension(path) == "js"
}

func IsJsonFile(path string) bool {
	return getFileExtension(path) == "json"
}

func IsJsMapFile(path string) bool {
	return getFileExtension(path) == "map"
}

func GetFileName(path string) string {
	_, fileName := pth.Split(path)
	return fileName
}
