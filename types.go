package cblib

func toStringArray(val interface{}) []string {
	result := make([]string, 0)

	arr := val.([]interface{})
	for _, item := range arr {
		result = append(result, item.(string))
	}

	return result
}
