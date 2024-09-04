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

func GetFileName(path string) string {
	_, fileName := pth.Split(path)
	return fileName
}

func IsClearbladePath(path string) bool {
	return IsAdaptorPath(path) ||
		IsBucketSetPath(path) ||
		IsCodePath(path) ||
		IsCollectionPath(path) ||
		IsDeploymentPath(path) ||
		IsDevicePath(path) ||
		IsEdgePath(path) ||
		IsExternalDbPath(path) ||
		IsMessageHistoryWhitelistPath(path) ||
		IsMessageTypeTriggerPath(path) ||
		IsPluginPath(path) ||
		IsPortalPath(path) ||
		IsRolePath(path) ||
		IsSecretPath(path) ||
		IsServiceCachePath(path) ||
		IsTimerPath(path) ||
		IsTriggerPath(path) ||
		IsUserPath(path) ||
		IsWebhookPath(path)
}
