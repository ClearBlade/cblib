package cblib

import (
	//"crypto/md5"
	"encoding/json"
	"fmt"
	cb "github.com/clearblade/Go-SDK"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	systemDotJSON = map[string]interface{}{}
	libCode = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	rolesInfo = []map[string]interface{}{}
}

func pullRoles(systemKey string, cli *cb.DevClient) ([]map[string]interface{}, error) {
	r, err := cli.GetAllRoles(systemKey)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, len(r))
	for idx, rIF := range r {
		rval[idx] = rIF.(map[string]interface{})
	}
	return rval, nil
}
func storeRoles(roles []map[string]interface{}) {
	rolesInfo = roles
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
		co := col.(map[string]interface{})
		id := co["collectionID"].(string)

		columnsResp, err := cli.GetColumns(id, sysMeta.Key, sysMeta.Secret)
		if err != nil {
			return nil, err
		}
		co["schema"] = columnsResp
		rval[i] = co
	}

	return rval, nil
}

/*

func pullCollections(sysMeta *System_meta, cli *cb.DevClient) ([]Collection_meta, error) {
	colls, err := cli.GetAllCollections(sysMeta.Key)
	if err != nil {
		return nil, err
	}
	collections := make([]Collection_meta, len(colls))
	for i, col := range colls {
		co := col.(map[string]interface{})
		id, ok := co["collectionID"].(string)
		if !ok {
			return nil, fmt.Errorf("collectionID is not a string")
		}
		columnsResp, err := cli.GetColumns(id, sysMeta.Key, sysMeta.Secret)

		if err != nil {
			return nil, err
		}
		columns := make([]Column, len(columnsResp))
		for j := 0; j < len(columnsResp); j++ {
			columns[j] = Column{
				ColumnName: columnsResp[j].(map[string]interface{})["ColumnName"].(string),
				ColumnType: columnsResp[j].(map[string]interface{})["ColumnType"].(string),
			}
		}

		collections[i] = Collection_meta{
			Name:          co["name"].(string),
			Collection_id: co["collectionID"].(string),
			Columns:       columns,
		}
	}

	return collections, nil
}

*/

func storeCollections(dir string, collections []Collection_meta) error {

	datadir := dir + "/data"
	if err := os.MkdirAll(datadir, 0777); err != nil {
		return err
	}
	meta_bytes, err := json.Marshal(collections)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(datadir+"/collections.json", meta_bytes, 0777); err != nil {
		return err
	}

	return nil
}

func pullUserColumns(systemKey string, cli *cb.DevClient) ([]map[string]interface{}, error) {
	resp, err := cli.GetUserColumns(systemKey)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, len(resp))
	for idx, colIF := range resp {
		col := colIF.(map[string]interface{})
		if col["ColumnName"] == "email" || col["ColumnName"] == "creation_date" {
			continue
		}
		rval[idx] = col
	}
	return rval, nil
}

func storeUserColumns(dir string, meta User_meta) error {

	userdir := dir + "/user"
	if err := os.MkdirAll(userdir, 0777); err != nil {
		return err
	}
	meta_bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(userdir+"/columns.json", meta_bytes, 0777); err != nil {
		return err
	}

	return nil
}

func pullServices(systemKey string, cli *cb.DevClient) ([]map[string]interface{}, error) {
	svcs, err := cli.GetServiceNames(systemKey)
	if err != nil {
		return nil, err
	}
	services := make([]map[string]interface{}, len(svcs))
	for i, svc := range svcs {
		service, err := cli.GetServiceRaw(systemKey, svc)
		if err != nil {
			return nil, err
		}
		service["code"] = strings.Replace(service["code"].(string), "\\n", "\n", -1)
		services[i] = service
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
		libCode[thisLib["name"].(string)] = thisLib["code"].(string)
		delete(thisLib, "code")
		delete(thisLib, "library_key")
		delete(thisLib, "system_key")
		libraries = append(libraries, thisLib)
	}
	return libraries, nil
}

func storeLibraries() error {
	for name, code := range libCode {
		fileName := libDir + "/" + name + ".js"
		if err := ioutil.WriteFile(fileName, []byte(code.(string)), 0666); err != nil {
			return fmt.Errorf("Could not store library %s: %s", fileName, err.Error())
		}
	}
	return nil
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
	}
	return timers, nil
}

func pullSystemMeta(systemKey string, cli *cb.DevClient) (*System_meta, error) {
	sys, err := cli.GetSystem(systemKey)
	if err != nil {
		return nil, err
	}
	serv_metas := make(map[string]Service_meta)
	/*
		svcs, err := pullServices(systemKey, cli)
		if err != nil {
			return nil, err
		}
		for _, svc := range svcs {
			serv_metas[svc.Name] = Service_meta{
				Name:    svc.Name,
				Version: svc.Version,
				Hash:    fmt.Sprintf("%x", md5.Sum([]byte(svc.Code))),
				Params:  svc.Params,
			}
		}
	*/
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
			if roleEnt["Name"].(string) != name {
				rval[roleName] = roleEnt["Level"]
			}
		}
	}
	return rval
}

func cleanServices(services []map[string]interface{}) []map[string]interface{} {
	for _, service := range services {
		service["source"] = service["name"].(string) + ".js"
		service["permissions"] = getRolesForThing(service["name"].(string), "CodeServices")
		delete(service, "code")
	}
	return services
}

func storeServices(dir string, services []map[string]interface{}, meta *System_meta) error {
	for _, service := range services {
		if err := ioutil.WriteFile(svcDir+"/"+service["name"].(string)+".js", []byte(service["code"].(string)), 0666); err != nil {
			return err
		}
	}
	return nil
}

func storeMeta(meta *System_meta) {
	systemDotJSON["platformURL"] = URL
	systemDotJSON["name"] = meta.Name
	systemDotJSON["description"] = meta.Description
	systemDotJSON["auth"] = true
}

func Export_cmd(sysKey, devToken string) error {
	fmt.Printf("Initializing...")
	cb.CB_ADDR = URL
	cli, err := auth(devToken)
	if err != nil {
		return err
	}
	fmt.Printf("Done.\nGetting System Info...")

	sysMeta, err := pullSystemMeta(sysKey, cli)
	if err != nil {
		return err
	}

	dir := rootDir
	if err := setupDirectoryStructure(sysMeta); err != nil {
		return err
	}
	storeMeta(sysMeta)
	fmt.Printf("Done.\nGetting Roles...")

	roles, err := pullRoles(sysKey, cli)
	if err != nil {
		return err
	}
	storeRoles(roles)

	fmt.Printf("Done.\nGetting Services...")
	services, err := pullServices(sysKey, cli)
	if err != nil {
		return err
	}
	if err := storeServices(dir, services, sysMeta); err != nil {
		return err
	}
	systemDotJSON["services"] = cleanServices(services)

	fmt.Printf("Done.\nGetting Libraries...")
	libraries, err := pullLibraries(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["libraries"] = libraries
	if err := storeLibraries(); err != nil {
		return err
	}

	fmt.Printf("Done.\nGetting Triggers...")
	if triggers, err := pullTriggers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["triggers"] = triggers
	}

	fmt.Printf("Done.\nGetting Timers...")
	if timers, err := pullTimers(sysMeta, cli); err != nil {
		return err
	} else {
		systemDotJSON["timers"] = timers
	}

	fmt.Printf("Done.\nGetting Collections...")
	colls, err := pullCollections(sysMeta, cli)
	if err != nil {
		return err
	}
	systemDotJSON["data"] = colls

	fmt.Printf("Done.\nGetting Users...")
	allUsers, err := cli.GetAllUsers(sysKey)
	if err != nil {
		return fmt.Errorf("GetAllUsers FAILED: %s", err.Error())
	}
	for _, aUser := range allUsers {
		userId := aUser["user_id"].(string)
		roles, err := cli.GetUserRoles(sysKey, userId)
		if err != nil {
			return fmt.Errorf("Could not get roles for %s: %s", userId, err.Error())
		}
		aUser["roles"] = roles
	}
	writeUsersFile(allUsers)

	userColumns, err := pullUserColumns(sysKey, cli)
	if err != nil {
		return err
	}
	systemDotJSON["users"] = userColumns
	fmt.Printf("Done.\n")

	if err = storeSystemDotJSON(); err != nil {
		return err
	}

	fmt.Printf("System '%s' has been exported into directory %s\n", sysMeta.Name, strings.Replace(sysMeta.Name, " ", "_", -1))
	return nil
}
