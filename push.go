package cblib

import (
	"bufio"
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
	"time"
)

func init() {
	pushCommand := &SubCommand{
		name:         "push",
		usage:        "push a specified resource to a system",
		needsAuth:    true,
		mustBeInRepo: true,
		run:          doPush,
	}

	pushCommand.flags.BoolVar(&UserSchema, "userschema", false, "push user table schema")
	pushCommand.flags.BoolVar(&AllServices, "all-services", false, "push all of the local services")
	pushCommand.flags.BoolVar(&AllLibraries, "all-libraries", false, "push all of the local libraries")
	pushCommand.flags.BoolVar(&AllDevices, "all-devices", false, "push all of the local devices")
	pushCommand.flags.BoolVar(&AllEdges, "all-edges", false, "push all of the local edges")
	pushCommand.flags.BoolVar(&AllPortals, "all-portals", false, "push all of the local portals")
	pushCommand.flags.BoolVar(&AllPlugins, "all-plugins", false, "push all of the local plugins")

	pushCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to push")
	pushCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to push")
	pushCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to push")
	pushCommand.flags.StringVar(&User, "user", "", "Name of user to push")
	pushCommand.flags.StringVar(&RoleName, "role", "", "Name of role to push")
	pushCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to push")
	pushCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to push")
	pushCommand.flags.StringVar(&DeviceName, "device", "", "Name of device to push")
	pushCommand.flags.StringVar(&EdgeName, "edge", "", "Name of edge to push")
	pushCommand.flags.StringVar(&PortalName, "portal", "", "Name of portal to push")
	pushCommand.flags.StringVar(&PluginName, "plugin", "", "Name of plugin to push")

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

func pushOneService(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing service %+s\n", ServiceName)
	service, err := getService(ServiceName)
	if err != nil {
		return err
	}
	return updateService(systemInfo.Key, service, cli)
}


/* Sample schema defintion - Keys not to be changed. Only the value
   e.g dont change "columns" or "ColumnName" etc tag names
{
    "columns": [
        {
            "ColumnName": "name",
            "ColumnType": "string",
            "PK": false
        },
        {
            "ColumnName": "city",
            "ColumnType": "string",
            "PK": false
        }
    ],
    "permissions": {}
} */
func pushUserSchema(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing user schema\n")
	exists := false
	userschema, err := getUserSchema()
	if err != nil {
		return err
	}
	userColumns, _ := cli.GetUserColumns(systemInfo.Key)
	typedSchema, ok := userschema["columns"].([]interface{})
	if !ok {
		return fmt.Errorf("Error in schema definition. Pls check the format of schema...\n")
	}
	// If user removes column from schema.json,
	// we check it by comparing length of columns in 
	// json file and no of columns on system.
	// len(userColumns) - 2 is done because there exist 2 columns
	// by default : Email & Date
	// We only want to check for columns added from schema
	if len(typedSchema) < len(userColumns)-2 {
		for i := 2; i < len(userColumns); i++ {
			exists = false
			existingColumn := userColumns[i].(map[string]interface{})["ColumnName"].(string)
			for j := 0; j < len(typedSchema); j++ {
				if typedSchema[j].(map[string]interface{})["ColumnName"].(string) == existingColumn {
					exists = true
					break
				}
			}
			if exists == false {
				if err := cli.DeleteUserColumn(systemInfo.Key, existingColumn); err != nil {
					return fmt.Errorf("User schema could not be updated. Deletion of column(s) failed\n", err)
				}
			}
		}
	} else {
		// Loop to add columns to system
		// Inner loop is used to check if column exists
		// If column exists, insertion does not work
		for i := 0; i < len(typedSchema); i++ {
			exists = false
			data := typedSchema[i].(map[string]interface{})
			columnName := data["ColumnName"].(string)
			columnType := data["ColumnType"].(string)
			for j:= 2; j < len(userColumns); j++ {
				existingColumn := userColumns[j].(map[string]interface{})["ColumnName"].(string)
				if existingColumn == columnName {
					exists = true
					break
				}
			}
			if exists == false {
				if err := cli.CreateUserColumn(systemInfo.Key, columnName, columnType); err != nil {
					return fmt.Errorf("User schema could not be updated\n", err)
				}
			}
		}
	}
	return nil
}

func pushOneCollection(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing collection %s\n", CollectionName)
	collection, err := getCollection(CollectionName)
	if err != nil {
		return err
	}
	return updateCollection(systemInfo.Key, collection, cli)
}

func pushOneCollectionById(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing collection with collectionID %s\n", CollectionId)
	collections, err := getCollections()
	if err != nil {
		return err
	}
	for _, collection := range collections {
		id, ok := collection["collectionID"].(string)
		if !ok {
			continue
		}
		if id == CollectionId {
			return updateCollection(systemInfo.Key, collection, cli)
		}
	}
	return fmt.Errorf("Collection with collectionID %+s not found.", CollectionId)
}

func pushOneUser(systemInfo *System_meta, cli *cb.DevClient) error {
	user, err := getUser(User)
	if err != nil {
		return err
	}
	return updateUser(systemInfo.Key, user, cli)
}

func pushOneUserById(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing user with user_id %s\n", UserId)
	users, err := getUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		id, ok := user["user_id"].(string)
		if !ok {
			continue
		}
		if id == UserId {
			return updateUser(systemInfo.Key, user, cli)
		}
	}
	return fmt.Errorf("User with user_id %+s not found.", UserId)
}

func pushOneRole(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing role %s\n", RoleName)
	role, err := getRole(RoleName)
	if err != nil {
		return err
	}
	return updateRole(systemInfo.Key, role, cli)
}

func pushOneTrigger(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing trigger %+s\n", TriggerName)
	trigger, err := getTrigger(TriggerName)
	if err != nil {
		return err
	}
	return updateTrigger(systemInfo.Key, trigger, cli)
}

func pushOneTimer(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing timer %+s\n", TimerName)
	timer, err := getTimer(TimerName)
	if err != nil {
		return err
	}
	return updateTimer(systemInfo.Key, timer, cli)
}

func pushOneDevice(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing device %+s\n", DeviceName)
	device, err := getDevice(DeviceName)
	if err != nil {
		return err
	}
	return updateDevice(systemInfo.Key, device, cli)
}

func pushAllDevices(systemInfo *System_meta, cli *cb.DevClient) error {
	devices, err := getDevices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		fmt.Printf("Pushing device %+s\n", device["name"].(string))
		if err := updateDevice(systemInfo.Key, device, cli); err != nil {
			return fmt.Errorf("Error updating device '%s': %s\n", device["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOneEdge(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing edge %+s\n", EdgeName)
	edge, err := getEdge(EdgeName)
	if err != nil {
		return err
	}
	return updateEdge(systemInfo.Key, edge, cli)
}

func pushAllEdges(systemInfo *System_meta, cli *cb.DevClient) error {
	edges, err := getEdges()
	if err != nil {
		return err
	}
	for _, edge := range edges {
		fmt.Printf("Pushing edge %+s\n", edge["name"].(string))
		if err := updateEdge(systemInfo.Key, edge, cli); err != nil {
			return fmt.Errorf("Error updating edge '%s': %s\n", edge["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOnePortal(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing portal %+s\n", PortalName)
	portal, err := getPortal(PortalName)
	if err != nil {
		return err
	}
	return updatePortal(systemInfo.Key, portal, cli)
}

func pushAllPortals(systemInfo *System_meta, cli *cb.DevClient) error {
	portals, err := getPortals()
	if err != nil {
		return err
	}
	for _, portal := range portals {
		fmt.Printf("Pushing portal %+s\n", portal["name"].(string))
		if err := updatePortal(systemInfo.Key, portal, cli); err != nil {
			return fmt.Errorf("Error updating portal '%s': %s\n", portal["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOnePlugin(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing portal %+s\n", PluginName)
	plugin, err := getPlugin(PluginName)
	if err != nil {
		return err
	}
	return updatePlugin(systemInfo.Key, plugin, cli)
}

func pushAllPlugins(systemInfo *System_meta, cli *cb.DevClient) error {
	plugins, err := getPlugins()
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		fmt.Printf("Pushing plugin %+s\n", plugin["name"].(string))
		if err := updatePlugin(systemInfo.Key, plugin, cli); err != nil {
			return fmt.Errorf("Error updating plugin '%s': %s\n", plugin["name"].(string), err.Error())
		}
	}
	return nil
}

func pushAllServices(systemInfo *System_meta, cli *cb.DevClient) error {
	services, err := getServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		fmt.Printf("Pushing service %+s\n", service["name"].(string))
		if err := updateService(systemInfo.Key, service, cli); err != nil {
			return fmt.Errorf("Error updating service '%s': %s\n", service["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOneLibrary(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing library %+s\n", LibraryName)

	library, err := getLibrary(LibraryName)
	if err != nil {
		return err
	}
	return updateLibrary(systemInfo.Key, library, cli)
}

func pushAllLibraries(systemInfo *System_meta, cli *cb.DevClient) error {
	libraries, err := getLibraries()
	if err != nil {
		return err
	}
	for _, library := range libraries {
		fmt.Printf("Pushing library %+s\n", library["name"].(string))
		if err := updateLibrary(systemInfo.Key, library, cli); err != nil {
			return fmt.Errorf("Error updating library '%s': %s\n", library["name"].(string), err.Error())
		}
	}
	return nil
}

func doPush(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	if err := checkPushArgsAndFlags(args); err != nil {
		return err
	}
	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}
	setRootDir(".")

	didSomething := false

	if AllServices {
		didSomething = true
		if err := pushAllServices(systemInfo, cli); err != nil {
			return err
		}
	}

	// Adding code to update user schema when pushed to system
	if UserSchema {
		didSomething = true
		if err := pushUserSchema(systemInfo, cli); err != nil {
			return err
		}
	}	

	if ServiceName != "" {
		didSomething = true
		if err := pushOneService(systemInfo, cli); err != nil {
			return err
		}
	}

	if AllLibraries {
		didSomething = true
		if err := pushAllLibraries(systemInfo, cli); err != nil {
			return err
		}
	}

	if LibraryName != "" {
		didSomething = true
		if err := pushOneLibrary(systemInfo, cli); err != nil {
			return err
		}
	}

	if CollectionName != "" {
		didSomething = true
		if err := pushOneCollection(systemInfo, cli); err != nil {
			return err
		}
	}

	if User != "" {
		didSomething = true
		if err := pushOneUser(systemInfo, cli); err != nil {
			return err
		}
	}

	if RoleName != "" {
		didSomething = true
		if err := pushOneRole(systemInfo, cli); err != nil {
			return err
		}
	}

	if TriggerName != "" {
		didSomething = true
		if err := pushOneTrigger(systemInfo, cli); err != nil {
			return err
		}
	}

	if TimerName != "" {
		didSomething = true
		if err := pushOneTimer(systemInfo, cli); err != nil {
			return err
		}
	}

	if AllDevices {
		didSomething = true
		if err := pushAllDevices(systemInfo, cli); err != nil {
			return err
		}
	}

	if DeviceName != "" {
		didSomething = true
		if err := pushOneDevice(systemInfo, cli); err != nil {
			return err
		}
	}

	if AllEdges {
		didSomething = true
		if err := pushAllEdges(systemInfo, cli); err != nil {
			return err
		}
	}

	if EdgeName != "" {
		didSomething = true
		if err := pushOneEdge(systemInfo, cli); err != nil {
			return err
		}
	}

	if AllPortals {
		didSomething = true
		if err := pushAllPortals(systemInfo, cli); err != nil {
			return err
		}
	}

	if PortalName != "" {
		didSomething = true
		if err := pushOnePortal(systemInfo, cli); err != nil {
			return err
		}
	}

	if AllPlugins {
		didSomething = true
		if err := pushAllPlugins(systemInfo, cli); err != nil {
			return err
		}
	}

	if PluginName != "" {
		didSomething = true
		if err := pushOnePlugin(systemInfo, cli); err != nil {
			return err
		}
	}

	if !didSomething {
		fmt.Printf("Nothing to push -- you must specify something to push (ie, -service=<svc_name>)\n")
	}

	return nil
}

func createRole(systemKey string, role map[string]interface{}, client *cb.DevClient) error {
	createIF, err := client.CreateRole(systemKey, role["Name"].(string))
	if err != nil {
		return err
	}
	createDict, ok := createIF.(map[string]interface{})
	if !ok {
		return fmt.Errorf("return value from CreateRole is not a map. It is %T", createIF)
	}
	fmt.Printf("Create Role returned: %+v\n", createDict)
	roleID, ok := createDict["role_id"].(string)
	if !ok {
		return fmt.Errorf("Did not get role_id key back from successful CreateRole call")
	}
	permissions, ok := role["Permissions"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Permissions for role do not exist or is not a map")
	}
	convertedPermissions := convertPermissionsStructure(permissions)
	convertedRole := map[string]interface{}{"ID": roleID, "Permissions": convertedPermissions}
	if err = client.UpdateRole(systemKey, role["Name"].(string), convertedRole); err != nil {
		return err
	}
	return nil
}

//
//  The roles structure we get back when we retrieve roles is different from
//  the format accepted for updating a role. Thus, we have this beauty of a
//  conversion function. -swm
//
//  THis is a gigantic cluster. We need to fix and learn from this. -swm
//
func convertPermissionsStructure(in map[string]interface{}) map[string]interface{} {
	fmt.Printf("convertPermissionStructure: %+v\n", in)
	out := map[string]interface{}{}
	for key, valIF := range in {
		switch key {
		case "CodeServices":
			services, err := getASliceOfMaps(valIF)
			if err != nil {
				fmt.Printf("Bad format for services permissions, not a slice of maps: %T\n", valIF)
				os.Exit(1)
			}
			svcs := make([]map[string]interface{}, len(services))
			for idx, mapVal := range services {
				svcs[idx] = map[string]interface{}{
					"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
					"permissions": mapVal["Level"],
				}
			}
			out["services"] = svcs
		case "Collections":
			collections, err := getASliceOfMaps(valIF)
			if err != nil {
				fmt.Printf("Bad format for collections permissions, not a slice of maps: %T\n", valIF)
				os.Exit(1)
			}
			cols := make([]map[string]interface{}, len(collections))
			for idx, mapVal := range collections {
				cols[idx] = map[string]interface{}{
					"itemInfo":    map[string]interface{}{"id": mapVal["ID"]},
					"permissions": mapVal["Level"],
				}
			}
			out["collections"] = cols
		case "DevicesList":
			val := getMap(valIF)
			out["devices"] = map[string]interface{}{"permissions": val["Level"]}
		case "MsgHistory":
			val := getMap(valIF)
			out["msgHistory"] = map[string]interface{}{"permissions": val["Level"]}
		case "Portals":
			portals, err := getASliceOfMaps(valIF)
			if err != nil {
				fmt.Printf("Bad format for portals permissions, not a slice of maps: %T\n", valIF)
				os.Exit(1)
			}
			ptls := make([]map[string]interface{}, len(portals))
			for idx, mapVal := range portals {
				ptls[idx] = map[string]interface{}{
					"itemInfo":    map[string]interface{}{"name": mapVal["Name"]},
					"permissions": mapVal["Level"],
				}
			}
			out["portals"] = ptls
		case "Push":
			val := getMap(valIF)
			out["push"] = map[string]interface{}{"permissions": val["Level"]}
		case "Topics":
			val, err := getASliceOfMaps(valIF)
			if err != nil {
				fmt.Printf("Bad format for topic permissions, not a slice of maps: %T\n", valIF)
				os.Exit(1)
			}
			out["topics"] = val
		case "UsersList":
			val := getMap(valIF)
			out["users"] = map[string]interface{}{"permissions": val["Level"]}

		default:
		}
	}
	fmt.Printf("convert results: %+v\n", out)
	return out
}

// The main thing I hate about go: type assertions
func getASliceOfMaps(val interface{}) ([]map[string]interface{}, error) {
	switch val.(type) {
	case []map[string]interface{}:
		return val.([]map[string]interface{}), nil
	case []interface{}:
		rval := make([]map[string]interface{}, len(val.([]interface{})))
		for idx, mapVal := range val.([]interface{}) {
			switch mapVal.(type) {
			case map[string]interface{}:
				rval[idx] = mapVal.(map[string]interface{})
			default:
				return nil, fmt.Errorf("slice values are not maps: %T\n", mapVal)
			}
		}
		return rval, nil
	default:
		return nil, fmt.Errorf("Bad type %T: expecting a slice", val)
	}
}

func getMap(val interface{}) map[string]interface{} {
	switch val.(type) {
	case map[string]interface{}:
		return val.(map[string]interface{})
	default:
		fmt.Printf("permissions type must be a map, not %T\n", val)
		os.Exit(1)
	}
	return map[string]interface{}{}
}

func updateUser(systemKey string, user map[string]interface{}, client *cb.DevClient) error {
	if id, ok := user["user_id"].(string); !ok {
		return fmt.Errorf("Missing user id %+v", user)
	} else {
		return client.UpdateUser(systemKey, id, user)
	}
}

func createUser(systemKey string, systemSecret string, user map[string]interface{}, client *cb.DevClient) (string, error) {
	email := user["email"].(string)
	password := "password"
	if pwd, ok := user["password"]; ok {
		password = pwd.(string)
	}
	newUser, err := client.RegisterNewUser(email, password, systemKey, systemSecret)
	if err != nil {
		return "", fmt.Errorf("Could not create user %s: %s", email, err.Error())
	}
	userId := newUser["user_id"].(string)
	niceRoles := mungeRoles(user["roles"].([]interface{}))
	if len(niceRoles) > 0 {
		if err := client.AddUserToRoles(systemKey, userId, niceRoles); err != nil {
			return "", err
		}
	}
	return userId, nil
}

func createTrigger(sysKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	triggerName := trigger["name"].(string)
	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = sysKey
	delete(trigger, "name")
	delete(trigger, "event_definition")
	if _, err := client.CreateEventHandler(sysKey, triggerName, trigger); err != nil {
		return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
	}
	return nil
}

func updateTrigger(systemKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	triggerName := trigger["name"].(string)
	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = systemKey
	delete(trigger, "name")
	delete(trigger, "event_definition")
	if _, err := client.UpdateEventHandler(systemKey, triggerName, trigger); err != nil {
		fmt.Printf("Could not find trigger %s\n", triggerName)
		fmt.Printf("Would you like to create a new trigger named %s? (Y/n)", triggerName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateEventHandler(systemKey, triggerName, trigger); err != nil {
					return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
				} else {
					fmt.Printf("Successfully created new trigger %s\n", triggerName)
				}
			} else {
				fmt.Printf("Trigger will not be created.\n")
			}
		}
	}
	return nil
}

func createTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) error {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}
	if _, err := client.CreateTimer(systemKey, timerName, timer); err != nil {
		return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
	}
	return nil
}

func updateTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) error {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}
	if _, err := client.UpdateTimer(systemKey, timerName, timer); err != nil {
		fmt.Printf("Could not find timer %s\n", timerName)
		fmt.Printf("Would you like to create a new timer named %s? (Y/n)", timerName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateEventHandler(systemKey, timerName, timer); err != nil {
					return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
				} else {
					fmt.Printf("Successfully created new timer %s\n", timerName)
				}
			} else {
				fmt.Printf("Timer will not be created.\n")
			}
		}
	}
	return nil
}

func updateDevice(systemKey string, device map[string]interface{}, client *cb.DevClient) error {
	deviceName := device["name"].(string)
	delete(device, "name")
	delete(device, "last_active_date")
	delete(device, "created_date")
	delete(device, "device_key")
	delete(device, "system_key")

	if _, err := client.UpdateDevice(systemKey, deviceName, device); err != nil {
		fmt.Printf("Could not find device %s\n", deviceName)
		fmt.Printf("Would you like to create a new device named %s? (Y/n)", deviceName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				device["name"] = deviceName
				if _, err := client.CreateDevice(systemKey, deviceName, device); err != nil {
					return fmt.Errorf("Could not create device %s: %s", deviceName, err.Error())
				} else {
					fmt.Printf("Successfully created new device %s\n", deviceName)
				}
			} else {
				fmt.Printf("Device will not be created.\n")
			}
		}
	}
	return nil
}

func updateEdge(systemKey string, edge map[string]interface{}, client *cb.DevClient) error {
	edgeName := edge["name"].(string)
	delete(edge, "name")
	delete(edge, "edge_key")
	delete(edge, "isConnected")
	delete(edge, "novi_system_key")
	delete(edge, "broker_auth_port")
	delete(edge, "broker_port")
	delete(edge, "broker_tls_port")
	delete(edge, "broker_ws_auth_port")
	delete(edge, "broker_ws_port")
	delete(edge, "broker_wss_port")
	delete(edge, "communication_style")
	delete(edge, "first_talked")
	delete(edge, "last_talked")
	delete(edge, "local_addr")
	delete(edge, "local_port")
	delete(edge, "public_addr")
	delete(edge, "public_port")
	delete(edge, "location")
	delete(edge, "mac_address")
	if edge["description"] == nil {
		edge["description"] = ""
	}

	_, err := client.GetEdge(systemKey, edgeName)
	if err != nil {
		// Edge does not exist
		fmt.Printf("Could not find edge %s\n", edgeName)
		fmt.Printf("Would you like to create a new edge named %s? (Y/n)", edgeName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateEdge(systemKey, edgeName, edge); err != nil {
					return fmt.Errorf("Could not create edge %s: %s", edgeName, err.Error())
				} else {
					fmt.Printf("Successfully created new edge %s\n", edgeName)
				}
			} else {
				fmt.Printf("Edge will not be created.\n")
			}
		}
	} else {
		client.UpdateEdge(systemKey, edgeName, edge)
	}
	return nil
}

func updatePortal(systemKey string, portal map[string]interface{}, client *cb.DevClient) error {
	portalName := portal["name"].(string)
	delete(portal, "system_key")
	if portal["description"] == nil {
		portal["description"] = ""
	}
	if portal["config"] == nil {
		portal["config"] = "{\"version\":1,\"allow_edit\":true,\"plugins\":[],\"panes\":[],\"datasources\":[],\"columns\":null}"
	}

	_, err := client.GetPortal(systemKey, portalName)
	if err != nil {
		// Portal DNE
		fmt.Printf("Could not find portal %s\n", portalName)
		fmt.Printf("Would you like to create a new portal named %s? (Y/n)", portalName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreatePortal(systemKey, portalName, portal); err != nil {
					return fmt.Errorf("Could not create portal %s: %s", portalName, err.Error())
				} else {
					fmt.Printf("Successfully created new portal %s\n", portalName)
				}
			} else {
				fmt.Printf("Portal will not be created.\n")
			}
		}
	} else {
		client.UpdatePortal(systemKey, portalName, portal)
	}

	return nil
}

func updatePlugin(systemKey string, plugin map[string]interface{}, client *cb.DevClient) error {
	pluginName := plugin["name"].(string)

	_, err := client.GetPlugin(systemKey, pluginName)
	if err != nil {
		// plugin DNE
		fmt.Printf("Could not find plugin %s\n", pluginName)
		fmt.Printf("Would you like to create a new plugin named %s? (Y/n)", pluginName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreatePlugin(systemKey, plugin); err != nil {
					return fmt.Errorf("Could not create plugin %s: %s", pluginName, err.Error())
				} else {
					fmt.Printf("Successfully created new plugin %s\n", pluginName)
				}
			} else {
				fmt.Printf("Plugin will not be created.\n")
			}
		}
	} else {
		client.UpdatePlugin(systemKey, pluginName, plugin)
	}

	return nil
}

func findService(systemKey, serviceName string) (map[string]interface{}, error) {
	services, err := getServices()
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if service["name"] == serviceName {
			return service, nil
		}
	}
	return nil, fmt.Errorf(NotExistErrorString)
}

func updateService(systemKey string, service map[string]interface{}, client *cb.DevClient) error {
	svcName := service["name"].(string)
	if ServiceName != "" {
		svcName = ServiceName
	}
	svcCode := service["code"].(string)
	svcDeps := service["dependencies"].(string)
	svcParams := []string{}
	for _, params := range service["params"].([]interface{}) {
		svcParams = append(svcParams, params.(string))
	}

	err, body := client.UpdateServiceWithLibraries(systemKey, svcName, svcCode, svcDeps, svcParams)
	if err != nil {
		fmt.Printf("Could not find service %s\n", svcName)
		fmt.Printf("Would you like to create a new service named %s? (Y/n)", svcName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if err := createService(systemKey, service, client); err != nil {
					return fmt.Errorf("Could not create service %s: %s", svcName, err.Error())
				} else {
					fmt.Printf("Successfully created new service %s\n", svcName)
				}
			} else {
				fmt.Printf("Service will not be created.\n")
			}
		}
	}
	if body != nil {
		service["current_version"] = body["version_number"]
		writeServiceVersion(svcName, service)
	}
	return nil
}

func createService(systemKey string, service map[string]interface{}, client *cb.DevClient) error {
	svcName := service["name"].(string)
	if ServiceName != "" {
		svcName = ServiceName
	}
	svcParams := mkSvcParams(service["params"].([]interface{}))
	svcDeps := service["dependencies"].(string)
	svcCode := service["code"].(string)
	if err := client.NewServiceWithLibraries(systemKey, svcName, svcCode, svcDeps, svcParams); err != nil {
		return err
	}
	if enableLogs(service) {
		if err := client.EnableLogsForService(systemKey, svcName); err != nil {
			return err
		}
	}
	permissions := service["permissions"].(map[string]interface{})
	//fetch roles again, find new id of role with same name
	roleIds := map[string]int{}
	for _, role := range rolesInfo {
		for roleName, level := range permissions {
			if role["Name"] == roleName {
				id := role["ID"].(string)
				roleIds[id] = int(level.(float64))
			}
		}
	}
	// now can iterate over ids instead of permission name
	for roleId, level := range roleIds {
		if err := client.AddServiceToRole(systemKey, svcName, roleId, level); err != nil {
			return err
		}
	}
	return nil
}

func updateLibrary(systemKey string, library map[string]interface{}, client *cb.DevClient) error {
	libName := library["name"].(string)
	if LibraryName != "" {
		libName = LibraryName
	}
	delete(library, "name")
	delete(library, "version")
	data, err := client.UpdateLibrary(systemKey, libName, library)
	if err != nil {
		fmt.Printf("Could not find library %s\n", libName)
		fmt.Printf("Would you like to create a new library named %s? (Y/n)", libName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				library["name"] = libName
				if err := createLibrary(systemKey, library, client); err != nil {
					return fmt.Errorf("Could not create library %s: %s", libName, err.Error())
				} else {
					fmt.Printf("Successfully created new library %s\n", libName)
				}
			} else {
				fmt.Printf("Library will not be created.\n")
			}
		}
	}
	delete(library, "code")
	library["version"] = data["version"]
	library["name"] = libName
	writeLibraryVersion(libName, library)
	return nil
}

func createLibrary(systemKey string, library map[string]interface{}, client *cb.DevClient) error {
	libName := library["name"].(string)
	if LibraryName != "" {
		libName = LibraryName
	}
	delete(library, "name")
	delete(library, "version")
	if _, err := client.CreateLibrary(systemKey, libName, library); err != nil {
		return fmt.Errorf("Could not create library %s: %s", libName, err.Error())
	}
	return nil
}

func updateCollection(systemKey string, collection map[string]interface{}, client *cb.DevClient) error {
	var err error
	collection_id := collection["collectionID"].(string)
	items := collection["items"].([]interface{})
	for _, row := range items {
		query := cb.NewQuery()
		query.EqualTo("item_id", row.(map[string]interface{})["item_id"])
		if err = client.UpdateData(collection_id, query, row.(map[string]interface{})); err != nil {
			break
		}
	}
	if err != nil {
		collName := collection["name"].(string)
		fmt.Printf("Error updating collection %s.\n", collName)
		collName = collName + "2"
		fmt.Printf("Would you like to create a new collection named %s? (Y/n)", collName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				collection["name"] = collName
				if err := createCollection(systemKey, collection, client); err != nil {
					return fmt.Errorf("Could not create collection %s: %s", collName, err.Error())
				} else {
					fmt.Printf("Successfully created new collection %s\n", collName)
				}
			} else {
				fmt.Printf("Collection will not be created.\n")
			}
		}
	}
	return nil
}

func createCollection(systemKey string, collection map[string]interface{}, client *cb.DevClient) error {
	collectionName := collection["name"].(string)
	isConnect := isConnectCollection(collection)
	var colId string
	var err error
	if isConnect {
		col, err := cb.GenerateConnectCollection(collection)
		if err != nil {
			return err
		}
		colId, err = client.NewConnectCollection(systemKey, col)
		if err != nil {
			return err
		}
	} else {
		colId, err = client.NewCollection(systemKey, collectionName)
		if err != nil {
			return err
		}
	}

	permissions := collection["permissions"].(map[string]interface{})

	roleIds := map[string]int{}
	for _, role := range rolesInfo {
		for roleName, level := range permissions {
			if role["Name"] == roleName {
				id := role["ID"].(string)
				roleIds[id] = int(level.(float64))
			}
		}
	}
	for roleId, level := range roleIds {
		if err := client.AddCollectionToRole(systemKey, colId, roleId, level); err != nil {
			return err
		}
	}

	if isConnect {
		return nil
	}

	columns := collection["schema"].([]interface{})
	for _, columnIF := range columns {
		column := columnIF.(map[string]interface{})
		colName := column["ColumnName"].(string)
		colType := column["ColumnType"].(string)
		if colName == "item_id" {
			continue
		}
		if err := client.AddColumn(colId, colName, colType); err != nil {
			return err
		}
	}
	items := collection["items"].([]interface{})
	if len(items) == 0 {
		return nil
	}
	for idx, itemIF := range items {
		items[idx] = itemIF.(map[string]interface{})
	}
	if _, err := client.CreateData(colId, items); err != nil {
		return err
	}
	return nil
}

func createEdge(systemKey, name string, edge map[string]interface{}, client *cb.DevClient) error {
	_, err := client.CreateEdge(systemKey, name, edge)
	if err != nil {
		return err
	}
	return nil
}

func createDevice(systemKey string, device map[string]interface{}, client *cb.DevClient) error {
	_, err := client.CreateDevice(systemKey, device["name"].(string), device)
	if err != nil {
		return err
	}
	return nil
}

func createPortal(systemKey string, port map[string]interface{}, client *cb.DevClient) error {
	/*
		<<<<<<< HEAD
	*/
	// Export stores config as dict, but import wants it as a string
	delete(port, "system_key")
	config, ok := port["config"]
	if ok {
		configStr := ""
		switch config.(type) {
		case string:
			configStr = config.(string)
		default:
			configBytes, err := json.Marshal(config)
			if err != nil {
				return err
			}
			configStr = string(configBytes)
		}
		port["config"] = configStr
		/*
			=======
				//remove any trash that the API doesn't like
				delete(port, "system_key")
				if port["description"] == nil {
					port["description"] = ""
			>>>>>>> f8d64d455cb9a80e4642021b1b949ec1b59dc2a0
		*/
	}
	_, err := client.CreatePortal(systemKey, port["name"].(string), port)
	if err != nil {
		return err
	}
	return nil
}

func createPlugin(systemKey string, plug map[string]interface{}, client *cb.DevClient) error {
	_, err := client.CreatePlugin(systemKey, plug)
	if err != nil {
		return err
	}
	return nil
}

func updateRole(systemKey string, role map[string]interface{}, cli *cb.DevClient) error {
	roleName := role["Name"].(string)
	if err := cli.UpdateRole(systemKey, roleName, role); err != nil {
		return fmt.Errorf("Role %s not updated\n", roleName)
	}
	return nil
}
