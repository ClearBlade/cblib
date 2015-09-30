package cblib

import (
	//"flag"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"time"
)

var (
	importRows  bool
	importUsers bool
)

func init() {
	myImportCommand := &SubCommand{
		name:  "import",
		usage: "just import stuff",
		run:   doImport,
	}
	myImportCommand.flags.BoolVar(&importRows, "importrows", false, "imports all data into all collections")
	myImportCommand.flags.BoolVar(&importUsers, "importusers", false, "imports all users into the system")
	AddCommand("import", myImportCommand)
	AddCommand("imp", myImportCommand)
	AddCommand("im", myImportCommand)
}

func createSystem(system map[string]interface{}, client *cb.DevClient) error {
	name := system["name"].(string)
	desc := system["description"].(string)
	auth := system["auth"].(bool)
	sysKey, sysErr := client.NewSystem(name, desc, auth)
	if sysErr != nil {
		return sysErr
	}
	realSystem, sysErr := client.GetSystem(sysKey)
	if sysErr != nil {
		return sysErr
	}
	system["systemKey"] = realSystem.Key
	system["systemSecret"] = realSystem.Secret
	return nil
}

//func createUsers(systemInfo map[string]interface{}, users []interface{}, client *cb.DevClient) error {
func createUsers(systemInfo map[string]interface{}, users []map[string]interface{}, client *cb.DevClient) error {

	//  Create user columns first -- if any
	sysKey := systemInfo["systemKey"].(string)
	sysSec := systemInfo["systemSecret"].(string)
	//userCols := systemInfo["users"].([]interface{})
	userSchema, err := getUserSchema()
	userCols := userSchema["columns"].([]interface{})
	if err != nil {
		return err
	}
	for _, columnIF := range userCols {
		column := columnIF.(map[string]interface{})
		columnName := column["ColumnName"].(string)
		columnType := column["ColumnType"].(string)
		if err := client.CreateUserColumn(sysKey, columnName, columnType); err != nil {
			return fmt.Errorf("Could not create user column %s: %s", columnName, err.Error())
		}
	}

	//  XXXSWM TODO -- add permissions to the users table:
	//  userTablePerms := userSchema["permissions"]

	if !importUsers {
		return nil
	}

	// Now, create users -- register, update roles, and update user-def colunms
	for _, user := range users {
		//user := userIF.(map[string]interface{})
		email := user["email"].(string)
		password := "password"
		if pwd, ok := user["password"]; ok {
			password = pwd.(string)
		}
		newUser, err := client.RegisterNewUser(email, password, sysKey, sysSec)
		if err != nil {
			return fmt.Errorf("Could not create user %s: %s", email, err.Error())
		}
		userId := newUser["user_id"].(string)
		niceRoles := mungeRoles(user["roles"].([]interface{}))
		if len(niceRoles) > 0 {
			if err := client.AddUserToRoles(sysKey, userId, niceRoles); err != nil {
				return err
			}
		}

		if len(userCols) == 0 {
			continue
		}

		updates := map[string]interface{}{}
		for _, columnIF := range userCols {
			column := columnIF.(map[string]interface{})
			columnName := column["ColumnName"].(string)
			if userVal, ok := user[columnName]; ok {
				if userVal != nil {
					updates[columnName] = userVal
				}
			}
		}

		if len(updates) == 0 {
			continue
		}

		if err := client.UpdateUser(sysKey, userId, updates); err != nil {
			return fmt.Errorf("Could not update user: %s", err.Error())
		}
	}

	return nil
}

func mungeRoles(roles []interface{}) []string {
	rval := []string{}
	for _, role := range roles {
		roleStr := role.(string)
		if roleStr == "Authenticated" { // This automagically happens when user auth'd
			continue
		}
		rval = append(rval, roleStr)
	}
	return rval
}

func createTriggers(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	//triggers := systemInfo["triggers"].([]interface{})
	triggers, err := getTriggers()
	if err != nil {
		return err
	}
	//for _, triggerIF := range triggers {
	for _, trigger := range triggers {
		//trigger := triggerIF.(map[string]interface{})
		triggerName := trigger["name"].(string)
		triggerDef := trigger["event_definition"].(map[string]interface{})
		trigger["def_module"] = triggerDef["def_module"]
		trigger["def_name"] = triggerDef["def_name"]
		trigger["system_key"] = systemInfo["systemKey"]
		delete(trigger, "name")
		delete(trigger, "event_definition")
		if _, err := client.CreateEventHandler(sysKey, triggerName, trigger); err != nil {
			return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
		}
	}
	return nil
}

func createTimers(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	//timers := systemInfo["timers"].([]interface{})
	timers, err := getTimers()
	if err != nil {
		return err
	}
	//for _, timerIF := range timers {
	for _, timer := range timers {
		//timer := timerIF.(map[string]interface{})
		timerName := timer["name"].(string)
		delete(timer, "name")
		startTime := timer["start_time"].(string)
		if startTime == "Now" {
			timer["start_time"] = time.Now().Format(time.RFC3339)
		}
		if _, err := client.CreateTimer(sysKey, timerName, timer); err != nil {
			return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
		}
	}
	return nil
}

func createServices(systemInfo map[string]interface{}, client *cb.DevClient) error {
	//services := systemInfo["services"].([]interface{})
	services, err := getServices()
	if err != nil {
		return err
	}
	sysKey := systemInfo["systemKey"].(string)
	//for _, serviceIF := range services {
	for _, service := range services {
		//service := serviceIF.(map[string]interface{})
		svcName := service["name"].(string)
		svcParams := mkSvcParams(service["params"].([]interface{}))
		svcDeps := service["dependencies"].(string)
		//svcCode, err := getServiceCode(svcName)
		svcCode := service["code"].(string)
		delete(service, "code")
		if err != nil {
			return err
		}
		if err := client.NewServiceWithLibraries(sysKey, svcName, svcCode, svcDeps, svcParams); err != nil {
			return err
		}
		if enableLogs(service) {
			if err := client.EnableLogsForService(sysKey, svcName); err != nil {
				return err
			}
		}
		permissions := service["permissions"].(map[string]interface{})
		for roleId, level := range permissions {
			if err := client.AddServiceToRole(sysKey, svcName, roleId, int(level.(float64))); err != nil {
				return err
			}
		}
	}
	return nil
}

func createLibraries(systemInfo map[string]interface{}, client *cb.DevClient) error {
	//libraries := systemInfo["libraries"].([]interface{})
	libraries, err := getLibraries()
	if err != nil {
		return err
	}
	sysKey := systemInfo["systemKey"].(string)
	//for _, libraryIF := range libraries {
	for _, library := range libraries {
		//library := libraryIF.(map[string]interface{})
		libName := library["name"].(string)
		/*
			libCode, err := getLibraryCode(libName)
			if err != nil {
				return err
			}
			library["code"] = libCode
		*/
		delete(library, "name")
		delete(library, "version")
		if _, err := client.CreateLibrary(sysKey, libName, library); err != nil {
			return fmt.Errorf("Could not create library %s: %s", libName, err.Error())
		}
	}
	return nil
}

func createCollections(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	//collections := systemInfo["data"].([]interface{})
	collections, err := getCollections()
	if err != nil {
		return err
	}
	for _, collection := range collections {
		//  Create the collection
		//collection := collectionIF.(map[string]interface{})
		collectionName := collection["name"].(string)
		colId, err := client.NewCollection(sysKey, collectionName)
		if err != nil {
			return err
		}

		permissions := collection["permissions"].(map[string]interface{})
		for roleId, level := range permissions {
			if err := client.AddCollectionToRole(sysKey, colId, roleId, int(level.(float64))); err != nil {
				return err
			}
		}

		//  Add the columns
		columns := collection["schema"].([]interface{})
		for _, columnIF := range columns {
			column := columnIF.(map[string]interface{})
			colName := column["ColumnName"].(string)
			colType := column["ColumnType"].(string)
			if colName == "item_id" {
				continue
			}
			if err := client.AddColumn(colId, colName, colType); err != nil {
				fmt.Printf("Add column: %s, %s, %s\n", collectionName, colName, colType)
				return err
			}
		}
		if !importRows {
			continue
		}

		//  Add the items
		//itemsIF, err := getCollectionItems(collectionName)
		itemsIF := collection["items"].([]interface{})
		if err != nil {
			return err
		}
		if len(itemsIF) == 0 {
			continue
		}
		items := make([]map[string]interface{}, len(itemsIF))
		for idx, itemIF := range itemsIF {
			items[idx] = itemIF.(map[string]interface{})
		}
		if _, err := client.CreateData(colId, items); err != nil {
			return err
		}
	}
	return nil
}

func enableLogs(service map[string]interface{}) bool {
	logVal, ok := service["logging_enabled"]
	if !ok {
		return false
	}
	switch logVal.(type) {
	case string:
		return logVal.(string) == "true"
	case bool:
		return logVal.(bool) == true
	}
	return false
}

func mkSvcParams(params []interface{}) []string {
	rval := []string{}
	for _, val := range params {
		rval = append(rval, val.(string))
	}
	return rval
}

func doImport(cmd *SubCommand, cli *cb.DevClient, args ...string) error {
	return importIt(cli)
}

func importIt(cli *cb.DevClient) error {
	cb.CB_ADDR = URL
	fmt.Printf("Reading system configuration files...")
	setRootDir(".")
	users, err := getUsers()
	if err != nil {
		return err
	}

	systemInfo, err := getDict("system.json")
	if err != nil {
		return err
	}
	fmt.Printf("Done.\nImporting system...")
	if err := createSystem(systemInfo, cli); err != nil {
		return fmt.Errorf("Could not create system %s: %s", systemInfo["name"], err.Error())
	}
	fmt.Printf("Done.\nImporting users...")
	if err := createUsers(systemInfo, users, cli); err != nil {
		return fmt.Errorf("Could not create users: %s", err.Error())
	}
	fmt.Printf("Done.\nImporting collections...")
	if err := createCollections(systemInfo, cli); err != nil {
		return fmt.Errorf("Could not create collections: %s", err.Error())
	}
	fmt.Printf("Done.\nImporting code services...")
	if err := createServices(systemInfo, cli); err != nil {
		return err
	}
	fmt.Printf("Done.\nImporting code libraries...")
	if err := createLibraries(systemInfo, cli); err != nil {
		return err
	}
	fmt.Printf("Done.\nImporting triggers...")
	if err := createTriggers(systemInfo, cli); err != nil {
		return fmt.Errorf("Could not create triggers: %s", err.Error())
	}
	fmt.Printf("Done.\nImporting timers...")
	if err := createTimers(systemInfo, cli); err != nil {
		return fmt.Errorf("Could not create timers: %s", err.Error())
	}

	fmt.Printf("Done\n")
	return nil
}
