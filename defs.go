package cblib

import (
	"github.com/clearblade/cblib/internal/types"
)

//
//  These are variables that can be used as
//  flags to a main package using this library, or
//  they can be set directly by unit tests, etc.
//  My, there are a lot of these...
//

const (
	NotExistErrorString    = "Does not exist"
	SpecialNoCBMetaError   = "No cbmeta file"
	ExportItemIdDefault    = true
	SortCollectionsDefault = false
	DataPageSizeDefault    = 100
)

var DefaultCollectionColumns = []string{"item_id"}

var (
	URL                        string
	MsgURL                     string
	SystemKey                  string
	DevToken                   string
	ShouldImportCollectionRows bool
	ImportRows                 bool
	ExportRows                 bool
	ExportItemId               bool
	ImportUsers                bool
	ExportUsers                bool
	CleanUp                    bool
	EdgeSchema                 bool
	DeviceSchema               bool
	UserSchema                 bool
	DataPageSize               int
	MaxRetries                 int
	Email                      string
	Password                   string
	CollectionSchema           string
	ServiceName                string
	LibraryName                string
	CollectionName             string
	CollectionId               string
	SortCollections            bool
	User                       string
	UserId                     string
	RoleName                   string
	TriggerName                string
	TimerName                  string
	DeviceName                 string
	EdgeName                   string
	PortalName                 string
	PluginName                 string
	AdaptorName                string
	DeploymentName             string
	ServiceCacheName           string
	WebhookName                string
	ExternalDatabaseName       string
	BucketSetName              string
	BucketSetFiles             string
	BucketSetBoxName           string
	BucketSetFileName          string
	SecretName                 string
	Message                    bool
	Topic                      string
	Payload                    string
	Help                       bool
	Params                     string
	Push                       bool
	AllServices                bool
	AllLibraries               bool
	AllDevices                 bool
	AllEdges                   bool
	AllPortals                 bool
	AllPlugins                 bool
	AllAdaptors                bool
	AllDeployments             bool
	AllCollections             bool
	AllRoles                   bool
	AllUsers                   bool
	AllAssets                  bool
	AllTriggers                bool
	AllTimers                  bool
	AllServiceCaches           bool
	AllWebhooks                bool
	AllExternalDatabases       bool
	AllBucketSets              bool
	AllBucketSetFiles          bool
	AllSecrets                 bool
	MessageHistoryStorage      bool
	MessageTypeTriggers        bool
	AutoApprove                bool
	TempDir                    string
	SkipUpdateMapNameToIdFiles bool
)

var (
	systemDotJSON map[string]interface{}
	libCode       map[string]interface{}
	svcCode       map[string]interface{}
	MetaInfo      map[string]interface{}
)

type AffectedAssets struct {
	AllAssets             bool
	AllServices           bool
	AllLibraries          bool
	AllEdges              bool
	AllDevices            bool
	AllPortals            bool
	AllPlugins            bool
	AllAdaptors           bool
	AllDeployments        bool
	AllCollections        bool
	AllRoles              bool
	AllUsers              bool
	AllTriggers           bool
	AllTimers             bool
	AllServiceCaches      bool
	AllWebhooks           bool
	AllExternalDatabases  bool
	AllBucketSets         bool
	AllSecrets            bool
	MessageTypeTriggers   bool
	DeviceSchema          bool
	UserSchema            bool
	EdgeSchema            bool
	MessageHistoryStorage bool
	CollectionSchema      string
	ServiceName           string
	LibraryName           string
	CollectionName        string
	User                  string
	RoleName              string
	TriggerName           string
	TimerName             string
	EdgeName              string
	DeviceName            string
	PortalName            string
	PluginName            string
	AdaptorName           string
	DeploymentName        string
	ServiceCacheName      string
	WebhookName           string
	ExternalDatabaseName  string
	BucketSetName         string
	BucketSetFiles        string
	AllBucketSetFiles     bool
	BucketSetBoxName      string
	BucketSetFileName     string
	SecretName            string
	ExportUsers           bool
	ExportRows            bool
	ExportItemId          bool
}

func createAffectedAssets() AffectedAssets {
	return AffectedAssets{
		AllAssets:             AllAssets,
		AllServices:           AllServices,
		AllLibraries:          AllLibraries,
		AllEdges:              AllEdges,
		AllDevices:            AllDevices,
		AllPortals:            AllPortals,
		AllPlugins:            AllPlugins,
		AllAdaptors:           AllAdaptors,
		AllDeployments:        AllDeployments,
		AllCollections:        AllCollections,
		AllRoles:              AllRoles,
		AllUsers:              AllUsers,
		AllServiceCaches:      AllServiceCaches,
		AllWebhooks:           AllWebhooks,
		AllExternalDatabases:  AllExternalDatabases,
		UserSchema:            UserSchema,
		DeviceSchema:          DeviceSchema,
		EdgeSchema:            EdgeSchema,
		AllTriggers:           AllTriggers,
		AllTimers:             AllTimers,
		AllBucketSets:         AllBucketSets,
		AllSecrets:            AllSecrets,
		CollectionSchema:      CollectionSchema,
		ServiceName:           ServiceName,
		LibraryName:           LibraryName,
		CollectionName:        CollectionName,
		User:                  User,
		RoleName:              RoleName,
		TriggerName:           TriggerName,
		TimerName:             TimerName,
		EdgeName:              EdgeName,
		DeviceName:            DeviceName,
		PortalName:            PortalName,
		PluginName:            PluginName,
		AdaptorName:           AdaptorName,
		DeploymentName:        DeploymentName,
		ServiceCacheName:      ServiceCacheName,
		WebhookName:           WebhookName,
		ExternalDatabaseName:  ExternalDatabaseName,
		BucketSetName:         BucketSetName,
		SecretName:            SecretName,
		ExportUsers:           ExportUsers,
		ExportRows:            ExportRows,
		ExportItemId:          ExportItemId,
		BucketSetFiles:        BucketSetFiles,
		AllBucketSetFiles:     AllBucketSetFiles,
		BucketSetBoxName:      BucketSetBoxName,
		BucketSetFileName:     BucketSetFileName,
		MessageHistoryStorage: MessageHistoryStorage,
		MessageTypeTriggers:   MessageTypeTriggers,
	}
}

type Role_meta struct {
	Name        string
	Description string
	Permission  []map[string]interface{}
}

type Column struct {
	ColumnName string
	ColumnType string
}

type Collection_meta struct {
	Name          string
	Collection_id string
	Columns       []Column
}

type User_meta struct {
	Columns []Column
}

func systemMetaToMap(meta *types.System_meta) map[string]interface{} {
	result := make(map[string]interface{})
	result["platform_url"] = meta.PlatformUrl
	result["messaging_url"] = meta.MessageUrl
	result["system_key"] = meta.Key
	result["system_secret"] = meta.Secret
	result["name"] = meta.Name
	result["description"] = meta.Description
	result["auth"] = true
	return result
}

func systemMetaFromMap(theMap map[string]interface{}) *types.System_meta {
	return &types.System_meta{
		Name:        theMap["name"].(string),
		Key:         theMap["system_key"].(string),
		Secret:      theMap["system_secret"].(string),
		Description: theMap["description"].(string),
		PlatformUrl: theMap["platform_url"].(string),
		MessageUrl:  theMap["messaging_url"].(string),
		Services:    map[string]types.Service_meta{},
	}
}
