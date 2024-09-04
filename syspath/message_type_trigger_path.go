package syspath

func IsMessageTypeTriggerPath(path string) bool {
	return topLevelDirectoryIs(path, "message-type-triggers")
}

func IsMessageTypeTriggersFile(path string) bool {
	return path == "message-type-triggers/triggers.json"
}
