package cblib

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/clearblade/cblib/syspath"
)

type systemPushOptions struct {
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

func DefaultPushOptions() *systemPushOptions {
	return &systemPushOptions{
		AllAssets:                 AllAssets,
		AllAdaptors:               AllAdaptors,
		AdaptorName:               AdaptorName,
		AllBucketSets:             AllBucketSets,
		BucketSetName:             BucketSetName,
		AllBucketSetFiles:         AllBucketSetFiles,
		BucketSetFiles:            BucketSetFiles,
		BucketSetBoxName:          BucketSetBoxName,
		BucketSetFileName:         BucketSetFileName,
		AllServiceCaches:          AllServiceCaches,
		ServiceCacheName:          ServiceCacheName,
		AllCollections:            AllCollections,
		CollectionName:            CollectionName,
		CollectionId:              CollectionId,
		CollectionSchema:          CollectionSchema,
		AllDeployments:            AllDeployments,
		DeploymentName:            DeploymentName,
		AllDevices:                AllDevices,
		DeviceName:                DeviceName,
		PushDeviceSchema:          DeviceSchema,
		AllEdges:                  AllEdges,
		EdgeName:                  EdgeName,
		PushEdgeSchema:            EdgeSchema,
		AllExternalDatabases:      AllExternalDatabases,
		ExternalDatabaseName:      ExternalDatabaseName,
		AllLibraries:              AllLibraries,
		LibraryName:               LibraryName,
		AllPlugins:                AllPlugins,
		PluginName:                PluginName,
		AllPortals:                AllPortals,
		PortalName:                PortalName,
		AllRoles:                  AllRoles,
		RoleName:                  RoleName,
		AllSecrets:                AllSecrets,
		SecretName:                SecretName,
		AllServices:               AllServices,
		ServiceName:               ServiceName,
		AllTimers:                 AllTimers,
		TimerName:                 TimerName,
		AllTriggers:               AllTriggers,
		TriggerName:               TriggerName,
		AllUsers:                  AllUsers,
		UserName:                  User,
		UserId:                    UserId,
		PushUserSchema:            UserSchema,
		AllWebhooks:               AllWebhooks,
		WebhookName:               WebhookName,
		PushMessageHistoryStorage: MessageHistoryStorage,
		PushMessageTypeTriggers:   MessageTypeTriggers,
	}
}

func (s *systemPushOptions) ShouldPushFile(relPath string) bool {
	if syspath.IsAdaptorPath(relPath) {
		return s.shouldPushAdaptorFile(relPath)
	}

	if syspath.IsBucketSetFilePath(relPath) {
		return s.shouldPushBucketSetDataFile(relPath)
	}

	if syspath.IsBucketSetMetaPath(relPath) {
		return s.shouldPushBucketSetMetaFile(relPath)
	}

	if syspath.IsCodePath(relPath) {
		return s.shouldPushCodeFile(relPath)
	}

	if syspath.IsCollectionPath(relPath) {
		return s.shouldPushCollectionFile(relPath)
	}

	if syspath.IsDeploymentPath(relPath) {
		return s.shouldPushDeploymentFile(relPath)
	}

	if syspath.IsDevicePath(relPath) {
		return s.shouldPushDeviceFile(relPath)
	}

	if syspath.IsEdgePath(relPath) {
		return s.shouldPushEdgeFile(relPath)
	}

	if syspath.IsExternalDbPath(relPath) {
		return s.shouldPushExternalDatabaseFile(relPath)
	}

	return false
}

func (s *systemPushOptions) shouldPushAdaptorFile(relPath string) bool {
	name, err := syspath.GetAdaptorNameFromPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushAdaptor(name)
}

func (s *systemPushOptions) shouldPushAdaptor(name string) bool {
	if s.AllAssets || s.AllAdaptors {
		return true
	}

	return s.AdaptorName == name
}

func (s *systemPushOptions) shouldPushBucketSetDataFile(relPath string) bool {
	data, err := syspath.ParseBucketPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushBucketSetData(data)
}

func (s *systemPushOptions) shouldPushBucketSetData(data *syspath.FullBucketPath) bool {
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

func (s *systemPushOptions) shouldPushBucketSetMetaFile(relPath string) bool {
	name, err := syspath.GetBucketSetNameFromPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushBucketSetMeta(name)
}

func (s *systemPushOptions) shouldPushBucketSetMeta(bucketName string) bool {
	if s.AllAssets || s.AllBucketSets || s.AllBucketSetFiles {
		return true
	}

	// Push the meta if any file for this bucket set is being pushed
	return s.BucketSetName == bucketName || s.BucketSetFiles == bucketName
}

func (s *systemPushOptions) shouldPushCodeFile(relPath string) bool {
	if service, err := syspath.GetServiceNameFromPath(relPath); err == nil {
		return s.shouldPushService(service)
	}

	if library, err := syspath.GetLibraryNameFromPath(relPath); err == nil {
		return s.shouldPushLibrary(library)
	}

	return false
}

func (s *systemPushOptions) shouldPushService(name string) bool {
	if s.AllAssets || s.AllServices {
		return true
	}

	return s.ServiceName == name
}

func (s *systemPushOptions) shouldPushLibrary(name string) bool {
	if s.AllAssets || s.AllLibraries {
		return true
	}

	return s.LibraryName == name
}

/**
 * The user can specify a collection name, collection id, or both.
 * This is a helper function to return all the collection names that
 * were specified, if any
 */
func (s *systemPushOptions) getCollectionNames() []string {
	names := []string{}
	if s.CollectionName != "" {
		names = append(names, s.CollectionName)
	}

	if s.CollectionId == "" {
		return names
	}

	name, err := getCollectionNameById(s.CollectionId)
	if err != nil {
		fmt.Printf("Not pushing collection id %q: %s", s.CollectionId, err)
	}

	names = append(names, name)
	return names
}

func (s *systemPushOptions) shouldPushCollectionFile(relPath string) bool {
	name, err := syspath.GetCollectionNameFromPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushCollection(name)
}

func (s *systemPushOptions) shouldPushCollection(name string) bool {
	if s.AllAssets || s.AllCollections {
		return true
	}

	return slices.Contains(s.getCollectionNames(), name)
}

func (s *systemPushOptions) shouldPushCollectionSchemaFileOnly(relPath string) bool {
	name, err := syspath.GetCollectionNameFromPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushCollectionSchemaOnly(name)
}

func (s *systemPushOptions) shouldPushCollectionSchemaOnly(name string) bool {
	// We're already pushing the entire collection
	if s.shouldPushCollection(name) {
		return false
	}

	return s.CollectionSchema == name
}

func (s *systemPushOptions) shouldPushDeploymentFile(relPath string) bool {
	name, err := syspath.GetDeploymentNameFromPath(relPath)
	if err != nil {
		return false
	}

	return s.shouldPushDeployment(name)
}

func (s *systemPushOptions) shouldPushDeployment(name string) bool {
	if s.AllAssets || s.AllDeployments {
		return true
	}

	return s.DeploymentName == name
}

func (s *systemPushOptions) shouldPushDeviceFile(relPath string) bool {
	if syspath.IsDeviceSchemaPath(relPath) {
		return s.shouldPushDeviceSchema()
	}

	if name, err := syspath.GetDeviceNameFromDataPath(relPath); err == nil {
		return s.shouldPushDevice(name)
	}

	if name, err := syspath.GetDeviceNameFromRolePath(relPath); err == nil {
		return s.shouldPushDevice(name)
	}

	return false
}

func (s *systemPushOptions) shouldPushDeviceSchema() bool {
	return s.AllAssets || s.AllDevices || s.PushDeviceSchema || s.DeviceName != ""
}

func (s *systemPushOptions) shouldPushDevice(name string) bool {
	if s.AllAssets || s.AllDevices {
		return true
	}

	return s.DeviceName == name
}

func (s *systemPushOptions) shouldPushEdgeFile(relPath string) bool {
	if syspath.IsEdgeSchemaPath(relPath) {
		return s.shouldPushEdgeSchema()
	}

	if name, err := syspath.GetEdgeNameFromPath(relPath); err == nil {
		return s.shouldPushEdge(name)
	}

	return false
}

func (s *systemPushOptions) shouldPushEdgeSchema() bool {
	return s.AllAssets || s.AllEdges || s.PushEdgeSchema || s.EdgeName != ""
}

func (s *systemPushOptions) shouldPushEdge(name string) bool {
	if s.AllAssets || s.AllEdges {
		return true
	}

	return s.EdgeName == name
}

func (s *systemPushOptions) shouldPushExternalDatabaseFile(relPath string) bool {
	if name, err := syspath.GetExternalDbNameFromPath(relPath); err == nil {
		return s.shouldPushExternalDatabase(name)
	}

	return false
}

func (s *systemPushOptions) shouldPushExternalDatabase(name string) bool {
	if s.AllAssets || s.AllExternalDatabases {
		return true
	}

	return s.ExternalDatabaseName == name
}

func (s *systemPushOptions) getMessageHistoryStorageRegex() string {
	if s.AllAssets || s.PushMessageHistoryStorage {
		return regexp.QuoteMeta("storage.json")
	}

	return ""
}

func (s *systemPushOptions) getMessageTypeTriggerRegex() string {
	if s.AllAssets || s.PushMessageTypeTriggers {
		return regexp.QuoteMeta("triggers.json")
	}

	return ""
}

func (s *systemPushOptions) getPluginsRegex() string {
	if s.AllAssets || s.AllPlugins {
		return ".*"
	}

	if s.PluginName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.PluginName))
}

func (s *systemPushOptions) getPortalsRegex() string {
	if s.AllAssets || s.AllPortals {
		return ".*"
	}

	if s.PortalName == "" {
		return ""
	}

	return fmt.Sprintf("%s/.*", regexp.QuoteMeta(s.PortalName))
}

func (s *systemPushOptions) getRolesRegex() string {
	if s.AllAssets || s.AllRoles {
		return ".*"
	}

	if s.RoleName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.RoleName))
}

func (s *systemPushOptions) getSecretsRegex() string {
	if s.AllAssets || s.AllSecrets {
		return ".*"
	}

	if s.SecretName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.SecretName))
}

func (s *systemPushOptions) getTimersRegex() string {
	if s.AllAssets || s.AllTimers {
		return ".*"
	}

	if s.TimerName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.TimerName))
}

func (s *systemPushOptions) getTriggerRegex() string {
	if s.AllAssets || s.AllTriggers {
		return ".*"
	}

	if s.TriggerName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.TriggerName))
}

func (s *systemPushOptions) getUserEmails() []string {
	emails := []string{}
	if s.UserName != "" {
		emails = append(emails, s.UserName)
	}

	if s.UserId == "" {
		return emails
	}

	email, err := getUserEmailByID(s.UserId)
	if err != nil {
		fmt.Printf("Ignoring user %q: %s", s.UserId, err)
		return emails
	}

	emails = append(emails, email)
	return emails
}

func (s *systemPushOptions) getUserRegex() string {
	if s.AllAssets || s.AllUsers {
		return ".*"
	}

	userRegex := strings.Builder{}
	emails := s.getUserEmails()
	if s.PushUserSchema || len(emails) > 0 {
		userRegex.WriteString("(" + regexp.QuoteMeta("schema.json") + ")")
	}

	for _, email := range emails {
		userRegex.WriteString("|(")
		userRegex.WriteString(regexp.QuoteMeta(fmt.Sprintf("%s.json", email)))
		userRegex.WriteString(")|(")
		userRegex.WriteString(regexp.QuoteMeta(fmt.Sprintf("roles/%s.json", email)))
		userRegex.WriteString(")")
	}

	return userRegex.String()
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

func (s *systemPushOptions) getWebhooksRegex() string {
	if s.AllAssets || s.AllWebhooks {
		return ".*"
	}

	if s.WebhookName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.WebhookName))
}

func (s *systemPushOptions) getCachesRegex() string {
	if s.AllAssets || s.AllServiceCaches {
		return ".*"
	}

	if s.ServiceCacheName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.ServiceCacheName))
}
