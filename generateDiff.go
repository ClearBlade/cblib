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

	dataMap := make(map[string]interface{});
	dataMap["services"] = diffServices;
	dataMap["libraries"] = diffLibraries;
	dataMap["devices"] = diffDevices;
	dataMap["devicesRoles"] = diffDeviceRoles;
	dataMap["deviceSchema"] = diffDeviceSchema;
	dataMap["edges"] = diffEdges;
	dataMap["edgesSchema"] = diffEdgesSchema;
	dataMap["sharedCaches"] = diffSharedCaches;
	dataMap["timers"] = diffTimers


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

		remoteEdge := getCorrespondingRemoteEdge(localEdgeName, remoteEdges)

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

func areLocalAndRemoteSchemaEqual(localSchema map[string]interface{}, remoteSchema map[string]interface{}) bool {
	// we need this util function because localSchema has entry like map[columns:<nil>] to show that there are no columns
  // whereas remoteSchema has entry like map[columns:[]] to show that there are no columns

	if localSchema["columns"] == nil && len(remoteSchema["columns"].([]interface{})) == 0 {
		return true;
	} else {
		return reflect.DeepEqual(localSchema["columns"], remoteSchema["columns"])
	}
}

func getCorrespondingRemoteEdge(name string, arr []interface{}) interface{} {
	for _, val := range arr {
		if val.(map[string]interface{})["name"] == name {
			return val;
		}
	}

	return nil;
}