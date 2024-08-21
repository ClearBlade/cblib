package cblib

import (
	"fmt"
	"regexp"
)

type systemPushOptions struct {
	AdaptorsRegex          string
	BucketSetsRegex        string
	BucketSetFilesRegex    string
	CachesRegex            string
	CollectionsRegex       string
	CollectionSchemasRegex string
	DeploymentsRegex       string
	DevicesRegex           string
	PushDeviceSchema       bool
	EdgesRegex             string
	PushEdgeSchema         bool
	ExternalDatabasesRegex string
	LibrariesRegex         string
	PluginsRegex           string
	PortalsRegex           string
	RolesRegex             string
	SecretsRegex           string
	ServicesRegex          string
	TimersRegex            string
	TriggersRegex          string
	UsersRegex             string
	PushUserSchema         bool
	WebhooksRegex          string
}

func NewPushOptions() *systemPushOptions {
	return &systemPushOptions{
		AdaptorsRegex:          getAdaptorsRegex(),
		BucketSetsRegex:        getBucketSetsRegex(),
		BucketSetFilesRegex:    getBucketSetFilesRegex(),
		CachesRegex:            getCachesRegex(),
		CollectionsRegex:       getCollectionsRegex(),
		CollectionSchemasRegex: getCollectionSchemasRegex(),
		DeploymentsRegex:       getDeploymentsRegex(),
		DevicesRegex:           getDevicesRegex(),
		PushDeviceSchema:       shouldPushDeviceSchema(),
		EdgesRegex:             getEdgesRegex(),
		PushEdgeSchema:         shouldPushEdgeSchema(),
		ExternalDatabasesRegex: getExternalDatabasesRegex(),
		LibrariesRegex:         getLibrariesRegex(),
		PluginsRegex:           getPluginsRegex(),
		PortalsRegex:           getPortalsRegex(),
		RolesRegex:             getRolesRegex(),
		SecretsRegex:           getSecretsRegex(),
		ServicesRegex:          getServicesRegex(),
		TimersRegex:            getTimersRegex(),
		TriggersRegex:          getTriggerRegex(),
		UsersRegex:             getUserRegex(),
		PushUserSchema:         shouldPushUserSchema(),
		WebhooksRegex:          getWebhooksRegex(),
	}
}

func getAdaptorsRegex() string {
	if AllAssets || AllAdaptors {
		return "*"
	}

	return regexp.QuoteMeta(AdaptorName)
}

func getBucketSetsRegex() string {
	if AllAssets || AllBucketSets {
		return "*"
	}

	return regexp.QuoteMeta(BucketSetName)
}

func getBucketSetFilesRegex() string {
	if AllAssets || AllBucketSetFiles {
		return "*"
	}

	if BucketSetFiles == "" || BucketSetBoxName == "" {
		return ""
	}

	if BucketSetFileName == "" {
		return regexp.QuoteMeta(fmt.Sprintf("%s/%s/*", BucketSetBoxName, BucketSetBoxName))
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s/%s/%s", BucketSetBoxName, BucketSetBoxName, BucketSetFileName))
}

func getCachesRegex() string {
	if AllAssets || AllServiceCaches {
		return "*"
	}

	return regexp.QuoteMeta(ServiceCacheName)
}

func getCollectionsRegex() string {
	if AllAssets || AllCollections {
		return "*"
	}

	// TODO: This OR logic is wrong!
	collectionsRegex := ""
	if CollectionName != "" {
		collectionsRegex = "(" + regexp.QuoteMeta(CollectionName) + ")"
	}

	if CollectionId != "" {
		name, err := getCollectionNameById(CollectionId)
		if err != nil {
			fmt.Printf("Could not determine name for collection with id %s: %s\n", CollectionId, err)
			return collectionsRegex
		}

		collectionsRegex += "|(" + regexp.QuoteMeta(name) + ")"
	}

	return regexp.QuoteMeta(collectionsRegex)
}

func getCollectionNameById(wantedId string) (string, error) {
	collections, err := getCollections()
	if err != nil {
		return "", err
	}
	for _, collection := range collections {
		id, ok := collection["collectionID"].(string)
		if !ok {
			continue
		}

		if id != wantedId {
			continue
		}

		name, ok := collection["name"].(string)
		if !ok {
			continue
		}

		return name, nil
	}
	return "", fmt.Errorf("collection with id %s not found", wantedId)
}

func getCollectionSchemasRegex() string {
	if AllAssets || AllCollections {
		return "*"
	}

	collectionSchemaRegex := ""
	if collection := getCollectionsRegex(); collection != "" {
		collectionSchemaRegex = ""
	}
	// TODO: Also push collection name

	return regexp.QuoteMeta(CollectionSchema)
}

func getDeploymentsRegex() string {
	if AllAssets || AllDeployments {
		return "*"
	}

	return regexp.QuoteMeta(DeploymentName)
}

func getDevicesRegex() string {
	if AllAssets || AllDevices {
		return "*"
	}

	return regexp.QuoteMeta(DeviceName)
}

func shouldPushDeviceSchema() bool {
	return DeviceSchema || getDevicesRegex() != ""
}

func getEdgesRegex() string {
	if AllAssets || AllEdges {
		return "*"
	}

	return regexp.QuoteMeta(EdgeName)
}

func shouldPushEdgeSchema() bool {
	return EdgeSchema || getEdgesRegex() != ""
}

func getExternalDatabasesRegex() string {
	if AllAssets || AllExternalDatabases {
		return "*"
	}

	return regexp.QuoteMeta(ExternalDatabaseName)
}

func getLibrariesRegex() string {
	if AllAssets || AllLibraries {
		return "*"
	}

	return regexp.QuoteMeta(LibraryName)
}

func getPluginsRegex() string {
	if AllAssets || AllPlugins {
		return "*"
	}

	return regexp.QuoteMeta(PluginName)
}

func getPortalsRegex() string {
	if AllAssets || AllPortals {
		return "*"
	}

	return regexp.QuoteMeta(PortalName)
}

func getRolesRegex() string {
	if AllAssets || AllRoles {
		return "*"
	}

	return regexp.QuoteMeta(RoleName)
}

func getSecretsRegex() string {
	if AllAssets || AllSecrets {
		return "*"
	}

	return regexp.QuoteMeta(SecretName)
}

func getServicesRegex() string {
	if AllAssets || AllServices {
		return "*"
	}

	return regexp.QuoteMeta(ServiceName)
}

func getTimersRegex() string {
	if AllAssets || AllTimers {
		return "*"
	}

	return regexp.QuoteMeta(TimerName)
}

func getTriggerRegex() string {
	if AllAssets || AllTriggers {
		return "*"
	}

	return regexp.QuoteMeta(TriggerName)
}

func getUserRegex() string {
	if AllAssets || AllUsers {
		return "*"
	}

	userRegex := ""
	if User != "" {
		userRegex = "(" + regexp.QuoteMeta(User) + ")"
	}

	if UserId != "" {
		email, err := getUserEmailByID(UserId)
		if err != nil {
			fmt.Printf("Could not determine email for user with id %s: %s\n", UserId, err)
			return userRegex
		}

		userRegex += "|(" + regexp.QuoteMeta(email) + ")"
	}

	return userRegex
}

func shouldPushUserSchema() bool {
	return UserSchema || getUserRegex() != ""
}

func getUserEmailById(wantedId string) (string, error) {
	users, err := getUsers()
	if err != nil {
		return "", err
	}

	for _, user := range users {
		id, ok := user["user_id"].(string)
		if !ok {
			continue
		}

		if id != wantedId {
			continue
		}

		email, ok := user["email"].(string)
		if !ok {
			continue
		}

		return email, nil
	}

	return "", fmt.Errorf("user with id %s not found", wantedId)
}

func getUserSchemaRegex() string {
	if AllAssets || AllUsers || UserSchema {
		return "*"
	}

	return ""
}

func getWebhooksRegex() string {
	if AllAssets || AllWebhooks {
		return "*"
	}

	return regexp.QuoteMeta(WebhookName)
}

// TODO: Message History Storage?
// TODO: Message Type Triggers
