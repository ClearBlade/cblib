package cblib

import (
	"fmt"
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

func NewDefaultPushOptions() *systemPushOptions {
	return &systemPushOptions{
		// TODO:
	}
}

func (s *systemPushOptions) GetFileRegex() *regexp.Regexp {
	regexBuilder := strings.Builder{}
	// TODO: handle collection schemas

	appendRegexForDirectory(&regexBuilder, adaptorsDir, s.getAdaptorsRegex())
	appendRegexForDirectory(&regexBuilder, bucketSetsDir, s.getBucketSetsRegex())
	appendRegexForDirectory(&regexBuilder, bucketSetFiles.BucketSetFilesDir, s.getBucketSetFilesRegex())
	appendRegexForDirectory(&regexBuilder, serviceCachesDir, s.getCachesRegex())
	appendRegexForDirectory(&regexBuilder, dataDir, s.getCollectionsRegex())
	appendRegexForDirectory(&regexBuilder, deploymentsDir, s.getDeploymentsRegex())
	appendRegexForDirectory(&regexBuilder, devicesDir, s.getDevicesRegex())
	appendRegexForDirectory(&regexBuilder, edgesDir, s.getEdgesRegex())
	appendRegexForDirectory(&regexBuilder, externalDatabasesDir, s.getExternalDatabasesRegex())
	appendRegexForDirectory(&regexBuilder, libDir, s.getLibrariesRegex())
	appendRegexForDirectory(&regexBuilder, messageHistoryStorageDir, s.getMessageHistoryStorageRegex())
	appendRegexForDirectory(&regexBuilder, messageTypeTriggersDir, s.getMessageTypeTriggerRegex())
	appendRegexForDirectory(&regexBuilder, pluginsDir, s.getPluginsRegex())
	appendRegexForDirectory(&regexBuilder, portalsDir, s.getPortalsRegex())
	appendRegexForDirectory(&regexBuilder, rolesDir, s.getRolesRegex())
	appendRegexForDirectory(&regexBuilder, secretsDir, s.getSecretsRegex())
	appendRegexForDirectory(&regexBuilder, svcDir, s.getServicesRegex())
	appendRegexForDirectory(&regexBuilder, timersDir, s.getTimersRegex())
	appendRegexForDirectory(&regexBuilder, triggersDir, s.getTriggerRegex())
	appendRegexForDirectory(&regexBuilder, usersDir, s.getUserRegex())
	appendRegexForDirectory(&regexBuilder, webhooksDir, s.getWebhooksRegex())

	// TODO: Remove trailing or
	return regexp.MustCompile(regexBuilder.String())
}

func appendRegexForDirectory(builder *strings.Builder, dir string, regex string) {
	if regex == "" {
		return
	}

	builder.WriteByte('(')
	builder.WriteString(regexp.QuoteMeta(dir))
	builder.WriteString("/(")
	builder.WriteString(regex)
	builder.WriteString("))|")
}

func (s *systemPushOptions) getAdaptorsRegex() string {
	if s.AllAssets || s.AllAdaptors {
		return "*+"
	}

	if s.AdaptorName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s/*+", s.AdaptorName))
}

func (s *systemPushOptions) getBucketSetsRegex() string {
	if s.AllAssets || s.AllBucketSets {
		return "*+"
	}

	if s.BucketSetName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.BucketSetName))
}

func (s *systemPushOptions) getBucketSetFilesRegex() string {
	if s.AllAssets || s.AllBucketSetFiles {
		return "*+"
	}

	if s.BucketSetFiles == "" || s.BucketSetBoxName == "" {
		return ""
	}

	if s.BucketSetFileName == "" {
		return regexp.QuoteMeta(fmt.Sprintf("%s/%s/*+", BucketSetFiles, BucketSetBoxName))
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s/%s/%s", BucketSetFiles, BucketSetBoxName, BucketSetFileName))
}

func (s *systemPushOptions) getCachesRegex() string {
	if s.AllAssets || s.AllServiceCaches {
		return "*+"
	}

	// TODO: Do a test where everything is false and check the regex
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
		return "*+"
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
		return "*+"
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
		return "*+"
	}

	if s.DeploymentName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.DeploymentName))
}

func (s *systemPushOptions) getDevicesRegex() string {
	if s.AllAssets || s.AllDevices {
		return "*+"
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
		return "*+"
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
		return "*+"
	}

	if s.ExternalDatabaseName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.ExternalDatabaseName))
}

func (s *systemPushOptions) getLibrariesRegex() string {
	if s.AllAssets || s.AllLibraries {
		return "*+"
	}

	if s.LibraryName == "" {
		return ""
	}

	return fmt.Sprintf("%s/*+", regexp.QuoteMeta(s.LibraryName))
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
		return "*+"
	}

	if s.PluginName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.PluginName))
}

func (s *systemPushOptions) getPortalsRegex() string {
	if s.AllAssets || s.AllPortals {
		return "*+"
	}

	if s.PortalName == "" {
		return ""
	}

	return fmt.Sprintf("%s/*+", regexp.QuoteMeta(s.PortalName))
}

func (s *systemPushOptions) getRolesRegex() string {
	if s.AllAssets || s.AllRoles {
		return "*+"
	}

	if s.RoleName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.RoleName))
}

func (s *systemPushOptions) getSecretsRegex() string {
	if s.AllAssets || s.AllSecrets {
		return "*+"
	}

	if s.SecretName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.SecretName))
}

func (s *systemPushOptions) getServicesRegex() string {
	if s.AllAssets || s.AllServices {
		return "*+"
	}

	if s.ServiceName == "" {
		return ""
	}

	return fmt.Sprintf("%s/*+", regexp.QuoteMeta(s.ServiceName))
}

func (s *systemPushOptions) getTimersRegex() string {
	if s.AllAssets || s.AllTimers {
		return "*+"
	}

	if s.TimerName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.TimerName))
}

func (s *systemPushOptions) getTriggerRegex() string {
	if s.AllAssets || s.AllTriggers {
		return "*+"
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
		fmt.Printf("Ignoring user %q: 5s", s.UserId, err)
		return emails
	}

	emails = append(emails, email)
	return emails
}

func (s *systemPushOptions) getUserRegex() string {
	if s.AllAssets || s.AllUsers {
		return "*+"
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
		return "*+"
	}

	if s.WebhookName == "" {
		return ""
	}

	return regexp.QuoteMeta(fmt.Sprintf("%s.json", s.WebhookName))
}
