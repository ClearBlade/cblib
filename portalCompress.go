package cblib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/totherme/unstructured"
)

func processDataSourceDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	return nil
}

func compressDatasources(portal *unstructured.Data, decompressedPortalDir string) error {
	myPayloadData, err := portal.GetByPointer("/config/datasources")

	if err != nil {
		return fmt.Errorf("Couldn't address into my own json")
	}
	datasourcesDir := filepath.Join(decompressedPortalDir, "datasources")
	filepath.Walk(datasourcesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		currDS, err := readFileAsString(path)
		if err != nil {
			return err
		}
		currDSObj, err := unstructured.ParseJSON(currDS)
		if err != nil {
			return err
		}
		dsID, err := currDSObj.GetByPointer("/id")
		if err != nil {
			return err
		}
		stringifiedDsID, err := dsID.StringValue()
		if err != nil {
			return err
		}
		updatedDsData, err := currDSObj.ObValue()
		if err != nil {
			return err
		}
		err = myPayloadData.SetField(stringifiedDsID, updatedDsData)

		if err != nil {
			return err
		}
		return nil
	})

	return nil
}

func extractUUiD(dirName string) string {
	re := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	uuidFromDir := re.Find([]byte(dirName))
	if uuidFromDir == nil {
		if strings.Contains(dirName, "TEXT_WIDGET_COMPONENT_title") {
			return "title"
		} else if strings.Contains(dirName, "TEXT_WIDGET_COMPONENT_flyoutTitle") {
			return "flyoutTitle"
		} else if strings.Contains(dirName, "HTML_WIDGET_COMPONENT_brand") {
			return "brand"
		}
		return ""
	}
	return string(uuidFromDir)
}

func recursivelyFindKeyPath(queryKey string, data map[string]interface{}, keysToIgnoreInData map[string]interface{}, keyPath string) string {
	for k, v := range data {
		if k == queryKey {
			return keyPath
		}
		switch v.(type) {
		case map[string]interface{}:
			if keysToIgnoreInData[k] != nil {
				continue
			}
			updatedKeyPath := keyPath + k + "/"
			val := recursivelyFindKeyPath(queryKey, v.(map[string]interface{}), keysToIgnoreInData, updatedKeyPath)
			if val != "" {
				return val
			}
		default:
			continue
		}
	}
	return ""
}

func updateObjUsingWebFiles(webData *unstructured.Data, currDir string) error {
	htmlFile := filepath.Join(currDir, outFile+".html")
	updateObjFromFile(webData, htmlFile, htmlKey)

	javascriptFile := filepath.Join(currDir, outFile+".js")
	updateObjFromFile(webData, javascriptFile, javascriptKey)

	cssFile := filepath.Join(currDir, outFile+".css")
	updateObjFromFile(webData, cssFile, cssKey)
	return nil
}

func updateObjFromFile(data *unstructured.Data, currFile string, fieldToSet string) error {
	s, err := readFileAsString(currFile)
	if err != nil {
		log.Println("Update obj from file error:", err)
		return err
	}
	data.SetField(fieldToSet, s)
	return nil
}

func processParser(currWidgetDir string, parserObj *unstructured.Data, parserType string) error {
	valueData, err := parserObj.GetByPointer("/value")
	if err != nil {
		return err
	}

	switch valueData.RawValue().(type) {
	case map[string]interface{}:
		currDir := filepath.Join(currWidgetDir, parserType)
		updateObjUsingWebFiles(&valueData, currDir)
	case string:
		currFile := filepath.Join(currWidgetDir, parserType, outFile+".js")
		updateObjFromFile(parserObj, currFile, "value")
	default:

	}
	return nil

}

func processCurrWidgetDir(path string, data *unstructured.Data) error {

	widgetSettings, err := data.ObValue()
	if err != nil {
		return err
	}

	return actOnParserSettings(widgetSettings, func(settingName, dataType string) error {
		settingDir := path + "/" + parsersDirectory + "/" + settingName

		if setting, err := data.GetByPointer("/props/" + settingName); err == nil {
			found := false
			if incoming, err := setting.GetByPointer("/" + incomingParserKey); err == nil {
				found = true
				if dataType != dynamicDataType {
					incoming = setting
				}
				if err := processParser(settingDir, &incoming, incomingParserKey); err != nil {
					return err
				}
			}

			if outgoing, err := setting.GetByPointer("/" + outgoingParserKey); err == nil {
				found = true
				if dataType != dynamicDataType {
					outgoing = setting
				}
				if err := processParser(settingDir, &outgoing, outgoingParserKey); err != nil {
					return err
				}
			}

			if !found {
				if setting.HasKey("value") {
					if err := processParser(settingDir, &setting, incomingParserKey); err != nil {
						return err
					}
				}
			}
		} else {
			return err
		}

		return nil
	})
}

func processOtherValues(currWidgetDir string, widgetsDataObj *unstructured.Data, keysToIgnoreInData map[string]interface{}) error {
	valueParent := recursivelyFindKeyPath("value", widgetsDataObj.RawValue().(map[string]interface{}), keysToIgnoreInData, "/")
	if valueParent == "" {
		return nil
	}
	valuePath := filepath.Join("/", valueParent, valueKey)
	valueParent = filepath.Join("/", valueParent)
	valueParentData, err := widgetsDataObj.GetByPointer(valueParent)
	if err != nil {
		log.Println("Got Err:", err)
		return err
	}
	valueData, err := widgetsDataObj.GetByPointer(valuePath)
	if err != nil {
		log.Println("Got Err:", err)
		return err
	}
	switch valueData.RawValue().(type) {
	case map[string]interface{}:
		updateObjUsingWebFiles(&valueData, currWidgetDir)
	case string:
		currFile := filepath.Join(currWidgetDir, outFile)
		updateObjFromFile(&valueParentData, currFile, valueKey)
	default:

	}
	return nil
}

func compressWidgets(portal *unstructured.Data, decompressedPortalDir string) error {
	widgetsDataObj, err := portal.GetByPointer("/config/widgets")
	if err != nil {
		return fmt.Errorf("Couldn't address into my own json")
	}

	widgetsDir := filepath.Join(decompressedPortalDir, "widgets")
	filepath.Walk(widgetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		split := strings.Split(path, "/")
		if split[len(split)-2] != widgetsDirectory {
			return nil
		}

		currUUID := extractUUiD(path)
		if currUUID == "" {
			return nil
		}
		currWidgetData, err := widgetsDataObj.GetByPointer("/" + currUUID)
		if err != nil {
			return err
		}

		return processCurrWidgetDir(path, &currWidgetData)
	})

	return nil
}

func getDecompressedPortalDir(portalName string) string {
	return filepath.Join(portalsDir, portalName, portalConfigDirectory)
}

func compressPortal(name string) (map[string]interface{}, error) {

	decompressedPortalDir := getDecompressedPortalDir(name)

	p, err := getPortal(name)
	if err != nil {
		return nil, err
	}
	portalConfig, err := convertPortalMapToUnstructured(p)
	if err != nil {
		return nil, err
	}

	if err := compressDatasources(portalConfig, decompressedPortalDir); err != nil {
		return nil, err
	}
	if err := compressWidgets(portalConfig, decompressedPortalDir); err != nil {
		return nil, err
	}

	return portalConfig.ObValue()
}
