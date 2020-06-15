package maputil

// --------------------------------
// Map utilities
// --------------------------------
// Utility functions for interacting with map[string]interface{} types.

// LookupKey looks for the first matching key in the given map and returns
// its value.
func LookupKey(m map[string]interface{}, keys ...string) (interface{}, bool) {
	for _, k := range keys {
		if value, ok := m[k]; ok {
			return value, true
		}
	}
	return nil, false
}

// LookupBool is similar to lookupKey but parses the value into an integer.
func LookupBool(m map[string]interface{}, keys ...string) (bool, bool) {
	var defaultFalse bool

	value, found := LookupKey(m, keys...)
	if !found {
		return defaultFalse, false
	}

	b, ok := value.(bool)
	if !ok {
		return defaultFalse, false
	}

	return b, true
}

// LookupInt is similar to lookupKey but parses the value into an integer.
func LookupInt(m map[string]interface{}, keys ...string) (int, bool) {
	var zero int

	value, found := LookupKey(m, keys...)
	if !found {
		return zero, false
	}

	n, ok := value.(int)
	if !ok {
		return zero, false
	}

	return n, true
}

// LookupFloat32 is similar to lookupKey but parses the value into a float32.
func LookupFloat32(m map[string]interface{}, keys ...string) (float32, bool) {
	var zero float32

	value, found := LookupKey(m, keys...)
	if !found {
		return zero, false
	}

	n, ok := value.(float32)
	if !ok {
		return zero, false
	}

	return n, true
}

// LookupFloat64 is similar to lookupKey but parses the value into a float64.
func LookupFloat64(m map[string]interface{}, keys ...string) (float64, bool) {
	var zero float64

	value, found := LookupKey(m, keys...)
	if !found {
		return zero, false
	}

	n, ok := value.(float64)
	if !ok {
		return zero, false
	}

	return n, true
}

// LookupString is similar to lookupKey but parses the value into a string.
func LookupString(m map[string]interface{}, keys ...string) (string, bool) {
	var empty string

	value, found := LookupKey(m, keys...)
	if !found {
		return empty, false
	}

	str, ok := value.(string)
	if !ok {
		return empty, false
	}

	return str, true
}

// LookupMap is similar to lookupKey but parses the value into a map[string]interface{} type.
func LookupMap(m map[string]interface{}, keys ...string) (map[string]interface{}, bool) {
	value, found := LookupKey(m, keys...)
	if !found {
		return nil, false
	}

	m, ok := value.(map[string]interface{})
	if !ok {
		return nil, false
	}

	return m, true
}

// SetIfMissing assigns the given value to the given key if the key is not
// present in the map.
func SetIfMissing(m map[string]interface{}, key string, value interface{}) bool {
	value, found := LookupKey(m, key)
	if !found {
		m[key] = value
		return true
	}
	return false
}
