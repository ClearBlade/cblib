package syspath

var (
	MessageHistoryStoragePath = "message-history-storage/storage.json"
)

func IsMessageHistoryWhitelistPath(path string) bool {
	return topLevelDirectoryIs(path, "message-history-storage")
}

func IsMessageHistoryStorageFile(path string) bool {
	return path == MessageHistoryStoragePath
}
