package cblib

import (
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strings"

	"github.com/clearblade/cblib/models/bucketSetFiles"
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

type CbSystemFile struct {
	fs.DirEntry
}

func (s *systemPushOptions) IsPathInUpload(path string) bool {
	// TODO: Check if any of our dirs are a prefix
}

// TODO: Does everything match the empty regex
func (s *systemPushOptions) GetFileRegex() *regexp.Regexp {
	regexBuilder := strings.Builder{}
	regexBuilder.WriteByte('^')

	regexDirs := []struct {
		dir   string
		regex string
	}{
		{adaptorsDir, s.getAdaptorsRegex()},
		{bucketSetsDir, s.getBucketSetsRegex()},
		{bucketSetFiles.BucketSetFilesDir, s.getBucketSetFilesRegex()},
		{serviceCachesDir, s.getCachesRegex()},
		{dataDir, s.getCollectionsRegex()},
		{deploymentsDir, s.getDeploymentsRegex()},
		{devicesDir, s.getDevicesRegex()},
		{edgesDir, s.getEdgesRegex()},
		{externalDatabasesDir, s.getExternalDatabasesRegex()},
		{libDir, s.getLibrariesRegex()},
		{messageHistoryStorageDir, s.getMessageHistoryStorageRegex()},
		{messageTypeTriggersDir, s.getMessageTypeTriggerRegex()},
		{pluginsDir, s.getPluginsRegex()},
		{portalsDir, s.getPortalsRegex()},
		{rolesDir, s.getRolesRegex()},
		{secretsDir, s.getSecretsRegex()},
		{svcDir, s.getServicesRegex()},
		{timersDir, s.getTimersRegex()},
		{triggersDir, s.getTriggerRegex()},
		{usersDir, s.getUserRegex()},
		{webhooksDir, s.getWebhooksRegex()},
	}

	lastWrittenIdx := -2
	for i, info := range regexDirs {
		if info.regex == "" || info.dir == "" {
			continue
		}

		if i == lastWrittenIdx+1 {
			regexBuilder.WriteByte('|')
		}

		regexBuilder.WriteString(makeRegexForDirectory(info.dir, info.regex))
		lastWrittenIdx = i
	}

	regexBuilder.WriteByte('$')
	return regexp.MustCompile(regexBuilder.String())
}

func (s *systemPushOptions) GetCollectionSchemaRegex() *regexp.Regexp {
	regexBuilder := strings.Builder{}
	regexBuilder.WriteByte('^')
	regexBuilder.WriteString(makeRegexForDirectory(dataDir, s.getCollectionSchemasRegex()))
	regexBuilder.WriteByte('$')
	return regexp.MustCompile(regexBuilder.String())
}

func makeRegexForDirectory(dir string, regex string) string {
	builder := strings.Builder{}
	builder.WriteByte('(')
	builder.WriteString(regexp.QuoteMeta(dir))
	builder.WriteString("/(")
	builder.WriteString(regex)
	builder.WriteString("))")
	return builder.String()
}

func (s *systemPushOptions) getAdaptorsRegex() string {
	if s.AllAssets || s.AllAdaptors {
		return ".*"
	}

	if s.AdaptorName == "" {
		return ""
	}

	return fmt.Sprintf("%s/.*", regexp.QuoteMeta(s.AdaptorName))
}

func (s *systemPushOptions) getBucketSetsRegex() string {
	if s.AllAssets || s.AllBucketSets {
		return ".*"
	}

	if s.BucketSetName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.BucketSetName))
}

func (s *systemPushOptions) getBucketSetFilesRegex() string {
	if s.AllAssets || s.AllBucketSetFiles {
		return ".*"
	}

	if s.BucketSetFiles == "" || s.BucketSetBoxName == "" {
		return ""
	}

	if s.BucketSetFileName == "" {
		return fmt.Sprintf("%s/%s/.*", regexp.QuoteMeta(BucketSetFiles), regexp.QuoteMeta(BucketSetBoxName))
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s/%s/%s", BucketSetFiles, BucketSetBoxName, BucketSetFileName))
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

func (s *systemPushOptions) getCollectionsRegex() string {
	if s.AllAssets || s.AllCollections {
		return ".*"
	}

	collectionsRegex := strings.Builder{}
	collections := s.getCollectionNames()
	for i, name := range collections {
		if i > 0 {
			collectionsRegex.WriteByte('|')
		}

		collectionsRegex.WriteString("(" + regexp.QuoteMeta(fmt.Sprintf("%s.json", name)) + ")")
	}

	return collectionsRegex.String()
}

// TODO: Move somewhere else?
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

func (s *systemPushOptions) getCollectionSchemasRegex() string {
	if s.AllAssets || s.AllCollections {
		return ".*"
	}

	if s.CollectionSchema == "" {
		return ""
	}

	// Don't need to push the schema if the collection is already being pushed
	collections := s.getCollectionNames()
	if slices.Contains(collections, s.CollectionSchema) {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.CollectionSchema))
}

func (s *systemPushOptions) getDeploymentsRegex() string {
	if s.AllAssets || s.AllDeployments {
		return ".*"
	}

	if s.DeploymentName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.DeploymentName))
}

func (s *systemPushOptions) getDevicesRegex() string {
	if s.AllAssets || s.AllDevices {
		return ".*"
	}

	devices := strings.Builder{}
	if s.PushDeviceSchema || s.DeviceName != "" {
		devices.WriteString("(" + regexp.QuoteMeta("schema.json") + ")")
	}

	if s.DeviceName == "" {
		return devices.String()
	}

	devices.WriteString("|(")
	devices.WriteString(regexp.QuoteMeta(fmt.Sprintf("%s.json", s.DeviceName)))
	devices.WriteString(")|(")
	devices.WriteString(regexp.QuoteMeta(fmt.Sprintf("roles/%s.json", s.DeviceName)))
	devices.WriteString(")")
	return devices.String()
}

func (s *systemPushOptions) getEdgesRegex() string {
	if s.AllAssets || s.AllEdges {
		return ".*"
	}

	edges := strings.Builder{}
	if s.PushEdgeSchema || s.EdgeName != "" {
		edges.WriteString("(" + regexp.QuoteMeta("schema.json") + ")")
	}

	if s.EdgeName == "" {
		return edges.String()
	}

	edges.WriteString("|(")
	edges.WriteString(regexp.QuoteMeta(fmt.Sprintf("%s.json", s.EdgeName)))
	edges.WriteString(")")
	return edges.String()
}

func (s *systemPushOptions) getExternalDatabasesRegex() string {
	if s.AllAssets || s.AllExternalDatabases {
		return ".*"
	}

	if s.ExternalDatabaseName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.ExternalDatabaseName))
}

func (s *systemPushOptions) getLibrariesRegex() string {
	if s.AllAssets || s.AllLibraries {
		return ".*"
	}

	if s.LibraryName == "" {
		return ""
	}

	return fmt.Sprintf("%s/.*", regexp.QuoteMeta(s.LibraryName))
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

func (s *systemPushOptions) getServicesRegex() string {
	if s.AllAssets || s.AllServices {
		return ".*"
	}

	if s.ServiceName == "" {
		return ""
	}

	return fmt.Sprintf("%s/.*", regexp.QuoteMeta(s.ServiceName))
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
