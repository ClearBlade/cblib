package cblib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
)

func init() {

	usage :=
		`
	Compresses or decompresses Portal code
	`

	example :=
		`
	cb-cli uncompress -portalName=portal1		#
	`

	uncompressCommand := &SubCommand{
		name:         "uncompress",
		usage:        usage,
		needsAuth:    false,
		mustBeInRepo: true,
		run:          doUncompress,
		example:      example,
	}
	uncompressCommand.flags.StringVar(&PortalName, "Portal", "", "Name of Portal to uncompress after editing")
	AddCommand("uncompress", uncompressCommand)
}

func checkPortalCodeManagerArgsAndFlags(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("There are no arguments to the update command, only command line options\n")
	}
	return nil
}

func writeDatasource(portalName, dataSourceName string, data map[string]interface{}) error {
	currentFileName := dataSourceName
	currDsDir := filepath.Join(portalsDir, portalName, "datasources")
	if err := os.MkdirAll(currDsDir, 0777); err != nil {
		return err
	}
	return writeEntity(currDsDir, currentFileName, data)
}

func uncompressDatasources(portal map[string]interface{}) error {
	var (
		portalName          string
		config, datasources map[string]interface{}
		ok                  bool
	)

	if portalName, ok = portal["name"].(string); !ok {
		return fmt.Errorf("Portal 'name' key missing in <Portal>.json file")
	}
	if config, ok = portal["config"].(map[string]interface{}); !ok {
		return fmt.Errorf("Portal 'config' key missing in <Portal>.json file")
	}
	if datasources, ok = config["datasources"].(map[string]interface{}); !ok {
		return fmt.Errorf("No Datasources defined in 'config' ")
	}

	for _, ds := range datasources {
		dataSourceData := ds.(map[string]interface{})
		dataSourceName := dataSourceData["name"].(string)
		if err := writeDatasource(portalName, dataSourceName, dataSourceData); err != nil {
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

func writeWidget(portalName, widgetName string, data map[string]interface{}) error {
	currWidgetDir := filepath.Join(portalsDir, portalName, "widgets", widgetName)

	//TODO see if widget name is passed, else write
	absFilePath := filepath.Join(currWidgetDir, "code")

	keysToIgnoreInData := map[string]interface{}{"incoming_parser": true, "outgoing_parser": true}
	if err := writeWebFiles(absFilePath, data, keysToIgnoreInData); err != nil {
		return err
	}

	if err := writeParserFiles("incoming_parser", currWidgetDir, data); err != nil {
		return err
	}
	if err := writeParserFiles("outgoing_parser", currWidgetDir, data); err != nil {
		return err
	}

	return nil
}

func writeWebFiles(absFilePath string, data, keysToIgnoreInData map[string]interface{}) error {

	outjs := resursivelyFindValueForKey("JavaScript", data, keysToIgnoreInData)
	outhtml := resursivelyFindValueForKey("HTML", data, keysToIgnoreInData)
	outcss := resursivelyFindValueForKey("CSS", data, keysToIgnoreInData)
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

func writeParserFiles(parserType, currWidgetDir string, data map[string]interface{}) error {
	keysToIgnoreInData := map[string]interface{}{}
	//log.Println("WriteParserFiles:: ", data)
	val := resursivelyFindValueForKey(parserType, data, keysToIgnoreInData)
	if val == nil {
		log.Println("Parser ", parserType, " does not exist in this widget")
		return nil
	}
	parserObj := val.(map[string]interface{})
	absFilePath := filepath.Join(currWidgetDir, parserType, "code")

	switch parserObj["value"].(type) {
	case string:
		if err := writeFile(absFilePath+".js", parserObj["value"].(string)); err != nil {
			return err
		}
	case map[string]interface{}:
		if err := writeWebFiles(absFilePath, parserObj["value"].(map[string]interface{}), keysToIgnoreInData); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func uncompressWidgets(portal map[string]interface{}) error {
	var (
		portalName      string
		config, widgets map[string]interface{}
		ok              bool
	)

	if portalName, ok = portal["name"].(string); !ok {
		return fmt.Errorf("Portal 'name' key missing in <Portal>.json file")
	}
	if config, ok = portal["config"].(map[string]interface{}); !ok {
		return fmt.Errorf("Portal 'config' key missing in <Portal>.json file")
	}
	if widgets, ok = config["widgets"].(map[string]interface{}); !ok {
		return fmt.Errorf("No widgets defined in 'config' ")
	}

	for _, ds := range widgets {
		widgetData := ds.(map[string]interface{})
		widgetName := getOrGenerateWidgetName(widgetData)
		if err := writeWidget(portalName, widgetName, widgetData); err != nil {
			return err
		}
	}
	return nil
}

func doUncompress(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if err := checkPortalCodeManagerArgsAndFlags(args); err != nil {
		return err
	}
	SetRootDir(".")
	portal, err := getPortal(PortalName)
	if err != nil {
		return err
	}

	if err = uncompressDatasources(portal); err != nil {
		return err
	}
	if err = uncompressWidgets(portal); err != nil {
		return err
	}
	return nil
}

func resursivelyFindValueForKey(queryKey string, data map[string]interface{}, keysToIgnoreInData map[string]interface{}) interface{} {
	for k, v := range data {
		if k == queryKey {
			return v
		}
		switch v.(type) {
		case map[string]interface{}:
			if keysToIgnoreInData[k] != nil {
				log.Println("key is", k)
				continue
			}
			val := resursivelyFindValueForKey(queryKey, v.(map[string]interface{}), keysToIgnoreInData)
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

func getOrGenerateWidgetName(widgetData map[string]interface{}) string {
	//  .widget_name_map check if name exists in the widget_name_map: id to name mapping
	// .widget_count_map if it does not exist, read the widget count map. widget_type_to_count_mapping
	// if these maps don't exist, the widgets are named same as Id
	//widgetConfigFileName := ".widget.config"
	widgetID := widgetData["id"].(string)
	widgetType := widgetData["type"].(string)
	name := fmt.Sprintf("%s"+"_"+"%v", widgetType, widgetID)
	return name

	// absFilePath := filepath.Join(rootDir, widgetConfigFileName)
	// log.Println("Widget ID", widgetID)
	// if !fileExists(absFilePath) || isFileEmpty(absFilePath) {
	// 	configStructData := map[string]interface{}{}
	// 	if errWriteFile := writeFile(absFilePath, configStructData); errWriteFile != nil {
	// 		log.Println("Error creating widget config file ", errWriteFile)
	// 		return widgetID, fmt.Errorf("Error creating widget config file %v ", errWriteFile)
	// 	}
	// }

	// widgetConfig, errConfig := getDict(absFilePath)
	// if errConfig != nil {
	// 	return widgetID, fmt.Errorf("Error with widget config %v", errConfig)
	// }

	// if widgetConfig == nil {
	// 	widgetConfig = map[string]interface{}{}
	// }
	// log.Println(widgetConfig)
	// widgetNameConfig, errNameConfig := widgetConfig["widget_name_config"].(map[string]interface{})
	// if !errNameConfig {
	// 	return string(widgetID), fmt.Errorf("Error with type conversion for widget name %v", errNameConfig)
	// }
	// if widgetNameConfig != nil {
	// 	// return if name exists in config
	// 	if widgetName := widgetNameConfig[widgetID]; widgetName != nil {
	// 		return widgetName.(string), nil
	// 	}
	// } else {
	// 	widgetConfig["widget_name_config"] = map[string]interface{}{}
	// 	widgetNameConfig = widgetConfig["widget_name_config"].(map[string]interface{})
	// }

	// widgetCountConfig, errCountConfig := widgetConfig["widget_count_config"].(map[string]int)
	// log.Println("Outside Error count config", errCountConfig)
	// if !errCountConfig || widgetConfig["widget_count_config"] == nil {
	// 	log.Println("In Error count config")
	// 	widgetConfig["widget_count_config"] = map[string]int{widgetType: 1}
	// 	widgetCountConfig = widgetConfig["widget_count_config"].(map[string]int)
	// }

	// // if count := widgetCountConfig[widgetType]; count == 0 {
	// // 	widgetCountConfig[widgetType] = 1
	// // } else {
	// widgetCountConfig[widgetType] = widgetCountConfig[widgetType] + 1
	// //}

	// widgetNameConfig[widgetID] = name

	// widgetConfig["widget_name_config"] = widgetNameConfig
	// widgetConfig["widget_count_config"] = widgetCountConfig
	// if errWriteFile := writeFile(absFilePath, widgetConfig); errWriteFile != nil {
	// 	log.Println("Error writing file...")
	// }

	// return name, nil
}
