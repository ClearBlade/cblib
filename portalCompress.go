package cblib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	cb "github.com/clearblade/Go-SDK"
	"github.com/totherme/unstructured"
)

const outgoingParserKey = "outgoing_parser"
const incomingParserKey = "incoming_parser"
const valueKey = "value"

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
	compressCommand.flags.StringVar(&PortalName, "portal", "", "Name of Portal to compress after editing")
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
	return nil
}

func compressDatasources(portalDotJSONAbsPath, portalUncompressedDir string) error {
	fmt.Println("Compressing Datasources...")
	portalJSONString, _ := readFileAsString(portalDotJSONAbsPath)
	portalData, _ := unstructured.ParseJSON(portalJSONString)
	myPayloadData, err := portalData.GetByPointer("/config/datasources")

	if err != nil {
		return fmt.Errorf("Couldn't address into my own json")
	}
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
	fmt.Println("Successfully Compressed Datasources")
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
	//fmt.Println("-------CurrDir in Update Obj using web files: ", currDir)
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

func processParser(currWidgetDir string, widgetsDataObj *unstructured.Data, parserType string) error {
	pathTillParserParent := recursivelyFindKeyPath(parserType, widgetsDataObj.RawValue().(map[string]interface{}), map[string]interface{}{}, "/")
	if pathTillParserParent == "" {
		return nil
	}
	parserPath := filepath.Join("/", pathTillParserParent, parserType)
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
		currFile := filepath.Join(currWidgetDir, parserType, outFile)
		updateObjFromFile(&parserData, currFile, "value")
	default:

	}
	return nil

}

func processCurrWidgetDir(path string, widgetsDataObj *unstructured.Data) error {
	fmt.Printf("process me: %+v\n", widgetsDataObj)
	return nil
	processParser(path, widgetsDataObj, incomingParserKey)
	processParser(path, widgetsDataObj, outgoingParserKey)
	keysToIgnoreInData := map[string]interface{}{"incoming_parser": true, "outgoing_parser": true}
	if err := processOtherValues(path, widgetsDataObj, keysToIgnoreInData); err != nil {
		return err
	}
	return nil
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

func compressWidgets(portalDotJSONAbsPath, portalUncompressedDir string) error {
	portalJSONString, _ := readFileAsString(portalDotJSONAbsPath)
	portalData, _ := unstructured.ParseJSON(portalJSONString)
	widgetsDataObj, err := portalData.GetByPointer("/config/widgets")
	if err != nil {
		return fmt.Errorf("Couldn't address into my own json")
	}

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
	fmt.Println("Successfully Compressed Widgets")

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

	if err := compressDatasources(portalDotJSONAbsPath, portalUncompressedDir); err != nil {
		return err
	}
	if err := compressWidgets(portalDotJSONAbsPath, portalUncompressedDir); err != nil {
		return err
	}

	return nil
}
