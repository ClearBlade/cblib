package cblib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	rootDir string
	dataDir string
	svcDir  string
	libDir  string
)

func setupDirectoryStructure(sys *System_meta) error {
	rootDir = strings.Replace(sys.Name, " ", "_", -1)
	svcDir = rootDir + "/code/services"
	libDir = rootDir + "/code/libraries"
	dataDir = rootDir + "/data"
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", rootDir, err.Error())
	}
	if err := os.MkdirAll(svcDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", svcDir, err.Error())
	}
	if err := os.MkdirAll(libDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", libDir, err.Error())
	}
	if err := os.MkdirAll(dataDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", dataDir, err.Error())
	}
	return nil
}

func storeSystemDotJSON() error {
	marshalled, err := json.MarshalIndent(systemDotJSON, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall system.json: %s", err.Error())
	}
	if err = ioutil.WriteFile(rootDir+"/system.json", marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to system.json: %s", err.Error())
	}
	return nil
}

func writeUsersFile(allUsers []map[string]interface{}) error {
	marshalled, err := json.MarshalIndent(allUsers, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall users.json: %s", err.Error())
	}
	if err = ioutil.WriteFile(rootDir+"/users.json", marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to users.json: %s", err.Error())
	}
	return nil
}
