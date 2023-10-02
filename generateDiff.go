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

	dataMap := make(map[string]interface{});
	dataMap["services"] = diffServices;
	dataMap["libraries"] = diffLibraries;

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

	fmt.Println(librariesDiff)

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

	fmt.Println(servicesDiff)
	fmt.Printf("\n")
	return servicesDiff, nil
}
