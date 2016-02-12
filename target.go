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
	myTargetCommand := &SubCommand{
		name:         "export",
		usage:        "Ain't no thing",
		needsAuth:    false,
		mustBeInRepo: true,
		run:          doTarget,
		//  TODO -- add help, usage, etc.
	}
	myTargetCommand.flags.StringVar(&URL, "url", "", "Clearblade platform url for target system")
	myTargetCommand.flags.StringVar(&SystemKey, "system-key", "", "System key for target system")
	myTargetCommand.flags.StringVar(&Email, "email", "", "Developer email for login")
	myTargetCommand.flags.StringVar(&Password, "password", "", "Developer password")
	AddCommand("target", myTargetCommand)
}

func doTarget(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) != 0 {
		fmt.Printf("init command takes no arguments; only options: '%+v'\n", args)
		os.Exit(1)
	}
	defaults := setupTargetDefaults()
	oldSysMeta, err := getSysMeta()
	if err != nil {
		return err
	}
	if err = os.Chdir(".."); err != nil {
		return fmt.Errorf("Could not move up to parent directory: %s", err.Error())
	}
	MetaInfo = nil
	client, err = Authorize(defaults)
	if err != nil {
		return err
	}
	return reallyTarget(client, SystemKey, oldSysMeta)
}

func fixSystemName(sysName string) string {
	return strings.Replace(sysName, " ", "_", -1)
}

func reallyTarget(cli *cb.DevClient, sysKey string, oldSysMeta *System_meta) error {
	sysMeta, err := pullSystemMeta(sysKey, cli)
	if err != nil {
		return err
	}

	fixo, fixn := fixSystemName(oldSysMeta.Name), fixSystemName(sysMeta.Name)
	if fixo != fixn {
		fmt.Printf("Renaming %s to %s\n", fixo, fixn)
		os.Rename(fixo, fixn)
	}

	setRootDir(fixn)
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
		"token":             cli.DevToken,
	}
	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been initialized into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}

//
//  This stuff is a hack -- it's used when initing inside a repo to give prompts
//  with defaults. see fillInTheBlanks(..) in newAuth.go
//

type DefaultInfo struct {
	url       string
	email     string
	systemKey string
}

func setupTargetDefaults() *DefaultInfo {
	meta, err := getSysMeta()
	if err != nil || MetaInfo == nil {
		return nil
	}
	return &DefaultInfo{
		url:       MetaInfo["platformURL"].(string),
		email:     MetaInfo["developerEmail"].(string),
		systemKey: meta.Key,
	}
}
