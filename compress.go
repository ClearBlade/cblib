package cblib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	cb "github.com/clearblade/Go-SDK"
	"github.com/totherme/unstructured"
)

const OUTGOING_PARSER_KEY = "outgoing_parser"
const INCOMING_PARSER_KEY = "incoming_parser"
const VALUE_KEY = "value"

func init() {

	usage :=
		`
	Compresses or decompresses Portal code
	`

	example :=
		`
	cb-cli compress -portalName=portal1		#
	`

	compressCommand := &SubCommand{
		name:         "compress",
		usage:        usage,
		needsAuth:    false,
		mustBeInRepo: true,
		run:          docompress,
		example:      example,
	}
	compressCommand.flags.StringVar(&PortalName, "Portal", "", "Name of Portal to compress after editing")
	AddCommand("compress", compressCommand)
}

func readFileAsString(absFilePath string) (string, error) {
	byts, err := ioutil.ReadFile(absFilePath)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

func processDataSourceDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	log.Println("Printing Path", path, info)

	return nil
}

func compressDatasources(portalDotJSONAbsPath, portalUncompressedDir string) error {

	portalJSONString, _ := readFileAsString(portalDotJSONAbsPath)
	portalData, _ := unstructured.ParseJSON(portalJSONString)
	//log.Println(portalData)
	myPayloadData, err := portalData.GetByPointer("/config/datasources")

	if err != nil {
		panic("Couldn't address into my own json")
	}

	marshaledValue, _ := json.Marshal(myPayloadData.RawValue())
	log.Println(string(marshaledValue))

	datasourcesDir := filepath.Join(portalUncompressedDir, "datasources")
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

	updatedPortalObject, _ := portalData.ObValue()
	err = writeFile(portalDotJSONAbsPath, updatedPortalObject)
	if err != nil {
		return err
	}
	return nil
}

func extractUUiD(dirName string) string {
	re := regexp.MustCompile("[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	uuidFromDir := re.Find([]byte(dirName))
	if uuidFromDir == nil {
		return ""
	}
	return string(uuidFromDir)
}

func resursivelyFindKeyPath(queryKey string, data map[string]interface{}, keysToIgnoreInData map[string]interface{}, keyPath string) string {
	for k, v := range data {
		if k == queryKey {
			return keyPath
		}
		switch v.(type) {
		case map[string]interface{}:
			if keysToIgnoreInData[k] != nil {
				log.Println("key is", k)
				continue
			}
			updatedKeyPath := keyPath + k + "/"
			val := resursivelyFindKeyPath(queryKey, v.(map[string]interface{}), keysToIgnoreInData, updatedKeyPath)
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

	htmlFile := filepath.Join(currDir, OUT_FILE+".html")
	updateObjFromFile(webData, htmlFile, HTML_KEY)

	javascriptFile := filepath.Join(currDir, OUT_FILE+".js")
	updateObjFromFile(webData, javascriptFile, JAVASCRIPT_KEY)

	cssFile := filepath.Join(currDir, OUT_FILE+".css")
	updateObjFromFile(webData, cssFile, CSS_KEY)
	return nil
}

func updateObjFromFile(data *unstructured.Data, currFile string, fieldToSet string) error {
	s, err := readFileAsString(currFile)
	if err != nil {
		return err
	}
	log.Println("IN UPdate Obj from File", s)

	data.SetField(fieldToSet, s)
	return nil
}

func processParser(currWidgetDir string, widgetsDataObj *unstructured.Data, parserType string) error {
	pathTillParserParent := resursivelyFindKeyPath(parserType, widgetsDataObj.RawValue().(map[string]interface{}), map[string]interface{}{}, "/")
	parserPath := filepath.Join("/", pathTillParserParent, parserType)
	log.Println("Path to parser", parserPath)
	parserData, err := widgetsDataObj.GetByPointer(parserPath)
	if err != nil {
		return err
	}
	valueData, err := parserData.GetByPointer("/value")
	if err != nil {
		return err
	}

	switch valueData.RawValue().(type) {
	case map[string]interface{}:
		currDir := filepath.Join(currWidgetDir, parserType)
		updateObjUsingWebFiles(&valueData, currDir)
	case string:
		currFile := filepath.Join(currWidgetDir, parserType, OUT_FILE)
		updateObjFromFile(&parserData, currFile, "value")
	default:

	}
	return nil

}

func processCurrWidgetDir(path string, widgetsDataObj *unstructured.Data) error {
	log.Println("Curr Widget Dir", path)
	processParser(path, widgetsDataObj, INCOMING_PARSER_KEY)
	processParser(path, widgetsDataObj, OUTGOING_PARSER_KEY)
	keysToIgnoreInData := map[string]interface{}{"incoming_parser": true, "outgoing_parser": true}
	processOtherValues(path, widgetsDataObj, keysToIgnoreInData)
	return nil
}

func processOtherValues(currWidgetDir string, widgetsDataObj *unstructured.Data, keysToIgnoreInData map[string]interface{}) error {
	valueParent := resursivelyFindKeyPath("value", widgetsDataObj.RawValue().(map[string]interface{}), keysToIgnoreInData, "/")
	log.Println("value parent", valueParent)
	valuePath := filepath.Join("/", valueParent, VALUE_KEY)
	valueParentData, err := widgetsDataObj.GetByPointer(valueParent)
	if err != nil {
		return err
	}
	valueData, err := widgetsDataObj.GetByPointer(valuePath)
	if err != nil {
		return err
	}
	switch valueData.RawValue().(type) {
	case map[string]interface{}:
		currDir := filepath.Join("/", currWidgetDir)
		updateObjUsingWebFiles(&valueData, currDir)
	case string:
		currFile := filepath.Join("/", currWidgetDir)
		updateObjFromFile(&valueParentData, currFile, VALUE_KEY)
		fmt.Println("valueParentData: ", valueParentData)
	default:

	}
	return nil
}

func compressWidgets(portalDotJSONAbsPath, portalUncompressedDir string) error {
	portalJSONString, _ := readFileAsString(portalDotJSONAbsPath)
	portalData, _ := unstructured.ParseJSON(portalJSONString)
	widgetsDataObj, err := portalData.GetByPointer("/config/widgets")
	fmt.Println("Compressing Widgets")
	if err != nil {
		panic("Couldn't address into my own json")
	}

	// marshaledValue, _ := json.Marshal(widgetsDataObj.RawValue())
	// log.Println(string(marshaledValue))

	widgetsDir := filepath.Join(portalUncompressedDir, "widgets")
	filepath.Walk(widgetsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
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

	updatedPortalObject, _ := portalData.ObValue()
	err = writeFile(portalDotJSONAbsPath, updatedPortalObject)
	if err != nil {
		return err
	}

	return nil
}

func getDecompressedPortalDir(portalName string) string {
	return filepath.Join(portalsDir, portalName)
}
func docompress(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if err := checkPortalCodeManagerArgsAndFlags(args); err != nil {
		return err
	}
	SetRootDir(".")
	portalDotJSONAbsPath := filepath.Join(portalsDir, PortalName+".json")

	portalUncompressedDir := getDecompressedPortalDir(PortalName)

	// if err := compressDatasources(portalDotJSONAbsPath, portalUncompressedDir); err != nil {
	// 	return err
	// }
	if err := compressWidgets(portalDotJSONAbsPath, portalUncompressedDir); err != nil {
		return err
	}
	// if err := compressWidgets(portal); err != nil {
	// 	return err
	// }
	return nil
}
