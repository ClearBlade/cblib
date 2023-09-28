package cblib

import (
	"fmt"
	"io/ioutil"

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
	
	generateDiffCommand.flags.StringVar(&Path, "path", "", "Path where diff file will be created")
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
	diffLibraries, err := getDiffEntity(systemInfo.Key, "libraries", client);

	if err != nil {
		return err;
	}

	logInfo("Diffing services:");
	diffServices, err := getDiffEntity(systemInfo.Key, "services", client);

	if err != nil {
		return err;
	}

	dataMap := make(map[string]interface{});
	dataMap["services"] = diffServices;
	dataMap["libraries"] = diffLibraries;

	err = storeDataInJSONFile(dataMap, Path, "diff.json");

	if err != nil {
		return err
	}

	logInfo("Created a diff.json file");
	return nil;
}

func getDiffEntity(systemKey string, entityType string, client *cb.DevClient) ([]string, error) {
	diffEntities := []string{}

	rootPath := "./code/" + entityType;

	entities, err := ioutil.ReadDir(rootPath);

	if err != nil {
		fmt.Println("error reading directory ", err);
		return  nil,err;
	}

	for _, entity := range entities {
		fmt.Print(entity.Name() + " ");
		entityPath := rootPath+"/"+entity.Name()+"/"+entity.Name()

		localEntityCode, err := ioutil.ReadFile(entityPath+".js");
		if err != nil {
			return nil,err;
		}

		var remoteEntityCode map[string]interface{}

		if(entityType == "libraries") {
			remoteEntityCode, err = pullLibrary(systemKey, entity.Name(), client)
		} else {
			remoteEntityCode, err = pullService(systemKey, entity.Name(), client)
		}

		if err != nil {
			return nil, err;
		}

		if string(localEntityCode) != remoteEntityCode["code"] {
			diffEntities = append(diffEntities, entity.Name())
		}

	}
	fmt.Printf("\n")

	return diffEntities, nil;
}
