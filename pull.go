package cblib

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/models"
)

var (
	PULL_ALL_USERS = "%PULL_ALL_USERS%"
)

func init() {
	usage :=
		`
	Pull a ClearBlade asset from the Platform to your local filesystem. Use -sort-collections for easier version controlling of datasets.

	Note: Collection rows are pulled by default.
	`

	example :=
		`
	cb-cli pull -service=Service1 									# Pulls Service1 from Platform to local filesystem
	cb-cli pull -collection=Collection1								# Pulls Collection1 from Platform to local filesystem, with all rows, unsorted
	cb-cli pull -collection=Collection1 -sort-collections=true		# Pulls Collection1 from Platform to local filesystem, with all rows, sorted
	`
	pullCommand := &SubCommand{
		name:         "pull",
		usage:        usage,
		needsAuth:    true,
		mustBeInRepo: true,
		run:          doPull,
		example:      example,
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
	pullCommand.flags.BoolVar(&AllAssets, "all", false, "pull all assets from system")
	pullCommand.flags.BoolVar(&AllTriggers, "all-triggers", false, "pull all triggers from system")
	pullCommand.flags.BoolVar(&AllTimers, "all-timers", false, "pull all timers from system")

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

	AddCommand("pull", pullCommand)
}

func doPull(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	SetRootDir(".")
	systemInfo, err := getSysMeta()
	setupDirectoryStructure(systemInfo)
	if err != nil {
		return err
	}

	// This is a hack to check if token has expired and auth again
	// since we dont have an endpoint to determine this
	client, err = checkIfTokenHasExpired(client, systemInfo.Key)
	if err != nil {
		return fmt.Errorf("Re-auth failed: %s\n", err)
	}

	// ??? we already have them locally
	if _, err := PullAndWriteRoles(systemInfo.Key, client, false); err != nil {
		return err
	}

	didSomething := false

	if AllServices || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all services:")
		if _, err := PullServices(systemInfo.Key, client); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if ServiceName != "" {
		didSomething = true
		fmt.Printf("Pulling service %+s\n", ServiceName)
		if err := PullAndWriteService(systemInfo.Key, ServiceName, client); err != nil {
			return err
		}
	}

	if AllLibraries || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all libraries:")
		if _, err := PullLibraries(systemInfo, client); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if LibraryName != "" {
		didSomething = true
		fmt.Printf("Pulling library %s\n", LibraryName)
		if lib, err := pullLibrary(systemInfo.Key, LibraryName, client); err != nil {
			return err
		} else {
			writeLibrary(lib["name"].(string), lib)
		}
	}

	if AllCollections || AllAssets {
		didSomething = true
		ExportRows = true
		fmt.Printf("Pulling all collections:")
		if _, err := pullCollections(systemInfo, client); err != nil {
			fmt.Printf("Error: Failed to pull all collections - %s\n", err.Error())
		}
	}

	if CollectionName != "" {
		didSomething = true
		ExportRows = true
		fmt.Printf("Pulling collection %+s\n", CollectionName)
		err := PullAndWriteCollection(systemInfo, CollectionName, client)
		if err != nil {
			return err
		}
	}

	if AllUsers || AllAssets {
		didSomething = true
		fmt.Println("Pulling all users:")
		if err := PullAndWriteUsers(systemInfo.Key, PULL_ALL_USERS, client); err != nil {
			fmt.Printf("Error: Failed to pull all users - %s\n", err.Error())
		}
		if _, err := pullUserSchemaInfo(systemInfo.Key, client, true); err != nil {
			fmt.Printf("Error: Failed to pull user schema - %s\n", err.Error())
		}
	}

	if User != "" {
		didSomething = true
		fmt.Printf("Pulling user %+s\n", User)
		err := PullAndWriteUsers(systemInfo.Key, User, client)
		if err != nil {
			return err
		}
		if _, err := pullUserSchemaInfo(systemInfo.Key, client, true); err != nil {
			return err
		}
	}

	if AllRoles || AllAssets {
		didSomething = true
		fmt.Println("Pulling all roles:")
		if _, err := PullAndWriteRoles(systemInfo.Key, client, true); err != nil {
			fmt.Printf("Error: Failed to pull all roles - %s\n", err.Error())
		}
	}

	if RoleName != "" {
		didSomething = true
		roles := make([]map[string]interface{}, 0)
		splitRoles := strings.Split(RoleName, ",")
		for _, role := range splitRoles {
			fmt.Printf("Pulling role %+s\n", role)
			if r, err := pullRole(systemInfo.Key, role, client); err != nil {
				return err
			} else {
				roles = append(roles, r)
				writeRole(role, r)
			}
		}
	}

	if AllTriggers || AllAssets {
		didSomething = true
		fmt.Println("Pulling all triggers:")
		if _, err := PullAndWriteTriggers(systemInfo, client); err != nil {
			fmt.Printf("Error: Failed to pull all triggers - %s\n", err.Error())
		}
	}

	if TriggerName != "" {
		didSomething = true
		fmt.Printf("Pulling trigger %+s\n", TriggerName)
		err := PullAndWriteTrigger(systemInfo.Key, TriggerName, client)
		if err != nil {
			return err
		}
	}

	if AllTimers || AllAssets {
		didSomething = true
		fmt.Println("Pulling all timers:")
		if _, err := PullAndWriteTimers(systemInfo, client); err != nil {
			fmt.Printf("Error: Failed to pull all timers - %s\n", err.Error())
		}
	}

	if TimerName != "" {
		didSomething = true
		fmt.Printf("Pulling timer %+s\n", TimerName)
		err := PullAndWriteTimer(systemInfo.Key, TimerName, client)
		if err != nil {
			return err
		}
	}

	if AllDevices || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all devices:")
		if _, err := PullDevices(systemInfo, client); err != nil {
			return err
		}
		if _, err := pullDevicesSchema(systemInfo.Key, client, true); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if DeviceName != "" {
		didSomething = true
		fmt.Printf("Pulling device %+s\n", DeviceName)
		if device, err := pullDevice(systemInfo.Key, DeviceName, client); err != nil {
			return err
		} else {
			if _, err := pullDevicesSchema(systemInfo.Key, client, true); err != nil {
				return err
			}
			writeDevice(DeviceName, device)
		}
	}

	if AllEdges || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all edges:")
		if _, err := PullEdges(systemInfo, client); err != nil {
			return err
		}
		if _, err := pullEdgesSchema(systemInfo.Key, client, true); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if EdgeName != "" {
		didSomething = true
		fmt.Printf("Pulling edge %+s\n", EdgeName)
		if edge, err := pullEdge(systemInfo.Key, EdgeName, client); err != nil {
			return err
		} else {
			writeEdge(EdgeName, edge)
		}
		if _, err := pullEdgesSchema(systemInfo.Key, client, true); err != nil {
			fmt.Printf("\nNo custom columns to pull and create schema.json from... Continuing...\n")
		}
	}

	if AllPortals || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all portals:")
		if _, err := PullPortals(systemInfo, client); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if PortalName != "" {
		didSomething = true
		fmt.Printf("Pulling portal %+s\n", PortalName)
		if err := PullAndWritePortal(systemInfo.Key, PortalName, client); err != nil {
			return err
		}
	}

	if AllPlugins || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all plugins:")
		if _, err := PullPlugins(systemInfo, client); err != nil {
			return err
		}
		fmt.Printf("\n")
	}

	if PluginName != "" {
		didSomething = true
		fmt.Printf("Pulling plugin %+s\n", PluginName)
		if err = PullAndWritePlugin(systemInfo.Key, PluginName, client); err != nil {
			return err
		}
	}

	if AllAdaptors || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all adaptors:")
		if err := backupAndCleanDirectory(adaptorsDir); err != nil {
			return err
		}
		if err := PullAdaptors(systemInfo, client); err != nil {
			if restoreErr := restoreBackupDirectory(adaptorsDir); restoreErr != nil {
				fmt.Printf("Failed to restore backup directory; %s\n", restoreErr.Error())
			}
			return err
		}
		if err := removeBackupDirectory(adaptorsDir); err != nil {
			fmt.Printf("Warning: Failed to remove backup directory for '%s'", adaptorsDir)
		}
		fmt.Printf("\n")
	}

	if AdaptorName != "" {
		didSomething = true
		fmt.Printf("Pulling adaptor %+s\n", AdaptorName)
		if err = PullAndWriteAdaptor(systemInfo.Key, AdaptorName, client); err != nil {
			return err
		}
	}

	if AllDeployments || AllAssets {
		didSomething = true
		fmt.Printf("Pulling all deployments:")
		if _, err = pullDeployments(systemInfo, client); err != nil {
			fmt.Printf("Error - Failed to pull all deployments: %s\n", err.Error())
		}
	}

	if DeploymentName != "" {
		didSomething = true
		fmt.Printf("Pulling deployment %+s\n", DeploymentName)
		if _, err = pullAndWriteDeployment(systemInfo, client, DeploymentName); err != nil {
			fmt.Printf("Error - Failed to pull deployment '%s': %s\n", DeploymentName, err.Error())
		}
	}

	if !didSomething {
		fmt.Printf("Nothing to pull -- you must specify something to pull (ie, -service=<svc_name>)\n")
	}
	return nil
}

func pullUserSchemaInfo(systemKey string, cli *cb.DevClient, writeThem bool) (map[string]interface{}, error) {
	resp, err := cli.GetUserColumns(systemKey)
	if err != nil {
		return nil, err
	}
	columns := []map[string]interface{}{}
	for _, colIF := range resp {
		col := colIF.(map[string]interface{})
		if col["ColumnName"] == "email" || col["ColumnName"] == "creation_date" {
			continue
		}
		columns = append(columns, col)
	}
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
	if service, err := client.GetServiceRaw(systemKey, serviceName); err != nil {
		return nil, err
	} else {
		service["code"] = strings.Replace(service["code"].(string), "\\n", "\n", -1)
		return service, nil
	}
}

func PullAndWriteLibrary(systemKey string, libraryName string, client *cb.DevClient) error {
	if svc, err := pullLibrary(systemKey, libraryName, client); err != nil {
		return err
	} else {
		return writeLibrary(libraryName, svc)
	}
}

func PullAndWriteUsers(systemKey string, userName string, client *cb.DevClient) error {
	if users, err := client.GetAllUsers(systemKey); err != nil {
		return err
	} else {
		ok := false
		for _, user := range users {
			if user["email"] == userName || userName == PULL_ALL_USERS {
				ok = true
				userId := user["user_id"].(string)
				if roles, err := client.GetUserRoles(systemKey, userId); err != nil {
					return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
				} else {
					user["roles"] = roles
				}
				err = writeUser(user["email"].(string), user)
				if err != nil {
					return err
				}
			}
		}
		if !ok {
			if userName == PULL_ALL_USERS {
				return fmt.Errorf("No users found")
			} else {
				return fmt.Errorf("User %+s not found\n", userName)
			}

		}
	}
	return nil
}

func PullAndWriteCollection(systemInfo *System_meta, collectionName string, client *cb.DevClient) error {
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
			if data, err := PullCollection(systemInfo, coll, client); err != nil {
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

func PullAndWriteCollections(sysMeta *System_meta, client *cb.DevClient) error {
	if allColls, err := client.GetAllCollections(sysMeta.Key); err != nil {
		return err
	} else {
		// iterate over allColls and find one with matching name
		for _, c := range allColls {
			coll := c.(map[string]interface{})
			if coll, err := client.GetCollectionInfo(coll["collectionID"].(string)); err != nil {
				return err
			} else {
				if data, err := PullCollection(sysMeta, coll, client); err != nil {
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

func PullAndWriteTrigger(systemKey, trigName string, client *cb.DevClient) error {
	if trigg, err := pullTrigger(systemKey, trigName, client); err != nil {
		return err
	} else {
		stripTriggerFields(trigg)
		err = writeTrigger(trigName, trigg)
		if err != nil {
			return err
		}
	}
	return nil
}

func PullAndWriteTriggers(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	trigs, err := cli.GetEventHandlers(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull triggers out of system %s: %s", sysMeta.Key, err.Error())
	}
	triggers := []map[string]interface{}{}
	for _, trig := range trigs {
		thisTrig := trig.(map[string]interface{})
		delete(thisTrig, "system_key")
		delete(thisTrig, "system_secret")
		triggers = append(triggers, thisTrig)
		err = writeTrigger(thisTrig["name"].(string), thisTrig)
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

func PullAndWriteTimers(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theTimers, err := cli.GetTimers(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull timers out of system %s: %s", sysMeta.Key, err.Error())
	}
	timers := []map[string]interface{}{}
	for _, timer := range theTimers {
		thisTimer := timer.(map[string]interface{})
		timers = append(timers, thisTimer)
		err = writeTimer(thisTimer["name"].(string), thisTimer)
		if err != nil {
			return nil, err
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
