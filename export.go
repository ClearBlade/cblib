package cblib

import (
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"os"
	"strings"
)

var (
	exportRows  bool
	exportUsers bool
)

func init() {
	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	rolesInfo = []map[string]interface{}{}
	myExportCommand := &SubCommand{
		name:  "export",
		usage: "Ain't no thing",
		run:   doExport,
		//  TODO -- add help, usage, etc.
	}
	myExportCommand.flags.BoolVar(&exportRows, "exportrows", false, "exports all data from all collections")
	myExportCommand.flags.BoolVar(&exportUsers, "exportusers", false, "exports user info")
	AddCommand("export", myExportCommand)
	AddCommand("ex", myExportCommand)
	AddCommand("exp", myExportCommand)
	ImportPageSize = 100 // TODO -- fix this
}

func pullRoles(systemKey string, cli *cb.DevClient, writeThem bool) ([]map[string]interface{}, error) {
	r, err := cli.GetAllRoles(systemKey)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, len(r))
	for idx, rIF := range r {
		thisRole := rIF.(map[string]interface{})
		rval[idx] = thisRole
		if writeThem {
			if err := writeRole(thisRole["Name"].(string), thisRole); err != nil {
				return nil, err
			}
		}
	}
	return rval, nil
}
func storeRoles(roles []map[string]interface{}) {
	roleList := make([]string, len(roles))
	for idx, role := range roles {
		roleList[idx] = role["Name"].(string)
	}
	systemDotJSON["roles"] = roleList
}

func pullCollections(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	colls, err := cli.GetAllCollections(sysMeta.Key)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, len(colls))
	for i, col := range colls {
		if r, err := pullCollection(sysMeta, col.(map[string]interface{}), cli); err != nil {
			return nil, err
		} else {
			writeCollection(r["name"].(string), r)
			rval[i] = r
		}
	}
	return rval, nil
}

func pullCollection(sysMeta *System_meta, co map[string]interface{}, cli *cb.DevClient) (map[string]interface{}, error) {
	id := co["collectionID"].(string)
	columnsResp, err := cli.GetColumns(id, sysMeta.Key, sysMeta.Secret)
	if err != nil {
		return nil, err
	}
	co["schema"] = columnsResp
	if err := getRolesForCollection(co); err != nil {
		return nil, err
	}
	co["items"] = []interface{}{}
	co["items"] = []interface{}{}
	if exportRows {
		items, err := pullCollectionData(co, cli)
		if err != nil {
			return nil, err
		}
		co["items"] = items
	}
	return co, nil
}

func getRolesForCollection(collection map[string]interface{}) error {
	colName := collection["name"].(string)
	perms := map[string]interface{}{}
	for _, role := range rolesInfo {
		roleName := role["Name"].(string)
		colPerms := role["Permissions"].(map[string]interface{})["Collections"].([]interface{})
		for _, colPermIF := range colPerms {
			colPerm := colPermIF.(map[string]interface{})
			if colPerm["Name"].(string) == colName {
				perms[roleName] = colPerm["Level"]
			}
		}
	}
	collection["permissions"] = perms
	return nil
}

func pullCollectionData(collection map[string]interface{}, client *cb.DevClient) ([]interface{}, error) {
	colId := collection["collectionID"].(string)

	totalItems, err := client.GetItemCount(colId)
	if err != nil {
		return nil, fmt.Errorf("GetItemCount Failed: %s", err.Error())
	}

	dataQuery := &cb.Query{}
	dataQuery.PageSize = ImportPageSize
	allData := []interface{}{}
	for j := 0; j < totalItems; j += ImportPageSize {
		dataQuery.PageNumber = (j / ImportPageSize) + 1
		data, err := client.GetData(colId, dataQuery)
		if err != nil {
			return nil, err
		}
		curData := data["DATA"].([]interface{})
		for _, oneItemIF := range curData {
			oneItem := oneItemIF.(map[string]interface{})
			delete(oneItem, "item_id")
		}
		allData = append(allData, curData...)
	}
	return allData, nil
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
	tablePerms := getUserTablePermissions()
	schema := map[string]interface{}{
		"columns":     columns,
		"permissions": tablePerms,
	}
	if writeThem {
		if err := writeUser("schema", schema); err != nil {
			return nil, err
		}
	}
	return schema, nil
}

func pullServices(systemKey string, cli *cb.DevClient) ([]map[string]interface{}, error) {
	svcs, err := cli.GetServiceNames(systemKey)
	if err != nil {
		return nil, err
	}
	services := make([]map[string]interface{}, len(svcs))
	for i, svc := range svcs {
		if s, err := pullService(systemKey, svc, cli); err != nil {
			return nil, err
		} else {
			services[i] = s
			writeService(s["name"].(string), s)
		}
	}
	return services, nil
}

func pullLibraries(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	libs, err := cli.GetLibraries(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull libraries out of system %s: %s", sysMeta.Key, err.Error())
	}
	libraries := []map[string]interface{}{}
	for _, lib := range libs {
		thisLib := lib.(map[string]interface{})
		if thisLib["visibility"] == "global" {
			continue
		}
		libraries = append(libraries, thisLib)
		writeLibrary(thisLib["name"].(string), thisLib)
	}
	return libraries, nil
}

func pullTriggers(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
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
		writeTrigger(thisTrig["name"].(string), thisTrig)
	}
	return triggers, nil
}

func pullTimers(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theTimers, err := cli.GetTimers(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull timers out of system %s: %s", sysMeta.Key, err.Error())
	}
	timers := []map[string]interface{}{}
	for _, timer := range theTimers {
		thisTimer := timer.(map[string]interface{})
		// lotsa system and user dependent stuff to get rid of...
		delete(thisTimer, "system_key")
		delete(thisTimer, "system_secret")
		delete(thisTimer, "timer_key")
		delete(thisTimer, "user_id")
		delete(thisTimer, "user_token")
		timers = append(timers, thisTimer)
		writeTimer(thisTimer["name"].(string), thisTimer)
	}
	return timers, nil
}

func pullSystemMeta(systemKey string, cli *cb.DevClient) (*System_meta, error) {
	sys, err := cli.GetSystem(systemKey)
	if err != nil {
		return nil, err
	}
	serv_metas := make(map[string]Service_meta)
	sysMeta := &System_meta{
		Name:        sys.Name,
		Key:         sys.Key,
		Secret:      sys.Secret,
		Description: sys.Description,
		Services:    serv_metas,
		PlatformUrl: URL,
	}
	return sysMeta, nil
}

func getRolesForThing(name, key string) map[string]interface{} {
	rval := map[string]interface{}{}
	for _, roleInfo := range rolesInfo {
		roleName := roleInfo["Name"].(string)
		roleSvcs := roleInfo["Permissions"].(map[string]interface{})[key].([]interface{}) // Mouthful
		for _, roleEntIF := range roleSvcs {
			roleEnt := roleEntIF.(map[string]interface{})
			if roleEnt["Name"].(string) == name {
				rval[roleName] = roleEnt["Level"]
			}
		}
	}
	return rval
}

func getUserTablePermissions() map[string]interface{} {
	rval := map[string]interface{}{}
	for _, roleInfo := range rolesInfo {
		roleName := roleInfo["Name"].(string)
		roleUsers := roleInfo["Permissions"].(map[string]interface{})["UsersList"].(map[string]interface{})
		level := int(roleUsers["Level"].(float64))
		if level != 0 {
			rval[roleName] = level
		}
	}
	return rval
}

func cleanService(service map[string]interface{}) {
	service["source"] = service["name"].(string) + ".js"
	service["permissions"] = getRolesForThing(service["name"].(string), "CodeServices")
	delete(service, "code")
}

func cleanServices(services []map[string]interface{}) []map[string]interface{} {
	for _, service := range services {
		cleanService(service)
	}
	return services
}

func storeMeta(meta *System_meta) {
	systemDotJSON["platformURL"] = URL
	systemDotJSON["systemKey"] = meta.Key
	systemDotJSON["systemSecret"] = meta.Secret
	systemDotJSON["name"] = meta.Name
	systemDotJSON["description"] = meta.Description
	systemDotJSON["auth"] = true
}

func pullUsers(sysMeta *System_meta, cli *cb.DevClient, saveThem bool) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	if !exportUsers {
		return []map[string]interface{}{}, nil
	}
	allUsers, err := cli.GetAllUsers(sysKey)
	if err != nil {
		return nil, fmt.Errorf("Could not get all users: %s", err.Error())
	}
	for _, aUser := range allUsers {
		userId := aUser["user_id"].(string)
		roles, err := cli.GetUserRoles(sysKey, userId)
		if err != nil {
			return nil, fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
		}
		aUser["roles"] = roles
		if saveThem {
			writeUser(aUser["email"].(string), aUser)
		}
	}
	return allUsers, nil
}

func doExport(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) == 0 {
		fmt.Printf("export command: missing system key\n")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Printf("export command: too many arguments\n")
		os.Exit(1)
	}
	return export(client, args[0])
}

func export(cli *cb.DevClient, sysKey string) error {
	cb.CB_ADDR = URL

	fmt.Printf("Exporting System Info...")
	sysMeta, err := pullSystemMeta(sysKey, cli)
	if err != nil {
		return err
	}

	//dir := rootDir
	setRootDir(strings.Replace(sysMeta.Name, " ", "_", -1))
	if err := setupDirectoryStructure(sysMeta); err != nil {
		return err
	}
	storeMeta(sysMeta)
	fmt.Printf("Done.\nExporting Roles...")

	roles, err := pullRoles(sysKey, cli, true)
	if err != nil {
		return err
	}
	rolesInfo = roles
	storeRoles(rolesInfo)

	fmt.Printf("Done.\nExporting Services...")
	services, err := pullServices(sysKey, cli)
	if err != nil {
		return err
	}
	/*
		if err := storeServices(dir, services, sysMeta); err != nil {
			return err
		}
	*/
	//systemDotJSON["services"] = cleanServices(services)
	systemDotJSON["services"] = services

	fmt.Printf("Done.\nExporting Libraries...")
	libraries, err := pullLibraries(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["libraries"] = libraries
	/*
		if err := storeLibraries(); err != nil {
			return err
		}
	*/

	fmt.Printf("Done.\nExporting Triggers...")
	if triggers, err := pullTriggers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["triggers"] = triggers
	}

	fmt.Printf("Done.\nExporting Timers...")
	if timers, err := pullTimers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["timers"] = timers
	}

	fmt.Printf("Done.\nExporting Collections...")
	colls, err := pullCollections(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["data"] = colls

	fmt.Printf("Done.\nExporting Users...")
	_, err = pullUsers(sysMeta, cli, true)
	if err != nil {
		return fmt.Errorf("GetAllUsers FAILED: %s", err.Error())
	}

	userSchema, err := pullUserSchemaInfo(sysKey, cli, true)
	if err != nil {
		return err
	}
	systemDotJSON["users"] = userSchema
	fmt.Printf("Done.\n")

	if err = storeSystemDotJSON(systemDotJSON); err != nil {
		return err
	}

	metaStuff := map[string]interface{}{
		"platformURL":       URL,
		"developerEmail":    Email,
		"assetRefreshDates": []interface{}{},
	}
	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been exported into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}
