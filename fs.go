package cblib

import (
	"bytes"
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/models"
	"io/ioutil"
	"os"
)

const SORT_KEY_CODE_SERVICE = "Name"
const SORT_KEY_COLLECTION_ITEM = "item_id"
const SORT_KEY_COLLECTION = "Name"

var (
	RootDirIsSet   bool
	
	rootDir        string
	dataDir        string
	svcDir         string
	libDir         string
	usersDir       string
	timersDir      string
	triggersDir    string
	rolesDir       string
	edgesDir       string
	devicesDir     string
	portalsDir     string
	pluginsDir     string
	adaptorsDir    string
	deploymentsDir string
	arrDir         [13]string //this is used to set up the directory structure for a system
)

func SetRootDir(theRootDir string) {

	RootDirIsSet = true

	rootDir = theRootDir
	svcDir = rootDir + "/code/services"
	libDir = rootDir + "/code/libraries"
	dataDir = rootDir + "/data"
	usersDir = rootDir + "/users"
	timersDir = rootDir + "/timers"
	triggersDir = rootDir + "/triggers"
	rolesDir = rootDir + "/roles"
	edgesDir = rootDir + "/edges"
	devicesDir = rootDir + "/devices"
	portalsDir = rootDir + "/portals"
	pluginsDir = rootDir + "/plugins"
	adaptorsDir = rootDir + "/adapters"
	deploymentsDir = rootDir + "/deployments"
	arrDir[0] = svcDir
	arrDir[1] = libDir
	arrDir[2] = dataDir
	arrDir[3] = usersDir
	arrDir[4] = timersDir
	arrDir[5] = triggersDir
	arrDir[6] = rolesDir
	arrDir[7] = edgesDir
	arrDir[8] = devicesDir
	arrDir[9] = portalsDir
	arrDir[10] = pluginsDir
	arrDir[11] = adaptorsDir
	arrDir[12] = deploymentsDir
	fmt.Println("Dirs:")
	fmt.Println(arrDir)
}

func setupDirectoryStructure(sys *System_meta) error {
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return fmt.Errorf("Could not make directory '%s': %s", rootDir, err.Error())
	}

	for i := 0; i < len(arrDir); i++ {
		if err := os.MkdirAll(arrDir[i], 0777); err != nil {
			return fmt.Errorf("Could not make directory '%s': %s", arrDir[i], err.Error())
		}
	}
	return nil
}

func cleanUpDirectories(sys *System_meta) error {
	fmt.Printf("CleaningUp Directories\n")
	for i := 0; i < len(arrDir); i++ {
		if err := os.RemoveAll(arrDir[i]); err != nil {
			return fmt.Errorf("Could not remove directory '%s': %s", arrDir[i], err.Error())
		}
	}
	return nil
}

func storeCBMeta(info map[string]interface{}) error {
	filename := ".cbmeta"
	marshalled, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshal .cbmeta info: %s", err.Error())
	}
	if err = ioutil.WriteFile(rootDir+"/"+filename, marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to .cbmeta: %s", err.Error())
	}
	return nil
}

func storeSystemDotJSON(systemDotJSON map[string]interface{}) error {
	delete(systemDotJSON, "services")
	delete(systemDotJSON, "libraries")
	delete(systemDotJSON, "timers")
	delete(systemDotJSON, "triggers")
	delete(systemDotJSON, "users")
	delete(systemDotJSON, "data")
	delete(systemDotJSON, "roles")
	delete(systemDotJSON, "edges")
	delete(systemDotJSON, "devices")
	delete(systemDotJSON, "portals")
	delete(systemDotJSON, "plugins")
	delete(systemDotJSON, "edge_deploy")
	marshalled, err := json.MarshalIndent(systemDotJSON, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall system.json: %s", err.Error())
	}
	if err = ioutil.WriteFile(rootDir+"/system.json", marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to system.json: %s", err.Error())
	}
	return nil
}

func storeDeployDotJSON(deployInfoList []map[string]interface{}) error {
	deployInfo := map[string]interface{}{"deployInfo": deployInfoList}
	marshalled, err := json.MarshalIndent(deployInfo, "", "    ")
	if err != nil {
		return fmt.Errorf("Could not marshall deploy.json: %s", err.Error())
	}
	if err = ioutil.WriteFile(rootDir+"/deploy.json", marshalled, 0666); err != nil {
		return fmt.Errorf("Could not write to deploy.json: %s", err.Error())
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
	byts, err := ioutil.ReadFile("code/" + dirName + "/" + fileName + "/" + fileName + ".js")
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

func getCollectionItems(collectionName string) ([]interface{}, error) {
	fileName := "data/" + collectionName + ".json"
	return getArray(fileName)
}

func getAdaptor(sysKey, adaptorName string, client *cb.DevClient) (*models.Adaptor, error) {
	currentDir := createFilePath(adaptorsDir, adaptorName)
	currentAdaptorInfo, err := getObject(currentDir, adaptorName+".json")
	if err != nil {
		return nil, err
	}

	adap := models.InitializeAdaptor(adaptorName, sysKey, client)
	adap.Info = currentAdaptorInfo

	adaptorFilesDir := createFilePath(currentDir, "files")
	adaptorFileDirList, err := getFileList(adaptorFilesDir, []string{})
	if err != nil {
		return nil, err
	}

	adap.InfoForFiles = make([]interface{}, 0)
	adap.ContentsForFiles = make([]map[string]interface{}, 0)

	for _, adaptorFileDirName := range adaptorFileDirList {
		currentFileDir := createFilePath(adaptorFilesDir, adaptorFileDirName)
		fileInfo, err := getObject(currentFileDir, adaptorFileDirName+".json")
		if err != nil {
			return nil, err
		}

		adap.InfoForFiles = append(adap.InfoForFiles, fileInfo)

		contentForFile := copyMap(fileInfo)

		fileContents, err := ioutil.ReadFile(createFilePath(currentFileDir, adaptorFileDirName))
		if err != nil {
			return nil, err
		}

		contentForFile["file"] = adap.EncodeFile(fileContents)

		adap.ContentsForFiles = append(adap.ContentsForFiles, contentForFile)

	}

	return adap, nil
}

func getAdaptors(sysKey string, client *cb.DevClient) ([]*models.Adaptor, error) {
	adaptorDirList, err := getFileList(adaptorsDir, []string{})
	if err != nil {
		// To ensure backwards-compatibility, we do not require
		// this folder to be present
		// As a result, let's log this error, but proceed
		fmt.Printf("Warning, could not read directory '%s' -- ignoring\n", adaptorsDir)
		return []*models.Adaptor{}, nil
	}
	rtn := make([]*models.Adaptor, 0)
	for _, adaptorDirName := range adaptorDirList {

		if adap, err := getAdaptor(sysKey, adaptorDirName, client); err != nil {
			return nil, err
		} else {
			rtn = append(rtn, adap)
		}

	}
	return rtn, nil
}

func writeServiceVersion(name string, data map[string]interface{}) error {
	mySvcDir := svcDir + "/" + name
	if err := os.MkdirAll(mySvcDir, 0777); err != nil {
		return err
	}
	cleanService(data)
	return writeEntity(mySvcDir, name, data)
}

func writeLibraryVersion(name string, data map[string]interface{}) error {
	myLibDir := libDir + "/" + name
	if err := os.MkdirAll(myLibDir, 0777); err != nil {
		return err
	}
	return writeEntity(myLibDir, name, data)
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
	if err := os.MkdirAll(dataDir, 0777); err != nil {
		return err
	}
	rawItemArray := data["items"]
	if rawItemArray == nil {
		return fmt.Errorf("Item array not found when accessing collection item array")
	}
	// Default value for items is an empty array, []
	itemArray, castSuccess := rawItemArray.([]interface{})
	if !castSuccess {
		return fmt.Errorf("Unable to process collection item array")
	}
	if SortCollections {
		fmt.Println("Note: Sorting collections by item_id. This may take time depending on collection size.")
		sortByFunction(&itemArray, compareCollectionItems)
	} else {
		fmt.Println("Note: Not sorting collections by item_id. Add sort-collection=true flag if desired.")
	}

	return writeEntity(dataDir, collectionName, data)
}

func writeUser(email string, data map[string]interface{}) error {
	if err := os.MkdirAll(usersDir, 0777); err != nil {
		return err
	}
	return writeEntity(usersDir, email, data)
}

func writeUserSchema(data map[string]interface{}) error {
	return writeEntity(usersDir, "schema", data)
}

func writeTrigger(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(triggersDir, 0777); err != nil {
		return err
	}
	return writeEntity(triggersDir, name, data)
}

func writeTimer(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(timersDir, 0777); err != nil {
		return err
	}
	return writeEntity(timersDir, name, data)
}

func writeDeployment(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(deploymentsDir, 0777); err != nil {
		return err
	}
	return writeEntity(deploymentsDir, name, data)
}

func writeRole(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(rolesDir, 0777); err != nil {
		return err
	}
	rawPermissions := data["Permissions"]
	if rawPermissions == nil {
		return fmt.Errorf("Permissions not found while processing role")
	}
	permissions, castSuccess := rawPermissions.(map[string]interface{})
	if !castSuccess {
		return fmt.Errorf("Unable to process role permissions: %v", rawPermissions)
	}
	// Default value for a role with no code services is null
	codeServices, castSuccess := permissions["CodeServices"].([]interface{})
	if castSuccess {
		sortByMapKey(&codeServices, SORT_KEY_CODE_SERVICE)
	}
	// Default value for a role with no collections is null
	collections, castSuccess := permissions["Collections"].([]interface{})
	if castSuccess {
		sortByMapKey(&collections, SORT_KEY_COLLECTION)
	}

	return writeEntity(rolesDir, name, data)
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

func writeEdge(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(edgesDir, 0777); err != nil {
		return err
	}
	return writeEntity(edgesDir, name, data)
}

func writeDevice(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(devicesDir, 0777); err != nil {
		return err
	}
	return writeEntity(devicesDir, name, data)
}

func writePortal(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(portalsDir, 0777); err != nil {
		return err
	}
	return writeEntity(portalsDir, name, data)
}

func writePlugin(name string, data map[string]interface{}) error {
	if err := os.MkdirAll(pluginsDir, 0777); err != nil {
		return err
	}
	return writeEntity(pluginsDir, name, data)
}

func writeAdaptor(a *models.Adaptor) error {
	myAdaptorDir := createFilePath(adaptorsDir, a.Name)
	if err := os.MkdirAll(myAdaptorDir, 0777); err != nil {
		return err
	}

	err := writeEntity(myAdaptorDir, a.Name, a.Info)
	if err != nil {
		return err
	}

	adaptorFilesDir := createFilePath(myAdaptorDir, "files")
	if err := os.MkdirAll(adaptorFilesDir, 0777); err != nil {
		return err
	}

	for i := 0; i < len(a.InfoForFiles); i++ {
		currentInfoForFile := a.InfoForFiles[i].(map[string]interface{})
		currentFileName := currentInfoForFile["name"].(string)
		currentAdaptorFileDir := createFilePath(myAdaptorDir, "files", currentFileName)
		if err := os.MkdirAll(currentAdaptorFileDir, 0777); err != nil {
			return err
		}
		if err := writeEntity(currentAdaptorFileDir, currentFileName, currentInfoForFile); err != nil {
			return err
		}
		fileContents, err := a.DecodeFileByName(currentFileName)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(createFilePath(currentAdaptorFileDir, currentFileName), fileContents, 0666); err != nil {
			return err
		}
	}

	return nil
}

func isException(name string, exceptions []string) bool {
	if name == "." || name == ".." || name == ".DS_Store" {
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
		// If the error is that the directory doesn't exist, this isn't an error per se,
		// so just return an empty list
		fmt.Printf("Warning, could not read directory '%s' -- ignoring\n", dirName)
		return []map[string]interface{}{}, nil
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
	dirList, err := getFileList(dirName, []string{".DS_Store", ".git", ".gitignore"}) // For starters
	if err != nil {
		fmt.Printf("getFileListFailed: %s, %s\n", dirName, err)
		return nil, err
	}
	rval := []map[string]interface{}{}
	for _, realDirName := range dirList {
		myRootDir := dirName + "/" + realDirName + "/"
		myObj, err := getObject(myRootDir, realDirName+".json")
		if err != nil {
			fmt.Printf("getObject failed: %s\n", err)
			return nil, err
		}
		byts, err := ioutil.ReadFile(myRootDir + "/" + realDirName + ".js")
		if err != nil {
			fmt.Printf("ioutil.ReadFile failed: %s\n", err)
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

func getRoles() ([]map[string]interface{}, error) {
	return getObjectList(rolesDir, []string{})
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

func getDeployments() ([]map[string]interface{}, error) {
	return getObjectList(deploymentsDir, []string{})
}

func getEdges() ([]map[string]interface{}, error) {
	return getObjectList(edgesDir, []string{"schema.json"})
}

func getEdgesSchema() (map[string]interface{}, error) {
	return getObject(edgesDir, "schema.json")
}

func getDevicesSchema() (map[string]interface{}, error) {
	return getObject(devicesDir, "schema.json")
}

func getDevices() ([]map[string]interface{}, error) {
	return getObjectList(devicesDir, []string{"schema.json"})
}

func getPortals() ([]map[string]interface{}, error) {
	return getObjectList(portalsDir, []string{})
}

func getPlugins() ([]map[string]interface{}, error) {
	return getObjectList(pluginsDir, []string{})
}

func getEdgeDeployInfo() (map[string]interface{}, error) {
	return getDict(rootDir + "/deploy.json")
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

func getRole(name string) (map[string]interface{}, error) {
	return getObject(rolesDir, name+".json")
}

func getUser(email string) (map[string]interface{}, error) {
	return getObject(usersDir, email+".json")
}

func getTrigger(name string) (map[string]interface{}, error) {
	return getObject(triggersDir, name+".json")
}

func getTimer(name string) (map[string]interface{}, error) {
	return getObject(timersDir, name+".json")
}

func getDevice(name string) (map[string]interface{}, error) {
	return getObject(devicesDir, name+".json")
}

func getEdge(name string) (map[string]interface{}, error) {
	return getObject(edgesDir, name+".json")
}

func getPortal(name string) (map[string]interface{}, error) {
	return getObject(portalsDir, name+".json")
}

func getPlugin(name string) (map[string]interface{}, error) {
	return getObject(pluginsDir, name+".json")
}

func getCollection(name string) (map[string]interface{}, error) {
	return getObject(dataDir, name+".json")
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

func getSysMeta() (*System_meta, error) {
	dict, err := getDict("system.json")
	if err != nil {
		return nil, err
	}
	platform_url, ok := dict["platformURL"].(string)
	if !ok {
		platform_url = dict["platform_url"].(string)
	}
	system_key, ok := dict["systemKey"].(string)
	if !ok {
		system_key = dict["system_key"].(string)
	}
	system_secret, ok := dict["systemSecret"].(string)
	if !ok {
		system_secret = dict["system_secret"].(string)
	}

	rval := &System_meta{
		Name:        dict["name"].(string),
		Key:         system_key,
		Secret:      system_secret,
		Description: dict["description"].(string),
		PlatformUrl: platform_url,
	}
	return rval, nil
}

func makeCollectionJsonConsistent(data map[string]interface{}) map[string]interface{} {
	data["collection_id"] = data["collectionID"].(string)
	data["app_id"] = data["appID"].(string)
	delete(data, "collectionID")
	delete(data, "appID")
	return data
}

// Although this is similar to utils.go's compareWithKey function,
// The logic in this function will diverge soon from the one below it in cb-cli v3
func compareCollectionItems(sliceOfItems *[]interface{}, i, j int) bool {

	sortKey := SORT_KEY_COLLECTION_ITEM

	slice := *sliceOfItems

	map1, castSuccess1 := slice[i].(map[string]interface{})
	map2, castSuccess2 := slice[j].(map[string]interface{})

	if !castSuccess1 || !castSuccess2 {
		return false
	}

	name1 := map1[sortKey]
	name2 := map2[sortKey]
	if !isString(name1) || !isString(name2) {
		return false
	}
	return name1.(string) < name2.(string)
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		// We have another error, but the file does exist
		// Just for the sake of this function, we ignore
		// the error
	}
	return true
}
