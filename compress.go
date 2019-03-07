package cblib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
	"github.com/totherme/unstructured"
)

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

func compressDatasources(portalDotJsonAbsPath, portalUncompressedDir string) error {

	portalJSONString, _ := readFileAsString(portalDotJsonAbsPath)
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

		currDS, _ := readFileAsString(path)
		currDSObj, _ := unstructured.ParseJSON(currDS)
		dsId, err := currDSObj.GetByPointer("/id")
		stringifiedDsId, _ := dsId.StringValue()
		//baseName := filepath.Base(path)
		// _, err = myPayloadData.GetByPointer("/" + stringifiedDsId)
		// if err != nil {
		// 	return err
		// }
		updatedDsData, _ := currDSObj.ObValue()
		err = myPayloadData.SetField(stringifiedDsId, updatedDsData)
		//log.Println(string(InnermarshaledValue))
		if err != nil {
			return err
		}
		log.Println("Printing Path", path)

		return nil
	})

	updatedPortalObject, _ := portalData.ObValue()
	err = writeFile(portalDotJsonAbsPath, updatedPortalObject)
	if err != nil {
		return err
	}
	// for _, ds := range datasources {
	// 	dataSourceData := ds.(map[string]interface{})
	// 	dataSourceName := dataSourceData["name"].(string)
	// 	if err := writeDatasource(portalName, dataSourceName, dataSourceData); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func compressWidgets(portal map[string]interface{}) error {
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

func getDecompressedPortalDir(portalName string) string {
	return filepath.Join(portalsDir, portalName)
}
func docompress(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if err := checkPortalCodeManagerArgsAndFlags(args); err != nil {
		return err
	}
	SetRootDir(".")
	portalDotJsonAbsPath := filepath.Join(portalsDir, PortalName+".json")

	portalUncompressedDir := getDecompressedPortalDir(PortalName)

	if err := compressDatasources(portalDotJsonAbsPath, portalUncompressedDir); err != nil {
		return err
	}
	// if err := compressWidgets(portal); err != nil {
	// 	return err
	// }
	return nil
}
