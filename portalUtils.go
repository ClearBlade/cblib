package cblib

import (
	"encoding/json"
	"fmt"

	"github.com/totherme/unstructured"
)

func actOnParserSettings(widgetConfig map[string]interface{}, cb func(string, string) error) error {
	widgetSettings := make(map[string]interface{})
	ok := true
	if widgetSettings, ok = widgetConfig["props"].(map[string]interface{}); !ok {
		return fmt.Errorf("No props key for widget config")
	}
	for settingName, v := range widgetSettings {
		switch v.(type) {
		case map[string]interface{}:
			// if there's a dataType property we know this setting is a parser
			if dataType, ok := v.(map[string]interface{})["dataType"].(string); ok {
				if err := cb(settingName, dataType); err != nil {
					return err
				}
			}
		default:
			continue
		}
	}
	return nil
}

func convertPortalMapToUnstructured(p map[string]interface{}) (*unstructured.Data, error) {
	jason, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	portalConfig, err := unstructured.ParseJSON(string(jason))
	if err != nil {
		return nil, err
	}
	return &portalConfig, nil
}
