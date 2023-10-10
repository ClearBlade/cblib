package cblib

import (
	"fmt"
	"reflect"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/types"
)

func init() {
	usage :=
		`
	Generates a diff between the code services and libraries present locally and in the system
	`

	example :=
		`
	  cb-cli diff # generate a diff file at the current path
		cb-cli diff -path='./diff/' # generate a diff file inside the diff folder
		`
	generateDiffCommand := &SubCommand{
		name:      "diff",
		usage:     usage,
		needsAuth: true,
		run:       doGenerateDiff,
		example:   example,
	}
	
	generateDiffCommand.flags.StringVar(&PathForDiffFile, "path", "", "Relative path where diff file will be created")
	AddCommand("diff", generateDiffCommand)
}

func doGenerateDiff(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	systemInfo, err := getSysMeta()
	if err != nil {
		return err;
	}

	client, err = checkIfTokenHasExpired(client, systemInfo.Key);
	if err != nil {
		return fmt.Errorf("Re-auth failed: %s\n", err)
	}

	logInfo("Diffing libraries:");
	diffLibraries, err := getLibrariesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing services:");
	diffServices, err := getServicesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing devices, deviceRoles and deviceSchema:");
	diffDevices, diffDeviceRoles, diffDeviceSchema, err := getDevicesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing edges:");
	diffEdges, diffEdgesSchema, err := getEdgesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing service caches: ");
	diffSharedCaches, err := getSharedCacheDiff(systemInfo, client)
	
	if err != nil {
		return err;
	}

	logInfo("Diffing timers:");
	diffTimers, err := getTimersDiff(systemInfo.Key, client)

	if err != nil {
		return err;
	}

	logInfo("Diffing users:");
	diffUsers, diffUserRoles, diffUserSchema, err := getUsersDiff(systemInfo.Key, client)

	if err != nil {
		return err
	}

	logInfo("Diffing webhooks:");
	diffWebhooks, err := getWebhooksDiff(systemInfo, client)

	if err != nil {
		return err
	}

	logInfo("Diffing messageHistoryStorage:");

	diffMessageHistoryStorage, err := getMessageHistoryStorageDiff(systemInfo, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing messageTypeTriggers:");
	diffMessageTypeTriggers, err := getMessageTypeTriggersDiff(systemInfo, client);

	if err != nil {
		return err;
	}

	logInfo("Diffing roles:");
	diffRoles, err := getRolesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	dataMap := make(map[string]interface{});
	dataMap["services"] = diffServices;
	dataMap["libraries"] = diffLibraries;
	dataMap["devices"] = diffDevices;
	dataMap["deviceRoles"] = diffDeviceRoles;
	dataMap["deviceSchema"] = diffDeviceSchema;
	dataMap["edges"] = diffEdges;
	dataMap["edgesSchema"] = diffEdgesSchema;
	dataMap["sharedCaches"] = diffSharedCaches;
	dataMap["timers"] = diffTimers
	dataMap["users"] = diffUsers;
	dataMap["userRoles"] = diffUserRoles;
	dataMap["userSchema"] = diffUserSchema;
	dataMap["webhooks"] = diffWebhooks;
	dataMap["messageHistoryStorage"] = diffMessageHistoryStorage
	dataMap["messageTypeTriggers"] = diffMessageTypeTriggers
	dataMap["roles"] = diffRoles;

	err = storeDataInJSONFile(dataMap, PathForDiffFile, "diff.json");

	if err != nil {
		return err
	}

	logInfo("Created a diff.json file");
	return nil;
}

func getLibrariesDiff(systemKey string, client *cb.DevClient) ([]string, error) {
	librariesDiff := []string{}

	localLibraries, err := getLibraries()

	if err != nil {
		return nil, err
	}

	for _, localLibrary := range localLibraries {
		localLibraryName := localLibrary["name"].(string)

		fmt.Printf(localLibraryName + " ");

		remoteLibrary, err := pullLibrary(systemKey, localLibraryName, client)

		if err != nil {
			librariesDiff = append(librariesDiff, localLibraryName)
			continue
		}

		localLibrary, remoteLibrary = keepCommonKeysFromMaps(localLibrary, remoteLibrary)

		if !reflect.DeepEqual(localLibrary, remoteLibrary) {
			librariesDiff = append(librariesDiff, localLibraryName)
		}
	}

	fmt.Printf("\n")
	return librariesDiff, nil
}

func getServicesDiff(systemKey string, client *cb.DevClient) ([]string, error) {
	servicesDiff := []string{}

	localServices, err := getServices()

	if err != nil {
		return nil, err
	}

	for _, localService := range localServices {
		localServiceName := localService["name"].(string)

		fmt.Printf(localServiceName + " ");

		remoteService, err := pullService(systemKey, localServiceName, client)

		if err != nil {
			servicesDiff = append(servicesDiff, localServiceName)
			continue
		}

		localService, remoteService = keepCommonKeysFromMaps(localService, remoteService)

		if !reflect.DeepEqual(localService, remoteService) {
			servicesDiff = append(servicesDiff, localServiceName)
		}
	}

	fmt.Printf("\n")
	return servicesDiff, nil
}

func getDevicesDiff(systemKey string, client *cb.DevClient) ([]string, []string, bool, error) {
	// does diffing for devices, deviceRoles as well as deviceSchema

	devicesDiff := []string{}
	deviceRolesDiff := []string{}

	localDevices, err := getDevices()
	
	if err != nil {
		return nil, nil, false, err
	}

	for _, localDevice := range localDevices {
		localDeviceName := (localDevice["name"]).(string)

		fmt.Printf(localDeviceName + " ");

		remoteDevice, err := pullDevice(systemKey, localDeviceName, client)

		if err != nil {
			// remoteDevice not found in the system, but is present locally. This should be added in our diff
			devicesDiff = append(devicesDiff, localDeviceName)
			deviceRolesDiff = append(deviceRolesDiff, localDeviceName)
			continue;
		}

		localDevice, remoteDevice = keepCommonKeysFromMaps(localDevice, remoteDevice)

		if !reflect.DeepEqual(localDevice, remoteDevice) {
			devicesDiff = append(devicesDiff, localDeviceName);
		}

		localDeviceRole, err := getDeviceRoles(localDeviceName)
		if err != nil {
			return nil, nil, false, err
		}

		remoteDeviceRole, err := pullDeviceRoles(systemKey, localDeviceName, client)
		if err != nil {
			return nil, nil, false, err
		}

		if !reflect.DeepEqual(convertInterfaceSliceToStringSlice(localDeviceRole), remoteDeviceRole) {
			deviceRolesDiff = append(deviceRolesDiff, localDeviceName)
		}
	}

	localDeviceSchema, err := getDevicesSchema()
	if err != nil {
		return nil, nil, false, err;
	}

	remoteDeviceSchema, err := pullDevicesSchema(systemKey, client, false)
	if err != nil {
		return nil, nil, false, err;
	}

	if !areLocalAndRemoteSchemaEqual(localDeviceSchema, remoteDeviceSchema) {
		return devicesDiff, deviceRolesDiff, true, nil;
	} else {
		return devicesDiff, deviceRolesDiff, false, err
	}
}

func getEdgesDiff(systemKey string, client *cb.DevClient) ([]string, bool, error) {
	edgesDiff := []string{}

	localEdges, err := getEdges()
	if err != nil {
		return nil, false, err;
	}

	// had to use pullAllEdges instead of pullEdge inside the for loop because pullEdge lacked some information
	// like last_connect, last_disconnect
	remoteEdges, err := pullAllEdges(systemKey, client)
	if err != nil {
		return nil, false, err;
	}

	for _, localEdge := range localEdges {
		localEdgeName := localEdge["name"].(string);
		fmt.Printf(localEdgeName + " ");

		remoteEdge := getCorrespondingRemoteMap("name", localEdgeName, remoteEdges)

		if remoteEdge == nil {
			// remoteEdge not present. Add this into diff
			edgesDiff = append(edgesDiff, localEdgeName);
			continue;
		}

		localEdge, remoteEdge = keepCommonKeysFromMaps(localEdge, remoteEdge.(map[string]interface{}))

		if !reflect.DeepEqual(localEdge, remoteEdge) {
			edgesDiff = append(edgesDiff, localEdgeName)
		}
	}

	localEdgesSchema, err := getEdgesSchema()
	if err != nil {
		return nil, false, err;
	}

	remoteEdgesSchema, err := pullEdgesSchema(systemKey, client, false)
	if err != nil {
		return nil, false, err;
	}

	if !areLocalAndRemoteSchemaEqual(localEdgesSchema, remoteEdgesSchema) {
		return edgesDiff, true, nil;
	} else {
		return edgesDiff, false, err
	}
}

func getMessageHistoryStorageDiff(systemInfo *types.System_meta, client *cb.DevClient) (bool, error) {
	localMessageHistoryStorage, err := getMessageHistoryStorage();

	if err != nil {
		return false, err;
	}

	remoteMessageHistoryStorage, err := pullMessageHistoryStorage(systemInfo, client)

	if err != nil {
		return false, err;
	}

	for ind, localMessageHistoryStorageItem := range localMessageHistoryStorage {
		remoteMessageHistoryStorageItem := remoteMessageHistoryStorage[ind].MessageHistoryStorageEntry
		if !reflect.DeepEqual(localMessageHistoryStorageItem, remoteMessageHistoryStorageItem) {
			return true, nil;
		}
	}

	return false, nil;
}

func getMessageTypeTriggersDiff(systemInfo *types.System_meta, client *cb.DevClient) (bool, error) {
	localMessageTypeTriggers, err := getMessageTypeTriggers();

	if err != nil {
		return false, err;
	}

	remoteMessageTypeTriggers, err := pullMessageTypeTriggers(systemInfo, client)

	if err != nil {
		return false, err;
	}

	for ind, localMessageTypeTrigger := range localMessageTypeTriggers {
		remoteMessageTypeTrigger:= remoteMessageTypeTriggers[ind]

		localMessageTypeTrigger, remoteMessageTypeTrigger = keepCommonKeysFromMaps(localMessageTypeTrigger, remoteMessageTypeTrigger)
		if !reflect.DeepEqual(localMessageTypeTrigger, remoteMessageTypeTrigger) {
			return true, nil;
		}
	}

	return false, nil;
}

func getSharedCacheDiff(systemInfo *types.System_meta, client *cb.DevClient) ([]string, error) {
	sharedCacheDiff := []string{}
	localSharedCaches, err := getServiceCaches()

	if err != nil {
		return nil, err
	}

	for _, localSharedCache := range localSharedCaches {
		localSharedCacheName := localSharedCache["name"].(string)
		
		fmt.Printf(localSharedCacheName + " ");
		
		remoteSharedCache, err := pullAndWriteServiceCache(systemInfo, client, localSharedCacheName, false)

		if err != nil {
			// this service cache is not present in the remote system. Add it to the diff
			sharedCacheDiff = append(sharedCacheDiff, localSharedCacheName)
			continue;
		}

		localSharedCache, remoteSharedCache = keepCommonKeysFromMaps(localSharedCache, remoteSharedCache)

		if !reflect.DeepEqual(localSharedCache, remoteSharedCache) {
			sharedCacheDiff = append(sharedCacheDiff, localSharedCacheName)
		}
	}

	return sharedCacheDiff, nil;
}

func getTimersDiff(systemKey string, client *cb.DevClient) ([]string, error) {
	timersDiff := []string{}
	localTimers, err := getTimers()

	if err != nil {
		return nil, err
	}

	for _, localTimer := range localTimers {
		localTimerName := localTimer["name"].(string)
		
		fmt.Printf(localTimerName + " ");
		
		remoteTimer, err := pullTimer(systemKey, localTimerName, client)

		if err != nil {
			// this timer is not present in the remote system. Add it to the diff
			timersDiff = append(timersDiff, localTimerName)
			continue;
		}

		localTimer, remoteTimer = keepCommonKeysFromMaps(localTimer, remoteTimer)

		if !reflect.DeepEqual(localTimer, remoteTimer) {
			timersDiff = append(timersDiff, localTimerName)
		}
	}

	return timersDiff, nil;
}

func getUsersDiff(systemKey string, client *cb.DevClient) ([]string, []string, bool, error) {
	// does diffing for users, userRoles as well as userSchema

	usersDiff := []string{}
	userRolesDiff := []string{}
	localUsers, err := getUsers()

	if err != nil {
		return nil, nil, false, err
	}

	remoteUsers, err := pullAllUsers(systemKey, client)

	if err != nil {
		return nil, nil, false, err
	}

	for _, localUser := range localUsers {
		localUserEmail := localUser["email"].(string)
		remoteUser := getCorrespondingRemoteMap("email", localUserEmail, remoteUsers)
		userId := remoteUser.(map[string]interface{})["user_id"].(string)

		fmt.Printf(userId + " ");

		localUser, remoteUser = keepCommonKeysFromMaps(localUser, remoteUser.(map[string]interface{}))

		if !reflect.DeepEqual(localUser, remoteUser) {
			usersDiff = append(usersDiff, localUserEmail)
		}

		localUserRoles, err := getUserRoles(localUserEmail)

		if err != nil {
			return nil, nil, false, err
		}

		convertedLocalUserRoles := convertInterfaceSliceToStringSlice(localUserRoles)

		remoteUserRoles, err := pullUserRoles(systemKey, userId, client)
		if err != nil {
			return nil, nil, false, err
		}

		if !reflect.DeepEqual(convertedLocalUserRoles, remoteUserRoles) {
			userRolesDiff = append(userRolesDiff, localUserEmail)
		}
	}

	localUserSchema, err := getUserSchema()

	if err != nil {
		return nil, nil, false, err
	}

	remoteUserSchema, err := pullUserSchemaInfo(systemKey, client, false)

	if err != nil {
		return nil, nil, false, err
	}

	diffUserSchema := !areLocalAndRemoteSchemaEqual(localUserSchema, remoteUserSchema)

	return usersDiff, userRolesDiff, diffUserSchema, nil;
}

func getWebhooksDiff(systemInfo *types.System_meta, client *cb.DevClient) ([]string, error) {
	webhooksDiff := []string{}
	localWebhooks, err := getWebhooks()

	if err != nil {
		return nil, err
	}

	for _, localWebhook := range localWebhooks {
		localWebhookName := localWebhook["name"].(string)
		
		fmt.Printf(localWebhookName + " ");
		
		remoteWebhook, err := pullAndWriteWebhook(systemInfo, client, localWebhookName, false)

		if err != nil {
			// this webhook is not present in the remote system. Add it to the diff
			webhooksDiff = append(webhooksDiff, localWebhookName)
			continue;
		}

		localWebhook, remoteWebhook = keepCommonKeysFromMaps(localWebhook, remoteWebhook)

		if !reflect.DeepEqual(localWebhook, remoteWebhook) {
			webhooksDiff = append(webhooksDiff, localWebhookName)
		}
	}

	return webhooksDiff, nil;

}

func getRolesDiff(systemKey string, client *cb.DevClient) ([]string, error) {
	rolesDiff := []string{}

	localRoles, err := getRoles()

	if err != nil {
		return nil, err
	}

	for _, localRole := range localRoles {
		localRoleName := localRole["Name"].(string)

		fmt.Printf(localRoleName + " ");

		remoteRole, err := pullRole(systemKey, localRoleName, client)

		localRolePerms, castSuccess := localRole["Permissions"].(map[string]interface{})
		if !castSuccess {
			return nil, fmt.Errorf("Unable to process role permissions: %v", localRole["Permissions"])
		}

		localRoleCodeServices, castSuccess := localRolePerms["CodeServices"].([]interface{})
		if castSuccess {
			sortByMapKey(&localRoleCodeServices, SORT_KEY_CODE_SERVICE)
			fmtServices := whitelistServicesPermissions(localRoleCodeServices)
			localRolePerms["CodeServices"] = fmtServices
		}

		localRoleCollections, castSuccess := localRolePerms["Collections"].([]interface{})
		if castSuccess {
			sortByMapKey(&localRoleCollections, SORT_KEY_COLLECTION)
			fmtCollections := whitelistCollectionsPermissions(localRoleCollections)
			localRolePerms["Collections"] = fmtCollections
		}

		localRolePortals, castSuccess := localRolePerms["Portals"].([]interface{})
		if castSuccess {
			sortByMapKey(&localRolePortals, SORT_KEY_PORTAL)
			fmtPortals := whitelistPortalsPermissions(localRolePortals)
			localRolePerms["Portals"] = fmtPortals
		}

		localRoleServiceCaches, castSuccess := localRolePerms["ServiceCaches"].([]interface{})
		if castSuccess {
			sortByMapKey(&localRoleServiceCaches, SORT_KEY_SERVICE_CACHES)
			localRolePerms["ServiceCaches"] = localRoleServiceCaches
		}

		localRoleTopics, castSuccess := localRolePerms["Topics"].([]interface{})
		if castSuccess {
			sortByMapKey(&localRoleTopics, SORT_KEY_TOPICS)
			localRolePerms["Topics"] = localRoleTopics
		}

		remoteRolePerms, castSuccess := remoteRole["Permissions"].(map[string]interface{})
		if !castSuccess {
			return nil, fmt.Errorf("Unable to process remote role permissions: %v", remoteRole["Permissions"])
		}

		remoteRoleCodeServices, castSuccess := remoteRolePerms["CodeServices"].([]interface{})
		if castSuccess {
			sortByMapKey(&remoteRoleCodeServices, SORT_KEY_CODE_SERVICE)
			fmtServices := whitelistServicesPermissions(remoteRoleCodeServices)
			remoteRolePerms["CodeServices"] = fmtServices
		}

		remoteRoleCollections, castSuccess := remoteRolePerms["Collections"].([]interface{})
		if castSuccess {
			sortByMapKey(&remoteRoleCollections, SORT_KEY_COLLECTION)
			fmtCollections := whitelistCollectionsPermissions(remoteRoleCollections)
			remoteRolePerms["Collections"] = fmtCollections
		}

		remoteRolePortals, castSuccess := remoteRolePerms["Portals"].([]interface{})
		if castSuccess {
			sortByMapKey(&remoteRolePortals, SORT_KEY_PORTAL)
			fmtPortals := whitelistPortalsPermissions(remoteRolePortals)
			remoteRolePerms["Portals"] = fmtPortals
		}

		remoteRoleServiceCaches, castSuccess := remoteRolePerms["ServiceCaches"].([]interface{})
		if castSuccess {
			sortByMapKey(&remoteRoleServiceCaches, "Name")
			remoteRolePerms["ServiceCaches"] = remoteRoleServiceCaches
		}

		remoteRoleTopics, castSuccess := remoteRolePerms["Topics"].([]interface{})
		if castSuccess {
			sortByMapKey(&remoteRoleTopics, "Name")
			remoteRolePerms["Topics"] = remoteRoleTopics
		}

		if err != nil {
			rolesDiff = append(rolesDiff, localRoleName)
			continue
		}

		localRole, remoteRole = keepCommonKeysFromMaps(localRole, remoteRole)

		if !reflect.DeepEqual(localRolePerms["Topics"], remoteRolePerms["Topics"]) {
			rolesDiff = append(rolesDiff, localRoleName);
			continue;
		}

		if !reflect.DeepEqual(localRolePerms["CodeServices"], remoteRolePerms["CodeServices"]) {
			rolesDiff = append(rolesDiff, localRoleName);
			continue;
		}

		if !reflect.DeepEqual(localRolePerms["Portals"], remoteRolePerms["Portals"]) {
			rolesDiff = append(rolesDiff, localRoleName);
			continue;
		}

		if !reflect.DeepEqual(localRolePerms["Collections"], remoteRolePerms["Collections"]) {
			rolesDiff = append(rolesDiff, localRoleName);
			continue;
		}

		if !reflect.DeepEqual(localRolePerms["ServiceCaches"], remoteRolePerms["ServiceCaches"]) {
			rolesDiff = append(rolesDiff, localRoleName);
			continue;
		}

		if !reflect.DeepEqual(localRole, remoteRole) {
			rolesDiff = append(rolesDiff, localRoleName)
		}
	}

	fmt.Printf("\n")
	return rolesDiff, nil
}

func areLocalAndRemoteSchemaEqual(localSchema map[string]interface{}, remoteSchema map[string]interface{}) bool {
	// we need this util function because localSchema has entry like map[columns:<nil>] to show that there are no columns
  // whereas remoteSchema has entry like map[columns:[]] to show that there are no columns

	if localSchema["columns"] == nil && len(remoteSchema["columns"].([]interface{})) == 0 {
		return true;
	} else {
		return reflect.DeepEqual(localSchema["columns"], remoteSchema["columns"])
	}
}

func getCorrespondingRemoteMap(property string, value string, arr []interface{}) interface{} {
	for _, val := range arr {
		if val.(map[string]interface{})[property] == value {
			return val;
		}
	}

	return nil;
}
