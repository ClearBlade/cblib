package cblib

import (
	"fmt"
	"os"
	"path"
	"strings"

	cb "github.com/clearblade/Go-SDK"

	"github.com/clearblade/cblib/maputil"
)

var (
	importRows  bool
	importUsers bool
)

func init() {

	usage :=
		`
	Import a system from your local filesystem to the ClearBlade Platform
	`

	example :=
		`
	cb-cli import 									# prompts for credentials
	cb-cli import -importrows=false -importusers=false			# prompts for credentials, excludes all collection-rows and users
	`
	myImportCommand := &SubCommand{
		name:      "import",
		usage:     usage,
		needsAuth: false,
		run:       doImport,
		example:   example,
	}
	DEFAULT_IMPORT_ROWS := true
	DEFAULT_IMPORT_USERS := true
	myImportCommand.flags.BoolVar(&importRows, "importrows", DEFAULT_IMPORT_ROWS, "imports all data into all collections")
	myImportCommand.flags.BoolVar(&importUsers, "importusers", DEFAULT_IMPORT_USERS, "imports all users into the system")
	myImportCommand.flags.StringVar(&URL, "url", "https://platform.clearblade.com", "Clearblade Platform URL where system is hosted, ex https://platform.clearblade.com")
	myImportCommand.flags.StringVar(&Email, "email", "", "Developer email for login to import destination")
	myImportCommand.flags.StringVar(&Password, "password", "", "Developer password at import destination")
	myImportCommand.flags.IntVar(&DataPageSize, "data-page-size", DataPageSizeDefault, "Number of rows in a collection to push/import at a time")
	myImportCommand.flags.IntVar(&MaxRetries, "max-retries", 3, "Number of retries to attempt if a request fails")
	AddCommand("import", myImportCommand)
	AddCommand("imp", myImportCommand)
	AddCommand("im", myImportCommand)
}

func doImport(cmd *SubCommand, cli *cb.DevClient, args ...string) error {

	systemPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// prompt and skip vaues we don't need (messaging URL and system key)
	promptAndFillMissingAuth(nil, PromptSkipMsgURL|PromptSkipSystemKey)

	// authorizes using global flags (import ignores cb meta)
	cli, err = authorizeUsingGlobalCLIFlags()
	if err != nil {
		return err
	}

	// creates import config and proceeds to import system
	config := MakeImportConfigFromGlobals()
	_, err = ImportSystemUsingConfig(config, systemPath, cli)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------
// Import config and other types
// --------------------------------
// We use an import config that is passed around as a parameter during the
// import process.

// ImportConfig contains configuration values for the import process.
// NOTE: Other configuration parameters can be added here. The idea is to pass
// them to the import process using an instance of this struct rather than using
// global variables. TRY TO KEEP ANY INSTANCE OF THIS STRUCTURE READ-ONLY.
type ImportConfig struct {
	SystemName        string // the name of the imported system
	SystemDescription string // the description of the imported system

	IntoExistingSystem   bool   // true if it should be imported on a system that already exists
	ExistingSystemKey    string // the system key of the existing system
	ExistingSystemSecret string // the system secret of the existing system

	ImportUsers            bool   // true if users should be imported
	ImportRows             bool   // true if collection rows should be imported
	DefaultUserPassword    string // default password for users that don't have one already
	DefaultDeviceActiveKey string // default active key for devices that don't have one already
}

// DefaultImportConfig contains the default configuration values for the import
// process. Note that this instance SHOULD NOT be updated and used as a global
// configuration object. If you wish to configure the import processs using the
// global variables, check the NewImportConfigFromGlobals function.
//
// To create your own configuration config just assign this one to your own
// and modify it:
//
// ```
// customImportConfig := DefaultImportConfig
// customImportConfig.DefaultUserPassword = "my-new-password"
// ````
var DefaultImportConfig = ImportConfig{
	SystemName:        "",
	SystemDescription: "",

	IntoExistingSystem:   false,
	ExistingSystemKey:    "",
	ExistingSystemSecret: "",

	ImportUsers:            false,
	ImportRows:             false,
	DefaultUserPassword:    "",
	DefaultDeviceActiveKey: "",
}

// MakeImportConfigFromGlobals creates a new ImportConfig instance from the
// GLOBAL variables in cblib. Use with caution. Note that this function starts
// with Make* and not with New* because it returns a normal instance, and not
// a pointer to an instance.
func MakeImportConfigFromGlobals() ImportConfig {
	config := DefaultImportConfig

	// TODO: confirm which global to use
	config.ImportUsers = importUsers // or ImportUsers global?
	config.ImportRows = importRows   // or ImportRows global?

	return config
}

// ImportResult holds relevant values resulting from a system import process.
type ImportResult struct {
	rawSystemInfo map[string]interface{}
	SystemName    string
	SystemKey     string
	SystemSecret  string
}

// --------------------------------
// Import process (creation, etc)
// --------------------------------
// Functions that focus on the creation of the system and other assets.

func createSystem(config ImportConfig, system map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	name := system["name"].(string)
	desc := system["description"].(string)
	auth := system["auth"].(bool)
	sysKey, sysErr := client.NewSystem(name, desc, auth)
	if sysErr != nil {
		return nil, sysErr
	}
	realSystem, sysErr := client.GetSystem(sysKey)
	if sysErr != nil {
		return nil, sysErr
	}
	system["systemKey"] = realSystem.Key
	system["systemSecret"] = realSystem.Secret
	return system, nil
}

func createRoles(config ImportConfig, systemInfo map[string]interface{}, collectionsInfo []CollectionInfo, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	roles, err := getRoles()
	if err != nil {
		return err
	}
	for _, role := range roles {
		name := role["Name"].(string)
		fmt.Printf(" %s", name)
		//if name != "Authenticated" && name != "Administrator" && name != "Anonymous" {
		if err := createRole(sysKey, role, collectionsInfo, client); err != nil {
			return err
		}
		//}
	}
	fmt.Println("\nUpdating local roles with newly created role IDs... ")
	// ids were created on import for the new roles, grab those
	_, err = PullAndWriteRoles(sysKey, client, true)
	if err != nil {
		return err
	}

	return nil
}

func createUsers(config ImportConfig, systemInfo map[string]interface{}, users []map[string]interface{}, client *cb.DevClient) ([]UserInfo, error) {
	//  Create user columns first -- if any
	sysKey := systemInfo["systemKey"].(string)
	sysSec := systemInfo["systemSecret"].(string)
	userCols := []interface{}{}
	userSchema, err := getUserSchema()
	if err == nil {
		userCols = userSchema["columns"].([]interface{})
	}
	for _, columnIF := range userCols {
		column := columnIF.(map[string]interface{})
		columnName := column["ColumnName"].(string)
		if columnName == "user_id" || columnName == "cb_service_account" || columnName == "cb_ttl_override" || columnName == "cb_token" {
			fmt.Printf("Warning: ignoring exported '%s' column\n", columnName)
			continue
		}
		columnType := column["ColumnType"].(string)
		if err := client.CreateUserColumn(sysKey, columnName, columnType); err != nil {
			return nil, fmt.Errorf("Could not create user column %s: %s", columnName, err.Error())
		}
	}

	if !config.ImportUsers {
		return nil, nil
	}

	rtn := make([]UserInfo, 0)
	// Now, create users -- register, update roles, and update user-def colunms
	for _, user := range users {

		// TODO: added if to keep it backwards-compatible
		if len(config.DefaultUserPassword) > 0 {
			password := randSeq(10)
			maputil.SetIfMissing(user, password, config.DefaultUserPassword)
		}

		fmt.Printf(" %s", user["email"].(string))
		userID, err := createUser(sysKey, sysSec, user, client)
		if err != nil {
			// don't return an error because we don't want to stop other users from being created
			fmt.Printf("Error: Failed to create user %s - %s", user["email"].(string), err.Error())
		}
		info := UserInfo{
			UserID: userID,
			Email:  user["email"].(string),
		}
		rtn = append(rtn, info)
		if err := updateUserEmailToId(info); err != nil {
			logErrorForUpdatingMapFile(getUserEmailToIdFullFilePath(), err)
		}

		if len(userCols) == 0 {
			continue
		}

		updates := map[string]interface{}{}
		for _, columnIF := range userCols {
			column := columnIF.(map[string]interface{})
			columnName := column["ColumnName"].(string)
			if columnName != "user_id" {
				if userVal, ok := user[columnName]; ok {
					if userVal != nil {
						updates[columnName] = userVal
					}
				}
			}
		}

		if len(updates) == 0 {
			continue
		}

		if err := client.UpdateUser(sysKey, userID, updates); err != nil {
			// don't return an error because we don't want to stop other users from being updated
			fmt.Printf("Could not update user: %s", err.Error())
		}
	}

	return rtn, nil
}

func unMungeRoles(roles []string) []interface{} {
	rval := []interface{}{}

	for _, role := range roles {
		rval = append(rval, role)
	}
	return rval
}

func updateTriggerInfo(trigger map[string]interface{}, usersInfo []UserInfo) {
	replaceEmailWithUserIdInTriggerKeyValuePairs(trigger, usersInfo)
}

func createTriggerWithUpdatedInfo(config ImportConfig, sysKey string, trigger map[string]interface{}, usersInfo []UserInfo, client *cb.DevClient) (map[string]interface{}, error) {
	updateTriggerInfo(trigger, usersInfo)
	return createTrigger(sysKey, trigger, client)
}

func createTriggers(config ImportConfig, systemInfo map[string]interface{}, usersInfo []UserInfo, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	triggers, err := getTriggers()
	if err != nil {
		return nil, err
	}
	triggersRval := make([]map[string]interface{}, len(triggers))
	for idx, trigger := range triggers {
		fmt.Printf(" %s", trigger["name"].(string))
		trigVal, err := createTriggerWithUpdatedInfo(config, sysKey, trigger, usersInfo, client)
		if err != nil {
			return nil, err
		}
		triggersRval[idx] = trigVal
	}
	return triggersRval, nil
}

func createTimers(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	timers, err := getTimers()
	if err != nil {
		return nil, err
	}
	timersRval := make([]map[string]interface{}, len(timers))
	for idx, timer := range timers {
		fmt.Printf(" %s", timer["name"].(string))
		timerVal, err := createTimer(sysKey, timer, client)
		if err != nil {
			return nil, err
		}
		timersRval[idx] = timerVal
	}
	return timersRval, nil
}

func createDeployments(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	deployments, err := getDeployments()
	if err != nil {
		return nil, err
	}
	deploymentsRval := make([]map[string]interface{}, len(deployments))
	for idx, deployment := range deployments {
		fmt.Printf(" %s", deployment["name"].(string))
		deploymentVal, err := createDeployment(sysKey, deployment, client)
		if err != nil {
			return nil, err
		}
		deploymentsRval[idx] = deploymentVal
	}
	return deploymentsRval, nil
}

func createServiceCaches(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	caches, err := getServiceCaches()
	if err != nil {
		return nil, err
	}
	for _, cache := range caches {
		fmt.Printf(" %s", cache["name"].(string))
		err := createServiceCache(sysKey, cache, client)
		if err != nil {
			return nil, err
		}
	}
	return caches, nil
}

func createWebhooks(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	hooks, err := getWebhooks()
	if err != nil {
		return nil, err
	}
	for _, hook := range hooks {
		fmt.Printf(" %s", hook["name"].(string))
		err := createWebhook(sysKey, hook, client)
		if err != nil {
			return nil, err
		}
	}
	return hooks, nil
}

func createExternalDatabases(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	externalDatabases, err := getExternalDatabases()
	if err != nil {
		return nil, err
	}
	for _, externalDB := range externalDatabases {
		fmt.Printf(" %s", externalDB["name"].(string))
		err := createExternalDatabase(sysKey, externalDB, client)
		if err != nil {
			return nil, err
		}
	}
	return externalDatabases, nil
}

func createServices(config ImportConfig, systemInfo map[string]interface{}, usersInfo []UserInfo, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	services, err := getServices()
	if err != nil {
		fmt.Printf("getServices Failed: %s\n", err)
		return err
	}
	for _, service := range services {
		fmt.Printf(" %s", service["name"].(string))
		if err := createServiceWithUpdatedInfo(sysKey, service, usersInfo, client); err != nil {
			fmt.Printf("createService Failed: %s\n", err)
			return err
		}
	}
	return nil
}

func createLibraries(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	libraries, err := getLibraries()
	if err != nil {
		fmt.Printf("getLibraries Failed: %s\n", err)
		return err
	}
	for _, library := range libraries {
		fmt.Printf(" %s", library["name"].(string))
		if err := createLibrary(sysKey, library, client); err != nil {
			fmt.Printf("createLibrary Failed: %s\n", err)
			return err
		}
	}
	return nil
}

func createAdaptors(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	adaptors, err := getAdaptors(sysKey, client)
	if err != nil {
		return err
	}
	for i := 0; i < len(adaptors); i++ {
		err := createAdaptor(adaptors[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func createCollections(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]CollectionInfo, error) {
	sysKey := systemInfo["systemKey"].(string)
	collections, err := getCollections()
	rtn := make([]CollectionInfo, 0)
	if err != nil {
		return rtn, err
	}

	for _, collection := range collections {
		fmt.Printf(" %s\n", collection["name"].(string))
		if info, err := CreateCollection(sysKey, collection, config.ImportRows, client); err != nil {
			return rtn, err
		} else {
			rtn = append(rtn, info)
		}
	}
	return rtn, nil
}

// Reads Filesystem and makes HTTP calls to platform to create edges and edge columns
// Note: Edge schemas are optional, so if it is not found, we log an error and continue
func createEdges(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	sysSecret := systemInfo["systemSecret"].(string)
	edgesSchema, err := getEdgesSchema()
	if err != nil {
		// To ensure backwards-compatibility, we do not require
		// this folder `edges` to be present
		// As a result, let's log this error, but proceed
		fmt.Printf("Warning, could not find optional edge schema -- ignoring\n")
		return nil
	}

	edgesCols, ok := edgesSchema["columns"].([]interface{})
	if ok {
		for _, columnIF := range edgesCols {
			column := columnIF.(map[string]interface{})
			columnName := column["ColumnName"].(string)
			columnType := column["ColumnType"].(string)
			if err := client.CreateEdgeColumn(sysKey, columnName, columnType); err != nil {
				return fmt.Errorf("Could not create edges column %s: %s", columnName, err.Error())
			}
		}
	}

	edges, err := getEdges()
	if err != nil {
		return err
	}
	for _, edge := range edges {
		fmt.Printf(" %s", edge["name"].(string))
		edgeName := edge["name"].(string)
		delete(edge, "name")
		edge["system_key"] = sysKey
		edge["system_secret"] = sysSecret
		if err := createEdge(sysKey, edgeName, edge, client); err != nil {
			return err
		}
	}
	return nil
}

func createDevices(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	schemaPresent := true
	sysKey := systemInfo["systemKey"].(string)
	devicesSchema, err := getDevicesSchema()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			schemaPresent = false
		} else {
			return nil, err
		}
	}
	if schemaPresent {
		deviceCols, ok := devicesSchema["columns"].([]interface{})
		if ok {
			for _, columnIF := range deviceCols {
				column := columnIF.(map[string]interface{})
				columnName := column["ColumnName"].(string)
				if columnName == "salt" || columnName == "cb_service_account" || columnName == "cb_ttl_override" || columnName == "cb_token" {
					fmt.Printf("Warning: ignoring exported '%s' column\n", columnName)
					continue
				}
				columnType := column["ColumnType"].(string)
				if err := client.CreateDeviceColumn(sysKey, columnName, columnType); err != nil {
					fmt.Printf("Failed Creating device column %s\n", columnName)
					return nil, fmt.Errorf("Could not create devices column %s: %s", columnName, err.Error())
				}
				fmt.Printf("Created device column %s\n", columnName)
			}
		} else {
			return nil, fmt.Errorf("columns key not present in schema.json for devices")
		}
	}
	devices, err := getDevices()
	if err != nil {
		return nil, err
	}
	devicesRval := make([]map[string]interface{}, len(devices))
	for idx, device := range devices {
		if !schemaPresent {
			if idx == 0 {
				for columnname, _ := range device {
					switch strings.ToLower(columnname) {
					case "device_key", "name", "system_key", "type", "state", "description", "enabled", "allow_key_auth", "active_key", "keys", "allow_certificate_auth", "certificate", "created_date", "last_active_date":
						continue
					default:
						err := client.CreateDeviceColumn(sysKey, columnname, "string")
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}

		// TODO: added if to keep it backwards-compatible
		if len(config.DefaultDeviceActiveKey) > 0 {
			maputil.SetIfMissing(device, "active_key", config.DefaultDeviceActiveKey)
		}

		deviceName := device["name"].(string)
		fmt.Printf(" %s", deviceName)
		deviceInfo, err := createDevice(sysKey, device, client)
		if err != nil {
			return nil, err
		}
		deviceRoles, err := getDeviceRoles(deviceName)
		if err != nil {
			// system is probably in the legacy format, let's just set the roles to the default
			deviceRoles = convertStringSliceToInterfaceSlice([]string{"Authenticated"})
			logWarning(fmt.Sprintf("Could not find roles for device with name '%s'. This device will be created with only the default 'Authenticated' role.", deviceName))
		}
		defaultRoles := convertStringSliceToInterfaceSlice([]string{"Authenticated"})
		roleDiff := diffRoles(deviceRoles, defaultRoles)
		if err := client.UpdateDeviceRoles(sysKey, deviceName, convertInterfaceSliceToStringSlice(roleDiff.Added), convertInterfaceSliceToStringSlice(roleDiff.Removed)); err != nil {
			return nil, err
		}
		devicesRval[idx] = deviceInfo
	}
	return devicesRval, nil
}

func createPortals(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	var portals []map[string]interface{}
	var err error
	if hasLegacyPortalDirectory() {
		portals, err = getLegacyPortals()
		if err != nil {
			return nil, err
		}
	} else {
		portals, err = getCompressedPortals()
		if err != nil {
			return nil, err
		}
	}
	portalsRval := make([]map[string]interface{}, len(portals))
	for idx, dash := range portals {
		fmt.Printf(" %s", dash["name"].(string))
		portalInfo, err := createPortal(sysKey, dash, client)
		if err != nil {
			return nil, err
		}
		portalsRval[idx] = portalInfo
	}
	return portalsRval, nil
}

func createAllEdgeDeployment(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) error {
	//  First, look for deploy.json file. This is the new way of doing edge
	//  deployment. If that fails try the old way.
	if fileExists(rootDir + "/deploy.json") {
		info, err := getEdgeDeployInfo()
		if err != nil {
			return err
		}
		return createEdgeDeployInfo(config, systemInfo, info, client)
	}
	return oldCreateEdgeDeployInfo(systemInfo, client) // old deprecated way
}

func createEdgeDeployInfo(config ImportConfig, systemInfo, deployInfo map[string]interface{}, client *cb.DevClient) error {
	deployList := deployInfo["deployInfo"].([]interface{})
	sysKey := systemInfo["systemKey"].(string)

	for _, deployOneIF := range deployList {
		deployOne, ok := deployOneIF.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Poorly structured edge deploy info")
		}
		platform := deployOne["platform"].(bool)
		resName := deployOne["resource_identifier"].(string)
		resType := deployOne["resource_type"].(string)

		//  Go sdk expects the edge query to be in the Query format, not a string
		edgeQueryStr := deployOne["edge"].(string)
		_, err := client.CreateDeployResourcesForSystem(sysKey, resName, resType, platform, edgeQueryStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func oldCreateEdgeDeployInfo(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	edgeInfo, ok := systemInfo["edgeSync"].(map[string]interface{})
	if !ok {
		fmt.Printf("Warning: Could not find any edge sync information\n")
		return nil
	}
	for edgeName, edgeSyncInfoIF := range edgeInfo {
		edgeSyncInfo, ok := edgeSyncInfoIF.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Poorly formed edge sync info")
		}
		converted, err := convertOldEdgeDeployInfo(edgeSyncInfo)
		if err != nil {
			return err
		}
		_, err = client.SyncResourceToEdge(sysKey, edgeName, converted, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func createPlugins(config ImportConfig, systemInfo map[string]interface{}, client *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := systemInfo["systemKey"].(string)
	plugins, err := getPlugins()
	if err != nil {
		return nil, err
	}
	pluginsRval := make([]map[string]interface{}, len(plugins))
	for idx, plug := range plugins {
		fmt.Printf(" %s", plug["name"].(string))
		pluginVal, err := createPlugin(sysKey, plug, client)
		if err != nil {
			return nil, err
		}
		pluginsRval[idx] = pluginVal
	}
	return pluginsRval, nil
}

func convertOldEdgeDeployInfo(info map[string]interface{}) (map[string][]string, error) {
	rval := map[string][]string{
		"service": []string{},
		"library": []string{},
		"trigger": []string{},
	}
	for resourceKey, _ := range info {
		stuff := strings.Split(resourceKey, "::")
		if len(stuff) != 2 {
			return nil, fmt.Errorf("Poorly formed edge sync info entry: '%s'", resourceKey)
		}
		rval[stuff[0]] = append(rval[stuff[0]], stuff[1])
	}
	return rval, nil
}

func enableLogs(service map[string]interface{}) bool {
	logVal, ok := service["logging_enabled"]
	if !ok {
		return false
	}
	switch logVal.(type) {
	case string:
		return logVal.(string) == "true"
	case bool:
		return logVal.(bool) == true
	}
	return false
}

// TODO Handle more specific error for if folder doesnt exist
// i.e. plugins folder not found vs plugins import failed due to syntax error
// https://clearblade.atlassian.net/browse/CBCOMM-227
func importAllAssets(config ImportConfig, systemInfo map[string]interface{}, users []map[string]interface{}, cli *cb.DevClient) error {

	// Common set of calls for a complete system import

	logInfo("Importing collections...")
	collectionsInfo, err := createCollections(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create collections: %s", err.Error())
	}

	logInfo("Importing roles...")
	err = createRoles(config, systemInfo, collectionsInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create roles: %s", err.Error())
	}

	logInfo("Importing users...")
	usersInfo, err := createUsers(config, systemInfo, users, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create users: %s", err.Error())
	}

	logInfo("Importing code services...")
	// Additonal modifications to the ImportIt functions
	if err := createServices(config, systemInfo, usersInfo, cli); err != nil {
		serr, _ := err.(*os.PathError)
		if err != serr {
			return err
		} else {
			fmt.Printf("Path Error importing services: Operation: %s Path %s, Error %s\n", serr.Op, serr.Path, serr.Err)
			fmt.Printf("Warning: Could not import code services... -- ignoring\n")
		}
	}

	logInfo("Importing code libraries...")
	if err := createLibraries(config, systemInfo, cli); err != nil {
		serr, _ := err.(*os.PathError)
		if err != serr {
			return err
		} else {
			fmt.Printf("Path Error importing libraries: Operation: %s Path %s, Error %s\n", serr.Op, serr.Path, serr.Err)
			fmt.Printf("Warning: Could not import code libraries... -- ignoring\n")
		}
	}

	logInfo("Importing triggers...")
	_, err = createTriggers(config, systemInfo, usersInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create triggers: %s", err.Error())
	}

	logInfo("Importing timers...")
	_, err = createTimers(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create timers: %s", err.Error())
	}

	logInfo("Importing edges...")
	if err := createEdges(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create edges: %s", err.Error())
	}

	logInfo("Importing devices...")
	_, err = createDevices(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create devices: %s", err.Error())
	}

	logInfo("Importing portals...")
	_, err = createPortals(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create portals: %s", err.Error())
	}

	logInfo("Importing plugins...")
	_, err = createPlugins(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create plugins: %s", err.Error())
	}

	logInfo("Importing adaptors...")
	if err := createAdaptors(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create adaptors: %s", err.Error())
	}

	logInfo("Importing deployments...")
	if _, err := createDeployments(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create deployments: %s", err.Error())
	}

	logInfo("Importing shared caches...")
	if _, err := createServiceCaches(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create shared caches: %s", err.Error())
	}

	logInfo("Importing webhooks...")
	if _, err := createWebhooks(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create webhooks: %s", err.Error())
	}

	logInfo("Importing external databases...")
	if _, err := createExternalDatabases(config, systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create external databases: %s", err.Error())
	}

	fmt.Printf(" Done\n")
	logInfo(fmt.Sprintf("Success! New system key is: %s", systemInfo["systemKey"].(string)))
	logInfo(fmt.Sprintf("New system secret is: %s", systemInfo["systemSecret"].(string)))
	return nil
}

// --------------------------------
// Import entrypoint and exposed functions
// --------------------------------

// importSystem will import the system rooted at the given path using the given
// config. Please that we assume that the given clearblade client is already
// authorized an ready to use.
func importSystem(config ImportConfig, systemPath string, cli *cb.DevClient) (map[string]interface{}, error) {

	// points the root directory to the system folder
	// WARNING: side-effect (changes globals)
	SetRootDir(systemPath)

	// sets up director strcuture
	// WARNING: side-effect (might change system)
	err := setupDirectoryStructure()
	if err != nil {
		return nil, err
	}

	// gets users from the system directory
	// WARNING: side-effect (reads filesystem)
	users, err := getUsers()
	if err != nil {
		return nil, err
	}

	// gets system info from the system directory
	// WARNING: side-effect (reads filesystem)
	systemInfoPath := path.Join(systemPath, "system.json")
	systemInfo, err := getDict(systemInfoPath)
	if err != nil {
		return nil, err
	}

	// creates system if we are not importing into an existing one
	if !config.IntoExistingSystem {

		if len(config.SystemName) > 0 {
			systemInfo["name"] = config.SystemName
		}

		if len(config.SystemDescription) > 0 {
			systemInfo["description"] = config.SystemDescription
		}

		// NOTE: createSystem will modify systemInfo map
		_, err := createSystem(config, systemInfo, cli)
		if err != nil {
			return nil, fmt.Errorf("could not create system named '%s': %s", config.SystemName, err)
		}

	} else {
		systemInfo["systemKey"] = config.ExistingSystemKey
		systemInfo["systemSecret"] = config.ExistingSystemSecret
	}

	// import assets into created/existing system
	err = importAllAssets(config, systemInfo, users, cli)
	if err != nil {
		return nil, err
	}

	return systemInfo, nil
}

// ImportSystem imports the system rooted at the given path, using the default
// import config.
// NOTE: this is a transitional function, refer to ImportSystemUsingConfig for a
// newer and safer implementation.
func ImportSystem(cli *cb.DevClient, systemPath string, userInfo map[string]interface{}) (map[string]interface{}, error) {

	// authorizes the client BEFORE going into the import process. The import
	// process SHOULD NOT care about authorization
	// TODO: If the cli passed above is already valid, we don't need to
	// authorize again. Can we try removing this?
	cli, err := authorizeUsingMetaInfo(userInfo)
	if err != nil {
		return nil, err
	}

	// refactored userInfo into custom ImportConfig. That way we get rid of
	// the weakly-typed userInfo object and use a strongly-typed ImportConfig
	// instance
	config := MakeImportConfigFromGlobals()
	config.SystemName, _ = maputil.LookupString(userInfo, "systemName", "system_name")
	config.IntoExistingSystem, _ = maputil.LookupBool(userInfo, "importIntoExistingSystem", "import_into_existing_system")
	config.ExistingSystemKey, _ = maputil.LookupString(userInfo, "systemKey", "system_key")
	config.ExistingSystemSecret, _ = maputil.LookupString(userInfo, "systemSecret", "system_secret")

	importResult, err := ImportSystemUsingConfig(config, systemPath, cli)
	if err != nil {
		return nil, err
	}

	// we return raw (rather than an ImportResult instance), to keep backward
	// compatibility with code that was already using this function
	return importResult.rawSystemInfo, nil
}

// ImportSystemUsingConfig imports the system rooted at the given path, using the
// given config for different values. The given client should already be
// authenticated and ready to go.
func ImportSystemUsingConfig(config ImportConfig, systemPath string, cli *cb.DevClient) (ImportResult, error) {

	blankImportResult := ImportResult{}

	rawSystemInfo, err := importSystem(config, systemPath, cli)
	if err != nil {
		return blankImportResult, err
	}

	systemName, _ := maputil.LookupString(rawSystemInfo, "systemName", "system_name")
	systemKey, systemKeyOk := maputil.LookupString(rawSystemInfo, "systemKey", "system_key")
	systemSecret, systemSecretOk := maputil.LookupString(rawSystemInfo, "systemSecret", "system_secret")

	if !systemKeyOk || !systemSecretOk {
		return blankImportResult, fmt.Errorf("unable to extract system information")
	}

	result := ImportResult{
		rawSystemInfo: rawSystemInfo,
		SystemName:    systemName,
		SystemKey:     systemKey,
		SystemSecret:  systemSecret,
	}

	return result, nil
}
