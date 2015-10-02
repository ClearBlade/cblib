package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
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

func createUsers(systemInfo map[string]interface{}, users []map[string]interface{}, client *cb.DevClient) error {

	//  Create user columns first -- if any
	sysKey := systemInfo["systemKey"].(string)
	sysSec := systemInfo["systemSecret"].(string)
	userSchema, err := getUserSchema()
	userCols := userSchema["columns"].([]interface{})
	userPerms := userSchema["permissions"].(map[string]interface{})
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

	for roleName, levelFloat := range userPerms {
		level := int(levelFloat.(float64))
		if err := client.AddGenericPermissionToRole(sysKey, roleName, "users", level); err != nil {
			return err
		}
	}

	if !importUsers {
		return nil
	}

	// Now, create users -- register, update roles, and update user-def colunms
	for _, user := range users {
		userId, err := createUser(sysKey, sysSec, user, client)
		if err != nil {
			return err
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

func unMungeRoles(roles []string) []interface{} {
	rval := []interface{}{}

	for _, role := range roles {
		rval = append(rval, role)
	}
	return rval
}

func createTriggers(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	triggers, err := getTriggers()
	if err != nil {
		return err
	}
	for _, trigger := range triggers {
		if err := createTrigger(sysKey, trigger, client); err != nil {
			return err
		}
	}
	return nil
}

func createTimers(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	timers, err := getTimers()
	if err != nil {
		return err
	}
	for _, timer := range timers {
		if err := createTimer(sysKey, timer, client); err != nil {
			return err
		}
	}
	return nil
}

func createServices(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	services, err := getServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		if err := createService(sysKey, service, client); err != nil {
			return err
		}
	}
	return nil
}

func createLibraries(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	libraries, err := getLibraries()
	if err != nil {
		return err
	}
	for _, library := range libraries {
		if err := createLibrary(sysKey, library, client); err != nil {
			return err
		}
	}
	return nil
}

func createCollections(systemInfo map[string]interface{}, client *cb.DevClient) error {
	sysKey := systemInfo["systemKey"].(string)
	collections, err := getCollections()
	if err != nil {
		return err
	}
	for _, collection := range collections {
		if err := createCollection(sysKey, collection, client); err != nil {
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
