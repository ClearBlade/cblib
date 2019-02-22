package cblib

import (
	"fmt"
	"io/ioutil"
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
	cb-cli compress -portalName=portal1		#
	`

	compressCommand := &SubCommand{
		name:         "compress",
		usage:        usage,
		needsAuth:    false,
		mustBeInRepo: true,
		run:          doCompress,
		example:      example,
	}
	compressCommand.flags.StringVar(&PortalName, "Portal", "", "Name of Portal to compress after editing")
	AddCommand("compress", compressCommand)
}

func checkPortalCodeManagerArgsAndFlags(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("There are no arguments to the update command, only command line options\n")
	}
	return nil
}

func writeDatasource(portalName, dsName string, data map[string]interface{}) error {
	currentFileName := dsName
	currDsDir := filepath.Join(portalsDir, portalName, "datasources")
	if err := os.MkdirAll(currDsDir, 0777); err != nil {
		return err
	}
	return writeEntity(currDsDir, currentFileName, data)
}

func compressDatasources(portal map[string]interface{}) error {
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
		dsName := dataSourceData["name"].(string)
		if err := writeDatasource(portalName, dsName, dataSourceData); err != nil {
			return err
		}
	}
	return nil
}

func writeWidget(portalName, widgetName string, data map[string]interface{}) error {
	currWidgetDir := filepath.Join(portalsDir, portalName, "widgets", widgetName)
	if err := os.MkdirAll(currWidgetDir, 0777); err != nil {
		return err
	}
	// see if widget name is passed, else write
	absFilePath := filepath.Join(currWidgetDir, "code")

	outjs := resursivelyFindValueForKey("JavaScript", data)
	if outjs != nil {
		if err := ioutil.WriteFile(absFilePath+".js", []byte(outjs.(string)), 0666); err != nil {
			return err
		}
	}

	outhtml := resursivelyFindValueForKey("HTML", data)
	if outhtml != nil {
		if err := ioutil.WriteFile(absFilePath+".html", []byte(outhtml.(string)), 0666); err != nil {
			return err
		}
	}

	outcss := resursivelyFindValueForKey("CSS", data)
	if outcss != nil {
		if err := ioutil.WriteFile(absFilePath+".css", []byte(outcss.(string)), 0666); err != nil {
			return err
		}
	}
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
		widgetName := widgetData["id"].(string)
		if err := writeWidget(portalName, widgetName, widgetData); err != nil {
			return err
		}
	}
	return nil
}

func doCompress(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if err := checkPortalCodeManagerArgsAndFlags(args); err != nil {
		return err
	}
	SetRootDir(".")
	portal, err := getPortal(PortalName)
	if err != nil {
		return err
	}

	if err = compressDatasources(portal); err != nil {
		return err
	}
	if err = compressWidgets(portal); err != nil {
		return err
	}
	return nil
}

func resursivelyFindValueForKey(queryKey string, data map[string]interface{}) interface{} {
	for k, v := range data {
		if k == queryKey {
			return v
		}
		switch v.(type) {
		case map[string]interface{}:
			val := resursivelyFindValueForKey(queryKey, v.(map[string]interface{}))
			if val != nil {
				return val
			}
		default:
			continue
		}
	}
	return nil
}
