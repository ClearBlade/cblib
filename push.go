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
	pushCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to push")
	pushCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to push")
	pushCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to push")
	pushCommand.flags.StringVar(&User, "user", "", "Name of user to push")
	pushCommand.flags.StringVar(&RoleName, "role", "", "Name of role to push")
	pushCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to push")
	pushCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to push")
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

	err, body := client.UpdateServiceWithLibraries(systemKey, svcName, svcCode, svcDeps, svcParams); 
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
	if (body != nil) {
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
	data, err := client.UpdateLibrary(systemKey, libName, library); 
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
