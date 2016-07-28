package cblib

import (
	"bufio"
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
	pushCommand.flags.BoolVar(&AllDashboards, "all-dashboards", false, "push all of the local dashboards")
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
	pushCommand.flags.StringVar(&DashboardName, "dashboard", "", "Name of dashboard to push")
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
	fmt.Printf("Pushing user %s\n", User)
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

func pushOneDashboard(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing dashboard %+s\n", DashboardName)
	dashboard, err := getDashboard(DashboardName)
	if err != nil {
		return err
	}
	return updateDashboard(systemInfo.Key, dashboard, cli)
}

func pushAllDashboards(systemInfo *System_meta, cli *cb.DevClient) error {
	dashboards, err := getDashboards()
	if err != nil {
		return err
	}
	for _, dashboard := range dashboards {
		fmt.Printf("Pushing dashboard %+s\n", dashboard["name"].(string))
		if err := updateDashboard(systemInfo.Key, dashboard, cli); err != nil {
			return fmt.Errorf("Error updating dashboard '%s': %s\n", dashboard["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOnePlugin(systemInfo *System_meta, cli *cb.DevClient) error {
	fmt.Printf("Pushing dashboard %+s\n", PluginName)
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

	if AllDashboards {
		didSomething = true
		if err := pushAllDashboards(systemInfo, cli); err != nil {
			return err
		}
	}

	if DashboardName != "" {
		didSomething = true
		if err := pushOneDashboard(systemInfo, cli); err != nil {
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
	if _, err := client.CreateRole(systemKey, role["Name"].(string)); err != nil {
		return err
	}
	return nil
}

func updateUser(systemKey string, user map[string]interface{}, client *cb.DevClient) error {
	if id, ok := user["Id"].(string); !ok {
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
	if(edge["description"] == nil){ edge["description"] = "" }

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

func updateDashboard(systemKey string, dashboard map[string]interface{}, client *cb.DevClient) error {
	dashboardName := dashboard["name"].(string)
	delete (dashboard, "system_key")
	if(dashboard["description"] == nil){ dashboard["description"] = "" }
	if(dashboard["config"] == nil){ dashboard["config"] = "{\"version\":1,\"allow_edit\":true,\"plugins\":[],\"panes\":[],\"datasources\":[],\"columns\":null}" }

	_, err := client.GetDashboard(systemKey, dashboardName)
	if err != nil {
		// Dashboard DNE
		fmt.Printf("Could not find dashboard %s\n", dashboardName)
		fmt.Printf("Would you like to create a new dashboard named %s? (Y/n)", dashboardName)
		reader := bufio.NewReader(os.Stdin)
		if text, err := reader.ReadString('\n'); err != nil {
			return err
		} else {
			if strings.Contains(strings.ToUpper(text), "Y") {
				if _, err := client.CreateDashboard(systemKey, dashboardName, dashboard); err != nil {
					return fmt.Errorf("Could not create dashboard %s: %s", dashboardName, err.Error())
				} else {
					fmt.Printf("Successfully created new dashboard %s\n", dashboardName)
				}
			} else {
				fmt.Printf("Dashboard will not be created.\n")
			}
		}
	} else {
		client.UpdateDashboard(systemKey, dashboardName, dashboard);
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
		client.UpdatePlugin(systemKey, pluginName, plugin);
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
	if err := client.UpdateServiceWithLibraries(systemKey, svcName, svcCode, svcDeps, svcParams); err != nil {
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
	if _, err := client.UpdateLibrary(systemKey, libName, library); err != nil {
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
		fmt.Println(err.Error())
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
	colId, err := client.NewCollection(systemKey, collectionName)
	if err != nil {
		return err
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

func createDashboard(systemKey string, dash map[string]interface{}, client *cb.DevClient) error {
	_, err := client.CreateDashboard(systemKey, dash["name"].(string), dash)
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
