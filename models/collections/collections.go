package collections

func IsConnectCollection(co map[string]interface{}) bool {
	if isConnect, ok := co["isConnect"]; ok {
		switch isCon := isConnect.(type) {
		case bool:
			return isCon
		case string:
			return isCon == "true"
		default:
			return false
		}
	}
	return false
}
