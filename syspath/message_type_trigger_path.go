package syspath

var (
	MessageTypeTriggersPath = "message-type-triggers/triggers.json"
)

func IsMessageTypeTriggerPath(path string) bool {
	return topLevelDirectoryIs(path, "message-type-triggers")
}

func IsMessageTypeTriggersFile(path string) bool {
	return path == MessageTypeTriggersPath
}
