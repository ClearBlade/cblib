package cblib

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	cb "github.com/clearblade/Go-SDK"

	"github.com/clearblade/cblib/models/collections"
	rt "github.com/clearblade/cblib/resourcetree"
	"github.com/clearblade/cblib/types"
)

var (
	inARepo bool
)

func init() {
	usage := `
	Export a System from a Platform to your local filesystem. By default, all assets are exported, except for Collection rows and Users.

	1) Exporting for first time - Run from any directory, will create a folder with same name as your system
	2) Exporting into an existing folder - 'cd' into the system's directory, and run 'cb-cli export' to export into that existing folder
	`

	example := `
	  cb-cli export                             # export default assets, omits db rows and users, Note: may prompt for remaining flags
	  cb-cli export -exportrows -exportusers    # export default asset, additionally rows and users, Note: may prompt for remaining flags
	  cb-cli export -url=https://platform.clearblade.com -messaging-url=platform.clearblade.com -system-key=9b9eea9c0bda8896a3dab5aeec9601 -email=MyDevEmail@dev.com   # Prompts for just password
	`

	systemDotJSON = map[string]interface{}{}
	svcCode = map[string]interface{}{}
	myExportCommand := &SubCommand{
		name:      "export",
		usage:     usage,
		needsAuth: false,
		run:       doExport,
		example:   example,
	}
	myExportCommand.flags.StringVar(&URL, "url", "https://platform.clearblade.com", "Clearblade Platform URL where system is hosted")
	myExportCommand.flags.StringVar(&MsgURL, "messaging-url", "platform.clearblade.com", "Clearblade messaging url for target system")
	myExportCommand.flags.StringVar(&SystemKey, "system-key", "", "System Key for target system, ex 9b9eea9c0bda8896a3dab5aeec9601")
	myExportCommand.flags.StringVar(&Email, "email", "", "Developer Email for login")
	myExportCommand.flags.StringVar(&DevToken, "dev-token", "", "Advanced: Developer Token for login")
	myExportCommand.flags.BoolVar(&CleanUp, "cleanup", false, "Clean up directories before export, recommended after having deleted assets on Platform")
	myExportCommand.flags.BoolVar(&ExportRows, "exportrows", false, "Exports all rows from all collections, Note: Large collections may take a long time")
	myExportCommand.flags.BoolVar(&ExportUsers, "exportusers", false, "exports user, Note: Passwords are not exported")
	myExportCommand.flags.BoolVar(&ExportItemId, "exportitemid", ExportItemIdDefault, "exports a collection rows' item_id column, Default: true")
	myExportCommand.flags.BoolVar(&SortCollections, "sort-collections", SortCollectionsDefault, "Sort collections version control ease, Note: exportitemid must be enabled")
	myExportCommand.flags.IntVar(&DataPageSize, "data-page-size", DataPageSizeDefault, "Number of rows in a collection to fetch at a time, Note: Large collections should increase up to 1000 rows")
	setBackoffFlags(myExportCommand.flags)
	AddCommand("export", myExportCommand)
}

func makeCollectionNameToIdMap(collections []map[string]interface{}) map[string]interface{} {
	rtn := make(map[string]interface{})
	for i := 0; i < len(collections); i++ {
		rtn[collections[i]["name"].(string)] = collections[i]["collection_id"]
	}
	return rtn
}

func makeRoleNameToIdMap(roles []map[string]interface{}) map[string]interface{} {
	rtn := make(map[string]interface{})
	for i := 0; i < len(roles); i++ {
		rtn[roles[i]["Name"].(string)] = roles[i]["ID"]
	}
	return rtn
}

func PullAndWriteCollections(sysMeta *types.System_meta, cli *cb.DevClient, saveThem, shouldExportRows, shouldExportItemID bool) ([]map[string]interface{}, error) {
	colls, err := cli.GetAllCollections(sysMeta.Key)
	if err != nil {
		return nil, err
	}
	rval := make([]map[string]interface{}, 0)
	for _, col := range colls {
		// Checking if collection is CB collection or different
		// Exporting only CB collections
		_, ok := col.(map[string]interface{})["dbtype"]
		if ok {
			continue
		}
		if r, err := PullCollection(sysMeta, cli, col.(map[string]interface{}), shouldExportRows, shouldExportItemID); err != nil {
			return nil, err
		} else {
			data := makeCollectionJsonConsistent(r)
			rval = append(rval, data)
			if saveThem {
				writeCollection(r["name"].(string), data)
			}
		}
	}
	return rval, nil
}

func pullAndWriteCollectionColumns(sysMeta *types.System_meta, cli *cb.DevClient, name string) ([]interface{}, error) {
	columnsResp, err := pullCollectionColumns(sysMeta, cli, name)
	if err != nil {
		return nil, err
	}

	err = updateCollectionSchema(name, columnsResp, cli, sysMeta)
	if err != nil {
		return nil, err
	}
	return columnsResp, nil
}

func pullCollectionColumns(sysMeta *types.System_meta, cli *cb.DevClient, name string) ([]interface{}, error) {
	return cli.GetColumnsByCollectionName(sysMeta.Key, name)
}

func pullAndWriteCollectionIndexes(sysMeta *types.System_meta, cli *cb.DevClient, name string) (*rt.Indexes, error) {

	indexes, err := pullCollectionIndexes(sysMeta, cli, name)
	if err != nil {
		return nil, err
	}

	err = updateCollectionIndexes(name, indexes, cli, sysMeta)
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func pullCollectionIndexes(sysMeta *types.System_meta, cli *cb.DevClient, name string) (*rt.Indexes, error) {

	data, err := cli.ListIndexes(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}

	indexes, err := rt.NewIndexesFromMap(data)
	if err != nil {
		return nil, err
	}

	return indexes, nil
}

func PullCollection(sysMeta *types.System_meta, cli *cb.DevClient, co map[string]interface{}, shouldExportRows, shouldExportItemId bool) (map[string]interface{}, error) {
	fmt.Printf(" %s", co["name"].(string))

	isConnect := collections.IsConnectCollection(co)

	var columnsResp []interface{}
	var err error
	var indexes *rt.Indexes

	if isConnect {

		columnsResp = []interface{}{}
		indexes = &rt.Indexes{}

	} else {

		columnsResp, err = pullCollectionColumns(sysMeta, cli, co["name"].(string))
		if err != nil {
			return nil, err
		}

		indexes, err = pullCollectionIndexes(sysMeta, cli, co["name"].(string))
		if err != nil {
			return nil, err
		}

	}

	//remove the item_id column if it is not supposed to be exported
	if !shouldExportItemId {
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
	co["indexes"] = indexes
	co["items"] = []interface{}{}

	if !isConnect && shouldExportRows {
		items, err := pullCollectionData(co, cli)
		if err != nil {
			return nil, err
		}
		co["items"] = items
	}

	return co, nil
}

func pullCollectionData(collection map[string]interface{}, client *cb.DevClient) ([]interface{}, error) {
	colId := collection["collectionID"].(string)
	totalItems, err := client.GetItemCount(colId)
	if err != nil {
		return nil, fmt.Errorf("GetItemCount Failed: %s", err.Error())
	}

	dataQuery := &cb.Query{}
	dataQuery.PageSize = DataPageSize

	//We have to add an orderby clause in order to ensure paging works. Without the orderby clause
	//The order returned for each page is not consistent and could therefore result in duplicate rows
	//
	//https://www.postgresql.org/docs/current/static/sql-select.html
	dataQuery.Order = []cb.Ordering{{OrderKey: "item_id", SortOrder: true}} // SortOrder: true means we are sorting item_id ascending
	allData := []interface{}{}
	itemIDs := make(map[string]interface{})
	totalDownloaded := 0

	if totalItems/DataPageSize > 1000 {
		fmt.Println("Large dataset detected. Recommend increasing page size. use flag: -data-page-size=1000 or -data-page-size=10000")
	}

	for j := 0; j < totalItems; j += DataPageSize {
		dataQuery.PageNumber = (j / DataPageSize) + 1

		data, err := retryRequest(func() (interface{}, error) {
			return client.GetData(colId, dataQuery)
		}, BackoffMaxRetries, BackoffInitialInterval, BackoffMaxInterval, BackoffRetryMultiplier)
		if err != nil {
			return nil, err
		}
		curData := data.(map[string]interface{})["DATA"].([]interface{})

		//Loop through the array of maps and store the value of the item_id column in
		//a map so that we can prevent adding duplicate rows
		//
		//Duplicate rows can occur when dealing with very large tables if rows are added
		//to the table while we are attempting to read pages of data. There currently is
		//no solution to remedy this.
		for _, rowMap := range curData {
			itemID := (rowMap.(map[string]interface{})["item_id"]).(string)

			if _, ok := itemIDs[itemID]; !ok {
				itemIDs[itemID] = ""

				//remove the item_id data if it is not supposed to be exported
				if !ExportItemId {
					delete(rowMap.(map[string]interface{}), "item_id")
				}
				allData = append(allData, rowMap)
				totalDownloaded++
			}
		}
		fmt.Printf("Downloaded: \tPage(s): %v / %v \tItem(s): %v / %v\n", dataQuery.PageNumber, (totalItems/DataPageSize)+1, totalDownloaded, totalItems)
	}
	return allData, nil
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

func PullLibraries(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
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
		// call the individual endpoint to retrieve the properly formatted code
		realLib, err := cli.GetLibrary(sysMeta.Key, thisLib["name"].(string))
		if err != nil {
			return nil, err
		}
		fmt.Printf(" %s", realLib["name"].(string))
		libraries = append(libraries, realLib)
		err = writeLibrary(realLib["name"].(string), realLib)
		if err != nil {
			return nil, err
		}
	}
	return libraries, nil
}

func pullAndWriteDeployment(sysMeta *types.System_meta, cli *cb.DevClient, name string) (map[string]interface{}, error) {
	deploymentDetails, err := cli.GetDeploymentByName(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}

	//sort the deployment "assets" array so that the deployment assets are always
	//returned in a predictable mannner
	sort.Slice(deploymentDetails["assets"], func(i, j int) bool {
		//deploymentDetails["assets"] = []interface{}
		return getAssetSortKey(deploymentDetails["assets"].([]interface{})[i].(map[string]interface{})["asset_class"].(string),
			deploymentDetails["assets"].([]interface{})[i].(map[string]interface{})["asset_id"].(string)) <
			getAssetSortKey(deploymentDetails["assets"].([]interface{})[j].(map[string]interface{})["asset_class"].(string),
				deploymentDetails["assets"].([]interface{})[j].(map[string]interface{})["asset_id"].(string))
	})

	if err = writeDeployment(deploymentDetails["name"].(string), deploymentDetails); err != nil {
		return nil, err
	}
	return deploymentDetails, nil
}

func getAssetSortKey(assetClass string, asset_id string) string {
	//
	//asset_id will be empty for collections
	//asset_id will contain the item_id for item_level sync with a collection
	//
	// From {"clearblade" repo}/registry/constants
	// Adaptors          = "adaptors"
	// BucketSets        = "bucketsets"
	// Devices           = "devices"
	// Users             = "users"
	// Services          = "services"
	// Libraries         = "libraries"
	// ServiceCaches     = "servicecaches"
	// Timers            = "timers"
	// Triggers          = "triggers"
	// Webhooks          = "webhooks"
	// Portals           = "portals"
	// Plugins           = "plugins"
	// RolesPerms        = "rolesperms"
	// UserSecrets       = "usersecrets"
	// Collections       = "collections" // This is just a placeholder. The actual asset class will be the collection name
	//
	switch assetClass {
	case "adaptors", "bucketsets", "devices", "users", "services", "libraries", "servicecaches", "timers", "triggers", "webhooks", "portals", "plugins", "usersecrets":
		return assetClass + asset_id
	default:
		return "collection" + assetClass + asset_id
	}
}

func pullDeployments(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theDeployments, err := cli.GetAllDeployments(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull deployments out of system %s: %s", sysMeta.Key, err)
	}
	deployments := []map[string]interface{}{}
	for _, deploymentIF := range theDeployments {

		deploymentSummary := deploymentIF.(map[string]interface{})
		deplName := deploymentSummary["name"].(string)
		fmt.Printf(" %s", deplName)
		deploymentDetails, err := pullAndWriteDeployment(sysMeta, cli, deplName)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, deploymentDetails)
	}
	return deployments, nil
}

func pullAndWriteServiceCache(sysMeta *types.System_meta, cli *cb.DevClient, name string) (map[string]interface{}, error) {
	cache, err := cli.GetServiceCacheMeta(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}
	if err = writeServiceCache(cache["name"].(string), cache); err != nil {
		return nil, err
	}
	return cache, nil
}

func pullServiceCaches(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theCaches, err := cli.GetAllServiceCacheMeta(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull shared caches out of system %s: %s", sysMeta.Key, err)
	}
	for _, cache := range theCaches {
		cacheName := cache["name"].(string)
		fmt.Printf(" %s", cacheName)
		err := writeServiceCache(cacheName, cache)
		if err != nil {
			return nil, err
		}
	}
	return theCaches, nil
}

func pullWebhooks(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theHooks, err := cli.GetAllWebhooks(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull webhooks out of system %s: %s", sysMeta.Key, err)
	}
	for _, hook := range theHooks {
		hookName := hook["name"].(string)
		fmt.Printf(" %s", hookName)
		err := writeWebhook(hookName, hook)
		if err != nil {
			return nil, err
		}
	}
	return theHooks, nil
}

func pullExternalDatabases(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	theExternalDatabases, err := cli.GetAllExternalDBConnections(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull external databases out of system %s: %s", sysMeta.Key, err)
	}
	rtn := make([]map[string]interface{}, 0)
	for _, db := range theExternalDatabases {
		dbName := db.(map[string]interface{})["name"].(string)
		fmt.Printf(" %s", dbName)
		fullDBMetadata, err := pullAndWriteExternalDatabase(sysMeta, cli, dbName)
		if err != nil {
			return nil, err
		}
		rtn = append(rtn, fullDBMetadata)
	}
	return rtn, nil
}

func pullAndWriteExternalDatabase(sysMeta *types.System_meta, cli *cb.DevClient, name string) (map[string]interface{}, error) {
	fullDBMetadata, err := cli.GetExternalDBConnection(sysMeta.Key, name)
	if err != nil {
		return nil, fmt.Errorf("Could not pull external database metadata for '%s': %s", name, err.Error())
	}
	if err := writeExternalDatabase(name, fullDBMetadata); err != nil {
		return nil, fmt.Errorf("Failed to write external database '%s' to file system: %s", name, err.Error())
	}
	return fullDBMetadata, nil
}

func pullBucketSets(sysMeta *types.System_meta, cli *cb.DevClient) ([]interface{}, error) {
	theBucketSets, err := cli.GetBucketSets(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull bucket sets out of system %s: %s", sysMeta.Key, err)
	}

	for _, bucketSet := range theBucketSets {
		bsMap := bucketSet.(map[string]interface{})
		bsName := bsMap["name"].(string)
		fmt.Printf(" %s", bsName)
		err := writeBucketSet(bsName, bsMap)
		if err != nil {
			return nil, err
		}
	}
	return theBucketSets, nil
}

func pullAndWriteBucketSet(sysMeta *types.System_meta, cli *cb.DevClient, name string) (map[string]interface{}, error) {
	bs, err := cli.GetBucketSet(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}
	if err = writeBucketSet(bs["name"].(string), bs); err != nil {
		return nil, err
	}
	return bs, nil
}

func pullSecrets(sysMeta *types.System_meta, cli *cb.DevClient) (map[string]interface{}, error) {
	theSecrets, err := cli.GetSecrets(sysMeta.Key)
	if err != nil {
		return nil, fmt.Errorf("Could not pull secrets out of system %s: %s", sysMeta.Key, err)
	}

	for secretName, secret := range theSecrets {
		err := writeSecret(secretName, map[string]interface{}{
			"name":   secretName,
			"secret": secret,
		})
		if err != nil {
			return nil, err
		}
	}
	return theSecrets, nil
}

func pullMessageHistoryStorage(sysMeta *types.System_meta, cli *cb.DevClient) error {
	storageEntries, err := cli.GetMessageHistoryStorage(sysMeta.Key)
	if err != nil {
		return fmt.Errorf("Could not pull message history storage out of system %s: %s", sysMeta.Key, err)
	}

	err = writeMessageHistoryStorage(storageEntries)
	if err != nil {
		return err
	}

	return nil
}

func pullMessageTypeTriggers(sysMeta *types.System_meta, cli *cb.DevClient) error {
	msgTypeTriggers, err := cli.GetMessageTypeTriggers(sysMeta.Key)
	if err != nil {
		return fmt.Errorf("Could not pull message type triggers out of system %s: %s", sysMeta.Key, err.Error())
	}

	err = writeMessageTypeTriggers(msgTypeTriggers)
	if err != nil {
		return err
	}

	return nil
}

func pullAndWriteSecret(sysMeta *types.System_meta, cli *cb.DevClient, name string) (interface{}, error) {
	sec, err := cli.GetSecret(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}

	if err = writeSecret(name, map[string]interface{}{"name": name, "secret": sec}); err != nil {
		return nil, err
	}
	return sec, nil
}

func pullAndWriteWebhook(sysMeta *types.System_meta, cli *cb.DevClient, name string) (map[string]interface{}, error) {
	hook, err := cli.GetWebhook(sysMeta.Key, name)
	if err != nil {
		return nil, err
	}
	if err = writeWebhook(hook["name"].(string), hook); err != nil {
		return nil, err
	}
	return hook, nil
}

func pullSystemMeta(systemKey string, cli *cb.DevClient) (*types.System_meta, error) {
	sys, err := cli.GetSystem(systemKey)
	if err != nil {
		return nil, err
	}

	serviceMetas := make(map[string]types.Service_meta)

	sysMeta := &types.System_meta{
		Name:        sys.Name,
		Key:         sys.Key,
		Secret:      sys.Secret,
		Description: sys.Description,
		Services:    serviceMetas,
		PlatformUrl: cli.HttpAddr,
		MessageUrl:  cli.MqttAddr,
	}

	return sysMeta, nil
}

func getUserTablePermissions(rolesInfo []map[string]interface{}) map[string]interface{} {
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

func pullAllEdges(systemKey string, cli *cb.DevClient) ([]interface{}, error) {
	return paginateRequests(systemKey, DataPageSize, cli.GetEdgesCountWithQuery, cli.GetEdgesWithQuery)
}

func PullEdges(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allEdges, err := pullAllEdges(sysKey, cli)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allEdges))
	for i := 0; i < len(allEdges); i++ {
		currentEdge := allEdges[i].(map[string]interface{})
		fmt.Printf(" %s", currentEdge["name"].(string))
		err = writeEdge(currentEdge["name"].(string), currentEdge)
		if err != nil {
			return nil, err
		}
		list = append(list, currentEdge)
	}

	return list, nil
}

func getUserDefinedColumns(columns []interface{}) []interface{} {
	var userDefinedColumns []interface{}
	for col := range columns {
		if columns[col].(map[string]interface{})["UserDefined"].(bool) {
			userDefinedColumns = append(userDefinedColumns, columns[col].(map[string]interface{}))
		}
	}
	return userDefinedColumns
}

func pullEdgesSchema(systemKey string, cli *cb.DevClient, writeThem bool) (map[string]interface{}, error) {
	resp, err := cli.GetEdgeColumns(systemKey)
	if err != nil {
		return nil, err
	}
	columns := getUserDefinedColumns(resp)
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

func pullDevicesSchema(systemKey string, cli *cb.DevClient, writeThem bool) (map[string]interface{}, error) {
	deviceCustomColumns, err := cli.GetDeviceColumns(systemKey)
	if err != nil {
		return nil, err
	}
	columns := getUserDefinedColumns(deviceCustomColumns)
	schema := map[string]interface{}{
		"columns": columns,
	}
	if writeThem {
		if err := writeDevice("schema", schema); err != nil {
			return nil, err
		}
	}
	return schema, nil
}

func pullAllDevices(systemKey string, cli *cb.DevClient) ([]interface{}, error) {
	return paginateRequests(systemKey, DataPageSize, cli.GetDevicesCount, cli.GetDevices)
}

func PullDevices(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allDevices, err := pullAllDevices(sysKey, cli)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allDevices))
	for i := 0; i < len(allDevices); i++ {
		currentDevice := allDevices[i].(map[string]interface{})
		name := currentDevice["name"].(string)
		fmt.Printf(" %s", name)
		roles, err := pullDeviceRoles(sysKey, name, cli)
		if err != nil {
			return nil, err
		}
		if err = writeDevice(name, currentDevice); err != nil {
			return nil, err
		}
		if err := writeDeviceRoles(name, roles); err != nil {
			return nil, err
		}
		list = append(list, currentDevice)
	}
	return list, nil
}

func pullDeviceRoles(sysKey, name string, cli *cb.DevClient) ([]string, error) {
	return cli.GetDeviceRoles(sysKey, name)
}

func pullEdgeDeployInfo(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	deployList, err := cli.GetDeployResourcesForSystem(sysKey)
	if err != nil {
		return nil, err
	}
	return deployList, nil
}

func PullPortals(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
	sysKey := sysMeta.Key
	allPortals, err := cli.GetPortals(sysKey)
	if err != nil {
		return nil, err
	}
	list := make([]map[string]interface{}, len(allPortals))
	for i := 0; i < len(allPortals); i++ {
		currentPortal := allPortals[i].(map[string]interface{})
		var err error
		if err := transformPortal(currentPortal); err != nil {
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

func PullPlugins(sysMeta *types.System_meta, cli *cb.DevClient) ([]map[string]interface{}, error) {
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

func PullAdaptors(sysMeta *types.System_meta, cli *cb.DevClient) error {
	sysKey := sysMeta.Key
	allAdaptors, err := cli.GetAdaptors(sysKey)
	if err != nil {
		return err
	}
	for i := 0; i < len(allAdaptors); i++ {
		currentAdaptorName := allAdaptors[i].(map[string]interface{})["name"].(string)
		currentAdaptor, err := pullAdaptor(sysKey, currentAdaptorName, cli)
		if err != nil {
			return err
		}

		if err = writeAdaptor(currentAdaptor); err != nil {
			return err
		}
	}

	return nil
}

func doExport(cmd *SubCommand, client *cb.DevClient, args ...string) error {
	parseBackoffFlags()
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

	client, err := Authorize(nil)
	if err != nil {
		return fmt.Errorf("Authorization failed: %s\n", err)
	}

	// This is a hack to check if token has expired and auth again
	// since we dont have an endpoint to determine this
	client, err = checkIfTokenHasExpired(client, SystemKey)
	if err != nil {
		return fmt.Errorf("Re-auth failed: %s", err)
	}

	return ExportSystem(client, SystemKey)
}

func exportOptionsExist() bool {
	return URL != "" || SystemKey != "" || Email != "" || DevToken != ""
}

func ExportSystem(cli *cb.DevClient, sysKey string) error {
	fmt.Printf("\nExporting System Info...\n")
	var sysMeta *types.System_meta
	var err error
	if inARepo {
		sysMeta, err = getSysMeta()
	} else {
		sysMeta, err = pullSystemMeta(sysKey, cli)
	}
	if err != nil {
		return err
	}
	// This was overwriting the rootdir set by cb_console
	// Only set if it has not already been set
	if !RootDirIsSet {
		SetRootDir(".")
	}

	if CleanUp {
		cleanUpDirectories(sysMeta)
	}

	if err := setupDirectoryStructure(); err != nil {
		return err
	}
	setGlobalSystemDotJSONFromSystemMeta(sysMeta)

	assetsToExport := createAffectedAssets()
	assetsToExport.AllAssets = true
	_, err = pullAssets(sysMeta, cli, assetsToExport)
	if err != nil {
		return err
	}

	if err = storeSystemDotJSON(systemDotJSON); err != nil {
		return err
	}

	// TODO: setting metaStuff using meta rather than globals
	// from here or the clearblade SDK. Might break.
	metaStuff := map[string]interface{}{
		// "platform_url":    cb.CB_ADDR,
		// "messaging_url":   cb.CB_MSG_ADDR,
		// "developer_email": Email,
		"platform_url":    sysMeta.PlatformUrl,
		"messaging_url":   sysMeta.MessageUrl,
		"developer_email": cli.Email,
		"token":           cli.DevToken,
	}

	if err = storeCBMeta(metaStuff); err != nil {
		return err
	}

	logInfo(fmt.Sprintf("System '%s' has been exported into the current directory\n", sysMeta.Name))
	return nil
}

func setupFromRepo() {
	var ok bool
	sysMeta, err := getSysMeta()
	if err != nil {
		fmt.Printf("Error getting sys meta: %s\n", err.Error())
		curDir, _ := os.Getwd()
		fmt.Printf("Current directory is %s\n", curDir)
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
