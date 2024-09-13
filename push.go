package cblib

import (
	"fmt"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/fs"
	"github.com/clearblade/cblib/models/systemUpload"
	"github.com/clearblade/cblib/models/systemUpload/dryRun"
	"github.com/clearblade/cblib/types"
)

func init() {

	usage :=
		`
	Push a ClearBlade asset from local filesystem to ClearBlade Platform
	`

	example :=
		`
	cb-cli push -all							# Push all assets up to Platform
	cb-cli push -all -auto-approve				# Push all assets up to Platform and automatically confirm any prompts for creating new assets
	cb-cli push -all-services -all-portals		# Push all services and all portals up to Platform
	cb-cli push -service=Service1				# Push a code service up to Platform
	cb-cli push -collection=Collection1			# Push a code service up to Platform
	`

	pushCommand := &SubCommand{
		name:      "push",
		usage:     usage,
		needsAuth: true,
		run:       doPush,
		example:   example,
	}

	pushCommand.flags.BoolVar(&UserSchema, "userschema", false, "push user table schema")
	pushCommand.flags.BoolVar(&EdgeSchema, "edgeschema", false, "push edges table schema")
	pushCommand.flags.BoolVar(&DeviceSchema, "deviceschema", false, "push devices table schema")
	pushCommand.flags.BoolVar(&AllServices, "all-services", false, "push all of the local services")
	pushCommand.flags.BoolVar(&AllLibraries, "all-libraries", false, "push all of the local libraries")
	pushCommand.flags.BoolVar(&AllDevices, "all-devices", false, "push all of the local devices")
	pushCommand.flags.BoolVar(&AllEdges, "all-edges", false, "push all of the local edges")
	pushCommand.flags.BoolVar(&AllPortals, "all-portals", false, "push all of the local portals")
	pushCommand.flags.BoolVar(&AllPlugins, "all-plugins", false, "push all of the local plugins")
	pushCommand.flags.BoolVar(&AllAdaptors, "all-adapters", false, "push all of the local adapters")
	pushCommand.flags.BoolVar(&AllCollections, "all-collections", false, "push all of the local collections")
	pushCommand.flags.BoolVar(&AllRoles, "all-roles", false, "push all of the local roles")
	pushCommand.flags.BoolVar(&AllUsers, "all-users", false, "push all of the local users")
	pushCommand.flags.BoolVar(&AllAssets, "all", false, "push all of the local assets")
	pushCommand.flags.BoolVar(&AllTriggers, "all-triggers", false, "push all of the local triggers")
	pushCommand.flags.BoolVar(&AllTimers, "all-timers", false, "push all of the local timers")
	pushCommand.flags.BoolVar(&AllDeployments, "all-deployments", false, "push all of the local deployments")
	pushCommand.flags.BoolVar(&AllServiceCaches, "all-shared-caches", false, "push all of the local shared caches")
	pushCommand.flags.BoolVar(&AllWebhooks, "all-webhooks", false, "push all of the local webhooks")
	pushCommand.flags.BoolVar(&AllExternalDatabases, "all-external-databases", false, "push all of the local external databases")
	pushCommand.flags.BoolVar(&AllBucketSets, "all-bucket-sets", false, "push all of the local bucket sets")
	pushCommand.flags.BoolVar(&AllBucketSetFiles, "all-bucket-set-files", false, "push all files from all local bucket sets")
	pushCommand.flags.BoolVar(&AutoApprove, "auto-approve", false, "automatically answer yes to all prompts. Useful for creating new entities when they aren't found in the platform")
	pushCommand.flags.BoolVar(&AllSecrets, "all-user-secrets", false, "push all user secrets")
	pushCommand.flags.BoolVar(&MessageHistoryStorage, "message-history-storage", false, "push message history storage")
	pushCommand.flags.BoolVar(&MessageTypeTriggers, "message-type-triggers", false, "push message type triggers")

	pushCommand.flags.StringVar(&CollectionSchema, "collectionschema", "", "Name of collection schema to push")
	pushCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to push")
	pushCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to push")
	pushCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to push")
	pushCommand.flags.StringVar(&CollectionId, "collectionID", "", "Unique id of collection to update. -collection flag is preferred")
	pushCommand.flags.StringVar(&User, "user", "", "Name of user to push")
	pushCommand.flags.StringVar(&UserId, "userID", "", "Unique id of user to update. -user flag is preferred")
	pushCommand.flags.StringVar(&RoleName, "role", "", "Name of role to push")
	pushCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to push")
	pushCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to push")
	pushCommand.flags.StringVar(&DeviceName, "device", "", "Name of device to push")
	pushCommand.flags.StringVar(&EdgeName, "edge", "", "Name of edge to push")
	pushCommand.flags.StringVar(&PortalName, "portal", "", "Name of portal to push")
	pushCommand.flags.StringVar(&PluginName, "plugin", "", "Name of plugin to push")
	pushCommand.flags.StringVar(&AdaptorName, "adapter", "", "Name of adapter to push")
	pushCommand.flags.StringVar(&DeploymentName, "deployment", "", "Name of deployment to push")
	pushCommand.flags.StringVar(&ServiceCacheName, "shared-cache", "", "Name of shared cache to push")
	pushCommand.flags.StringVar(&WebhookName, "webhook", "", "Name of webhook to push")
	pushCommand.flags.StringVar(&ExternalDatabaseName, "external-database", "", "Name of external database to push")
	pushCommand.flags.StringVar(&BucketSetName, "bucket-set", "", "Name of bucket set to push")
	pushCommand.flags.StringVar(&BucketSetFiles, "bucket-set-files", "", "Name of bucket set to push files to. Can be used in conjunction with -box and -file")
	pushCommand.flags.StringVar(&BucketSetBoxName, "box", "", "Name of box to search in bucket set")
	pushCommand.flags.StringVar(&BucketSetFileName, "file", "", "Name of file to push from bucket set box")
	pushCommand.flags.StringVar(&SecretName, "user-secret", "", "Name of user secret to push")

	setBackoffFlags(pushCommand.flags)

	AddCommand("push", pushCommand)
}

func checkPushArgsAndFlags(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("There are no arguments to the push command, only command line options\n")
	}
	if AllServices && ServiceName != "" {
		return fmt.Errorf("Cannot specify both -all-services and -service=<service_name>\n")
	}
	if AllLibraries && LibraryName != "" {
		return fmt.Errorf("Cannot specify both -all-libraries and -library=<library_name>\n")
	}
	return nil
}

func doPush(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	parseBackoffFlags()
	if err := checkPushArgsAndFlags(args); err != nil {
		return err
	}
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}
	SetRootDir(".")

	// This is a hack to check if token has expired and auth again
	// since we dont have an endpoint to determine this
	client, err = checkIfTokenHasExpired(client, systemInfo.Key)
	if err != nil {
		return fmt.Errorf("Re-auth failed...")
	}

	version, err := systemUpload.GetSystemUploadVersion(systemInfo, client)
	if err != nil {
		return err
	}

	// Below version 5 we only support code services, so we need to do the legacy push
	if version < 5 {
		return doLegacyPush(cmd, client, systemInfo)
	}

	return pushSystemZip(systemInfo, client, defaultZipOptions())
}

type prompter struct{}

func (p prompter) PromptForSecret(prompt string) string {
	return getOneItem(prompt, true)
}

func pushSystemZip(systemInfo *types.System_meta, client *cb.DevClient, options *fs.ZipOptions) error {
	fmt.Printf("Preparing to push system %s\n", systemInfo.Name)
	buffer, err := fs.GetSystemZipBytes(rootDir, prompter{}, options)
	if err != nil {
		return err
	}

	fmt.Println("Doing dry run")
	result, err := client.UploadToSystemDryRun(systemInfo.Key, buffer)
	if err != nil {
		return err
	}

	dryRun, err := dryRun.New(result)
	if err != nil {
		return err
	}

	if dryRun.HasErrors() {
		return fmt.Errorf(dryRun.String())
	}

	if !dryRun.HasChanges() {
		fmt.Println("Nothing to push")
		return nil
	}

	fmt.Print(dryRun.String())
	changesAccepted, err := confirmPrompt(fmt.Sprintln("Would you like to accept these changes?"))
	if err != nil {
		return err
	}

	if !changesAccepted {
		fmt.Println("Changes will not be pushed")
		return nil
	}

	fmt.Println("Pushing changes")
	r, err := client.UploadToSystem(systemInfo.Key, buffer)
	if err != nil {
		return err
	}

	updateIdMap(r)
	return r.Error()
}

func updateIdMap(result *cb.SystemUploadChanges) {
	updateCollectionMap(result)
	updateUserMap(result)
	updateRoleMap(result)
}

func updateCollectionMap(result *cb.SystemUploadChanges) {
	for name, id := range result.CollectionNameToId {
		if err := updateCollectionNameToId(CollectionInfo{
			ID:   id,
			Name: name,
		}); err != nil {
			fmt.Printf("Could not update collection map entry (%s to %s): %s\n", name, id, err)
		}
	}
}

func updateUserMap(result *cb.SystemUploadChanges) {
	for email, id := range result.UserEmailToId {
		if err := updateUserEmailToId(UserInfo{
			UserID: id,
			Email:  email,
		}); err != nil {
			fmt.Printf("Could not update user map entry (%s to %s): %s\n", email, id, err)
		}
	}
}

func updateRoleMap(result *cb.SystemUploadChanges) {
	for name, id := range result.RoleNameToId {
		if err := updateRoleNameToId(RoleInfo{
			ID:   id,
			Name: name,
		}); err != nil {
			fmt.Printf("Could not update role map entry (%s to %s): %s\n", name, id, err)
		}
	}
}
