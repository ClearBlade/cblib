package cblib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	rootDir     string
	dataDir     string
	svcDir      string
	libDir      string
	usersDir    string
	timersDir   string
	triggersDir string
)

func setupDirectoryStructure(sys *System_meta) error {
	rootDir = strings.Replace(sys.Name, " ", "_", -1)
	svcDir = rootDir + "/code/services"
	libDir = rootDir + "/code/libraries"
	dataDir = rootDir + "/data"
	usersDir = rootDir + "/users"
	timersDir = rootDir + "/timers"
	triggersDir = rootDir + "/triggers"
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
	if err := os.MkdirAll(usersDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", usersDir, err.Error())
	}
	if err := os.MkdirAll(timersDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", timersDir, err.Error())
	}
	if err := os.MkdirAll(triggersDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", triggersDir, err.Error())
	}
	return nil
}

func storeSystemDotJSON() error {
	delete(systemDotJSON, "services")
	delete(systemDotJSON, "libraries")
	delete(systemDotJSON, "timers")
	delete(systemDotJSON, "triggers")
	delete(systemDotJSON, "users")
	delete(systemDotJSON, "data")
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

/*
func writeCollection(collection map[string]interface{}, allData []interface{}) error {
	colName := collection["name"].(string)
	fileName := dataDir + "/" + colName + ".json"
	marshalled, err := json.MarshalIndent(allData, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall collection data for: %s", colName)
	}
	if err = ioutil.WriteFile(fileName, marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to %s: %s", fileName, err.Error())
	}
	return nil
}
*/

func getDict(filename string) (map[string]interface{}, error) {
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parsed := map[string]interface{}{}
	err = json.Unmarshal(jsonStr, &parsed)
	if err != nil {
		jsonErr := err.(*json.SyntaxError)
		return nil, fmt.Errorf("%s: (%s, line %d)\n", err.Error(), filename,
			bytes.Count(jsonStr[:jsonErr.Offset], []byte("\n"))+1)
	}
	return parsed, nil
}

func getArray(filename string) ([]interface{}, error) {
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	parsed := []interface{}{}
	err = json.Unmarshal(jsonStr, &parsed)
	if err != nil {
		jsonErr := err.(*json.SyntaxError)
		return nil, fmt.Errorf("%s: (%s, line %d)\n", err.Error(), filename,
			bytes.Count(jsonStr[:jsonErr.Offset], []byte("\n"))+1)
	}
	return parsed, nil
}

func getServiceCode(serviceName string) (string, error) {
	return getCode("services", serviceName)
}

func getLibraryCode(libraryName string) (string, error) {
	return getCode("libraries", libraryName)
}

func getCode(dirName, fileName string) (string, error) {
	byts, err := ioutil.ReadFile("code/" + dirName + "/" + fileName + ".js")
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

func getCollectionItems(collectionName string) ([]interface{}, error) {
	fileName := "data/" + collectionName + ".json"
	return getArray(fileName)
}

func writeEntity(dirName, fileName string, stuff interface{}) error {
	marshalled, err := json.MarshalIndent(stuff, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall %s: %s", fileName, err.Error())
	}
	if err = ioutil.WriteFile(dirName+"/"+fileName+".json", marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to %s: %s", fileName, err.Error())
	}
	return nil
}

func writeCollection(collectionName string, data map[string]interface{}) error {
	return writeEntity(dataDir, collectionName, data)
}

func writeUser(email string, data map[string]interface{}) error {
	return writeEntity(usersDir, email, data)
}

func writeUserSchema(data []map[string]interface{}) error {
	return writeEntity(usersDir, "schema", data)
}

func writeTrigger(name string, data map[string]interface{}) error {
	return writeEntity(triggersDir, name, data)
}

func writeTimer(name string, data map[string]interface{}) error {
	return writeEntity(timersDir, name, data)
}

func writeService(name string, data map[string]interface{}) error {
	mySvcDir := svcDir + "/" + name
	if err := os.MkdirAll(mySvcDir, 0777); err != nil {
		return err
	}

	if err := ioutil.WriteFile(mySvcDir+"/"+name+".js", []byte(data["code"].(string)), 0666); err != nil {
		return err
	}

	cleanService(data)
	return writeEntity(mySvcDir, name, data)
}

func writeLibrary(name string, data map[string]interface{}) error {
	return writeEntity(libDir, name, data)
}
