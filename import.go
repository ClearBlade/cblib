package cblib

import (
	"fmt"
	"os"
	"path"

	cb "github.com/clearblade/Go-SDK"

	"github.com/clearblade/cblib/fs"
	"github.com/clearblade/cblib/models/systemUpload"
	"github.com/clearblade/cblib/types"
)

var (
	importPiecemeal bool
	importRows      bool
	importUsers     bool
)

func init() {

	usage :=
		`
	Import a system from your local filesystem to the ClearBlade Platform
	`

	example :=
		`
	cb-cli import 									# prompts for credentials
	cb-cli import -importrows=false -importusers=false			# prompts for credentials, excludes all collection-rows and users
	`
	myImportCommand := &SubCommand{
		name:      "import",
		usage:     usage,
		needsAuth: false,
		run:       doImport,
		example:   example,
	}
	myImportCommand.flags.BoolVar(&importPiecemeal, "piecemeal", false, "perform push through many individual http requests instead of uploading a single zip")
	myImportCommand.flags.BoolVar(&importRows, "importrows", true, "imports all data into all collections")
	myImportCommand.flags.BoolVar(&importUsers, "importusers", true, "imports all users into the system")
	myImportCommand.flags.StringVar(&URL, "url", "https://platform.clearblade.com", "Clearblade Platform URL where system is hosted, ex https://platform.clearblade.com")
	myImportCommand.flags.StringVar(&Email, "email", "", "Developer email for login to import destination")
	myImportCommand.flags.StringVar(&Password, "password", "", "Developer password at import destination")
	myImportCommand.flags.StringVar(&DevToken, "dev-token", "", "Developer token to use instead of email/password")
	myImportCommand.flags.IntVar(&DataPageSize, "data-page-size", DataPageSizeDefault, "Number of rows in a collection to push/import at a time")
	setBackoffFlags(myImportCommand.flags)
	AddCommand("import", myImportCommand)
	AddCommand("imp", myImportCommand)
	AddCommand("im", myImportCommand)
}

func doImport(cmd *SubCommand, _ *cb.DevClient, _ ...string) error {
	parseBackoffFlags()
	systemPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// prompt and skip values we don't need
	skips := PromptSkipMsgURL | PromptSkipSystemKey
	if DevToken != "" {
		skips |= PromptSkipEmail
		skips |= PromptSkipPassword
	}
	promptAndFillMissingAuth(nil, skips)

	// authorizes using global flags (import ignores cb meta)
	cli, err := authorizeUsingGlobalCLIFlags()
	if err != nil {
		return err
	}

	// creates import config and proceeds to import system
	config := MakeImportConfigFromGlobals()
	AutoApprove = true
	_, err = ImportSystemUsingConfig(config, systemPath, cli)
	if err != nil {
		return err
	}

	return nil
}

// --------------------------------
// Import config and other types
// --------------------------------
// We use an import config that is passed around as a parameter during the
// import process.

// ImportConfig contains configuration values for the import process.
// NOTE: Other configuration parameters can be added here. The idea is to pass
// them to the import process using an instance of this struct rather than using
// global variables. TRY TO KEEP ANY INSTANCE OF THIS STRUCTURE READ-ONLY.
type ImportConfig struct {
	SystemName        string // the name of the imported system
	SystemDescription string // the description of the imported system

	IntoExistingSystem   bool   // true if it should be imported on a system that already exists
	ExistingSystemKey    string // the system key of the existing system
	ExistingSystemSecret string // the system secret of the existing system

	ImportUsers     bool // true if users should be imported
	ImportRows      bool // true if collection rows should be imported
	ImportPiecemeal bool // true if system upload endpoint should not be used
}

// DefaultImportConfig contains the default configuration values for the import
// process. Note that this instance SHOULD NOT be updated and used as a global
// configuration object. If you wish to configure the import processs using the
// global variables, check the NewImportConfigFromGlobals function.
//
// To create your own configuration config just assign this one to your own
// and modify it:
//
// ```
// customImportConfig := DefaultImportConfig
// customImportConfig.DefaultUserPassword = "my-new-password"
// ````
var DefaultImportConfig = ImportConfig{
	SystemName:        "",
	SystemDescription: "",

	IntoExistingSystem:   false,
	ExistingSystemKey:    "",
	ExistingSystemSecret: "",

	ImportUsers:     false,
	ImportRows:      false,
	ImportPiecemeal: false,
}

// MakeImportConfigFromGlobals creates a new ImportConfig instance from the
// GLOBAL variables in cblib. Use with caution. Note that this function starts
// with Make* and not with New* because it returns a normal instance, and not
// a pointer to an instance.
func MakeImportConfigFromGlobals() ImportConfig {
	config := DefaultImportConfig
	config.ImportUsers = importUsers
	config.ImportRows = importRows
	config.ImportPiecemeal = importPiecemeal

	return config
}

// ImportResult holds relevant values resulting from a system import process.
type ImportResult struct {
	rawSystemInfo map[string]interface{}
	SystemName    string
	SystemKey     string
	SystemSecret  string
}

// --------------------------------
// Import process (creation, etc)
// --------------------------------
// Functions that focus on the creation of the system and other assets.

func createSystem(config ImportConfig, system *types.System_meta, client *cb.DevClient) (*types.System_meta, error) {
	name := system.Name
	desc := system.Description
	auth := true
	sysKey, sysErr := client.NewSystem(name, desc, auth)
	if sysErr != nil {
		return nil, sysErr
	}
	realSystem, sysErr := client.GetSystem(sysKey)
	if sysErr != nil {
		return nil, sysErr
	}
	system.Key = realSystem.Key
	system.Secret = realSystem.Secret
	return system, nil
}

func importAllAssets(config ImportConfig, systemInfo *types.System_meta, users []map[string]interface{}, cli *cb.DevClient) error {
	version, err := systemUpload.GetSystemUploadVersion(systemInfo, cli)
	if err != nil {
		return err
	}

	// Below version 5 we only support code services, so we need to do the legacy push
	if version < 5 || config.ImportPiecemeal {
		return importAllAssetsLegacy(config, systemInfo, users, cli)
	}

	opts := fs.NewZipOptions(&mapper{})
	opts.AllAdaptors = true
	opts.AllBucketSets = true
	opts.AllBucketSetFiles = true
	opts.AllCollections = config.ImportRows
	opts.AllCollectionSchemas = true
	opts.AllServiceCaches = true
	opts.AllDeployments = true
	opts.AllDevices = true
	opts.AllEdges = true
	opts.AllExternalDatabases = true
	opts.AllLibraries = true
	opts.AllPlugins = true
	opts.AllPortals = true
	opts.AllRoles = true
	opts.AllSecrets = true
	opts.AllServices = true
	opts.AllTimers = true
	opts.AllTriggers = true
	opts.AllUsers = config.ImportUsers
	opts.AllWebhooks = true
	opts.PushMessageHistoryStorage = true
	opts.PushMessageTypeTriggers = true
	opts.PushUserSchema = true

	if err := pushSystemZip(systemInfo, cli, opts); err != nil {
		return err
	}

	fmt.Printf(" Done\n")
	logInfo(fmt.Sprintf("Success! New system key is: %s", systemInfo.Key))
	logInfo(fmt.Sprintf("New system secret is: %s", systemInfo.Secret))
	return nil
}

// --------------------------------
// Import entrypoint and exposed functions
// --------------------------------

// importSystem will import the system rooted at the given path using the given
// config. Please that we assume that the given clearblade client is already
// authorized an ready to use.
func importSystem(config ImportConfig, systemPath string, cli *cb.DevClient) (*types.System_meta, error) {

	// points the root directory to the system folder
	// WARNING: side-effect (changes globals)
	SetRootDir(systemPath)

	// sets up director strcuture
	// WARNING: side-effect (might change system)
	err := setupDirectoryStructure()
	if err != nil {
		return nil, err
	}

	// gets users from the system directory
	// WARNING: side-effect (reads filesystem)
	users, err := getUsers()
	if err != nil {
		return nil, err
	}

	// gets system info from the system directory
	// WARNING: side-effect (reads filesystem)
	systemInfoPath := path.Join(systemPath, "system.json")
	systemInfoMap, err := getDict(systemInfoPath)
	if err != nil {
		return nil, err
	}
	systemInfo := systemMetaFromMap(systemInfoMap)

	// creates system if we are not importing into an existing one
	if !config.IntoExistingSystem {

		if len(config.SystemName) > 0 {
			systemInfo.Name = config.SystemName
		}

		if len(config.SystemDescription) > 0 {
			systemInfo.Description = config.SystemDescription
		}

		// NOTE: createSystem will modify systemInfo map
		_, err := createSystem(config, systemInfo, cli)
		if err != nil {
			return nil, fmt.Errorf("could not create system named '%s': %s", config.SystemName, err)
		}

	} else {
		systemInfo.Key = config.ExistingSystemKey
		systemInfo.Secret = config.ExistingSystemSecret
	}

	// import assets into created/existing system
	err = importAllAssets(config, systemInfo, users, cli)
	if err != nil {
		return nil, err
	}

	return systemInfo, nil
}

// ImportSystemUsingConfig imports the system rooted at the given path, using the
// given config for different values. The given client should already be
// authenticated and ready to go.
func ImportSystemUsingConfig(config ImportConfig, systemPath string, cli *cb.DevClient) (*types.System_meta, error) {
	systemInfo, err := importSystem(config, systemPath, cli)
	if err != nil {
		return nil, err
	}

	return systemInfo, nil
}
