package cblib

import (
	//"flag"
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
		name:            "init",
		usage:           "Ain't no thing",
		needsAuth:       true,
		mustBeInRepo:    false,
		mustNotBeInRepo: true,
		run:             doInit,
		//  TODO -- add help, usage, etc.
	}
	myInitCommand.flags.StringVar(&URL, "url", "", "Clearblade platform url for target system")
	myInitCommand.flags.StringVar(&MsgURL, "messaging-url", "", "Clearblade messaging url for target system")
	myInitCommand.flags.StringVar(&SystemKey, "system-key", "", "System key for target system")
	myInitCommand.flags.StringVar(&Email, "email", "", "Developer email for login")
	myInitCommand.flags.StringVar(&Password, "password", "", "Developer password")
	AddCommand("init", myInitCommand)
}

func doInit(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) != 0 {
		fmt.Printf("init command takes no arguments; only options: '%+v'\n", args)
		os.Exit(1)
	}
	return reallyInit(client, SystemKey)
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
		"platform_url":       cb.CB_ADDR,
		"messaging_url":      cb.CB_MSG_ADDR,
		"developer_email":    Email,
		"asset_refresh_dates": []interface{}{},
		"token":             cli.DevToken,
	}
	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been initialized into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}
