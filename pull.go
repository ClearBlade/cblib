package cblib

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/models"
	"github.com/clearblade/cblib/types"
)

var (
	PULL_ALL_USERS           = "%PULL_ALL_USERS%"
	FORBIDDEN_CHARS_IN_NAMES = "/"
)

func init() {
	usage :=
		`
	Pull a ClearBlade asset from the Platform to your local filesystem. Use -sort-collections for easier version controlling of datasets.

	Note: Collection rows are pulled by default.
	`

	example :=
		`
	cb-cli pull -all												# Pull all assets from Platform to local filesystem
	cb-cli pull -all-services -all-portals							# Pull all services and all portals from Platform to local filesystem
	cb-cli pull -service=Service1 									# Pulls Service1 from Platform to local filesystem
	cb-cli pull -collection=Collection1								# Pulls Collection1 from Platform to local filesystem, with all rows, unsorted
	cb-cli pull -collection=Collection1 -sort-collections=true		# Pulls Collection1 from Platform to local filesystem, with all rows, sorted
	`
	pullCommand := &SubCommand{
		name:      "pull",
		usage:     usage,
		needsAuth: true,
		run:       doPull,
		example:   example,
	}

	pullCommand.flags.BoolVar(&AllServices, "all-services", false, "pull all services from system")
	pullCommand.flags.BoolVar(&AllLibraries, "all-libraries", false, "pull all libraries from system")
	pullCommand.flags.BoolVar(&AllEdges, "all-edges", false, "pull all edges from system")
	pullCommand.flags.BoolVar(&AllDevices, "all-devices", false, "pull all devices from system")
	pullCommand.flags.BoolVar(&AllPortals, "all-portals", false, "pull all portals from system")
	pullCommand.flags.BoolVar(&AllPlugins, "all-plugins", false, "pull all plugins from system")
	pullCommand.flags.BoolVar(&AllAdaptors, "all-adapters", false, "pull all adapters from system")
	pullCommand.flags.BoolVar(&AllDeployments, "all-deployments", false, "pull all deployments from system")
	pullCommand.flags.BoolVar(&AllCollections, "all-collections", false, "pull all collections from system")
	pullCommand.flags.BoolVar(&AllRoles, "all-roles", false, "pull all roles from system")
	pullCommand.flags.BoolVar(&AllUsers, "all-users", false, "pull all users from system")
	pullCommand.flags.BoolVar(&UserSchema, "userschema", false, "pull user table schema")
	pullCommand.flags.BoolVar(&DeviceSchema, "deviceschema", false, "pull device table schema")
	pullCommand.flags.BoolVar(&EdgeSchema, "edgeschema", false, "pull edges table schema")
	pullCommand.flags.BoolVar(&AllAssets, "all", false, "pull all assets from system")
	pullCommand.flags.BoolVar(&AllTriggers, "all-triggers", false, "pull all triggers from system")
	pullCommand.flags.BoolVar(&AllTimers, "all-timers", false, "pull all timers from system")
	pullCommand.flags.BoolVar(&AllServiceCaches, "all-shared-caches", false, "pull all shared caches from system")
	pullCommand.flags.BoolVar(&AllWebhooks, "all-webhooks", false, "pull all webhooks from system")
	pullCommand.flags.BoolVar(&AllExternalDatabases, "all-external-databases", false, "pull all external databases from system")
	pullCommand.flags.BoolVar(&AllFileStores, "all-file-stores", false, "pull all file stores from system")
	pullCommand.flags.BoolVar(&AllFileStoreFiles, "all-file-store-files", false, "pull all files from all file stores from system")
	pullCommand.flags.BoolVar(&AllBucketSets, "all-bucket-sets", false, "pull all bucket sets from system")
	pullCommand.flags.BoolVar(&AllBucketSetFiles, "all-bucket-set-files", false, "pull all files from all bucket sets from system")
	pullCommand.flags.BoolVar(&AllSecrets, "all-user-secrets", false, "pull all user secrets from system")
	pullCommand.flags.BoolVar(&MessageHistoryStorage, "message-history-storage", false, "pull message history storage from system")
	pullCommand.flags.BoolVar(&MessageTypeTriggers, "message-type-triggers", false, "pull message type triggers from system")

	pullCommand.flags.StringVar(&CollectionSchema, "collectionschema", "", "Name of collection schema to pull")
	pullCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to pull")
	pullCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to pull")
	pullCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to pull")
	pullCommand.flags.BoolVar(&SortCollections, "sort-collections", SortCollectionsDefault, "Sort collections by item id, for version control ease")
	pullCommand.flags.IntVar(&DataPageSize, "page-size", DataPageSizeDefault, "Number of rows in a collection to request at a time")
	pullCommand.flags.StringVar(&User, "user", "", "Name of user to pull")
	pullCommand.flags.StringVar(&RoleName, "role", "", "Name of role to pull")
	pullCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to pull")
	pullCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to pull")
	pullCommand.flags.StringVar(&EdgeName, "edge", "", "Name of edge to pull")
	pullCommand.flags.StringVar(&DeviceName, "device", "", "Name of device to pull")
	pullCommand.flags.StringVar(&PortalName, "portal", "", "Name of portal to pull")
	pullCommand.flags.StringVar(&PluginName, "plugin", "", "Name of plugin to pull")
	pullCommand.flags.StringVar(&AdaptorName, "adapter", "", "Name of adapter to pull")
	pullCommand.flags.StringVar(&DeploymentName, "deployment", "", "Name of deployment to pull")
	pullCommand.flags.StringVar(&ServiceCacheName, "shared-cache", "", "Name of shared cache to pull")
	pullCommand.flags.StringVar(&WebhookName, "webhook", "", "Name of webhook to pull")
	pullCommand.flags.StringVar(&ExternalDatabaseName, "external-database", "", "Name of external database to pull")
	pullCommand.flags.StringVar(&BucketSetName, "bucket-set", "", "Name of bucket set to pull")
	pullCommand.flags.StringVar(&BucketSetFiles, "bucket-set-files", "", "Name of bucket set to pull files from. Can be used in conjunction with -box and -file")
	pullCommand.flags.StringVar(&BucketSetBoxName, "box", "", "Name of box to search in bucket set")
	pullCommand.flags.StringVar(&BucketSetFileName, "file", "", "Name of file to pull from bucket set box")
	pullCommand.flags.StringVar(&FileStoreName, "file-store", "", "Name of file store to pull")
	pullCommand.flags.StringVar(&FileStoreFiles, "file-store-files", "", "Name of file store to pull files from. Can be used in conjunction with -file-store-file")
	pullCommand.flags.StringVar(&FileStoreFileName, "file-store-file", "", "Name of file to pull from file store specified with -file-store-files")
	pullCommand.flags.StringVar(&SecretName, "user-secret", "", "Name of user secret to pull")

	setBackoffFlags(pullCommand.flags)

	AddCommand("pull", pullCommand)
}

func doPull(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	parseBackoffFlags()
	SetRootDir(".")
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}

	if err := setupDirectoryStructure(); err != nil {
		return err
	}

	// This is a hack to check if token has expired and auth again
	// since we dont have an endpoint to determine this
	client, err = checkIfTokenHasExpired(client, systemInfo.Key)
	if err != nil {
		return fmt.Errorf("Re-auth failed: %s\n", err)
	}

	assetsToPull := createAffectedAssets()
	assetsToPull.ExportItemId = true
	assetsToPull.ExportRows = true
	assetsToPull.ExportUsers = true
	didSomething, err := pullAssets(systemInfo, client, assetsToPull)

	if !didSomething {
		fmt.Printf("Nothing to pull -- you must specify something to pull (ie, -service=<svc_name>)\n")
	}
	return nil
}

var userColumnsToSkip = []string{"email", "creation_date", "cb_service_account", "cb_ttl_override", "cb_token"}

func isInList(list []string, item string) bool {
	for i := 0; i < len(list); i++ {
		if list[i] == item {
			return true
		}
	}
	return false
}

func pullUserSchemaInfo(systemKey string, cli *cb.DevClient, writeThem bool) (map[string]interface{}, error) {
	resp, err := cli.GetUserColumns(systemKey)
	if err != nil {
		return nil, err
	}
	columns := getUserDefinedColumns(resp)
	schema := map[string]interface{}{
		"columns": columns,
	}
	if writeThem {
		if err := writeUserSchema(schema); err != nil {
			return nil, err
		}
	}
	return schema, nil
}

func pullRole(systemKey string, roleName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetRole(systemKey, roleName)
}

func PullAndWriteRoles(systemKey string, cli *cb.DevClient, writeThem bool) ([]map[string]interface{}, error) {
	r, err := cli.GetAllRoles(systemKey)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, 0)
	for _, rIF := range r {
		thisRole := rIF.(map[string]interface{})
		fmt.Printf(" %s", thisRole["Name"].(string))
		rval = append(rval, thisRole)
		if writeThem {
			if err := writeRole(thisRole["Name"].(string), thisRole); err != nil {
				return nil, err
			}
		}
	}
	return rval, nil
}

func PullAndWriteService(systemKey string, serviceName string, client *cb.DevClient) error {
	if svc, err := pullService(systemKey, serviceName, client); err != nil {
		return err
	} else {
		return writeService(serviceName, svc)
	}
}

func pullService(systemKey string, serviceName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetServiceRaw(systemKey, serviceName)
}

func PullAndWriteLibrary(systemKey string, libraryName string, client *cb.DevClient) error {
	if svc, err := pullLibrary(systemKey, libraryName, client); err != nil {
		return err
	} else {
		return writeLibrary(libraryName, svc)
	}
}

func pullAllUsers(systemKey string, client *cb.DevClient) ([]interface{}, error) {
	return paginateRequests(systemKey, DataPageSize, client.GetUserCountWithQuery, client.GetUsersWithQuery)
}

func PullAndWriteUsers(systemKey string, userName string, client *cb.DevClient, saveThem bool) ([]map[string]interface{}, error) {
	if users, err := pullAllUsers(systemKey, client); err != nil {
		return nil, err
	} else {
		ok := false
		rtn := make([]map[string]interface{}, 0)
		for _, u := range users {
			user := u.(map[string]interface{})
			if user["email"] == userName || userName == PULL_ALL_USERS {
				email := user["email"].(string)
				fmt.Printf(" %s", email)
				ok = true
				userId := user["user_id"].(string)
				roles, err := client.GetUserRoles(systemKey, userId)
				if err != nil {
					return nil, fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
				}
				rtn = append(rtn, user)
				if saveThem {
					err := writeUser(email, user)
					if err != nil {
						return nil, err
					}
					err = writeUserRoles(email, roles)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if !ok {
			if userName == PULL_ALL_USERS {
				return nil, fmt.Errorf("No users found")
			} else {
				return nil, fmt.Errorf("User %+s not found\n", userName)
			}

		} else {
			return rtn, nil
		}

	}
}

func PullAndWriteCollection(systemInfo *types.System_meta, collectionName string, client *cb.DevClient, shouldExportRows, shouldExportItemId bool) error {
	if allColls, err := client.GetAllCollections(systemInfo.Key); err != nil {
		return err
	} else {
		var collID string
		// iterate over allColls and find one with matching name
		for _, c := range allColls {
			coll := c.(map[string]interface{})
			if collectionName == coll["name"] {
				collID = coll["collectionID"].(string)
			}
		}
		if len(collID) < 1 {
			return fmt.Errorf("Collection %s not found.", collectionName)
		}
		if coll, err := client.GetCollectionInfo(collID); err != nil {
			return err
		} else {
			if data, err := PullCollection(systemInfo, client, coll, shouldExportRows, shouldExportItemId); err != nil {
				return err
			} else {
				d := makeCollectionJsonConsistent(data)
				err = writeCollection(d["name"].(string), d)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func pullLibrary(systemKey string, libraryName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetLibrary(systemKey, libraryName)
}

func stripTriggerFields(trig map[string]interface{}) {
	delete(trig, "system_key")
	delete(trig, "system_secret")
	return
}

func writeTriggerWithUserInfo(name string, trig map[string]interface{}) error {
	users, err := getUserEmailToId()
	if err != nil {
		logWarning(fmt.Sprintf("Unable to fetch user email map when writing trigger. This can be ignored if your system doesn't have users or doesn't have any user triggers; Any user triggers in the system will be stored with userId rather than email which will affect their portability between systems. Any user triggers will need to be recreated after importing into a new system. Message: %s", err.Error()))
	} else {
		replaceUserIdWithEmailInTriggerKeyValuePairs(trig, users)
	}
	return writeTrigger(name, trig)
}

func PullAndWriteTrigger(systemKey, trigName string, client *cb.DevClient) error {
	if trigg, err := pullTrigger(systemKey, trigName, client); err != nil {
		return err
	} else {
		stripTriggerFields(trigg)
		err = writeTriggerWithUserInfo(trigName, trigg)
		if err != nil {
			return err
		}
	}
	return nil
}

func PullAndWriteTriggers(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	trigs, err := cli.GetEventHandlers(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull triggers out of system %s: %s", sysMeta.Key, err.Error())
	}
	triggers := []map[string]interface{}{}
	for _, trig := range trigs {
		thisTrig := trig.(map[string]interface{})
		fmt.Printf(" %s", thisTrig["name"].(string))
		stripTriggerFields(thisTrig)
		triggers = append(triggers, thisTrig)
		err = writeTriggerWithUserInfo(thisTrig["name"].(string), thisTrig)
		if err != nil {
			return nil, err
		}
	}
	return triggers, nil
}

func PullAndWriteTimer(systemKey, timerName string, client *cb.DevClient) error {
	if timer, err := pullTimer(systemKey, timerName, client); err != nil {
		return err
	} else {
		err = writeTimer(timerName, timer)
		if err != nil {
			return err
		}
	}
	return nil
}

func PullAndWriteTimers(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theTimers, err := cli.GetTimers(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull timers out of system %s: %s", sysMeta.Key, err.Error())
	}
	timers := []map[string]interface{}{}
	for _, timer := range theTimers {
		thisTimer := timer.(map[string]interface{})
		thisTimerName := thisTimer["name"].(string)
		fmt.Printf("\n%s ", thisTimerName)
		if strings.ContainsAny(thisTimerName, FORBIDDEN_CHARS_IN_NAMES) {
			fmt.Printf(" (skipped: forbidden char in name)")
		} else {
			timers = append(timers, thisTimer)
			err = writeTimer(thisTimerName, thisTimer)
			if err != nil {
				return nil, err
			}
		}
	}
	return timers, nil
}

func PullAndWritePortal(systemKey, name string, client *cb.DevClient) error {
	if portal, err := pullPortal(systemKey, name, client); err != nil {
		return err
	} else {
		return writePortal(name, portal)
	}
}

func PullAndWritePlugin(systemKey, name string, client *cb.DevClient) error {
	if plugin, err := pullPlugin(systemKey, name, client); err != nil {
		return err
	} else {
		if err = writePlugin(name, plugin); err != nil {
			return err
		}
	}
	return nil
}

func PullAndWriteAdaptor(systemKey, name string, client *cb.DevClient) error {
	if adaptor, err := pullAdaptor(systemKey, name, client); err != nil {
		return err
	} else {
		if err = writeAdaptor(adaptor); err != nil {
			return err
		}
	}
	return nil
}

func pullTrigger(systemKey string, triggerName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetEventHandler(systemKey, triggerName)
}

func pullTimer(systemKey string, timerName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetTimer(systemKey, timerName)
}

func pullDevice(systemKey string, deviceName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetDevice(systemKey, deviceName)
}

func pullEdge(systemKey string, edgeName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetEdge(systemKey, edgeName)
}

func transformPortal(portal map[string]interface{}) error {
	portal = removeBlacklistedPortalKeys(portal)
	if parsed, err := parseIfNeeded(portal["config"]); err != nil {
		return err
	} else {
		portal["config"] = parsed
	}
	return nil
}

func pullPortal(systemKey string, portalName string, client *cb.DevClient) (map[string]interface{}, error) {
	portal, err := client.GetPortal(systemKey, portalName)
	if err != nil {
		return nil, err
	}
	if err := transformPortal(portal); err != nil {
		return nil, err
	}
	return portal, nil
}

func pullPlugin(systemKey string, pluginName string, client *cb.DevClient) (map[string]interface{}, error) {
	return client.GetPlugin(systemKey, pluginName)
}

func pullAdaptor(systemKey, adaptorName string, client *cb.DevClient) (*models.Adaptor, error) {
	fmt.Printf("\n %s", adaptorName)
	currentAdaptor := models.InitializeAdaptor(adaptorName, systemKey, client)

	if err := currentAdaptor.FetchAllInfo(); err != nil {
		return nil, err
	}

	return currentAdaptor, nil
}

func updateMapNameToIDFiles(systemInfo *types.System_meta, client *cb.DevClient) {
	logInfo("Updating roles...")
	if roles, err := getRoles(); err != nil {
		logError(fmt.Sprintf("Failed to get roles %s", err.Error()))
	} else {
		for i := 0; i < len(roles); i++ {
			roleName := roles[i]["Name"].(string)
			fmt.Printf(" %s", roleName)
			role, err := pullRole(systemInfo.Key, roleName, client)
			if err != nil {
				logError(fmt.Sprintf("Failed to pull role '%s'. %s", roleName, err.Error()))
			} else {
				updateRoleNameToId(RoleInfo{
					ID:   role["ID"].(string),
					Name: role["Name"].(string),
				})
			}
		}
	}

	logInfo("\nUpdating collections...")
	if collections, err := getCollections(); err != nil {
		logError(fmt.Sprintf("Failed to get collections %s", err.Error()))
	} else {
		if data, err := client.GetAllCollections(systemInfo.Key); err != nil {
			logError(fmt.Sprintf("Failed to get all collections metadata. %s", err.Error()))
		} else {
			for i := 0; i < len(collections); i++ {
				collectionName := collections[i]["name"].(string)
				fmt.Printf(" %s", collectionName)
				if found, collectionID := findCollectionID(data, collectionName); found {
					updateCollectionNameToId(CollectionInfo{
						ID:   collectionID,
						Name: collectionName,
					})
				} else {
					logWarning(fmt.Sprintf("Could not find collection '%s'", collectionName))
				}
			}
		}
	}

	logInfo("Updating users...")
	if users, err := getUsers(); err != nil {
		logError(fmt.Sprintf("Failed to get users %s", err.Error()))
	} else if len(users) > 0 {
		query := cb.NewQuery()
		query.Columns = []string{"user_id"}
		userEmails := []string{}
		for i := 0; i < len(users); i++ {
			userEmails = append(userEmails, users[i]["email"].(string))
		}
		query.Filters[0] = append(query.Filters[0], cb.Filter{
			Field:    "email",
			Value:    userEmails,
			Operator: "IN",
		})
		client.GetUsersWithQuery(systemInfo.Key, query)
		for i := 0; i < len(users); i++ {
			userEmail := users[i]["email"].(string)
			data, err := PullAndWriteUsers(systemInfo.Key, userEmail, client, false)
			if err != nil {
				logError(fmt.Sprintf("Failed to pull user '%s'. %s", userEmail, err.Error()))
			} else {
				updateUserEmailToId(UserInfo{
					Email:  userEmail,
					UserID: data[0]["user_id"].(string),
				})
			}
		}
	}
}

func findCollectionID(collections []interface{}, collectionName string) (bool, string) {
	for i := 0; i < len(collections); i++ {
		if collections[i].(map[string]interface{})["name"] == collectionName {
			return true, collections[i].(map[string]interface{})["collectionID"].(string)
		}
	}
	return false, ""
}
