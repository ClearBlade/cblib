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
		name:  "push",
		usage: "push a specified resource to a system",
		run:   doPush,
	}
	pushCommand.flags.BoolVar(&UserSchema, "userschema", false, "diff user table schema")
	pushCommand.flags.StringVar(&ServiceName, "service", "", "Name of service to diff")
	pushCommand.flags.StringVar(&LibraryName, "library", "", "Name of library to diff")
	pushCommand.flags.StringVar(&CollectionName, "collection", "", "Name of collection to diff")
	pushCommand.flags.StringVar(&User, "user", "", "Name of user to diff")
	pushCommand.flags.StringVar(&RoleName, "role", "", "Name of role to diff")
	pushCommand.flags.StringVar(&TriggerName, "trigger", "", "Name of trigger to diff")
	pushCommand.flags.StringVar(&TimerName, "timer", "", "Name of timer to diff")
	AddCommand("push", pushCommand)
}

func doPush(cmd *SubCommand, cli *cb.DevClient, args ...string) error {

	systemInfo, err := getSysMeta()
	if err != nil {
		return err
	}
	setRootDir(".")

	if ServiceName != "" {
		fmt.Printf("Pushing service %+s\n", ServiceName)
		services, err := getServices()
		if err != nil {
			return err
		}
		ok := false
		for _, service := range services {
			if service["name"] == ServiceName {
				ok = true
				if err := updateService(systemInfo.Key, service, cli); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Service %+s not found\n", ServiceName)
		}
	}

	if LibraryName != "" {
		fmt.Printf("Pushing library %+s\n", LibraryName)
		libraries, err := getLibraries()
		if err != nil {
			return err
		}
		ok := false
		for _, library := range libraries {
			if library["name"] == LibraryName {
				ok = true
				if err := updateLibrary(systemInfo.Key, library, cli); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Library %+s not found\n", LibraryName)
		}
	}

	if CollectionName != "" {
		fmt.Printf("Pushing collection %+s\n", CollectionName)
		if collections, err := getCollections(); err != nil {
			return err
		} else {
			var (
				ok   = false
				coll map[string]interface{}
			)
			for _, c := range collections {
				if c["name"] == CollectionName {
					ok = true
					coll = c
				}
			}
			if !ok {
				return fmt.Errorf("Collection %+s not found\n", CollectionName)
			}
			if err := updateCollection(systemInfo.Key, coll, cli); err != nil {
				return err
			}
		}
	}

	/***** XXXSWM -- need to fix this
	if User != "" {
		fmt.Printf("Pushing user %+s\n", User)
		sysSecret := systemInfo.Secret
		if users, err := getUsers(); err != nil {
			return err
		} else {
			ok := false
			for _, user := range users {
				userMap := user.(map[string]interface{})
				if userMap["email"] == User {
					ok = true
					for _, userCol := range p.SysInfo["users"].([]interface{}) {
						column := userCol.(map[string]interface{})
						columnName := column["ColumnName"].(string)
						columnType := column["ColumnType"].(string)
						if err := cli.CreateUserColumn(systemInfo.Key, columnName, columnType); err != nil {
							return fmt.Errorf("Could not create user column %s: %s", columnName, err.Error())
						}
					}
					userId := userMap["user_id"].(string)
					if roles, err := cli.GetUserRoles(systemInfo.Key, userId); err != nil {
						return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
					} else {
						userMap["roles"] = roles
					}
					if _, err := createUser(systemInfo.Key, sysSecret, userMap, cli); err != nil {
						return fmt.Errorf("Could not create user %s: %s", User, err.Error())
					}
				}
			}
			if !ok {
				return fmt.Errorf("User %+s not found\n", User)
			}
		}
	}
	************************/

	if RoleName != "" {
		ok := false
		roles, err := getRoles()
		if err != nil {
			return fmt.Errorf("Could not get local roles: %s", err.Error())
		}
		for _, role := range roles {
			//role := roleIF.(map[string]interface{})
			if role["Name"] == RoleName {
				ok = true
				fmt.Printf("Pushing role %s\n", RoleName)
				if err := updateRole(systemInfo.Key, role, cli); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Role %s not found\n", RoleName)
		}
	}

	if TriggerName != "" {
		fmt.Printf("Pushing trigger %+s\n", TriggerName)
		triggers, err := getTriggers()
		if err != nil {
			return fmt.Errorf("Could not get local triggers: %s", err.Error())
		}
		ok := false
		for _, trigger := range triggers {
			//trigger := triggIF.(map[string]interface{})
			if trigger["name"] == TriggerName {
				ok = true
				if err := updateTrigger(systemInfo.Key, trigger, cli); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Trigger %+s not found\n", TriggerName)
		}
	}

	if TimerName != "" {
		fmt.Printf("Pushing timer %+s\n", TimerName)
		timers, err := getTimers()
		if err != nil {
			return fmt.Errorf("Could not get local timers: %s", err.Error())
		}
		ok := false
		for _, timer := range timers {
			if timer["name"] == TimerName {
				ok = true
				if err := updateTimer(systemInfo.Key, timer, cli); err != nil {
					return err
				}
			}
		}
		if !ok {
			return fmt.Errorf("Timer %+s not found\n", TimerName)
		}
	}

	return nil
}

func createRole(systemKey string, role map[string]interface{}, client *cb.DevClient) error {
	if _, err := client.CreateRole(systemKey, role["Name"].(string)); err != nil {
		return err
	}
	return nil
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
	svcCode := service["code"].(string)
	svcParams := []string{}
	for _, params := range service["params"].([]interface{}) {
		svcParams = append(svcParams, params.(string))
	}
	if err := client.UpdateService(systemKey, svcName, svcCode, svcParams); err != nil {
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
		// query.EqualTo(field, value)
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

func updateRole(systemKey string, role map[string]interface{}, cli *cb.DevClient) error {
	roleName := role["Name"].(string)
	if err := cli.UpdateRole(systemKey, roleName, role); err != nil {
		return fmt.Errorf("Role %s not updated\n", roleName)
	}
	return nil
}
