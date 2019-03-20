package cblib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/totherme/unstructured"
)

func cleanUpAndDecompress(name string, portal map[string]interface{}) (map[string]interface{}, error) {
	if err := os.RemoveAll(filepath.Join(portalsDir, name, portalConfigDirectory)); err != nil {
		return nil, err
	}

	portalConfig, err := convertPortalMapToUnstructured(portal)
	if err != nil {
		return nil, err
	}

	if err = decompressDatasources(portalConfig); err != nil {
		return nil, err
	}
	if err = decompressWidgets(portalConfig); err != nil {
		return nil, err
	}
	if err = decompressInternalResources(portalConfig); err != nil {
		return nil, err
	}

	return portalConfig.ObValue()
}

func checkPortalCodeManagerArgsAndFlags(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("There are no arguments to the update command, only command line options")
	}
	return nil
}

func decompressInternalResources(portal *unstructured.Data) error {
	portalName, err := portal.UnsafeGetField("name").StringValue()
	if err != nil {
		return err
	}

	resources, err := portal.GetByPointer(portalInternalResourcesPath)
	if err != nil {
		return err
	}

	keys, err := resources.Keys()
	if err != nil {
		return err
	}
	for _, id := range keys {
		resourceData, err := resources.GetByPointer("/" + id)
		if err != nil {
			return err
		}
		resourceName, err := resourceData.UnsafeGetField("name").StringValue()
		if err != nil {
			return err
		}
		if err := writeInternalResource(portalName, resourceName, &resourceData); err != nil {
			return err
		}

	}

	portalConfig, err := portal.GetByPointer(portalConfigPath)
	if err != nil {
		return err
	}
	if err = portalConfig.SetField("internalResources", "___placeholder___"); err != nil {
		return err
	}

	return nil
}

func decompressDatasources(portal *unstructured.Data) error {

	portalName, err := portal.UnsafeGetField("name").StringValue()
	if err != nil {
		return err
	}

	d, err := portal.GetByPointer(portalDatasourcesPath)
	if err != nil {
		return err
	}
	datasources, err := d.ObValue()
	if err != nil {
		return err
	}

	for _, ds := range datasources {
		dataSourceData := ds.(map[string]interface{})
		dataSourceName := dataSourceData["name"].(string)
		if err := writeDatasource(portalName, dataSourceName, dataSourceData); err != nil {
			return err
		}
	}

	portalConfig, err := portal.GetByPointer(portalConfigPath)
	if err != nil {
		return err
	}
	if err = portalConfig.SetField("datasources", "___placeholder___"); err != nil {
		return err
	}
	return nil
}

func getDatasourceParser(settings map[string]interface{}) string {
	useParser, ok := settings[datasourceUseParserKey].(bool)
	if ok && useParser == true {
		return settings[datasourceParserKey].(string)
	}
	return ""
}

func writeDatasource(portalName, dataSourceName string, data map[string]interface{}) error {
	myDatasourceDir := filepath.Join(portalsDir, portalName, portalConfigDirectory, datasourceDirectory, dataSourceName)
	if err := os.MkdirAll(myDatasourceDir, 0777); err != nil {
		return err
	}

	settings, ok := data["settings"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("No datasource settings for '%s'", dataSourceName)
	}

	dsParser := getDatasourceParser(settings)
	if dsParser != "" {
		settings[datasourceParserKey] = "./" + datasourceParserFileName
		if err := writeFile(myDatasourceDir+"/"+datasourceParserFileName, dsParser); err != nil {
			return err
		}
	}

	return writeEntity(myDatasourceDir, "meta", data)
}

func decompressWidgets(portal *unstructured.Data) error {

	portalName, err := portal.UnsafeGetField("name").StringValue()
	if err != nil {
		return err
	}

	w, err := portal.GetByPointer(portalWidgetsPath)
	if err != nil {
		return err
	}
	widgets, err := w.ObValue()
	if err != nil {
		return err
	}

	for id := range widgets {
		widgetData, err := portal.GetByPointer(portalWidgetsPath + "/" + id)
		if err != nil {
			return err
		}
		widgetName := getOrGenerateWidgetName(widgetData)
		if err := writeWidget(portalName, widgetName, &widgetData); err != nil {
			return err
		}
	}

	portalConfig, err := portal.GetByPointer(portalConfigPath)
	if err != nil {
		return err
	}
	portalConfig.SetField("widgets", "___placeholder___")
	return nil
}

func getOrGenerateWidgetName(widgetData unstructured.Data) string {
	widgetID, _ := widgetData.UnsafeGetField("id").StringValue()
	widgetType, _ := widgetData.UnsafeGetField("type").StringValue()
	name := fmt.Sprintf("%s"+"_"+"%v", widgetType, widgetID)
	return name
}

func writeParserBasedOnDataType(dataType string, setting *unstructured.Data, filePath string) error {
	found := false
	if setting.HasKey(incomingParserKey) {
		raw, _ := setting.GetByPointer("/" + incomingParserKey)
		ip := &raw
		found = true
		if dataType != dynamicDataType {
			ip = setting
		}
		if err := writeParserFiles(ip, filePath+"/"+incomingParserKey); err != nil {
			return err
		}
	}

	if setting.HasKey(outgoingParserKey) {
		found = true
		raw, _ := setting.GetByPointer("/" + outgoingParserKey)
		op := &raw
		if dataType != dynamicDataType {
			op = setting
		}
		if err := writeParserFiles(op, filePath+"/"+outgoingParserKey); err != nil {
			return err
		}
	}

	if !found {
		if setting.HasKey("value") {
			if err := writeParserFiles(setting, filePath+"/"+incomingParserKey); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeWidgetMeta(widgetDir string, widgetConfig *unstructured.Data) error {
	keys, err := widgetConfig.Keys()
	if err != nil {
		return err
	}
	meta := make(map[string]interface{})
	// grab all the keys except for "props" aka settings
	for _, k := range keys {
		if k != "props" {
			meta[k] = widgetConfig.UnsafeGetField(k).RawValue()
		}
	}
	return writeFile(filepath.Join(widgetDir, portalWidgetMetaFile), meta)
}

func writeWidgetSettings(widgetDir string, widgetConfig *unstructured.Data) error {
	return writeFile(filepath.Join(widgetDir, portalWidgetSettingsFile), widgetConfig.UnsafeGetField("props").RawValue())
}

func createInternalResourceMeta(resourceData *unstructured.Data) (map[string]interface{}, error) {
	keys, err := resourceData.Keys()
	if err != nil {
		return nil, err
	}
	rtn := make(map[string]interface{})
	for _, k := range keys {
		if k == "file" {
			rtn[k] = "___placeholder___"
		} else {
			rtn[k] = resourceData.UnsafeGetField(k).RawValue()
		}
	}
	return rtn, nil
}

func writeInternalResource(portalName, resourceName string, resourceData *unstructured.Data) error {
	// write the parser file
	currResourceDir := filepath.Join(portalsDir, portalName, portalInternalResourcesPath, resourceName)

	file := resourceData.UnsafeGetField("file")
	fileStr, err := file.StringValue()
	if err != nil {
		return err
	}

	if err := writeFile(currResourceDir+"/"+resourceName, fileStr); err != nil {
		return err
	}

	meta, err := createInternalResourceMeta(resourceData)
	if err != nil {
		return err
	}

	if err := writeFile(currResourceDir+"/"+portalInternalResourceMetaFile, meta); err != nil {
		return err
	}

	return nil
}

func writeWidget(portalName, widgetName string, widgetData *unstructured.Data) error {
	currWidgetDir := filepath.Join(portalsDir, portalName, portalConfigDirectory, widgetsDirectory, widgetName)

	widgetDataMap, err := widgetData.UnsafeGetField("props").ObValue()
	if err != nil {
		return err
	}
	if err := actOnParserSettings(widgetDataMap, func(settingName, dataType string) error {
		parserSetting, err := widgetData.GetByPointer("/props/" + settingName)
		if err != nil {
			return err
		}
		if err := writeParserBasedOnDataType(dataType, &parserSetting, currWidgetDir+"/"+parsersDirectory+"/"+settingName); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	if err := writeWidgetMeta(currWidgetDir, widgetData); err != nil {
		return err
	}

	if err := writeWidgetSettings(currWidgetDir, widgetData); err != nil {
		return err
	}
	return nil
}

func writeParserFiles(parserData *unstructured.Data, currWidgetDir string) error {
	keysToIgnoreInData := map[string]interface{}{}
	absFilePath := filepath.Join(currWidgetDir, outFile)

	value := parserData.UnsafeGetField("value")

	switch value.RawValue().(type) {
	case string:
		str, _ := value.StringValue()
		if err := writeFile(absFilePath+".js", str); err != nil {
			return err
		}
		if err := parserData.SetField("value", "___placeholder___"); err != nil {
			return err
		}
	case map[string]interface{}:
		mapp, _ := value.ObValue()
		if err := writeWebFiles(absFilePath, mapp, keysToIgnoreInData); err != nil {
			return err
		}
		if err := parserData.SetField("value", map[string]interface{}{"placeholder": map[string]interface{}{}}); err != nil {
			return err
		}
	default:
		return nil
	}

	return nil
}

func writeWebFiles(absFilePath string, data, keysToIgnoreInData map[string]interface{}) error {

	outjs := recursivelyFindValueForKey(javascriptKey, data, keysToIgnoreInData)
	outhtml := recursivelyFindValueForKey(htmlKey, data, keysToIgnoreInData)
	outcss := recursivelyFindValueForKey(cssKey, data, keysToIgnoreInData)
	if outhtml != nil {
		if err := writeFile(absFilePath+".html", outhtml.(interface{})); err != nil {
			return err
		}
	}

	if outjs != nil {
		if err := writeFile(absFilePath+".js", outjs.(interface{})); err != nil {
			return err
		}
	}

	if outcss != nil {
		if err := writeFile(absFilePath+".css", outcss.(interface{})); err != nil {
			return err
		}
	}
	return nil
}

func writeFile(absFilePath string, data interface{}) error {
	if data == nil {
		return nil
	}
	outDir := filepath.Dir(absFilePath)
	if err := os.MkdirAll(outDir, 0777); err != nil {
		return err
	}
	switch data.(type) {
	case string:
		if err := ioutil.WriteFile(absFilePath, []byte(data.(string)), 0666); err != nil {
			return err
		}
	case map[string]interface{}:
		marshalled, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return fmt.Errorf("Could not marshall object: %s", err.Error())
		}
		if err := ioutil.WriteFile(absFilePath, []byte(marshalled), 0666); err != nil {
			return err
		}
	}
	return nil
}

func recursivelyFindValueForKey(queryKey string, data map[string]interface{}, keysToIgnoreInData map[string]interface{}) interface{} {
	for k, v := range data {
		if k == queryKey {
			return v
		}
		switch v.(type) {
		case map[string]interface{}:
			if keysToIgnoreInData[k] != nil {
				continue
			}
			val := recursivelyFindValueForKey(queryKey, v.(map[string]interface{}), keysToIgnoreInData)
			if val != nil {
				return val
			}
		default:
			continue
		}
	}
	return nil
}

func isFileEmpty(absFilePath string) bool {
	if fileInfo, err := os.Stat(absFilePath); err == nil {
		log.Println("Is File Empty", fileInfo.Size())
		if fileInfo.Size() == 0 {
			return true
		}
	}
	return false
}
