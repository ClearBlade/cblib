package cblib

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/clearblade/cblib/colutil"
	"github.com/clearblade/cblib/listutil"
	"github.com/clearblade/cblib/models"
	"github.com/clearblade/cblib/models/bucketSetFiles"
	"github.com/clearblade/cblib/models/collections"
	"github.com/clearblade/cblib/models/index"
	libPkg "github.com/clearblade/cblib/models/libraries"
	"github.com/clearblade/cblib/models/roles"
	"github.com/clearblade/cblib/models/systemUpload"
	"github.com/clearblade/cblib/types"
	"github.com/nsf/jsondiff"

	cb "github.com/clearblade/Go-SDK"
	cbfs "github.com/clearblade/cblib/fs"

	rt "github.com/clearblade/cblib/resourcetree"
)

func doLegacyPush(client *cb.DevClient, systemInfo *types.System_meta) error {
	didSomething := false
	if AllRoles || AllAssets {
		didSomething = true
		if err := pushRoles(systemInfo, client); err != nil {
			return err
		}
	}

	if RoleName != "" {
		didSomething = true
		if err := pushOneRole(systemInfo, client, RoleName); err != nil {
			return err
		}
	}

	if UserSchema || AllAssets {
		didSomething = true
		if err := pushUserSchema(systemInfo, client); err != nil {
			return err
		}
	}

	if AllUsers || AllAssets {
		didSomething = true
		if err := pushUsers(systemInfo, client); err != nil {
			return err
		}
	}

	if User != "" {
		didSomething = true
		if err := pushOneUser(systemInfo, client, User); err != nil {
			return err
		}
	}

	if UserId != "" {
		didSomething = true
		if err := pushOneUserById(systemInfo, client, UserId); err != nil {
			return err
		}
	}

	if AllServiceCaches || AllAssets {
		didSomething = true
		if err := pushAllServiceCaches(systemInfo, client); err != nil {
			return err
		}
	}

	if ServiceCacheName != "" {
		didSomething = true
		if err := pushOneServiceCache(systemInfo, client, ServiceCacheName); err != nil {
			return err
		}
	}

	if AllLibraries || AllServices || AllAssets {
		didSomething = true
		opts := cbfs.NewZipOptions(&mapper{})
		opts.AllServices = AllServices || AllAssets
		opts.AllLibraries = AllLibraries || AllAssets
		if err := pushCode(systemInfo, client, opts); err != nil {
			return err
		}
	}

	if LibraryName != "" {
		didSomething = true
		if err := pushOneLibrary(systemInfo, client, LibraryName); err != nil {
			return err
		}
	}

	if ServiceName != "" {
		didSomething = true
		if err := pushOneService(systemInfo, client, ServiceName); err != nil {
			return err
		}
	}

	if AllCollections || AllAssets {
		didSomething = true
		if err := pushAllCollections(systemInfo, client); err != nil {
			return err
		}
	}

	if CollectionSchema != "" {
		didSomething = true
		if err := pushOneCollectionSchema(systemInfo, client, CollectionSchema); err != nil {
			return err
		}
		if err := pushCollectionIndexes(systemInfo, client, CollectionSchema); err != nil {
			return err
		}
	}

	if CollectionName != "" {
		didSomething = true
		if err := pushOneCollection(systemInfo, client, CollectionName); err != nil {
			return err
		}
	}

	if CollectionId != "" {
		didSomething = true
		if err := pushOneCollectionById(systemInfo, client, CollectionId); err != nil {
			return err
		}
	}

	if AllTriggers || AllAssets {
		didSomething = true
		if err := pushTriggers(systemInfo, client); err != nil {
			return err
		}
	}

	if TriggerName != "" {
		didSomething = true
		if err := pushOneTrigger(systemInfo, client, TriggerName); err != nil {
			return err
		}
	}

	if AllTimers || AllAssets {
		didSomething = true
		if err := pushTimers(systemInfo, client); err != nil {
			return err
		}
	}

	if TimerName != "" {
		didSomething = true
		if err := pushOneTimer(systemInfo, client, TimerName); err != nil {
			return err
		}
	}

	if DeviceSchema || AllAssets {
		didSomething = true
		if err := pushDevicesSchema(systemInfo, client); err != nil {
			return err
		}
	}

	if AllDevices || AllAssets {
		didSomething = true
		if err := pushAllDevices(systemInfo, client); err != nil {
			return err
		}
	}

	if DeviceName != "" {
		didSomething = true
		if err := pushOneDevice(systemInfo, client, DeviceName); err != nil {
			return err
		}
	}

	if EdgeSchema || AllAssets {
		didSomething = true
		if err := pushEdgesSchema(systemInfo, client); err != nil {
			return err
		}
	}

	if AllEdges || AllAssets {
		didSomething = true
		if err := pushAllEdges(systemInfo, client); err != nil {
			return err
		}
	}

	if EdgeName != "" {
		didSomething = true
		if err := pushOneEdge(systemInfo, client, EdgeName); err != nil {
			return err
		}
	}

	if AllPortals || AllAssets {
		didSomething = true
		if err := pushAllPortals(systemInfo, client); err != nil {
			return err
		}
	}

	if PortalName != "" {
		didSomething = true
		if err := pushOnePortal(systemInfo, client, PortalName); err != nil {
			return err
		}
	}

	if AllPlugins || AllAssets {
		didSomething = true
		if err := pushAllPlugins(systemInfo, client); err != nil {
			return err
		}
	}

	if PluginName != "" {
		didSomething = true
		if err := pushOnePlugin(systemInfo, client, PluginName); err != nil {
			return err
		}
	}

	if AllAdaptors || AllAssets {
		didSomething = true
		if err := pushAllAdaptors(systemInfo, client); err != nil {
			return err
		}
	}

	if AdaptorName != "" {
		didSomething = true
		if err := pushOneAdaptor(systemInfo, client, AdaptorName); err != nil {
			return err
		}
	}

	if AllDeployments || AllAssets {
		didSomething = true
		if err := pushDeployments(systemInfo, client); err != nil {
			return err
		}
	}

	if DeploymentName != "" {
		didSomething = true
		if err := pushDeployment(systemInfo, client, DeploymentName); err != nil {
			return err
		}
	}

	if AllWebhooks || AllAssets {
		didSomething = true
		if err := pushAllWebhooks(systemInfo, client); err != nil {
			return err
		}
	}

	if WebhookName != "" {
		didSomething = true
		if err := pushOneWebhook(systemInfo, client, WebhookName); err != nil {
			return err
		}
	}

	if AllExternalDatabases || AllAssets {
		didSomething = true
		if err := pushAllExternalDatabases(systemInfo, client); err != nil {
			return err
		}
	}

	if ExternalDatabaseName != "" {
		didSomething = true
		if err := pushOneExternalDatabase(systemInfo, client, ExternalDatabaseName); err != nil {
			return err
		}
	}

	if AllBucketSets || AllAssets {
		didSomething = true
		if err := pushAllBucketSets(systemInfo, client); err != nil {
			return err
		}
	}

	if BucketSetName != "" {
		didSomething = true
		if err := pushOneBucketSet(systemInfo, client, BucketSetName); err != nil {
			return err
		}
	}

	if BucketSetFiles != "" {
		didSomething = true
		if BucketSetBoxName != "" && BucketSetFileName != "" {
			// push individual file within bucket set's box
			if err := bucketSetFiles.PushFile(systemInfo, client, BucketSetFiles, BucketSetBoxName, BucketSetFileName); err != nil {
				return err
			}
		} else {
			// push all files within bucket set
			if err := bucketSetFiles.PushFiles(systemInfo, client, BucketSetFiles, BucketSetBoxName); err != nil {
				return err
			}
		}
	}

	if AllBucketSetFiles || AllAssets {
		didSomething = true
		if err := bucketSetFiles.PushFilesForAllBucketSets(systemInfo, client); err != nil {
			return err
		}
	}

	if AllSecrets || AllAssets {
		didSomething = true
		if err := pushAllSecrets(systemInfo, client); err != nil {
			return err
		}
	}

	if SecretName != "" {
		didSomething = true
		if err := pushOneSecret(systemInfo, client, SecretName); err != nil {
			return err
		}
	}

	if MessageHistoryStorage || AllAssets {
		didSomething = true
		if err := pushMessageHistoryStorage(systemInfo, client); err != nil {
			return err
		}
	}

	if MessageTypeTriggers || AllAssets {
		didSomething = true
		if err := pushMessageTypeTriggers(systemInfo, client); err != nil {
			// ignoring error while pushing message type triggers since older platforms don't support this endpoint
			fmt.Printf("Ignoring error while pushing message type triggers. Assuming that you're pointed at an older platform. Error: %s\n", err.Error())
		}
	}

	if !didSomething {
		fmt.Printf("Nothing to push -- you must specify something to push (ie, -service=<svc_name>)\n")
	}

	return nil
}

func pushOneService(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing service %+s\n", name)
	service, err := getService(name)
	if err != nil {
		return err
	}
	return updateServiceWithRunAs(systemInfo.Key, name, service, client)
}

func pushUserSchema(systemInfo *types.System_meta, client *cb.DevClient) error {
	fmt.Printf("Pushing user schema\n")
	userschema, err := getUserSchema()
	if err != nil {
		return err
	}
	userColumns, err := client.GetUserColumns(systemInfo.Key)
	if err != nil {
		return fmt.Errorf("Error fetching user columns: %s", err.Error())
	}

	localSchema, ok := userschema["columns"].([]interface{})
	if !ok {
		return fmt.Errorf("Error in schema definition. Pls check the format of schema...\n")
	}

	diff := colutil.GetDiffForColumnsWithDynamicListOfDefaultColumns(convertInterfaceSlice[map[string]interface{}](localSchema), convertInterfaceSlice[map[string]interface{}](userColumns))
	for i := 0; i < len(diff.Removed); i++ {
		if err := client.DeleteUserColumn(systemInfo.Key, diff.Removed[i]["ColumnName"].(string)); err != nil {
			return fmt.Errorf("User schema could not be updated. Deletion of column(s) failed: %s", err)
		}
	}
	for i := 0; i < len(diff.Added); i++ {
		if err := client.CreateUserColumn(systemInfo.Key, diff.Added[i]["ColumnName"].(string), diff.Added[i]["ColumnType"].(string)); err != nil {
			return fmt.Errorf("Failed to create user column '%s': %s", diff.Added[i]["ColumnName"].(string), err.Error())
		}
	}
	return nil
}

func pushEdgesSchema(systemInfo *types.System_meta, client *cb.DevClient) error {
	fmt.Println("Pushing edge schema")
	edgeschema, err := getEdgesSchema()
	if err != nil {
		return err
	}
	allEdgeColumns, err := client.GetEdgeColumns(systemInfo.Key)
	if err != nil {
		return err
	}

	typedLocalSchema, ok := edgeschema["columns"].([]interface{})
	if !ok {
		return fmt.Errorf("Error in schema definition. Please verify the format of the schema.json. Value is: %+v - %+v\n", edgeschema["columns"], ok)
	}

	diff := colutil.GetDiffForColumnsWithDynamicListOfDefaultColumns(convertInterfaceSlice[map[string]interface{}](typedLocalSchema), convertInterfaceSlice[map[string]interface{}](allEdgeColumns))
	for i := 0; i < len(diff.Removed); i++ {
		if err := client.DeleteEdgeColumn(systemInfo.Key, diff.Removed[i]["ColumnName"].(string)); err != nil {
			return fmt.Errorf("Unable to delete column '%s': %s", diff.Removed[i]["ColumnName"].(string), err.Error())
		}
	}
	for i := 0; i < len(diff.Added); i++ {
		if err := client.CreateEdgeColumn(systemInfo.Key, diff.Added[i]["ColumnName"].(string), diff.Added[i]["ColumnType"].(string)); err != nil {
			return fmt.Errorf("Unable to create column '%s': %s", diff.Added[i]["ColumnName"].(string), err.Error())
		}
	}

	return nil

}

func pushDevicesSchema(systemInfo *types.System_meta, client *cb.DevClient) error {
	fmt.Println("Pushing device schema")
	deviceSchema, err := getDevicesSchema()
	if err != nil {
		return err
	}
	allDeviceColumns, err := client.GetDeviceColumns(systemInfo.Key)
	if err != nil {
		return err
	}
	localSchema, ok := deviceSchema["columns"].([]interface{})
	if !ok {
		return fmt.Errorf("Error in schema definition. Please verify the format of the schema.json\n")
	}

	diff := colutil.GetDiffForColumnsWithDynamicListOfDefaultColumns(convertInterfaceSlice[map[string]interface{}](localSchema), convertInterfaceSlice[map[string]interface{}](allDeviceColumns))
	for i := 0; i < len(diff.Removed); i++ {
		if err := client.DeleteDeviceColumn(systemInfo.Key, diff.Removed[i]["ColumnName"].(string)); err != nil {
			return fmt.Errorf("Unable to delete column '%s': %s", diff.Removed[i]["ColumnName"].(string), err.Error())
		}
	}
	for i := 0; i < len(diff.Added); i++ {
		if err := client.CreateDeviceColumn(systemInfo.Key, diff.Added[i]["ColumnName"].(string), diff.Added[i]["ColumnType"].(string)); err != nil {
			return fmt.Errorf("Unable to create column '%s': %s", diff.Added[i]["ColumnName"].(string), err.Error())
		}
	}

	return nil

}

func pushAllCollections(systemInfo *types.System_meta, client *cb.DevClient) error {
	allColls, err := getCollections()
	if err != nil {
		return err
	}
	for i := 0; i < len(allColls); i++ {
		err := pushOneCollection(systemInfo, client, allColls[i]["name"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func pushOneCollectionSchema(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing collection schema %s\n", name)
	collection, err := getCollection(name)
	if err != nil {
		fmt.Printf("error is %+v\n", err)
		return err
	}
	if out, err := createCollectionIfNecessary(systemInfo, collection, client, CreateCollectionIfNecessaryOptions{pullItems: false, pushItems: false}); err != nil {
		return err
	} else if !out.collectionExistsOrWasCreated {
		return nil
	}
	return pushCollectionSchema(systemInfo, collection, client)
}

func pushOneCollection(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing collection %s\n", name)
	collection, err := getCollection(name)
	if err != nil {
		fmt.Printf("error is %+v\n", err)
		return err
	}
	return updateCollection(systemInfo, collection, client)
}

func pushOneCollectionById(systemInfo *types.System_meta, client *cb.DevClient, wantedId string) error {
	fmt.Printf("Pushing collection with collectionID %s\n", wantedId)
	collections, err := getCollections()
	if err != nil {
		return err
	}
	for _, collection := range collections {
		id, ok := collection["collectionID"].(string)
		if !ok {
			continue
		}
		if id == wantedId {
			return updateCollection(systemInfo, collection, client)
		}
	}
	return fmt.Errorf("Collection with collectionID %+s not found.", wantedId)
}

func pushUsers(systemInfo *types.System_meta, client *cb.DevClient) error {
	users, err := getUsers()
	if err != nil {
		return err
	}
	for i := 0; i < len(users); i++ {
		// todo: make getUser accept user object so that it doesn't refetch from the FS
		if err := pushOneUser(systemInfo, client, users[i]["email"].(string)); err != nil {
			return err
		}
	}
	return nil
}

func pushOneUser(systemInfo *types.System_meta, client *cb.DevClient, email string) error {
	user, err := getFullUserObject(email)
	if err != nil {
		return err
	}
	return updateUser(systemInfo, user, client)
}

func pushOneUserById(systemInfo *types.System_meta, client *cb.DevClient, wantedId string) error {
	fmt.Printf("Pushing user with user_id %s\n", wantedId)
	users, err := getUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		id, ok := user["user_id"].(string)
		if !ok {
			continue
		}
		if id == wantedId {
			return updateUser(systemInfo, user, client)
		}
	}
	return fmt.Errorf("User with user_id %+s not found.", wantedId)
}

func pushRoles(systemInfo *types.System_meta, client *cb.DevClient) error {
	allRoles, err := getRoles()
	if err != nil {
		return err
	}
	for i := 0; i < len(allRoles); i++ {
		if err := pushOneRole(systemInfo, client, allRoles[i]["Name"].(string)); err != nil {
			return err
		}
	}
	return nil
}

func pushOneRole(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing role %s\n", name)
	role, err := getRole(name)
	if err != nil {
		return err
	}
	return updateRole(systemInfo, role, client)
}

func pushTriggers(systemInfo *types.System_meta, client *cb.DevClient) error {
	allTriggers, err := getTriggers()
	if err != nil {
		return err
	}
	for i := 0; i < len(allTriggers); i++ {
		if err := pushOneTrigger(systemInfo, client, allTriggers[i]["name"].(string)); err != nil {
			return err
		}
	}
	return nil
}

func pushOneTrigger(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing trigger %+s\n", name)
	trigger, err := getTrigger(name)
	if err != nil {
		return err
	}
	return updateTriggerWithUpdatedInfo(systemInfo.Key, trigger, client)
}

func pushTimers(systemInfo *types.System_meta, client *cb.DevClient) error {
	allTimers, err := getTimers()
	if err != nil {
		return err
	}
	for i := 0; i < len(allTimers); i++ {
		if err := pushOneTimer(systemInfo, client, allTimers[i]["name"].(string)); err != nil {
			return err
		}
	}
	return nil
}

func pushOneTimer(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing timer %+s\n", name)
	timer, err := getTimer(name)
	if err != nil {
		return err
	}
	return updateTimer(systemInfo.Key, timer, client)
}

func pushOneDevice(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing device %+s\n", name)
	device, err := getDevice(name)
	if err != nil {
		return err
	}
	return updateDevice(systemInfo.Key, device, client)
}

func pushAllDevices(systemInfo *types.System_meta, client *cb.DevClient) error {
	devices, err := getDevices()
	if err != nil {
		return err
	}
	for _, device := range devices {
		fmt.Printf("Pushing device %+s\n", device["name"].(string))
		if err := updateDevice(systemInfo.Key, device, client); err != nil {
			return fmt.Errorf("Error updating device '%s': %s\n", device["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOneEdge(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing edge %+s\n", name)
	edge, err := getEdge(name)
	if err != nil {
		return err
	}
	return updateEdge(systemInfo.Key, edge, client)
}

func pushAllEdges(systemInfo *types.System_meta, client *cb.DevClient) error {
	edges, err := getEdges()
	if err != nil {
		return err
	}
	for _, edge := range edges {
		edgeName := edge["name"].(string) // storing this here since updateEdge deletes "name" from the map
		fmt.Printf("Pushing edge %+s\n", edgeName)
		if err := updateEdge(systemInfo.Key, edge, client); err != nil {
			return fmt.Errorf("Error updating edge '%s': %s\n", edgeName, err.Error())
		}
	}
	return nil
}

func pushOnePortal(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing portal %+s\n", name)
	compressedPortal, err := compressPortal(name)
	if err != nil {
		return err
	}
	return updatePortal(systemInfo.Key, compressedPortal, client)
}

func pushAllPortals(systemInfo *types.System_meta, client *cb.DevClient) error {
	portals, err := getCompressedPortals()
	if err != nil {
		return err
	}
	for _, portal := range portals {
		name := portal["name"].(string)
		fmt.Printf("Pushing portal %+s\n", name)
		if err := updatePortal(systemInfo.Key, portal, client); err != nil {
			return fmt.Errorf("Error updating portal '%s': %s\n", name, err.Error())
		}
	}
	return nil
}

func pushOnePlugin(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing portal %+s\n", name)
	plugin, err := getPlugin(name)
	if err != nil {
		return err
	}
	return updatePlugin(systemInfo.Key, plugin, client)
}

func pushAllPlugins(systemInfo *types.System_meta, client *cb.DevClient) error {
	plugins, err := getPlugins()
	if err != nil {
		return err
	}
	for _, plugin := range plugins {
		fmt.Printf("Pushing plugin %+s\n", plugin["name"].(string))
		if err := updatePlugin(systemInfo.Key, plugin, client); err != nil {
			return fmt.Errorf("Error updating plugin '%s': %s\n", plugin["name"].(string), err.Error())
		}
	}
	return nil
}

func pushAllServiceCaches(systemInfo *types.System_meta, client *cb.DevClient) error {
	caches, err := getServiceCaches()
	if err != nil {
		return err
	}
	for _, cache := range caches {
		fmt.Printf("Pushing shared cache %+s\n", cache["name"].(string))
		if err := updateServiceCache(systemInfo.Key, cache, client); err != nil {
			return fmt.Errorf("Error updating shared cache '%s': %s\n", cache["name"].(string), err.Error())
		}
	}
	return nil
}

func pushOneServiceCache(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing shared cache %+s\n", name)
	cache, err := getServiceCache(name)
	if err != nil {
		return err
	}
	return updateServiceCache(systemInfo.Key, cache, client)
}

func pushAllWebhooks(systemInfo *types.System_meta, client *cb.DevClient) error {
	hooks, err := getWebhooks()
	if err != nil {
		return err
	}
	for _, hook := range hooks {
		fmt.Printf("Pushing webhook %+s\n", hook["name"].(string))
		if err := updateWebhook(systemInfo.Key, hook, client); err != nil {
			return fmt.Errorf("Error updating webhook '%+v': %s\n", hook, err.Error())
		}
	}
	return nil
}

func pushAllExternalDatabases(systemInfo *types.System_meta, client *cb.DevClient) error {
	extDBs, err := getExternalDatabases()
	if err != nil {
		return err
	}
	for _, extDB := range extDBs {
		fmt.Printf("Pushing external database %+s\n", extDB["name"].(string))
		if err := updateExternalDatabase(systemInfo.Key, extDB, client); err != nil {
			return fmt.Errorf("Error updating external database '%+v': %s\n", extDB, err.Error())
		}
	}
	return nil
}

func pushAllBucketSets(systemInfo *types.System_meta, client *cb.DevClient) error {
	bucketSets, err := getBucketSets()
	if err != nil {
		return err
	}
	for _, bucketSet := range bucketSets {
		fmt.Printf("Pushing bucket set %+s\n", bucketSet["name"].(string))
		if err := updateBucketSet(systemInfo.Key, bucketSet, client); err != nil {
			return fmt.Errorf("Error updating bucket set '%+v': %s\n", bucketSet, err.Error())
		}
	}
	return nil
}

func pushOneBucketSet(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing bucket set %+s\n", name)

	bucketSet, err := getBucketSet(name)
	if err != nil {
		return err
	}

	return updateBucketSet(systemInfo.Key, bucketSet, client)
}

func pushAllSecrets(systemInfo *types.System_meta, client *cb.DevClient) error {
	secrets, err := getSecrets()
	if err != nil {
		return err
	}
	for _, secret := range secrets {
		fmt.Printf("Pushing user secret %+s\n", secret)
		if err := updateSecret(systemInfo.Key, secret, client); err != nil {
			return fmt.Errorf("Error updating user secret '%+v': %s\n", secret, err.Error())
		}
	}
	return nil
}

func pushOneSecret(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing user secret %+s\n", name)

	secret, err := getSecret(name)
	if err != nil {
		return err
	}

	return updateSecret(systemInfo.Key, secret, client)
}

func pushOneWebhook(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing webhook %+s\n", name)
	hook, err := getWebhook(name)
	if err != nil {
		return err
	}
	return updateWebhook(systemInfo.Key, hook, client)
}

func pushOneExternalDatabase(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing external database %+s\n", name)
	db, err := getExternalDatabase(name)
	if err != nil {
		return err
	}
	return updateExternalDatabase(systemInfo.Key, db, client)
}

func updateServiceCache(systemKey string, cache map[string]interface{}, cli *cb.DevClient) error {
	cacheName := cache["name"].(string)

	_, err := cli.GetServiceCacheMeta(systemKey, cacheName)
	if err != nil {
		// shared cache DNE
		fmt.Printf("Could not find shared cache %s\n", cacheName)
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new shared cache named %s?", cacheName))
		if err != nil {
			return err
		} else {
			if c {
				if err := createServiceCache(systemKey, cache, cli); err != nil {
					return fmt.Errorf("Could not create shared cache %s: %s", cacheName, err.Error())
				} else {
					fmt.Printf("Successfully created new shared cache %s\n", cacheName)
				}
			} else {
				fmt.Printf("Shared cache will not be created.\n")
			}
		}
	} else {
		delete(cache, "name")
		return cli.UpdateServiceCacheMeta(systemKey, cacheName, cache)
	}

	return nil
}

func updateWebhook(systemKey string, hook map[string]interface{}, cli *cb.DevClient) error {
	hookName := hook["name"].(string)

	_, err := cli.GetWebhook(systemKey, hookName)
	if err != nil {
		// webhook DNE
		fmt.Printf("Could not find webhook %s\n", hookName)
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new webhook named %s?", hookName))
		if err != nil {
			return err
		} else {
			if c {
				if err := createWebhook(systemKey, hook, cli); err != nil {
					return fmt.Errorf("Could not create webhook %s: %s", hookName, err.Error())
				} else {
					fmt.Printf("Successfully created new webhook %s\n", hookName)
				}
			} else {
				fmt.Printf("Webhook will not be created.\n")
			}
		}
	} else {
		// not allowed to update these fields
		delete(hook, "name")
		delete(hook, "service_name")
		return cli.UpdateWebhook(systemKey, hookName, hook)
	}

	return nil
}

func updateExternalDatabase(systemKey string, obj map[string]interface{}, cli *cb.DevClient) error {
	name := obj["name"].(string)

	_, err := cli.GetExternalDBConnection(systemKey, name)
	if err != nil {
		// external database DNE
		fmt.Printf("Could not find external database %s\n", name)
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new external database named %s?", name))
		if err != nil {
			return err
		} else {
			if c {
				if err := createExternalDatabase(systemKey, obj, cli); err != nil {
					return fmt.Errorf("Could not create external database %s: %s", name, err.Error())
				} else {
					fmt.Printf("Successfully created new external database %s\n", name)
				}
			} else {
				fmt.Printf("External database will not be created.\n")
			}
		}
	} else {
		// not allowed to update these fields
		delete(obj, "name")
		delete(obj, "dbtype")
		return cli.UpdateExternalDBConnection(systemKey, name, obj)
	}

	return nil
}

func updateBucketSet(systemKey string, obj map[string]interface{}, cli *cb.DevClient) error {
	name := obj["name"].(string)

	bucketSet, err := cli.GetBucketSet(systemKey, name)
	if err != nil {
		// bucket set DNE
		fmt.Printf("Could not find bucket set %s\n", name)
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new bucket set named %s?", name))
		if err != nil {
			return err
		} else {
			if c {
				if err := createBucketSet(systemKey, obj, cli); err != nil {
					return fmt.Errorf("Could not create bucket set %s: %s", name, err.Error())
				} else {
					fmt.Printf("Successfully created new bucket set %s\n", name)
				}
			} else {
				fmt.Printf("Bucket set will not be created.\n")
			}
		}
	} else {
		ourBucket, err := json.Marshal(obj)
		if err != nil {
			return err
		}
		theirBucket, err := json.Marshal(whitelistBucketSet(bucketSet))
		if err != nil {
			return err
		}
		diff, _ := jsondiff.Compare(ourBucket, theirBucket, &jsondiff.Options{})
		if diff == jsondiff.FullMatch {
			fmt.Println("Bucket set hasn't changed, not updating")
			// no changes have been made, exit
			return nil
		}

		// since there is no UpdateBucketSet we must first delete and then create
		err = cli.DeleteBucketSet(systemKey, name)
		if err != nil {
			return err
		}

		return createBucketSet(systemKey, obj, cli)
	}

	return nil
}

func updateSecret(systemKey string, obj map[string]interface{}, cli *cb.DevClient) error {
	secret, err := cli.GetSecret(systemKey, obj["name"].(string))
	if err != nil {
		// secret DNE
		fmt.Printf("Could not find user secret %s\n", obj["name"].(string))
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new user secret named %s?", obj["name"].(string)))
		if err != nil {
			return err
		} else {
			if c {
				if err := createSecret(systemKey, obj, cli); err != nil {
					return fmt.Errorf("Could not create user secret %s: %s", obj["name"].(string), err.Error())
				} else {
					fmt.Printf("Successfully created new user secret %s\n", obj["name"].(string))
				}
			} else {
				fmt.Printf("User secret will not be created.\n")
			}
		}
	} else {
		if obj["secret"].(string) == secret {
			fmt.Println("User secret hasn't changed, not updating")
			// no changes have been made, exit
			return nil
		}

		//Update the secret since there are changes
		_, err = cli.UpdateSecret(systemKey, obj["name"].(string), obj["secret"].(string))
		if err != nil {
			return err
		}
	}

	return nil
}

func pushOneAdaptor(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing adaptor %+s\n", name)
	sysKey := systemInfo.Key
	adaptor, err := getAdaptor(sysKey, name, client)
	if err != nil {
		return err
	}
	return handleUpdateAdaptor(systemInfo.Key, adaptor, client)
}

func pushAllAdaptors(systemInfo *types.System_meta, client *cb.DevClient) error {
	sysKey := systemInfo.Key
	adaptors, err := getAdaptors(sysKey, client)
	if err != nil {
		return err
	}
	for i := 0; i < len(adaptors); i++ {
		currentAdaptor := adaptors[i]
		fmt.Printf("Pushing adaptor %+s\n", currentAdaptor.Name)
		if err := handleUpdateAdaptor(sysKey, currentAdaptor, client); err != nil {
			return fmt.Errorf("Error updating adaptor '%s': %s\n", currentAdaptor.Name, err.Error())
		}
	}
	return nil
}

func pushAllServices(systemInfo *types.System_meta, client *cb.DevClient) error {
	services, err := getServices()
	if err != nil {
		return err
	}
	for _, service := range services {
		name := service["name"].(string)
		fmt.Printf("Pushing service %+s\n", name)
		if err := updateServiceWithRunAs(systemInfo.Key, name, service, client); err != nil {
			return fmt.Errorf("Error updating service '%s': %s\n", service["name"].(string), err.Error())
		}
	}
	return nil
}

/**
 * Legacy behavior of push where each object is pushed individually.
 * This should be used on systems that do not support system upload.
 */
func pushCodeLegacy(systemInfo *types.System_meta, client *cb.DevClient, options *cbfs.ZipOptions) error {
	if options.AllAssets || options.AllLibraries {
		if err := pushAllLibraries(systemInfo, client); err != nil {
			return err
		}
	}

	if options.AllAssets || options.AllServices {
		if err := pushAllServices(systemInfo, client); err != nil {
			return err
		}
	}

	return nil
}

func pushOneLibrary(systemInfo *types.System_meta, client *cb.DevClient, name string) error {
	fmt.Printf("Pushing library %+s\n", name)

	library, err := getLibrary(name)
	if err != nil {
		return err
	}
	return updateLibrary(systemInfo.Key, library, client)
}

func pushMessageHistoryStorage(systemInfo *types.System_meta, client *cb.DevClient) error {
	storage, err := getMessageHistoryStorage()
	if err != nil {
		return err
	}

	err = client.UpdateMessageHistoryStorage(systemInfo.Key, storage)
	if err != nil {
		return err
	}
	return nil
}

func pushMessageTypeTriggers(systemInfo *types.System_meta, client *cb.DevClient) error {
	fmt.Println("Pushing message type triggers")
	msgTypeTriggers, err := getMessageTypeTriggers()
	if err != nil {
		return err
	}

	err = client.DeleteMessageTypeTriggers(systemInfo.Key)
	if err != nil {
		return err
	}

	if len(msgTypeTriggers) == 0 {
		// if there aren't any message type triggers, just return. the platform returns an error an empty array is POSTed
		return nil
	}

	err = client.AddMessageTypeTriggers(systemInfo.Key, msgTypeTriggers)
	if err != nil {
		return err
	}
	return nil
}

func pushAllLibraries(systemInfo *types.System_meta, client *cb.DevClient) error {
	rawLibraries, err := getLibraries()
	if err != nil {
		return err
	}

	libraries := make([]libPkg.Library, 0)
	for _, rawLib := range rawLibraries {
		libraries = append(libraries, libPkg.NewLibraryFromMap(rawLib))
	}

	orderedLibraries := libPkg.PostorderLibraries(libraries)

	for _, library := range orderedLibraries {
		fmt.Printf("Pushing library %+s\n", library.GetName())
		if err := updateLibrary(systemInfo.Key, library.GetMap(), client); err != nil {
			return fmt.Errorf("Error updating library '%s': %s\n", library.GetName(), err.Error())
		}
	}
	return nil
}

func pushCollectionSchema(systemInfo *types.System_meta, collection map[string]interface{}, cli *cb.DevClient) error {
	name, err := getCollectionName(collection)
	if err != nil {
		return err
	}

	collID, err := getCollectionIdByName(name, cli, systemInfo)
	if err != nil {
		return err
	}

	backendSchema, err := cli.GetColumnsByCollectionName(systemInfo.Key, name)
	if err != nil {
		return err
	}
	localSchema, ok := collection["schema"].([]interface{})
	if !ok {
		return fmt.Errorf("Error in schema definition. Please verify the format of the schema.json\n")
	}

	diff := colutil.GetDiffForColumnsWithStaticListOfDefaultColumns(convertInterfaceSlice[map[string]interface{}](localSchema), convertInterfaceSlice[map[string]interface{}](backendSchema), DefaultCollectionColumns)
	for i := 0; i < len(diff.Removed); i++ {
		if err := cli.DeleteColumn(collID, diff.Removed[i]["ColumnName"].(string)); err != nil {
			return fmt.Errorf("Unable to delete column '%s': %s", diff.Removed[i]["ColumnName"].(string), err.Error())
		}
	}
	for i := 0; i < len(diff.Added); i++ {
		if err := cli.AddColumn(collID, diff.Added[i]["ColumnName"].(string), diff.Added[i]["ColumnType"].(string)); err != nil {
			return fmt.Errorf("Unable to create column '%s': %s", diff.Added[i]["ColumnName"].(string), err.Error())
		}
	}

	// check if collection map has a hypertable_properties key
	hypertablePropertiesMap, ok := collection["hypertable_properties"].(map[string]interface{})
	if ok {
		localHypertableProperties, err := NewHypertablePropertiesFromMap(hypertablePropertiesMap)
		if err != nil {
			return err
		}

		allCollections, err := getAllCollectionsInfo(cli, systemInfo)
		if err != nil {
			return err
		}

		var collectionInfo CollectionInfo
		found := false
		for _, coll := range allCollections {
			if coll.Name == name {
				collectionInfo = coll
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("Unable to update hypertable info. Collection '%s' not found", name)
		}

		if collectionInfo.IsHypertable {
			if collectionInfo.HyperTableProperties.ChunkTimeInterval.IntervalString == localHypertableProperties.ChunkTimeInterval.IntervalString && collectionInfo.HyperTableProperties.DataRetentionPolicy.IntervalString == localHypertableProperties.DataRetentionPolicy.IntervalString {
				return nil
			}

			// update the hypertable properties
			err = cli.UpdateHypertableProperties(systemInfo.Key, name, map[string]interface{}{
				"chunk_time_interval": map[string]interface{}{
					"interval_string": localHypertableProperties.ChunkTimeInterval.IntervalString,
				},
				"data_retention_policy": map[string]interface{}{
					"interval_string": localHypertableProperties.DataRetentionPolicy.IntervalString,
				},
			})
			if err != nil {
				return err
			}
		} else {
			// create the hypertable
			err = cli.ConvertCollectionToHypertable(systemInfo.Key, name, map[string]interface{}{
				"migrate_data": true,
				"time_column":  localHypertableProperties.TimeColumn,
				"chunk_time_interval": map[string]interface{}{
					"interval_string": localHypertableProperties.ChunkTimeInterval.IntervalString,
				},
				"data_retention_policy": map[string]interface{}{
					"interval_string": localHypertableProperties.DataRetentionPolicy.IntervalString,
				},
			})
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func NewHypertablePropertiesFromMap(hypertablePropertiesMap map[string]interface{}) (HypertableProperties, error) {
	timeColumn, ok := hypertablePropertiesMap["time_column"].(string)
	if !ok {
		return HypertableProperties{}, fmt.Errorf("time_column is not a string")
	}
	chunkTimeIntervalMap, ok := hypertablePropertiesMap["chunk_time_interval"].(map[string]interface{})
	if !ok {
		return HypertableProperties{}, fmt.Errorf("chunk_time_interval is not a map")
	}
	chunkTimeIntervalString, ok := chunkTimeIntervalMap["interval_string"].(string)
	if !ok {
		return HypertableProperties{}, fmt.Errorf("chunk_time_interval.interval_string is not a string")
	}
	chunkTimeInterval := HypertableChunkTimeInterval{
		IntervalString: chunkTimeIntervalString,
	}
	dataRetentionPolicyMap, ok := hypertablePropertiesMap["data_retention_policy"].(map[string]interface{})
	if !ok {
		return HypertableProperties{}, fmt.Errorf("data_retention_policy is not a map")
	}
	dataRetentionPolicyString, ok := dataRetentionPolicyMap["interval_string"].(string)
	if !ok {
		return HypertableProperties{}, fmt.Errorf("data_retention_policy.interval_string is not a string")
	}
	dataRetentionPolicy := HypertableDataRetentionPolicy{
		IntervalString: dataRetentionPolicyString,
	}

	return HypertableProperties{
		TimeColumn:          timeColumn,
		ChunkTimeInterval:   chunkTimeInterval,
		DataRetentionPolicy: dataRetentionPolicy,
	}, nil
}

func pushCollectionIndexes(systemInfo *types.System_meta, cli *cb.DevClient, name string) error {

	fmt.Printf("Pushing collection indexes for '%s'\n", name)

	localColl, err := getCollection(name)
	if err != nil {
		return err
	}

	// NOTE: cast will default to nil if no indexes entry is available
	maybeIndexes, _ := localColl["indexes"].(map[string]interface{})
	localIndexes, err := rt.NewIndexesFromMap(maybeIndexes)
	if err != nil {
		return err
	}

	remoteIndexesData, err := cli.ListIndexes(systemInfo.Key, name)
	if err != nil {
		return err
	}

	remoteIndexes, err := rt.NewIndexesFromMap(remoteIndexesData)
	if err != nil {
		return err
	}

	diff := index.DiffIndexesFull(localIndexes.Data, remoteIndexes.Data)

	for _, index := range diff.Removed {
		err = doDropIndex(
			index,
			func() error { return cli.DropUniqueIndex(systemInfo.Key, name, index.Name) },
			func() error { return cli.DropIndex(systemInfo.Key, name, index.Name) },
		)
		if err != nil {
			return err
		}
	}

	for _, index := range diff.Added {
		err = doCreateIndex(
			index,
			func() error { return cli.CreateUniqueIndex(systemInfo.Key, name, index.Name) },
			func() error { return cli.CreateIndex(systemInfo.Key, name, index.Name) },
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func pushDeployments(systemInfo *types.System_meta, cli *cb.DevClient) error {
	deps, err := getDeployments()
	if err != nil {
		return err
	}
	for i := 0; i < len(deps); i++ {
		err := pushDeployment(systemInfo, cli, deps[i]["name"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func pushDeployment(systemInfo *types.System_meta, cli *cb.DevClient, name string) error {
	dep, err := getDeployment(name)
	if err != nil {
		return err
	}
	return updateDeployment(systemInfo, cli, name, dep, PreserveEdges)
}

func updateDeployment(systemInfo *types.System_meta, cli *cb.DevClient, name string, dep map[string]interface{}, preserveEdges bool) error {
	// fetch deployment
	backendDep, err := cli.GetDeploymentByName(systemInfo.Key, name)
	if err != nil {
		fmt.Printf("Could not find deployment '%s'. Error is - %s\n", name, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new deployment named %s?", name))
		if err != nil {
			return err
		} else {
			if c {
				if _, err := createDeployment(systemInfo.Key, dep, cli); err != nil {
					return fmt.Errorf("Could not create deployment %s: %s", name, err.Error())
				} else {
					fmt.Printf("Successfully created new deployment %s\n", name)
				}
			} else {
				fmt.Printf("Deployment will not be created.\n")
			}
		}
	} else {

		// diff backend deployment and local deployment
		theDiff := diffDeployments(dep, backendDep, preserveEdges)
		if _, err := cli.UpdateDeploymentByName(systemInfo.Key, name, theDiff); err != nil {
			return err
		}
	}

	return nil
}

func diffDeployments(localDep map[string]interface{}, backendDep map[string]interface{}, preserveEdges bool) map[string]interface{} {
	assetDiff := listutil.CompareLists(localDep["assets"].([]interface{}), backendDep["assets"].([]interface{}), isAssetMatch)
	edgeDiff := listutil.CompareLists(localDep["edges"].([]interface{}), backendDep["edges"].([]interface{}), isEdgeMatch)
	edgesToAdd := edgeDiff.Added
	edgesToRemove := edgeDiff.Removed
	if preserveEdges {
		edgesToAdd = []interface{}{}
		edgesToRemove = []interface{}{}
	}
	return map[string]interface{}{
		"assets": map[string]interface{}{
			"add":    assetDiff.Added,
			"remove": assetDiff.Removed,
		},
		"edges": map[string]interface{}{
			"add":    edgesToAdd,
			"remove": edgesToRemove,
		},
	}
}

func isEdgeMatch(edgeA interface{}, edgeB interface{}) bool {
	return edgeA.(string) == edgeB.(string)
}

func isAssetMatch(assetA interface{}, assetB interface{}) bool {
	typedA := assetA.(map[string]interface{})
	typedB := assetB.(map[string]interface{})
	return typedA["asset_class"].(string) == typedB["asset_class"].(string) && typedA["asset_id"].(string) == typedB["asset_id"].(string) && typedA["sync_to_edge"].(bool) == typedB["sync_to_edge"].(bool) && typedA["sync_to_platform"].(bool) == typedB["sync_to_platform"].(bool)
}

func createRole(systemInfo *types.System_meta, role map[string]interface{}, client *cb.DevClient) error {
	roleName := role["Name"].(string)
	var roleID string
	if roleName != "Authenticated" && roleName != "Anonymous" && roleName != "Administrator" {
		createIF, err := client.CreateRole(systemInfo.Key, role["Name"].(string))
		if err != nil {
			return err
		}
		createDict, ok := createIF.(map[string]interface{})
		if !ok {
			return fmt.Errorf("return value from CreateRole is not a map. It is %T", createIF)
		}
		roleID, ok = createDict["role_id"].(string)
		if !ok {
			return fmt.Errorf("Did not get role_id key back from successful CreateRole call")
		}
	} else {
		roleID = roleName // Administrator, Authorized, Anonymous
	}
	updateRoleBody, err := roles.PackageRoleForUpdate(roleID, role, networkCollectionFetcher{client: client, systemInfo: systemInfo})
	if err != nil {
		return err
	}
	if err := client.UpdateRole(systemInfo.Key, role["Name"].(string), updateRoleBody); err != nil {
		return err
	}
	if err := updateRoleNameToId(RoleInfo{ID: roleID, Name: roleName}); err != nil {
		logErrorForUpdatingMapFile(getRoleNameToIdFullFilePath(), err)
	}
	return nil
}

func lookupCollectionIdByName(theNameWeWant string, collectionsInfo []CollectionInfo) (string, bool) {
	for i := 0; i < len(collectionsInfo); i++ {
		if collectionsInfo[i].Name == theNameWeWant {
			return collectionsInfo[i].ID, true
		}
	}
	return "", false
}

type networkCollectionFetcher struct {
	client     *cb.DevClient
	systemInfo *types.System_meta
}

func (f networkCollectionFetcher) GetCollectionIdByName(theNameWeWant string) (string, error) {
	return getCollectionIdByName(theNameWeWant, f.client, f.systemInfo)
}

func getCollectionIdByName(theNameWeWant string, client *cb.DevClient, systemInfo *types.System_meta) (string, error) {
	collectionsInfo, err := getCollectionNameToIdAsSlice()
	if os.IsNotExist(err) {
		collectionsInfo = make([]CollectionInfo, 0)
	}
	maybeCollectionId, found := lookupCollectionIdByName(theNameWeWant, collectionsInfo)
	if found {
		return maybeCollectionId, nil
	}
	fmt.Printf("Couldn't find ID for collection name '%s'. Fetching IDs from platform...\n", theNameWeWant)
	collections, err := getAllCollectionsInfo(client, systemInfo)
	if err != nil {
		return "", err
	}
	for i := 0; i < len(collections); i++ {
		updateCollectionNameToId(collections[i])
	}
	maybeCollectionId, found = lookupCollectionIdByName(theNameWeWant, collections)
	if found {
		return maybeCollectionId, nil
	}

	return "", fmt.Errorf("Couldn't find ID for collection name '%s'\n", theNameWeWant)
}

func getAllCollectionsInfo(client *cb.DevClient, systemInfo *types.System_meta) ([]CollectionInfo, error) {
	collections, err := client.GetAllCollections(systemInfo.Key)
	if err != nil {
		return nil, err
	}
	var infoList []CollectionInfo
	for i := 0; i < len(collections); i++ {
		// check if the hypertable_properties key exists, if it does, cast it to a HyperTableProperties
		var hypertableProperties HypertableProperties
		isHypertable := false
		if hypertablePropertiesInterface, ok := collections[i].(map[string]interface{})["hypertable_properties"]; ok {
			hypertableProperties, err = NewHypertablePropertiesFromMap(hypertablePropertiesInterface.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			isHypertable = true
		}
		infoList = append(infoList, CollectionInfo{
			ID:                   collections[i].(map[string]interface{})["collectionID"].(string),
			Name:                 collections[i].(map[string]interface{})["name"].(string),
			HyperTableProperties: hypertableProperties,
			IsHypertable:         isHypertable,
		})
	}
	return infoList, nil
}

func updateUser(meta *types.System_meta, user map[string]interface{}, client *cb.DevClient) error {
	delete(user, "cb_token")
	if email, ok := user["email"].(string); !ok {
		return fmt.Errorf("Missing user email %+v", user)
	} else {
		_, err := client.GetUserInfo(meta.Key, email)
		if err != nil {
			fmt.Printf("Could not update user '%s'. Error is - %s\n", email, err.Error())
			c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new user with email %s?", email))
			if err != nil {
				return err
			} else {
				if c {
					id, err := createUser(meta.Key, meta.Secret, user, client)
					if err != nil {
						return fmt.Errorf("Could not create user %s: %s", email, err.Error())
					} else {
						// tack the new user id onto the user object so it can be used in subsequent requests
						user["user_id"] = id
						fmt.Printf("Successfully created new user %s\n", email)
					}
				} else {
					fmt.Printf("User will not be created.\n")
					return nil
				}
			}
		}
	}

	userRoles, err := getUserRoles(user["email"].(string))
	if err != nil {
		return err
	}

	userID := user["user_id"].(string)
	backendUserRoles, err := client.GetUserRoles(meta.Key, userID)
	if err != nil {
		return err
	}
	roleDiff := roles.DiffRoles(userRoles, backendUserRoles)
	user["roles"] = map[string]interface{}{
		"add":    roleDiff.Added,
		"delete": roleDiff.Removed,
	}

	delete(user, "user_id")
	return client.UpdateUser(meta.Key, userID, user)
}

func createUser(systemKey string, systemSecret string, user map[string]interface{}, client *cb.DevClient) (string, error) {
	email := user["email"].(string)
	password := randSeq(10)
	if pwd, ok := user["password"]; ok {
		password = pwd.(string)
	}
	newUser, err := client.RegisterNewUser(email, password, systemKey, systemSecret)
	if err != nil {
		return "", fmt.Errorf("Could not create user %s: %s", email, err.Error())
	}
	userId := newUser["user_id"].(string)
	if err := updateUserEmailToId(UserInfo{
		UserID: userId,
		Email:  email,
	}); err != nil {
		logErrorForUpdatingMapFile(getUserEmailToIdFullFilePath(), err)
	}
	userRoles, err := getUserRoles(email)
	if err != nil {
		// couldn't get user roles, let's see if they're on the user map (legacy format)
		if r, ok := user["roles"].([]interface{}); ok {
			userRoles = convertInterfaceSliceToStringSlice(r)
		} else {
			logWarning(fmt.Sprintf("Could not find roles for user with email '%s'. This user will be created with only the default 'Authenticated' role.", email))
			userRoles = []string{"Authenticated"}
		}
	}
	defaultRoles := []string{"Authenticated"}
	roleDiff := roles.DiffRoles(userRoles, defaultRoles)
	if len(roleDiff.Added) > 0 || len(roleDiff.Removed) > 0 {
		if err := client.UpdateUserRoles(systemKey, userId, roleDiff.Added, roleDiff.Removed); err != nil {
			return "", err
		}
	}
	return userId, nil
}

func createTrigger(sysKey string, trigger map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	triggerName := trigger["name"].(string)
	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = sysKey
	delete(trigger, "name")
	delete(trigger, "event_definition")
	stuff, err := client.CreateEventHandler(sysKey, triggerName, trigger)
	if err != nil {
		return nil, fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
	}
	return stuff, nil
}

func updateUserTriggerInfo(trigger map[string]interface{}) {
	if email, _, ok := isTriggerForSpecificUser(trigger); ok {
		if id, err := getUserIdByEmail(email); err == nil {
			replaceEmailWithUserIdInTriggerKeyValuePairs(trigger, []UserInfo{{Email: email, UserID: id}})
		}
	}
}

func updateTriggerWithUpdatedInfo(systemKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	updateUserTriggerInfo(trigger)
	return updateTrigger(systemKey, trigger, client)
}

func updateTrigger(systemKey string, trigger map[string]interface{}, client *cb.DevClient) error {
	triggerName := trigger["name"].(string)

	triggerDef := trigger["event_definition"].(map[string]interface{})
	trigger["def_module"] = triggerDef["def_module"]
	trigger["def_name"] = triggerDef["def_name"]
	trigger["system_key"] = systemKey
	delete(trigger, "event_definition")

	if _, err := pullTrigger(systemKey, triggerName, client); err != nil {
		fmt.Printf("Could not find trigger '%s'. Error is - %s\n", triggerName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new trigger named %s?", triggerName))
		if err != nil {
			return err
		} else {
			if c {
				if _, err := client.CreateEventHandler(systemKey, triggerName, trigger); err != nil {
					return fmt.Errorf("Could not create trigger %s: %s", triggerName, err.Error())
				} else {
					fmt.Printf("Successfully created new trigger %s\n", triggerName)
				}
			} else {
				fmt.Printf("Trigger will not be created.\n")
			}
		}
	} else {

		delete(trigger, "name")

		if _, err := client.UpdateEventHandler(systemKey, triggerName, trigger); err != nil {
			return err
		}
	}
	return nil
}

func createTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}
	if _, err := client.CreateTimer(systemKey, timerName, timer); err != nil {
		return nil, fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
	}
	return timer, nil
}

func updateTimer(systemKey string, timer map[string]interface{}, client *cb.DevClient) error {
	timerName := timer["name"].(string)
	delete(timer, "name")
	startTime := timer["start_time"].(string)
	if startTime == "Now" {
		timer["start_time"] = time.Now().Format(time.RFC3339)
	}

	if _, err := pullTimer(systemKey, timerName, client); err != nil {
		fmt.Printf("Could not find timer '%s'. Error is - %s\n", timerName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new timer named %s?", timerName))
		if err != nil {
			return err
		} else {
			if c {
				if _, err := client.CreateTimer(systemKey, timerName, timer); err != nil {
					return fmt.Errorf("Could not create timer %s: %s", timerName, err.Error())
				} else {
					fmt.Printf("Successfully created new timer %s\n", timerName)
				}
			} else {
				fmt.Printf("Timer will not be created.\n")
			}
		}
	} else {
		if _, err := client.UpdateTimer(systemKey, timerName, timer); err != nil {
			return err
		}
	}
	return nil
}

func createDeployment(systemKey string, deployment map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	deploymentName := deployment["name"].(string)
	//delete(deployment, "name")
	if _, err := client.CreateDeploymentByName(systemKey, deploymentName, deployment); err != nil {
		return nil, fmt.Errorf("Could not create deployment %s: %s", deploymentName, err.Error())
	}
	return deployment, nil
}

func createServiceCache(systemKey string, cache map[string]interface{}, client *cb.DevClient) error {
	cacheName := cache["name"].(string)
	if err := client.CreateServiceCacheMeta(systemKey, cacheName, cache); err != nil {
		return fmt.Errorf("Could not create cache %s: %s", cacheName, err.Error())
	}
	return nil
}

func createWebhook(systemKey string, hook map[string]interface{}, client *cb.DevClient) error {
	hookName := hook["name"].(string)
	if err := client.CreateWebhook(systemKey, hookName, hook); err != nil {
		return fmt.Errorf("Could not create webhook %s: %s", hookName, err.Error())
	}
	return nil
}

func createExternalDatabase(systemKey string, obj map[string]interface{}, client *cb.DevClient) error {
	name := obj["name"].(string)
	// add a new line before prompting for password
	password := getOneItem(fmt.Sprintf("Password for external database '%s'", name), true)
	obj["credentials"].(map[string]interface{})["password"] = password

	if err := client.AddExternalDBConnection(systemKey, obj); err != nil {
		return fmt.Errorf("Could not create external database %s: %s", name, err.Error())
	}
	return nil
}

func createBucketSet(systemKey string, bucketSet map[string]interface{}, client *cb.DevClient) error {
	bucketSetName := bucketSet["name"].(string)
	if _, err := client.CreateBucketSet(systemKey, bucketSetName, bucketSet); err != nil {
		return fmt.Errorf("Could not create bucket set %s: %s", bucketSetName, err.Error())
	}
	return nil
}

func createSecret(systemKey string, secret map[string]interface{}, client *cb.DevClient) error {
	if _, err := client.AddSecret(systemKey, secret["name"].(string), secret["secret"].(string)); err != nil {
		return fmt.Errorf("Could not add secret %s: %s", secret["name"].(string), err.Error())
	}
	return nil
}

func updateDevice(systemKey string, device map[string]interface{}, client *cb.DevClient) error {
	deviceName := device["name"].(string)
	delete(device, "last_active_date")
	delete(device, "created_date")
	delete(device, "device_key")
	delete(device, "system_key")

	if _, err := pullDevice(systemKey, deviceName, client); err != nil {
		fmt.Printf("Could not find device '%s'. Error is - %s\n", deviceName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new device named %s?", deviceName))
		if err != nil {
			return err
		} else {
			if c {
				if _, err := createDevice(systemKey, device, client); err != nil {
					return err
				}
			} else {
				fmt.Printf("Device will not be created.\n")
				return nil
			}
		}
	} else {
		delete(device, "name")
		if _, err := client.UpdateDevice(systemKey, deviceName, device); err != nil {
			return err
		}
	}
	deviceRoles, err := getDeviceRoles(deviceName)
	if err != nil {
		return err
	}
	backendDeviceRoles, err := pullDeviceRoles(systemKey, deviceName, client)
	if err != nil {
		return err
	}
	roleDiff := roles.DiffRoles(deviceRoles, backendDeviceRoles)
	return client.UpdateDeviceRoles(
		systemKey,
		deviceName,
		roleDiff.Added,
		roleDiff.Removed)
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

	originalColumns := make(map[string]interface{})
	customColumns := make(map[string]interface{})
	for columnName, value := range edge {
		switch strings.ToLower(columnName) {
		case "system_key", "system_secret", "token", "description", "location", "mac_address", "policy_name", "resolver_func", "sync_edge_tables", "last_seen_version":
			originalColumns[columnName] = value
			break
		default:
			customColumns[columnName] = value
			break
		}
	}

	_, err := client.GetEdge(systemKey, edgeName)
	if err != nil {
		// Edge does not exist
		fmt.Printf("Could not update edge '%s'. Error is - %s\n", edgeName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new edge named %s?", edgeName))
		if err != nil {
			return err
		} else {
			if c {
				if _, err := client.CreateEdge(systemKey, edgeName, originalColumns); err != nil {
					return fmt.Errorf("Could not create edge %s: %s", edgeName, err.Error())
				} else {
					fmt.Printf("Successfully created new edge %s\n", edgeName)
				}
				_, err = client.UpdateEdge(systemKey, edgeName, customColumns)
				if err != nil {
					return err
				} else {
					return nil
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
	} else {
		rawConfig, _ := json.Marshal(portal["config"])
		portal["config"] = string(rawConfig)
	}

	_, err := client.GetPortal(systemKey, portalName)
	if err != nil {
		// Portal DNE
		fmt.Printf("Could not update portal '%s'. Error is - %s\n", portalName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new portal named %s?", portalName))
		if err != nil {
			return err
		}

		if c {
			if _, err := client.CreatePortal(systemKey, portalName, portal); err != nil {
				return fmt.Errorf("Could not create portal %s: %s", portalName, err.Error())
			}
			fmt.Printf("Successfully created new portal %s\n", portalName)
		} else {
			fmt.Printf("Portal will not be created.\n")
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
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new plugin named %s?", pluginName))
		if err != nil {
			return err
		} else {
			if c {
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

func updateAdaptor(adaptor *models.Adaptor) error {
	return adaptor.UpdateAllInfo()
}

func handleUpdateAdaptor(systemKey string, adaptor *models.Adaptor, client *cb.DevClient) error {
	adaptorName := adaptor.Name

	_, err := client.GetAdaptor(systemKey, adaptorName)
	if err != nil {
		// adaptor DNE
		fmt.Printf("Could not update adapter '%s'. Error is - %s\n", adaptorName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new adapter named %s?", adaptorName))
		if err != nil {
			return err
		} else {
			if c {
				if err := createAdaptor(adaptor); err != nil {
					return fmt.Errorf("Could not create adapter %s: %s", adaptorName, err.Error())
				} else {
					fmt.Printf("Successfully created new adapter %s\n", adaptorName)
				}
			} else {
				fmt.Printf("Adapter will not be created.\n")
			}
		}
	} else {
		return updateAdaptor(adaptor)
	}

	return nil
}

func findService(serviceName string) (map[string]interface{}, error) {
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

func updateServiceWithRunAs(systemKey, name string, service map[string]interface{}, client *cb.DevClient) error {
	// if savedRunAs, ok := service[runUserKey].(string); ok {
	// 	if id, err := getUserIdByEmail(savedRunAs); err == nil {
	// 		service[runUserKey] = id
	// 	} else if savedRunAs != "" {
	// 		service[runUserKey] = ""
	// 		logWarning(fmt.Sprintf("Failed to retrieve run_user ID for email '%s'. Please check to make sure that the user exists and that there is a matching entry in .cb-cli/map-name-to-id/users.json. Empty value will be used for run_user", savedRunAs))
	// 	}
	// }

	return updateService(systemKey, name, service, client)
}

func updateService(systemKey, name string, service map[string]interface{}, client *cb.DevClient) error {
	if _, err := pullService(systemKey, name, client); err != nil {
		fmt.Printf("Could not find service '%s'. Error is - %s\n", name, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new service named %s?", name))
		if err != nil {
			return err
		} else {
			if c {
				if err := createService(systemKey, service, client); err != nil {
					return fmt.Errorf("Could not create service %s: %s", name, err.Error())
				} else {
					fmt.Printf("Successfully created new service %s\n", name)
				}
			} else {
				fmt.Printf("Service will not be created.\n")
			}
		}
	} else {

		svcCode := service["code"].(string)

		extra := getServiceBody(service)
		_, err := client.UpdateServiceWithBody(systemKey, name, svcCode, extra)
		if err != nil {
			return err
		}

	}
	return nil
}

func mkSvcParams(params []interface{}) []string {
	rval := []string{}
	for _, val := range params {
		rval = append(rval, val.(string))
	}
	return rval
}

func getServiceBody(service map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{
		"logging_enabled":   false,
		"execution_timeout": 60,
		"parameters":        make([]interface{}, 0),
		"auto_balance":      false,
		"auto_restart":      false,
		"concurrency":       0,
		"dependencies":      "",
		"run_user":          "",
		"log_ttl_minutes":   10080,
		"run_on_edge":       true,
		"run_on_platform":   true,
		"log_level":         "debug",
		"engine_type":       0,
	}
	if loggingEnabled, ok := service["logging_enabled"]; ok {
		ret["logging_enabled"] = loggingEnabled
	}
	if executionTimeout, ok := service["execution_timeout"].(float64); ok {
		ret["execution_timeout"] = executionTimeout
	}
	if parameters, ok := service["params"].([]interface{}); ok { // GET for a service returns 'params' but POST/PUT expect 'parameters'
		ret["parameters"] = mkSvcParams(parameters)
	}
	if dependencies, ok := service["dependencies"].(string); ok {
		ret["dependencies"] = dependencies
	}
	if autoBalance, ok := service["auto_balance"].(bool); ok {
		ret["auto_balance"] = autoBalance
	}
	if autoRestart, ok := service["auto_restart"].(bool); ok {
		ret["auto_restart"] = autoRestart
	}
	if concurrency, ok := service["concurrency"].(float64); ok {
		ret["concurrency"] = concurrency
	}
	if runUser, ok := service["run_user"].(string); ok {
		ret["run_user"] = runUser
	}
	if log_ttl_minutes, ok := service["log_ttl_minutes"].(float64); ok {
		ret["log_ttl_minutes"] = log_ttl_minutes
	}
	if run_on_edge, ok := service["run_on_edge"].(bool); ok {
		ret["run_on_edge"] = run_on_edge
	}
	if run_on_platform, ok := service["run_on_platform"].(bool); ok {
		ret["run_on_platform"] = run_on_platform
	}
	if log_level, ok := service["log_level"].(string); ok {
		ret["log_level"] = log_level
	}
	if engine_type, ok := service["engine_type"].(float64); ok {
		ret["engine_type"] = engine_type
	}
	return ret
}

func createService(systemKey string, service map[string]interface{}, client *cb.DevClient) error {
	svcName := service["name"].(string)
	if ServiceName != "" {
		svcName = ServiceName
	}
	svcCode := service["code"].(string)
	extra := getServiceBody(service)
	if err := client.NewServiceWithBody(systemKey, svcName, svcCode, extra); err != nil {
		return err
	}
	if enableLogs(service) {
		if err := client.EnableLogsForService(systemKey, svcName); err != nil {
			return err
		}
	}
	return nil
}

func updateLibrary(systemKey string, library map[string]interface{}, client *cb.DevClient) error {
	libName := library["name"].(string)

	if _, err := pullLibrary(systemKey, libName, client); err != nil {
		fmt.Printf("Could not find library '%s'. Error is - %s\n", libName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new library named %s?", libName))
		if err != nil {
			return err
		} else {
			if c {
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
	} else {

		delete(library, "name")
		delete(library, "version")
		if _, err := client.UpdateLibrary(systemKey, libName, library); err != nil {
			return err
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

func updateCollection(meta *types.System_meta, collection map[string]interface{}, client *cb.DevClient) error {
	out, err := createCollectionIfNecessary(meta, collection, client, CreateCollectionIfNecessaryOptions{pullItems: true, pushItems: true})
	if err != nil {
		return err
	} else if !out.collectionExistsOrWasCreated {
		return nil
	}

	// here's our workflow for updating a collection:
	// 1) diff and update the collection schema
	// 2) diff and update the collection indexes
	// 3) attempt to update all of our items
	// 4) if update fails, we assume the item doesn't exist, so we create the item
	if err := pushCollectionSchema(meta, collection, client); err != nil {
		return err
	}

	collection_name, err := getCollectionName(collection)
	if err != nil {
		return err
	}

	if err := pushCollectionIndexes(meta, client, collection_name); err != nil {
		return err
	}

	fmt.Printf("Pushing collection data for '%s'", collection_name)
	items := collection["items"].([]interface{})
	for _, row := range items {
		query := cb.NewQuery()
		query.EqualTo("item_id", row.(map[string]interface{})["item_id"])

		if row.(map[string]interface{})["item_id"] != nil {
			if resp, err := client.UpdateDataByName(meta.Key, collection_name, query, row.(map[string]interface{})); err != nil {
				fmt.Printf("Error updating item '%s'. Skipping. Error is - %s\n", row.(map[string]interface{})["item_id"], err.Error())
			} else if resp.Count == 0 {
				if _, err := client.CreateDataByName(meta.Key, collection_name, row.(map[string]interface{})); err != nil {
					return fmt.Errorf("Failed to create item. Error is - %s", err.Error())
				}
			}
		} else {
			if _, err := client.CreateDataByName(meta.Key, collection_name, row.(map[string]interface{})); err != nil {
				return fmt.Errorf("Failed to create item. Error is - %s", err.Error())
			}
		}

	}
	return nil
}

type CollectionInfo struct {
	ID                   string
	Name                 string
	HyperTableProperties HypertableProperties
	IsHypertable         bool
}

type HypertableProperties struct {
	TimeColumn          string                        `json:"time_column"`
	ChunkTimeInterval   HypertableChunkTimeInterval   `json:"chunk_time_interval"`
	DataRetentionPolicy HypertableDataRetentionPolicy `json:"data_retention_policy"`
}

type HypertableChunkTimeInterval struct {
	IntervalString string `json:"interval_string"`
}

type HypertableDataRetentionPolicy struct {
	IntervalString string `json:"interval_string"`
}

type RoleInfo struct {
	ID   string
	Name string
}

func CreateCollection(systemKey string, collection map[string]interface{}, pushItems bool, client *cb.DevClient) (CollectionInfo, error) {
	collectionName := collection["name"].(string)
	isConnect := collections.IsConnectCollection(collection)
	var colId string
	var err error
	if isConnect {
		col, err := cb.GenerateConnectCollection(collection)
		if err != nil {
			return CollectionInfo{}, err
		}
		colId, err = client.NewConnectCollection(systemKey, col)
		if err != nil {
			return CollectionInfo{}, err
		}
	} else {
		colId, err = client.NewCollection(systemKey, collectionName)
		if err != nil {
			return CollectionInfo{}, err
		}
	}

	myInfo := CollectionInfo{
		ID:   colId,
		Name: collectionName,
	}
	if isConnect {
		return myInfo, nil
	}

	if err := updateCollectionNameToId(myInfo); err != nil {
		logErrorForUpdatingMapFile(getCollectionNameToIdFullFilePath(), err)
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
			return CollectionInfo{}, err
		}
	}

	indexes, ok := collection["indexes"].(map[string]interface{})
	if ok {
		indexInfo, err := rt.NewIndexesFromMap(indexes)
		if err != nil {
			return CollectionInfo{}, err
		}
		for _, index := range indexInfo.Data {
			err := doCreateIndex(
				index,
				func() error { return client.CreateUniqueIndex(systemKey, collectionName, index.Name) },
				func() error { return client.CreateIndex(systemKey, collectionName, index.Name) },
			)
			if err != nil {
				return CollectionInfo{}, err
			}
		}
	}

	if !pushItems {
		return myInfo, nil
	}

	allItems := collection["items"].([]interface{})
	totalItems := len(allItems)
	totalPushed := 0

	if totalItems == 0 {
		return myInfo, nil
	}
	if totalItems/DataPageSize > 1000 {
		fmt.Println("Large dataset detected. Recommend increasing page size. Use flag: -data-page-size=1000")
	}

	for i := 0; i < totalItems; i += DataPageSize {
		pageNumber := (i / DataPageSize) + 1
		beginningOfRange := i

		// this will be equal to max index + 1
		// to account for golang #slice conventions
		endOfRange := i + DataPageSize

		// if this is last page, and items on this page are fewer than page size
		if totalItems < endOfRange {
			endOfRange = totalItems
		}

		itemsInThisPage := allItems[beginningOfRange:endOfRange]

		for i, item := range itemsInThisPage {
			itemsInThisPage[i] = item.(map[string]interface{})
		}

		if _, err := retryRequest(func() (interface{}, error) { return client.CreateData(colId, itemsInThisPage) }, BackoffMaxRetries, BackoffInitialInterval, BackoffMaxInterval, BackoffRetryMultiplier); err != nil {
			return CollectionInfo{}, err
		}
		totalPushed += len(itemsInThisPage)
		fmt.Printf("Pushed: \tPage(s): %v / %v \tItem(s): %v / %v\n", pageNumber, (totalItems/DataPageSize)+1, totalPushed, totalItems)
	}
	return myInfo, nil
}

func createEdge(systemKey, name string, edge map[string]interface{}, client *cb.DevClient) error {
	originalColumns := make(map[string]interface{})
	customColumns := make(map[string]interface{})
	for columnName, value := range edge {
		switch strings.ToLower(columnName) {
		case "system_key", "system_secret", "token", "description", "location", "mac_address", "policy_name", "resolver_func", "sync_edge_tables", "last_seen_version":
			originalColumns[columnName] = value
		default:
			if value != nil {
				customColumns[columnName] = value
			}
		}
	}
	_, err := client.CreateEdge(systemKey, name, originalColumns)
	if err != nil {
		return err
	}
	if len(customColumns) == 0 {
		return nil
	}

	//  We only do this if there ARE custom columns to create
	_, err = client.UpdateEdge(systemKey, name, customColumns)
	if err != nil {
		return err
	}
	return nil
}

func createDevice(systemKey string, device map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	var randomActiveKey string
	activeKey, ok := device["active_key"].(string)
	if !ok {
		// Active key not present in json file. Creating a random one
		fmt.Printf(" Active key not present. Creating a random one for device creation. Please update the active key from the ClearBlade Console after creation\n")
		randomActiveKey = randSeq(8)
		device["active_key"] = randomActiveKey
	} else {
		if activeKey == "" || len(activeKey) < 6 {
			fmt.Printf("Active is either an empty string or less than 6 characters. Creating a random one for device creation. Please update the active key from the ClearBlade Console after creation\n")
			randomActiveKey = randSeq(8)
			device["active_key"] = randomActiveKey
		}
	}

	originalColumns := make(map[string]interface{})
	customColumns := make(map[string]interface{})
	for columnName, value := range device {
		switch strings.ToLower(columnName) {
		case "name", "type", "state", "description", "enabled", "allow_key_auth", "keys", "active_key", "allow_certificate_auth", "certificate":
			originalColumns[columnName] = value
			break
		case "cb_token", "cb_ttl_override":
			break
		default:
			customColumns[columnName] = value
			break
		}
	}
	deviceStuff, err := client.CreateDevice(systemKey, device["name"].(string), originalColumns)
	if err != nil {
		fmt.Printf("CREATE DEVICE ERROR: %s\n", err)
		return nil, err
	}
	_, err = client.UpdateDevice(systemKey, device["name"].(string), customColumns)
	if err != nil {
		fmt.Printf("UPDATE DEVICE ERROR: %s\n", err)
		return nil, err
	}
	return deviceStuff, nil
}

func createPortal(systemKey string, port map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	delete(port, "system_key")
	if port["description"] == nil {
		port["description"] = ""
	}
	if port["last_updated"] == nil {
		port["last_updated"] = ""
	}
	// Export stores config as dict, but import wants it as a string
	config, ok := port["config"]
	if ok {
		configStr := ""
		switch config.(type) {
		case string:
			configStr = config.(string)
		default:
			configBytes, err := json.Marshal(config)
			if err != nil {
				return nil, err
			}
			configStr = string(configBytes)
		}
		port["config"] = configStr
	}
	portalStuff, err := client.CreatePortal(systemKey, port["name"].(string), port)
	if err != nil {
		return nil, err
	}
	return portalStuff, nil
}

func createPlugin(systemKey string, plug map[string]interface{}, client *cb.DevClient) (map[string]interface{}, error) {
	return client.CreatePlugin(systemKey, plug)
}

func createAdaptor(adap *models.Adaptor) error {
	return adap.UploadAllInfo()
}

func updateRole(systemInfo *types.System_meta, role map[string]interface{}, client *cb.DevClient) error {
	roleName := role["Name"].(string)

	if _, err := pullRole(systemInfo.Key, roleName, client); err != nil {
		fmt.Printf("Could not find role '%s'. Error is - %s\n", roleName, err.Error())
		c, err := confirmPrompt(fmt.Sprintf("Would you like to create a new role named %s?", roleName))
		if err != nil {
			return err
		} else {
			if c {
				if err := createRole(systemInfo, role, client); err != nil {
					return fmt.Errorf("Could not create role %s: %s", roleName, err.Error())
				} else {
					fmt.Printf("Successfully created new role %s\n", roleName)
				}
			} else {
				fmt.Printf("Role will not be created.\n")
			}
		}
	} else {
		roleID, err := getRoleIdByName(roleName)
		if err != nil {
			return fmt.Errorf("Error updating role: %s", err.Error())
		}
		updateRoleBody, err := roles.PackageRoleForUpdate(roleID, role, networkCollectionFetcher{client: client, systemInfo: systemInfo})
		if err != nil {
			return err
		}
		if err := client.UpdateRole(systemInfo.Key, roleName, updateRoleBody); err != nil {
			if byts, err := json.Marshal(updateRoleBody); err == nil {
				fmt.Printf("Failed to update role '%s'. Request body is - \n%s\n", roleName, string(byts))
			}
			return err
		}
	}
	return nil
}

func pushCode(systemInfo *types.System_meta, client *cb.DevClient, options *cbfs.ZipOptions) error {
	if systemUpload.DoesBackendSupportSystemUploadForCode(systemInfo, client) {
		return pushSystemZip(systemInfo, client, options)
	} else {
		return pushCodeLegacy(systemInfo, client, options)
	}
}
