package fs

import (
	"fmt"
	"slices"

	"github.com/clearblade/cblib/syspath"
)

type ZipOptions struct {
	mapper IdMapper

	AllAssets   bool
	AllAdaptors bool
	AdaptorName string

	AllBucketSets bool
	BucketSetName string

	AllBucketSetFiles bool
	BucketSetFiles    string
	BucketSetBoxName  string
	BucketSetFileName string

	AllServiceCaches bool
	ServiceCacheName string

	AllCollections   bool
	CollectionName   string
	CollectionId     string
	CollectionSchema string

	AllDeployments bool
	DeploymentName string

	AllDevices       bool
	DeviceName       string
	PushDeviceSchema bool

	AllEdges       bool
	EdgeName       string
	PushEdgeSchema bool

	AllExternalDatabases bool
	ExternalDatabaseName string

	AllLibraries bool
	LibraryName  string

	AllPlugins bool
	PluginName string

	AllPortals bool
	PortalName string

	AllRoles bool
	RoleName string

	AllSecrets bool
	SecretName string

	AllServices bool
	ServiceName string

	AllTimers bool
	TimerName string

	AllTriggers bool
	TriggerName string

	AllUsers       bool
	UserName       string
	UserId         string
	PushUserSchema bool

	AllWebhooks bool
	WebhookName string

	PushMessageHistoryStorage bool
	PushMessageTypeTriggers   bool
}

type IdMapper interface {
	GetCollectionNameById(id string) (string, error)
	GetUserEmailById(id string) (string, error)
}

func NewZipOptions(idMapper IdMapper) *ZipOptions {
	return &ZipOptions{mapper: idMapper}
}

func (s *ZipOptions) shouldPushAdaptor(name string) bool {
	if s.AllAssets || s.AllAdaptors {
		return true
	}

	return s.AdaptorName == name
}

func (s *ZipOptions) shouldPushBucketSetFile(data *syspath.FullBucketPath) bool {
	if s.AllAssets || s.AllBucketSetFiles {
		return true
	}

	if s.BucketSetFiles != data.BucketName || s.BucketSetBoxName != data.Box {
		return false
	}

	// Empty means push all files in the box
	if s.BucketSetFileName == "" {
		return true
	}

	return s.BucketSetFileName == data.RelativePath
}

func (s *ZipOptions) shouldPushBucketSetMeta(bucketName string) bool {
	if s.AllAssets || s.AllBucketSets || s.AllBucketSetFiles {
		return true
	}

	// Push the meta if any file for this bucket set is being pushed
	return s.BucketSetName == bucketName || s.BucketSetFiles == bucketName
}

func (s *ZipOptions) shouldPushService(name string) bool {
	if s.AllAssets || s.AllServices {
		return true
	}

	return s.ServiceName == name
}

func (s *ZipOptions) shouldPushLibrary(name string) bool {
	if s.AllAssets || s.AllLibraries {
		return true
	}

	return s.LibraryName == name
}

func (s *ZipOptions) shouldPushCollection(name string) bool {
	if s.AllAssets || s.AllCollections {
		return true
	}

	return slices.Contains(s.getCollectionNames(), name)
}

/**
 * The user can specify a collection name, collection id, or both.
 * This is a helper function to return all the collection names that
 * were specified, if any
 */
func (s *ZipOptions) getCollectionNames() []string {
	names := []string{}
	if s.CollectionName != "" {
		names = append(names, s.CollectionName)
	}

	if s.CollectionId == "" {
		return names
	}

	name, err := s.mapper.GetCollectionNameById(s.CollectionId)
	if err != nil {
		fmt.Printf("Not pushing collection id %q: %s", s.CollectionId, err)
		return names
	}

	names = append(names, name)
	return names
}

func (s *ZipOptions) shouldPushCollectionSchemaOnly(name string) bool {
	// We're already pushing the entire collection
	if s.shouldPushCollection(name) {
		return false
	}

	return s.CollectionSchema == name
}

func (s *ZipOptions) shouldPushDeployment(name string) bool {
	if s.AllAssets || s.AllDeployments {
		return true
	}

	return s.DeploymentName == name
}

func (s *ZipOptions) shouldPushDeviceSchema() bool {
	return s.AllAssets || s.AllDevices || s.PushDeviceSchema || s.DeviceName != ""
}

func (s *ZipOptions) shouldPushDevice(name string) bool {
	if s.AllAssets || s.AllDevices {
		return true
	}

	return s.DeviceName == name
}

func (s *ZipOptions) shouldPushEdgeSchema() bool {
	return s.AllAssets || s.AllEdges || s.PushEdgeSchema || s.EdgeName != ""
}

func (s *ZipOptions) shouldPushEdge(name string) bool {
	if s.AllAssets || s.AllEdges {
		return true
	}

	return s.EdgeName == name
}

func (s *ZipOptions) shouldPushExternalDatabase(name string) bool {
	if s.AllAssets || s.AllExternalDatabases {
		return true
	}

	return s.ExternalDatabaseName == name
}

func (s *ZipOptions) shouldPushMessageHistoryStorage() bool {
	return s.AllAssets || s.PushMessageHistoryStorage
}

func (s *ZipOptions) shouldPushMessageTypeTriggers() bool {
	return s.AllAssets || s.PushMessageTypeTriggers
}

func (s *ZipOptions) shouldPushPlugin(name string) bool {
	if s.AllAssets || s.AllPlugins {
		return true
	}

	return s.PluginName == name
}

func (s *ZipOptions) shouldPushPortal(name string) bool {
	if s.AllAssets || s.AllPortals {
		return true
	}

	return s.PortalName == name
}

func (s *ZipOptions) shouldPushRole(name string) bool {
	if s.AllAssets || s.AllRoles {
		return true
	}

	return s.RoleName == name
}

func (s *ZipOptions) shouldPushSecret(name string) bool {
	if s.AllAssets || s.AllSecrets {
		return true
	}

	return s.SecretName == name
}

func (s *ZipOptions) shouldPushServiceCache(name string) bool {
	if s.AllAssets || s.AllServiceCaches {
		return true
	}

	return s.ServiceCacheName == name
}

func (s *ZipOptions) shouldPushTimer(name string) bool {
	if s.AllAssets || s.AllTimers {
		return true
	}

	return s.TimerName == name
}

func (s *ZipOptions) shouldPushTrigger(name string) bool {
	if s.AllAssets || s.AllTriggers {
		return true
	}

	return s.TriggerName == name
}

func (s *ZipOptions) shouldPushUser(email string) bool {
	if s.AllAssets || s.AllUsers {
		return true
	}

	emails := s.getUserEmails()
	return slices.Contains(emails, email)
}

func (s *ZipOptions) shouldPushUserSchema() bool {
	return s.AllAssets || s.AllUsers || s.PushUserSchema || len(s.getUserEmails()) > 0
}

func (s *ZipOptions) getUserEmails() []string {
	emails := []string{}
	if s.UserName != "" {
		emails = append(emails, s.UserName)
	}

	if s.UserId == "" {
		return emails
	}

	email, err := s.mapper.GetUserEmailById(s.UserId)
	if err != nil {
		fmt.Printf("Ignoring user %q: %s", s.UserId, err)
		return emails
	}

	emails = append(emails, email)
	return emails
}

func (s *ZipOptions) shouldPushWebhook(name string) bool {
	if s.AllAssets || s.AllWebhooks {
		return true
	}

	return s.WebhookName == name
}
