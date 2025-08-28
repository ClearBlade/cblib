package cblib

import (
	"fmt"

	"github.com/clearblade/cblib/fs"
)

func defaultZipOptions() *fs.ZipOptions {
	opts := fs.NewZipOptions(&mapper{})
	opts.AllAssets = AllAssets
	opts.AllAdaptors = AllAdaptors
	opts.AdaptorName = AdaptorName
	opts.AllBucketSets = AllBucketSets
	opts.BucketSetName = BucketSetName
	opts.AllBucketSetFiles = AllBucketSetFiles
	opts.BucketSetFiles = BucketSetFiles
	opts.BucketSetBoxName = BucketSetBoxName
	opts.BucketSetFileName = BucketSetFileName
	opts.AllServiceCaches = AllServiceCaches
	opts.ServiceCacheName = ServiceCacheName
	opts.AllCollections = AllCollections
	opts.AllCollectionSchemas = false
	opts.CollectionName = CollectionName
	opts.CollectionId = CollectionId
	opts.CollectionSchema = CollectionSchema
	opts.AllDeployments = AllDeployments
	opts.DeploymentName = DeploymentName
	opts.AllDevices = AllDevices
	opts.DeviceName = DeviceName
	opts.PushDeviceSchema = DeviceSchema
	opts.AllEdges = AllEdges
	opts.EdgeName = EdgeName
	opts.PushEdgeSchema = EdgeSchema
	opts.AllFileStores = AllFileStores
	opts.FileStoreName = FileStoreName
	opts.AllFileStoreFiles = AllFileStoreFiles
	opts.FileStoreFiles = FileStoreFiles
	opts.FileStoreFileName = FileStoreFileName
	opts.AllExternalDatabases = AllExternalDatabases
	opts.ExternalDatabaseName = ExternalDatabaseName
	opts.AllLibraries = AllLibraries
	opts.LibraryName = LibraryName
	opts.AllPlugins = AllPlugins
	opts.PluginName = PluginName
	opts.AllPortals = AllPortals
	opts.PortalName = PortalName
	opts.AllRoles = AllRoles
	opts.RoleName = RoleName
	opts.AllSecrets = AllSecrets
	opts.SecretName = SecretName
	opts.AllServices = AllServices
	opts.ServiceName = ServiceName
	opts.AllTimers = AllTimers
	opts.TimerName = TimerName
	opts.AllTriggers = AllTriggers
	opts.TriggerName = TriggerName
	opts.AllUsers = AllUsers
	opts.UserName = User
	opts.UserId = UserId
	opts.PushUserSchema = UserSchema
	opts.AllWebhooks = AllWebhooks
	opts.WebhookName = WebhookName
	opts.PushMessageHistoryStorage = MessageHistoryStorage
	opts.PushMessageTypeTriggers = MessageTypeTriggers
	return opts
}

type mapper struct{}

func (m *mapper) GetCollectionNameById(id string) (string, error) {
	return getCollectionNameById(id)
}

func (m *mapper) GetUserEmailById(wantedId string) (string, error) {
	users, err := getUserEmailToId()
	if err != nil {
		return "", err
	}

	for email, id := range users {
		if id == wantedId {
			return email, nil
		}
	}

	return "", fmt.Errorf("user with id %s not found", wantedId)
}
