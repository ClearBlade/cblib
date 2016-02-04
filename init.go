package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
)

func init() {
	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	rolesInfo = []map[string]interface{}{}
	myInitCommand := &SubCommand{
		name:  "export",
		usage: "Ain't no thing",
		run:   doInit,
		//  TODO -- add help, usage, etc.
	}
	AddCommand("init", myInitCommand)
}

func doInit(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) == 0 {
		fmt.Printf("init command: missing system key\n")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Printf("init command: too many arguments\n")
		os.Exit(1)
	}
	return reallyInit(client, args[0])
}

func reallyInit(cli *cb.DevClient, sysKey string) error {
	sysMeta, err := pullSystemMeta(sysKey, cli)
	if err != nil {
		return err
	}

	setRootDir(strings.Replace(sysMeta.Name, " ", "_", -1))
	if err := setupDirectoryStructure(sysMeta); err != nil {
		return err
	}
	storeMeta(sysMeta)

	if err = storeSystemDotJSON(systemDotJSON); err != nil {
		return err
	}

	metaStuff := map[string]interface{}{
		"platformURL":       URL,
		"developerEmail":    Email,
		"assetRefreshDates": []interface{}{},
		"token":             DevToken,
	}
	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been initialized into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}
