package cblib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//"strings"
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

func setRootDir(theRootDir string) {
	rootDir = theRootDir
	svcDir = rootDir + "/code/services"
	libDir = rootDir + "/code/libraries"
	dataDir = rootDir + "/data"
	usersDir = rootDir + "/users"
	timersDir = rootDir + "/timers"
	triggersDir = rootDir + "/triggers"
}
func setupDirectoryStructure(sys *System_meta) error {
	/*
		rootDir = strings.Replace(sys.Name, " ", "_", -1)
		svcDir = rootDir + "/code/services"
		libDir = rootDir + "/code/libraries"
		dataDir = rootDir + "/data"
		usersDir = rootDir + "/users"
		timersDir = rootDir + "/timers"
		triggersDir = rootDir + "/triggers"
	*/
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
	myLibDir := libDir + "/" + name
	if err := os.MkdirAll(myLibDir, 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(myLibDir+"/"+name+".js", []byte(data["code"].(string)), 0666); err != nil {
		return err
	}
	delete(data, "code")
	delete(data, "library_key")
	delete(data, "system_key")
	return writeEntity(myLibDir, name, data)
}

func isException(name string, exceptions []string) bool {
	if name == "." || name == ".." {
		return false
	}
	for _, exception := range exceptions {
		if name == exception {
			return true
		}
	}
	return false
}

func getFileList(dirName string, exceptions []string) ([]string, error) {
	rval := []string{}
	fileList, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, oneFile := range fileList {
		if isException(oneFile.Name(), exceptions) {
			continue
		}
		rval = append(rval, oneFile.Name())
	}
	return rval, nil
}

func getObjectList(dirName string, exceptions []string) ([]map[string]interface{}, error) {
	rval := []map[string]interface{}{}
	fileList, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, err
	}
	for _, oneFile := range fileList {
		if isException(oneFile.Name(), exceptions) {
			continue
		}
		objMap, err := getObject(dirName, oneFile.Name())
		if err != nil {
			return nil, err
		}
		rval = append(rval, objMap)
	}
	return rval, nil
}

func getCodeStuff(dirName string) ([]map[string]interface{}, error) {
	dirList, err := getFileList(dirName, []string{})
	rval := []map[string]interface{}{}
	if err != nil {
		return nil, err
	}
	for _, realDirName := range dirList {
		myRootDir := dirName + "/" + realDirName + "/"
		myObj, err := getObject(myRootDir, realDirName+".json")
		if err != nil {
			return nil, err
		}
		byts, err := ioutil.ReadFile(myRootDir + "/" + realDirName + ".js")
		if err != nil {
			return nil, err
		}
		myObj["code"] = string(byts)
		delete(myObj, "source")
		rval = append(rval, myObj)
	}
	return rval, nil
}

func getLibraries() ([]map[string]interface{}, error) {
	return getCodeStuff(libDir)
}

func getServices() ([]map[string]interface{}, error) {
	return getCodeStuff(svcDir)
}

func getUsers() ([]map[string]interface{}, error) {
	return getObjectList(usersDir, []string{"schema.json"})
}

func getCollections() ([]map[string]interface{}, error) {
	return getObjectList(dataDir, []string{})
}

func getTriggers() ([]map[string]interface{}, error) {
	return getObjectList(triggersDir, []string{})
}

func getTimers() ([]map[string]interface{}, error) {
	return getObjectList(timersDir, []string{})
}

//  For most of these calls below (getUser, etc) the second arg
//  is really the filename as obtained by ReadDir, not the actual object
//  name -- it is <object name>.json

func getObject(dirName, objName string) (map[string]interface{}, error) {
	return getDict(dirName + "/" + objName)
}

func getUserSchema() (map[string]interface{}, error) {
	return getObject(usersDir, "schema.json")
}

func getUser(email string) (map[string]interface{}, error) {
	return getObject(usersDir, email)
}

func getTrigger(name string) (map[string]interface{}, error) {
	return getObject(triggersDir, name)
}

func getTimer(name string) (map[string]interface{}, error) {
	return getObject(timersDir, name)
}

func getCollection(name string) (map[string]interface{}, error) {
	return getObject(dataDir, name)
}

func getService(name string) (map[string]interface{}, error) {
	svcRootDir := svcDir + "/" + name
	codeFile := name + ".js"
	schemaFile := name + ".json"

	svcMap, err := getObject(svcRootDir, schemaFile)
	if err != nil {
		return nil, err
	}
	byts, err := ioutil.ReadFile(svcRootDir + "/" + codeFile)
	if err != nil {
		return nil, err
	}
	svcMap["code"] = string(byts)
	return svcMap, nil
}

func getLibrary(name string) (map[string]interface{}, error) {
	libRootDir := libDir + "/" + name
	codeFile := name + ".js"
	schemaFile := name + ".json"

	libMap, err := getObject(libRootDir, schemaFile)
	if err != nil {
		return nil, err
	}
	byts, err := ioutil.ReadFile(libRootDir + "/" + codeFile)
	if err != nil {
		return nil, err
	}
	libMap["code"] = string(byts)
	return libMap, nil
}
