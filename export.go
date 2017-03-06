package cblib

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

var (
	exportUsers bool
	inARepo     bool
)

func init() {
	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	rolesInfo = []map[string]interface{}{}
	myExportCommand := &SubCommand{
		name:         "export",
		usage:        "Ain't no thing",
		needsAuth:    false,
		mustBeInRepo: false,
		run:          doExport,
		//  TODO -- add help, usage, etc.
	}
	myExportCommand.flags.StringVar(&URL, "url", "", "Clearblade platform url for target system")
	myExportCommand.flags.StringVar(&MsgURL, "messaging-url", "", "Clearblade messaging url for target system")
	myExportCommand.flags.StringVar(&SystemKey, "system-key", "", "System key for target system")
	myExportCommand.flags.StringVar(&Email, "email", "", "Developer email for login")
	myExportCommand.flags.StringVar(&DevToken, "dev-token", "", "Dev token for login")
	myExportCommand.flags.BoolVar(&ExportRows, "exportrows", false, "exports all data from all collections")
	myExportCommand.flags.BoolVar(&exportUsers, "exportusers", false, "exports user info")
	AddCommand("export", myExportCommand)
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
		// Checking if collection is CB collection or different
		// Exporting only CB collections
		_, ok := col.(map[string]interface{})["dbtype"]
		if ok {
			continue
		}
		if r, err := PullCollection(sysMeta, col.(map[string]interface{}), cli); err != nil {
			return nil, err
		} else {
			data := makeCollectionJsonConsistent(r)
			writeCollection(r["name"].(string), data)
			rval[i] = data
		}
	}
	return rval, nil
}

func PullCollection(sysMeta *System_meta, co map[string]interface{}, cli *cb.DevClient) (map[string]interface{}, error) {
	id := co["collectionID"].(string)
	isConnect := isConnectCollection(co)
	var columnsResp []interface{}
	var err error
	if isConnect {
		columnsResp = []interface{}{}
	} else {
		columnsResp, err = cli.GetColumns(id, sysMeta.Key, sysMeta.Secret)
		if err != nil {
			return nil, err
		}
	}

	//remove the item_id column if it is not supposed to be exported
	if !ExportItemId {
		//Loop through the array of maps and find the one where ColumnName = item_id
		//Remove it from the slice
		for ndx, columnMap := range columnsResp {
			if columnMap.(map[string]interface{})["ColumnName"] == "item_id" {
				columnsResp = append(columnsResp[:ndx], columnsResp[ndx+1:]...)
				break
			}
		}
	}

	co["schema"] = columnsResp
	if err = getRolesForCollection(co); err != nil {
		return nil, err
	}
	co["items"] = []interface{}{}
	if !isConnect && ExportRows {
		items, err := pullCollectionData(co, cli)
		if err != nil {
			return nil, err
		}
		co["items"] = items
	}
	return co, nil
}

func isConnectCollection(co map[string]interface{}) bool {
	if isConnect, ok := co["isConnect"]; ok {
		switch isConnect.(type) {
		case bool:
			return isConnect.(bool)
		case string:
			return isConnect.(string) == "true"
		default:
			return false
		}
	}
	return false
}

func pullCollectionAndInfo(sysMeta *System_meta, id string, cli *cb.DevClient) (map[string]interface{}, error) {

	colInfo, err := cli.GetCollectionInfo(id)
	if err != nil {
		return nil, err
	}
	return PullCollection(sysMeta, colInfo, cli)
}

func getRolesForCollection(collection map[string]interface{}) error {
	colName := collection["name"].(string)
	perms := map[string]interface{}{}
	for _, role := range rolesInfo {
		roleName := role["Name"].(string)

		if _, ok := role["Permissions"].(map[string]interface{}); !ok {
			continue
		}
		rolePerms := role["Permissions"].(map[string]interface{})

		if _, ok := rolePerms["Collections"].([]interface{}); !ok {
			continue
		}
		colPerms := rolePerms["Collections"].([]interface{})

		//colPerms := role["Permissions"].(map[string]interface{})["Collections"].([]interface{})
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

		//remove the item_id data if it is not supposed to be exported
		if !ExportItemId {
			//Loop through the array of maps and find the one where ColumnName = item_id
			//Remove it from the slice
			for _, rowMap := range curData {
				delete(rowMap.(map[string]interface{}), "item_id")
			}
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

func PullServices(systemKey string, cli *cb.DevClient) ([]map[string]interface{}, error) {
	svcs, err := cli.GetServiceNames(systemKey)
	if err != nil {
		return nil, err
	}
	services := make([]map[string]interface{}, len(svcs))
	for i, svc := range svcs {
		fmt.Printf(" %s", svc)
		if s, err := pullService(systemKey, svc, cli); err != nil {
			return nil, err
		} else {
			services[i] = s
			err = writeService(s["name"].(string), s)
			if err != nil {
				return nil, err
			}
		}
	}
	return services, nil
}

func PullLibraries(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
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
		fmt.Printf(" %s", thisLib["name"].(string))
		libraries = append(libraries, thisLib)
		err = writeLibrary(thisLib["name"].(string), thisLib)
		if err != nil {
			return nil, err
		}
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
		err = writeTrigger(thisTrig["name"].(string), thisTrig)
		if err != nil {
			return nil, err
		}
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
		err = writeTimer(thisTimer["name"].(string), thisTimer)
		if err != nil {
			return nil, err
		}
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
		perms := roleInfo["Permissions"].(map[string]interface{})
		svcPerms := perms[key]

		if roleSvcs, ok := svcPerms.([]interface{}); ok {
			for _, roleEntIF := range roleSvcs {
				roleEnt := roleEntIF.(map[string]interface{})
				if roleEnt["Name"].(string) == name {
					rval[roleName] = roleEnt["Level"]
				}
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
	systemDotJSON["platform_url"] = cb.CB_ADDR
	systemDotJSON["messaging_url"] = cb.CB_MSG_ADDR
	systemDotJSON["system_key"] = meta.Key
	systemDotJSON["system_secret"] = meta.Secret
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

func PullEdges(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allEdges, err := cli.GetEdges(sysKey)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allEdges))
	for i := 0; i < len(allEdges); i++ {
		currentEdge := allEdges[i].(map[string]interface{})
		fmt.Printf(" %s", currentEdge["name"].(string))
		delete(currentEdge, "edge_key")
		delete(currentEdge, "isConnected")
		delete(currentEdge, "novi_system_key")
		delete(currentEdge, "broker_auth_port")
		delete(currentEdge, "broker_port")
		delete(currentEdge, "broker_tls_port")
		delete(currentEdge, "broker_ws_auth_port")
		delete(currentEdge, "broker_ws_port")
		delete(currentEdge, "broker_wss_port")
		delete(currentEdge, "communication_style")
		delete(currentEdge, "first_talked")
		delete(currentEdge, "last_talked")
		delete(currentEdge, "local_addr")
		delete(currentEdge, "local_port")
		delete(currentEdge, "public_addr")
		delete(currentEdge, "public_port")
		err = writeEdge(currentEdge["name"].(string), currentEdge)
		if err != nil {
			return nil, err
		}
		list = append(list, currentEdge)
	}

	return list, nil
}

func pullEdgesSchema(systemKey string, cli *cb.DevClient, writeThem bool) (map[string]interface{}, error) {
	resp, err := cli.GetEdgeColumns(systemKey)
	if err != nil {
		return nil, err
	}
	columns := []map[string]interface{}{}
	sort.Strings(DefaultEdgeColumns)
	for _, colIF := range resp {
		col := colIF.(map[string]interface{})
		if i := sort.SearchStrings(DefaultEdgeColumns, col["ColumnName"].(string)); i < len(DefaultEdgeColumns) && DefaultEdgeColumns[i] != col["ColumnName"].(string) {
			columns = append(columns, col)
		}
	}
	schema := map[string]interface{}{
		"columns": columns,
	}
	if writeThem {
		if err := writeEdge("schema", schema); err != nil {
			return nil, err
		}
	}
	return schema, nil
}

func PullDevices(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allDevices, err := cli.GetDevices(sysKey)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allDevices))
	for i := 0; i < len(allDevices); i++ {
		currentDevice := allDevices[i].(map[string]interface{})
		fmt.Printf(" %s", currentDevice["name"].(string))
		delete(currentDevice, "device_key")
		delete(currentDevice, "system_key")
		delete(currentDevice, "last_active_date")
		delete(currentDevice, "__HostId__")
		delete(currentDevice, "created_date")
		delete(currentDevice, "last_active_date")
		err = writeDevice(currentDevice["name"].(string), currentDevice)
		if err != nil {
			return nil, err
		}
		list = append(list, currentDevice)
	}
	return list, nil
}

func pullEdgeSyncInfo(sysMeta *System_meta, cli *cb.DevClient) (map[string]interface{}, error) {
	sysKey := sysMeta.Key
	syncMap, err := cli.GetSyncResourcesForEdge(sysKey)
	if err != nil {
		return nil, err
	}
	return syncMap, nil
}

func PullPortals(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allPortals, err := cli.GetPortals(sysKey)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allPortals))
	for i := 0; i < len(allPortals); i++ {
		currentPortal := allPortals[i].(map[string]interface{})
		var err error
		if currentPortal["config"], err = parseIfNeeded(currentPortal["config"]); err != nil {
			return nil, err
		}
		fmt.Printf(" %s", currentPortal["name"].(string))
		err = writePortal(currentPortal["name"].(string), currentPortal)
		if err != nil {
			return nil, err
		}
		list = append(list, currentPortal)
	}
	return list, nil
}

func PullPlugins(sysMeta *System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allPlugins, err := cli.GetPlugins(sysKey)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allPlugins))
	for i := 0; i < len(allPlugins); i++ {
		currentPlugin := allPlugins[i].(map[string]interface{})
		fmt.Printf(" %s", currentPlugin["name"].(string))
		if err = writePlugin(currentPlugin["name"].(string), currentPlugin); err != nil {
			return nil, err
		}
		list = append(list, currentPlugin)
	}

	return list, nil
}

func doExport(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	if len(args) != 0 {
		return fmt.Errorf("export command takes no arguments; only options\n")
	}
	inARepo = MetaInfo != nil
	if inARepo {
		if exportOptionsExist() {
			return fmt.Errorf("When in a repo, you cannot have command line options")
		}
		/*
			if err := os.Chdir(".."); err != nil {
				return fmt.Errorf("Could not change to parent directory: %s", err.Error())
			}
		*/
		setupFromRepo()
	}
	if exportOptionsExist() {
		client = cb.NewDevClientWithToken(DevToken, Email)
	} else {
		client, _ = Authorize(nil) // This is a hack for now. Need to handle error returned by Authorize
	}

	SetRootDir(".")

	// This is a hack to check if token has expired and auth again
	// since we dont have an endpoint to determine this
	client, err := checkIfTokenHasExpired(client, SystemKey)
	if err != nil {
		return fmt.Errorf("Re-auth failed...", err)
	}
	return ExportSystem(client, SystemKey)
}

func exportOptionsExist() bool {
	return URL != "" || SystemKey != "" || Email != "" || DevToken != ""
}

func ExportSystem(cli *cb.DevClient, sysKey string) error {
	fmt.Printf("Exporting System Info...")
	var sysMeta *System_meta
	var err error
	if inARepo {
		sysMeta, err = getSysMeta()
		os.Chdir("..")
	} else {
		sysMeta, err = pullSystemMeta(sysKey, cli)
	}
	if err != nil {
		return err
	}

	//dir := rootDir
	SetRootDir(strings.Replace(sysMeta.Name, " ", "_", -1))
	if err := setupDirectoryStructure(sysMeta); err != nil {
		return err
	}
	storeMeta(sysMeta)
	fmt.Printf(" Done.\nExporting Roles...")

	roles, err := pullRoles(sysKey, cli, true)
	if err != nil {
		return err
	}
	rolesInfo = roles
	storeRoles(rolesInfo)

	fmt.Printf(" Done.\nExporting Services...")
	services, err := PullServices(sysKey, cli)
	if err != nil {
		return err
	}
	systemDotJSON["services"] = services

	fmt.Printf(" Done.\nExporting Libraries...")
	libraries, err := PullLibraries(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["libraries"] = libraries

	fmt.Printf(" Done.\nExporting Triggers...")
	if triggers, err := pullTriggers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["triggers"] = triggers
	}

	fmt.Printf(" Done.\nExporting Timers...")
	if timers, err := pullTimers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["timers"] = timers
	}

	fmt.Printf(" Done.\nExporting Collections...")
	colls, err := pullCollections(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["data"] = colls

	fmt.Printf(" Done.\nExporting Users...")
	_, err = pullUsers(sysMeta, cli, true)
	if err != nil {
		return fmt.Errorf("GetAllUsers FAILED: %s", err.Error())
	}

	userSchema, err := pullUserSchemaInfo(sysKey, cli, true)
	if err != nil {
		return err
	}
	systemDotJSON["users"] = userSchema

	fmt.Printf(" Done.\nExporting Edges...")
	edges, err := PullEdges(sysMeta, cli)
	if err != nil {
		return err
	}
	if _, err := pullEdgesSchema(sysKey, cli, true); err != nil {
		fmt.Printf("\nNo custom columns to pull and create schema.json from... Continuing...\n")
	}
	systemDotJSON["edges"] = edges

	fmt.Printf(" Done.\nExporting Devices...")
	devices, err := PullDevices(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["devices"] = devices

	fmt.Printf(" Done.\nExporting Edge Sync Information...")
	syncInfo, err := pullEdgeSyncInfo(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["edge_sync"] = syncInfo

	fmt.Printf(" Done.\nExporting Portals...")
	portals, err := PullPortals(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["portals"] = portals

	fmt.Printf(" Done.\nExporting Plugins...")
	plugins, err := PullPlugins(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["plugins"] = plugins

	fmt.Printf(" Done.\n")

	if err = storeSystemDotJSON(systemDotJSON); err != nil {
		return err
	}

	metaStuff := map[string]interface{}{
		"platform_url":        cb.CB_ADDR,
		"messaging_url":       cb.CB_MSG_ADDR,
		"developer_email":     Email,
		"asset_refresh_dates": []interface{}{},
		"token":               cli.DevToken,
	}
	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been exported into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}

func setupFromRepo() {
	var ok bool
	sysMeta, err := getSysMeta()
	if err != nil {
		fmt.Printf("Error getting sys meta: %s\n", err.Error())
		curDir, _ := os.Getwd()
		fmt.Printf("WORKING DIRECTORY: %s\n", curDir)
	}
	Email, ok = MetaInfo["developerEmail"].(string)
	if !ok {
		Email = MetaInfo["developer_email"].(string)
	}
	URL, ok = MetaInfo["platformURL"].(string)
	if !ok {
		URL = MetaInfo["platform_url"].(string)
	}
	DevToken = MetaInfo["token"].(string)
	SystemKey = sysMeta.Key
}

func parseIfNeeded(stuff interface{}) (map[string]interface{}, error) {
	switch stuff.(type) {
	case map[string]interface{}:
		return stuff.(map[string]interface{}), nil
	case string:
		parsed := map[string]interface{}{}
		if err := json.Unmarshal([]byte(stuff.(string)), &parsed); err != nil {
			return nil, err
		}
		return parsed, nil
	default:
		return nil, fmt.Errorf("Invalid type passed into parseIfNeeded. Must be string or map[string]interface{}")
	}
}
