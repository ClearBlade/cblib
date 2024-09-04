package syspath

func IsMessageHistoryWhitelistPath(path string) bool {
	return topLevelDirectoryIs(path, "message-history-storage")
}

func IsMessageHistoryStorageFile(path string) bool {
	return path == "message-history-storage/storage.json"
}
