package cblib

import (
	//"flag"
	"fmt"
	"os"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/remote"
)

func init() {
	usage :=
		`
		Initializes your filesystem with your ClearBlade Platform System or targets your local system to a different remote system within a ClearBlade Platform
	`

	example :=
		`
	cb-cli init
	cb-cli init -url=https://platform.clearblade.com -messaging-url=platform.clearblade.com -system-key=8abcd6aa0baadcd8bbe3fabca29301 -email=dev@dev.com -password=pw
	`
	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	myInitCommand := &SubCommand{
		name:      "init",
		usage:     usage,
		needsAuth: false,
		run:       doInit,
		example:   example,
	}
	myInitCommand.flags.StringVar(&URL, "url", "", "Clearblade platform url for target system")
	myInitCommand.flags.StringVar(&MsgURL, "messaging-url", "", "Clearblade messaging url for target system")
	myInitCommand.flags.StringVar(&SystemKey, "system-key", "", "System Key for target system, ex 9b9eea9c0bda8896a3dab5aeec9601")
	myInitCommand.flags.StringVar(&Email, "email", "", "Developer email for login")
	myInitCommand.flags.StringVar(&Password, "password", "", "Developer password")
	myInitCommand.flags.StringVar(&DevToken, "dev-token", "", "Developer token to use instead of email/password")
	myInitCommand.flags.BoolVar(&SkipUpdateMapNameToIdFiles, "skip-update-map-name-to-id", false, "Set this to true to skip pulling the IDs for roles, collections, and users. This is useful if the system has lots of these types of assets and the goal is to retrieve the schema for the tables after initialization.")
	AddCommand("init", myInitCommand)
}

func doInit(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) != 0 {
		fmt.Printf("init command takes no arguments; only options: '%+v'\n", args)
		os.Exit(1)
	}
	defaults := setupInitDefaults()
	MetaInfo = nil
	client, err := Authorize(defaults)
	if err != nil {
		return err
	}
	return reallyInit(cmd, client, SystemKey)
}

func reallyInit(cmd *SubCommand, cli *cb.DevClient, sysKey string) error {
	SetRootDir(".")

	if err := setupDirectoryStructure(); err != nil {
		return err
	}

	systemMeta, err := pullSystemMeta(sysKey, cli)
	if err != nil {
		return err
	}

	err = initRemote(cmd, systemMeta, cli)
	if err != nil {
		return err
	}

	if !SkipUpdateMapNameToIdFiles {
		logInfo("Updating map name to ID files...")
		updateMapNameToIDFiles(systemMeta, cli)
	}

	fmt.Printf("System '%s' has been initialized in the current directory.\n", systemMeta.Name)
	return nil
}

func initRemote(cmd *SubCommand, systemMeta *System_meta, cli *cb.DevClient) error {
	cmd.remotes = remote.NewRemotes()

	initRemote := &remote.Remote{
		Name:         "init",
		PlatformURL:  cli.HttpAddr,
		MessagingURL: cli.MqttAddr,
		SystemKey:    systemMeta.Key,
		SystemSecret: systemMeta.Secret,
		Token:        cli.DevToken,
	}

	cmd.remotes.Put(initRemote)

	systemJSON := systemMetaToMap(systemMeta)

	cbmeta := map[string]interface{}{
		"platform_url":    cli.HttpAddr,
		"messaging_url":   cli.MqttAddr,
		"developer_email": cli.Email,
		"token":           cli.DevToken,
	}

	return useRemoteByMerging(systemJSON, cbmeta, initRemote)
}

type DefaultInfo struct {
	url       string
	email     string
	systemKey string
	msgUrl    string
}

func setupInitDefaults() *DefaultInfo {
	meta, err := getSysMeta()
	if err != nil || MetaInfo == nil {
		return nil
	}

	platform_url, ok := MetaInfo["platformURL"].(string)
	if !ok {
		platform_url = MetaInfo["platform_url"].(string)
	}
	email, ok := MetaInfo["developerEmail"].(string)
	if !ok {
		email = MetaInfo["developer_email"].(string)
	}
	messaging_url, ok := MetaInfo["messagingURL"].(string)
	if !ok {
		messaging_url = MetaInfo["messaging_url"].(string)
	}

	return &DefaultInfo{
		url:       platform_url,
		email:     email,
		systemKey: meta.Key,
		msgUrl:    messaging_url,
	}
}
