package cblib

import (
	"fmt"
	"reflect"

	cb "github.com/clearblade/Go-SDK"
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

	logInfo("Diffing devices and deviceRoles:");
	diffDevices, diffDeviceRoles, diffDeviceSchema, err := getDevicesDiff(systemInfo.Key, client);

	if err != nil {
		return err;
	}

	dataMap := make(map[string]interface{});
	dataMap["services"] = diffServices;
	dataMap["libraries"] = diffLibraries;
	dataMap["devices"] = diffDevices;
	dataMap["devicesRoles"] = diffDeviceRoles;
	dataMap["deviceSchema"] = diffDeviceSchema;

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

	if !areLocalAndRemoteDeviceSchemaEqual(localDeviceSchema, remoteDeviceSchema) {
		return devicesDiff, deviceRolesDiff, true, nil;
	} else {
		return devicesDiff, deviceRolesDiff, false, err
	}
}

func areLocalAndRemoteDeviceSchemaEqual(localDeviceSchema map[string]interface{}, remoteDeviceSchema map[string]interface{}) bool {
	// we need this util function because localDeviceSchema has entry like map[columns:<nil>] to show that there are no columns
  // whereas remoteDeviceSchema has entry like map[columns:[]] to show that there are no columns

	if localDeviceSchema["columns"] == nil && len(remoteDeviceSchema["columns"].([]interface{})) == 0 {
		return true;
	} else {
		return reflect.DeepEqual(localDeviceSchema["columns"], remoteDeviceSchema["columns"])
	}
}