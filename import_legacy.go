package cblib

import (
	"fmt"
	"os"
	"strings"

	cb "github.com/clearblade/Go-SDK"

	"github.com/clearblade/cblib/models/bucketSetFiles"
	libPkg "github.com/clearblade/cblib/models/libraries"
	"github.com/clearblade/cblib/models/roles"
	"github.com/clearblade/cblib/types"
)

func createRoles(systemInfo *types.System_meta, client *cb.DevClient) error {

	roles, err := getRoles()
	if err != nil {
		return err
	}
	for _, role := range roles {
		name := role["Name"].(string)
		fmt.Printf(" %s", name)
		//if name != "Authenticated" && name != "Administrator" && name != "Anonymous" {
		if err := createRole(systemInfo, role, client); err != nil {
			return err
		}
		//}
	}
	fmt.Println("\nUpdating local roles with newly created role IDs... ")
	// ids were created on import for the new roles, grab those
	_, err = PullAndWriteRoles(systemInfo.Key, client, true)
	if err != nil {
		return err
	}

	return nil
}

func createUsers(config ImportConfig, systemInfo *types.System_meta, users []map[string]interface{}, client *cb.DevClient) ([]UserInfo, error) {
	//  Create user columns first -- if any
	userCols := []interface{}{}
	userSchema, err := getUserSchema()
	if err == nil {
		userColsIF, ok := userSchema["columns"]
		if ok && userColsIF != nil {
			userCols = userColsIF.([]interface{})
		}
	}
	for _, columnIF := range userCols {
		column := columnIF.(map[string]interface{})
		columnName := column["ColumnName"].(string)
		if columnName == "user_id" || columnName == "cb_service_account" || columnName == "cb_ttl_override" || columnName == "cb_token" {
			fmt.Printf("Warning: ignoring exported '%s' column\n", columnName)
			continue
		}
		columnType := column["ColumnType"].(string)
		if err := client.CreateUserColumn(systemInfo.Key, columnName, columnType); err != nil {
			return nil, fmt.Errorf("Could not create user column %s: %s", columnName, err.Error())
		}
	}

	if !config.ImportUsers {
		return nil, nil
	}

	rtn := make([]UserInfo, 0)
	// Now, create users -- register, update roles, and update user-def colunms
	for _, user := range users {
		fmt.Printf(" %s", user["email"].(string))
		userID, err := createUser(systemInfo.Key, systemInfo.Secret, user, client)
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
		if isServiceAccount, ok := user["cb_service_account"]; ok {
			updates["cb_service_account"] = isServiceAccount
		}

		if len(updates) == 0 {
			continue
		}

		if err := client.UpdateUser(systemInfo.Key, userID, updates); err != nil {
			// don't return an error because we don't want to stop other users from being updated
			fmt.Printf("Could not update user: %s", err.Error())
		}
	}

	return rtn, nil
}

func updateTriggerInfo(trigger map[string]interface{}, usersInfo []UserInfo) {
	replaceEmailWithUserIdInTriggerKeyValuePairs(trigger, usersInfo)
}

func createTriggerWithUpdatedInfo(sysKey string, trigger map[string]interface{}, usersInfo []UserInfo, client *cb.DevClient) (map[string]interface{}, error) {
	updateTriggerInfo(trigger, usersInfo)
	return createTrigger(sysKey, trigger, client)
}

func createTriggers(systemInfo *types.System_meta, usersInfo []UserInfo, client *cb.DevClient) ([]map[string]interface{}, error) {
	triggers, err := getTriggers()
	if err != nil {
		return nil, err
	}
	triggersRval := make([]map[string]interface{}, len(triggers))
	for idx, trigger := range triggers {
		fmt.Printf(" %s", trigger["name"].(string))
		trigVal, err := createTriggerWithUpdatedInfo(systemInfo.Key, trigger, usersInfo, client)
		if err != nil {
			return nil, err
		}
		triggersRval[idx] = trigVal
	}
	return triggersRval, nil
}

func createTimers(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	timers, err := getTimers()
	if err != nil {
		return nil, err
	}
	timersRval := make([]map[string]interface{}, len(timers))
	for idx, timer := range timers {
		fmt.Printf(" %s", timer["name"].(string))
		timerVal, err := createTimer(systemInfo.Key, timer, client)
		if err != nil {
			return nil, err
		}
		timersRval[idx] = timerVal
	}
	return timersRval, nil
}

func createDeployments(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	deployments, err := getDeployments()
	if err != nil {
		return nil, err
	}
	deploymentsRval := make([]map[string]interface{}, len(deployments))
	for idx, deployment := range deployments {
		fmt.Printf(" %s", deployment["name"].(string))
		deploymentVal, err := createDeployment(systemInfo.Key, deployment, client)
		if err != nil {
			return nil, err
		}
		deploymentsRval[idx] = deploymentVal
	}
	return deploymentsRval, nil
}

func createServiceCaches(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	caches, err := getServiceCaches()
	if err != nil {
		return nil, err
	}
	for _, cache := range caches {
		fmt.Printf(" %s", cache["name"].(string))
		err := createServiceCache(systemInfo.Key, cache, client)
		if err != nil {
			return nil, err
		}
	}
	return caches, nil
}

func createWebhooks(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	hooks, err := getWebhooks()
	if err != nil {
		return nil, err
	}
	for _, hook := range hooks {
		fmt.Printf(" %s", hook["name"].(string))
		err := createWebhook(systemInfo.Key, hook, client)
		if err != nil {
			return nil, err
		}
	}
	return hooks, nil
}

func createExternalDatabases(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	externalDatabases, err := getExternalDatabases()
	if err != nil {
		return nil, err
	}
	for _, externalDB := range externalDatabases {
		fmt.Printf(" %s", externalDB["name"].(string))
		err := createExternalDatabase(systemInfo.Key, externalDB, client)
		if err != nil {
			return nil, err
		}
	}
	return externalDatabases, nil
}

func createBucketSets(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	bucketSets, err := getBucketSets()
	if err != nil {
		return nil, err
	}
	for _, bucketSet := range bucketSets {
		fmt.Printf(" %s", bucketSet["name"].(string))
		err := createBucketSet(systemInfo.Key, bucketSet, client)
		if err != nil {
			return nil, err
		}
	}
	return bucketSets, nil
}

func createSecrets(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	secrets, err := getSecrets()
	if err != nil {
		return nil, err
	}
	for _, secret := range secrets {
		fmt.Printf(" %s", secret["name"].(string))
		err := createSecret(systemInfo.Key, secret, client)
		if err != nil {
			return nil, err
		}
	}
	return secrets, nil
}

func createServices(systemInfo *types.System_meta, client *cb.DevClient) error {
	services, err := getServices()
	if err != nil {
		fmt.Printf("getServices Failed: %s\n", err)
		return err
	}
	for _, service := range services {
		fmt.Printf(" %s", service["name"].(string))
		if err := createService(systemInfo.Key, service, client); err != nil {
			fmt.Printf("createService Failed: %s\n", err)
			return err
		}
	}
	return nil
}

func createLibraries(systemInfo *types.System_meta, client *cb.DevClient) error {
	rawLibraries, err := getLibraries()
	if err != nil {
		return err
	}

	libraries := make([]libPkg.Library, 0)
	for _, rawLib := range rawLibraries {
		libraries = append(libraries, libPkg.NewLibraryFromMap(rawLib))
	}

	orderedLibraries := libPkg.PostorderLibraries(libraries)

	for _, library := range orderedLibraries {
		fmt.Printf(" %s", library.GetName())
		if err := createLibrary(systemInfo.Key, library.GetMap(), client); err != nil {
			fmt.Printf("createLibrary Failed: %s\n", err)
			return err
		}
	}
	return nil
}

func createAdaptors(systemInfo *types.System_meta, client *cb.DevClient) error {
	adaptors, err := getAdaptors(systemInfo.Key, client)
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

func createCollections(config ImportConfig, systemInfo *types.System_meta, client *cb.DevClient) ([]CollectionInfo, error) {
	collections, err := getCollections()
	rtn := make([]CollectionInfo, 0)
	if err != nil {
		return rtn, err
	}

	for _, collection := range collections {
		fmt.Printf(" %s\n", collection["name"].(string))
		if info, err := CreateCollection(systemInfo.Key, collection, config.ImportRows, client); err != nil {
			return rtn, err
		} else {
			rtn = append(rtn, info)
		}
	}
	return rtn, nil
}

// Reads Filesystem and makes HTTP calls to platform to create edges and edge columns
// Note: Edge schemas are optional, so if it is not found, we log an error and continue
func createEdges(systemInfo *types.System_meta, client *cb.DevClient) error {
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
			if err := client.CreateEdgeColumn(systemInfo.Key, columnName, columnType); err != nil {
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
		edge["system_key"] = systemInfo.Key
		edge["system_secret"] = systemInfo.Secret
		if err := createEdge(systemInfo.Key, edgeName, edge, client); err != nil {
			return err
		}
	}
	return nil
}

func createDevices(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	schemaPresent := true
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
				if err := client.CreateDeviceColumn(systemInfo.Key, columnName, columnType); err != nil {
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
				for columnname := range device {
					switch strings.ToLower(columnname) {
					case "device_key", "name", "system_key", "type", "state", "description", "enabled", "allow_key_auth", "active_key", "keys", "allow_certificate_auth", "certificate", "created_date", "last_active_date":
						continue
					default:
						err := client.CreateDeviceColumn(systemInfo.Key, columnname, "string")
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}

		deviceName := device["name"].(string)
		fmt.Printf(" %s", deviceName)
		deviceInfo, err := createDevice(systemInfo.Key, device, client)
		if err != nil {
			return nil, err
		}
		deviceRoles, err := getDeviceRoles(deviceName)
		if err != nil {
			// system is probably in the legacy format, let's just set the roles to the default
			deviceRoles = []string{"Authenticated"}
			logWarning(fmt.Sprintf("Could not find roles for device with name '%s'. This device will be created with only the default 'Authenticated' role.", deviceName))
		}
		defaultRoles := []string{"Authenticated"}
		roleDiff := roles.DiffRoles(deviceRoles, defaultRoles)
		if err := client.UpdateDeviceRoles(systemInfo.Key, deviceName, roleDiff.Added, roleDiff.Removed); err != nil {
			return nil, err
		}
		devicesRval[idx] = deviceInfo
	}
	return devicesRval, nil
}

func createPortals(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
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
		portalInfo, err := createPortal(systemInfo.Key, dash, client)
		if err != nil {
			return nil, err
		}
		portalsRval[idx] = portalInfo
	}
	return portalsRval, nil
}

func createPlugins(systemInfo *types.System_meta, client *cb.DevClient) ([]map[string]interface{}, error) {
	plugins, err := getPlugins()
	if err != nil {
		return nil, err
	}
	pluginsRval := make([]map[string]interface{}, len(plugins))
	for idx, plug := range plugins {
		fmt.Printf(" %s", plug["name"].(string))
		pluginVal, err := createPlugin(systemInfo.Key, plug, client)
		if err != nil {
			return nil, err
		}
		pluginsRval[idx] = pluginVal
	}
	return pluginsRval, nil
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
func importAllAssetsLegacy(config ImportConfig, systemInfo *types.System_meta, users []map[string]interface{}, cli *cb.DevClient) error {
	// Common set of calls for a complete system import
	logInfo("Importing collections...")
	_, err := createCollections(config, systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create collections: %s", err.Error())
	}

	logInfo("Importing roles...")
	err = createRoles(systemInfo, cli)
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

	logInfo("Importing code libraries...")
	if err := createLibraries(systemInfo, cli); err != nil {
		serr, _ := err.(*os.PathError)
		if err != serr {
			return err
		} else {
			fmt.Printf("Path Error importing libraries: Operation: %s Path %s, Error %s\n", serr.Op, serr.Path, serr.Err)
			fmt.Printf("Warning: Could not import code libraries... -- ignoring\n")
		}
	}

	logInfo("Importing shared caches...")
	if _, err := createServiceCaches(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create shared caches: %s", err.Error())
	}

	logInfo("Importing code services...")
	// Additonal modifications to the ImportIt functions
	if err := createServices(systemInfo, cli); err != nil {
		serr, _ := err.(*os.PathError)
		if err != serr {
			return err
		} else {
			fmt.Printf("Path Error importing services: Operation: %s Path %s, Error %s\n", serr.Op, serr.Path, serr.Err)
			fmt.Printf("Warning: Could not import code services... -- ignoring\n")
		}
	}

	logInfo("Importing triggers...")
	_, err = createTriggers(systemInfo, usersInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create triggers: %s", err.Error())
	}

	logInfo("Importing timers...")
	_, err = createTimers(systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create timers: %s", err.Error())
	}

	logInfo("Importing edges...")
	if err := createEdges(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create edges: %s", err.Error())
	}

	logInfo("Importing devices...")
	_, err = createDevices(systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create devices: %s", err.Error())
	}

	logInfo("Importing portals...")
	_, err = createPortals(systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create portals: %s", err.Error())
	}

	logInfo("Importing plugins...")
	_, err = createPlugins(systemInfo, cli)
	if err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create plugins: %s", err.Error())
	}

	logInfo("Importing adaptors...")
	if err := createAdaptors(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create adaptors: %s", err.Error())
	}

	logInfo("Importing deployments...")
	if _, err := createDeployments(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create deployments: %s", err.Error())
	}

	logInfo("Importing webhooks...")
	if _, err := createWebhooks(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create webhooks: %s", err.Error())
	}

	logInfo("Importing external databases...")
	if _, err := createExternalDatabases(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create external databases: %s", err.Error())
	}

	logInfo("Importing bucket sets...")
	if _, err := createBucketSets(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create bucket sets: %s", err.Error())
	}

	logInfo("Importing bucket set files...")
	if err := bucketSetFiles.PushFilesForAllBucketSets(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not import bucket set files: %s", err.Error())
	}

	logInfo("Importing secrets...")
	if _, err := createSecrets(systemInfo, cli); err != nil {
		//  Don't return an err, just warn -- so we keep back compat with old systems
		fmt.Printf("Could not create secrets: %s", err.Error())
	}

	logInfo("Importing message history storage...")
	if err := pushMessageHistoryStorage(systemInfo, cli); err != nil {
		fmt.Printf("Could not import message history storage: %s", err.Error())
	}

	logInfo("Importing message type triggers...")
	if err := pushMessageTypeTriggers(systemInfo, cli); err != nil {
		fmt.Printf("Could not import message type triggers: %s", err.Error())
	}

	fmt.Printf(" Done\n")
	logInfo(fmt.Sprintf("Success! New system key is: %s", systemInfo.Key))
	logInfo(fmt.Sprintf("New system secret is: %s", systemInfo.Secret))
	return nil
}
