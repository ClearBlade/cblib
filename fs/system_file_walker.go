package fs

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/clearblade/cblib/syspath"
)

type systemFileHandler interface {
	WalkAdaptor(path, relPath string, adaptorName string)
	WalkAdaptorFile(path, relPath string, adaptorName string)
	WalkAdaptorFileMeta(path, relPath string, adaptorName string)
	WalkBucketSetMeta(path, relPath string, bucketSetName string)
	WalkBucketSetFile(path, relPath string, bucketFile *syspath.FullBucketPath)
	WalkService(path, relPath string, serviceName string)
	WalkLibrary(path, relPath string, libraryName string)
	WalkCollection(path, relPath string, collectionName string)
	WalkDeployment(path, relPath string, deploymentName string)
	WalkDevice(path, relPath string, deviceName string)
	WalkDeviceRole(path, relPath string, deviceName string)
	WalkDeviceSchema(path string)
	WalkEdge(path, relPath string, edgeName string)
	WalkEdgeSchema(path string)
	WalkExternalDatabase(path, relPath string, externalDatabaseName string)
	WalkMessageHistoryStorage(path string)
	WalkMessageTypeTriggers(path string)
	WalkPlugin(path, relPath string, pluginName string)
	WalkPortal(path, relPath string, portalName string)
	WalkPortalDatasource(path, relPath string, portalName string)
	WalkPortalInternalResources(path, relPath string, portalName string)
	WalkPortalWidget(path, relPath string, portalName string)
	WalkPortalWidgetParser(path, relPath string, portalName string)
	WalkRole(path, relPath string, roleName string)
	WalkSecret(path, relPath string, secretName string)
	WalkServiceCache(path, relPath string, serviceCacheName string)
	WalkTimer(path, relPath string, timerName string)
	WalkTrigger(path, relPath string, triggerName string)
	WalkUser(path, relPath string, email string)
	WalkUserRole(path, relPath string, email string)
	WalkUserSchema(path string)
	WalkWebhook(path, relPath string, webhookName string)
}

func walkSystemFiles(rootDir string, handler systemFileHandler) error {
	return filepath.WalkDir(rootDir, func(absolutePath string, d fs.DirEntry, err error) error {
		path, pathErr := filepath.Rel(rootDir, absolutePath)
		if pathErr != nil {
			return fmt.Errorf("could not make %s relative to %s: %w", absolutePath, rootDir, pathErr)
		}

		// Skip directories we don't care about
		if d.IsDir() && !isAssetPath(path) && path != "." {
			return filepath.SkipDir
		}

		// Only call handlers on files
		if !d.IsDir() && err == nil {
			callHandler(handler, absolutePath, path)
		}

		return err
	})
}

type assetPathHandler struct {
	isAssetPath       func(relPath string) bool
	handleAssetAtPath func(systemFileHandler systemFileHandler, absPath, relPath string)
}

var assetHandlers = []assetPathHandler{
	{syspath.IsAdaptorPath, callAdaptorHandlers},
	{syspath.IsBucketSetMetaPath, callBucketSetMetaHandlers},
	{syspath.IsBucketSetFilePath, callBucketSetFileHandlers},
	{syspath.IsCodePath, callCodeHandlers},
	{syspath.IsCollectionPath, callCollectionHandlers},
	{syspath.IsDeploymentPath, callDeploymentHandlers},
	{syspath.IsDevicePath, callDeviceHandlers},
	{syspath.IsEdgePath, callEdgeHandlers},
	{syspath.IsExternalDbPath, callExternalDatabaseHandlers},
	{syspath.IsMessageHistoryStorageFile, callMessageHistoryStorageHandlers},
	{syspath.IsMessageTypeTriggerPath, callMessageTypeTriggersHandlers},
	{syspath.IsPluginPath, callPluginHandlers},
	{syspath.IsPortalPath, callPortalHandlers},
	{syspath.IsRolePath, callRoleHandlers},
	{syspath.IsSecretPath, callSecretHandlers},
	{syspath.IsServiceCachePath, callServiceCacheHandlers},
	{syspath.IsTimerPath, callTimerHandlers},
	{syspath.IsTriggerPath, callTriggerHandlers},
	{syspath.IsUserPath, callUserHandlers},
	{syspath.IsWebhookPath, callWebhookHandlers},
}

func isAssetPath(relPath string) bool {
	for _, assetHandler := range assetHandlers {
		if assetHandler.isAssetPath(relPath) {
			return true
		}
	}

	return false
}

func callHandler(handler systemFileHandler, absPath, relPath string) {
	for _, assetHandler := range assetHandlers {
		if assetHandler.isAssetPath(relPath) {
			assetHandler.handleAssetAtPath(handler, absPath, relPath)
			return
		}
	}
}

func callAdaptorHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetAdaptorNameFromPath(relPath); err == nil {
		handler.WalkAdaptor(absPath, relPath, name)
	} else if name, _, err := syspath.GetAdaptorFileMetaNameFromPath(relPath); err == nil {
		handler.WalkAdaptorFileMeta(absPath, relPath, name)
	} else if name, _, err := syspath.GetAdaptorFileDataNameFromPath(relPath); err == nil {
		handler.WalkAdaptorFile(absPath, relPath, name)
	}
}

func callBucketSetMetaHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetBucketSetNameFromPath(relPath); err == nil {
		handler.WalkBucketSetMeta(absPath, relPath, name)
	}
}

func callBucketSetFileHandlers(handler systemFileHandler, absPath, relPath string) {
	if parsedPath, err := syspath.ParseBucketPath(relPath); err == nil {
		handler.WalkBucketSetFile(absPath, relPath, parsedPath)
	}
}

func callCodeHandlers(handler systemFileHandler, absPath, relPath string) {
	if service, err := syspath.GetServiceNameFromPath(relPath); err == nil {
		handler.WalkService(absPath, relPath, service)
	} else if library, err := syspath.GetLibraryNameFromPath(relPath); err == nil {
		handler.WalkLibrary(absPath, relPath, library)
	}
}

func callCollectionHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetCollectionNameFromPath(relPath); err == nil {
		handler.WalkCollection(absPath, relPath, name)
	}
}

func callDeploymentHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetDeploymentNameFromPath(relPath); err == nil {
		handler.WalkDeployment(absPath, relPath, name)
	}
}

func callDeviceHandlers(handler systemFileHandler, absPath, relPath string) {
	if syspath.IsDeviceSchemaPath(relPath) {
		handler.WalkDeviceSchema(absPath)
	} else if name, err := syspath.GetDeviceNameFromDataPath(relPath); err == nil {
		handler.WalkDevice(absPath, relPath, name)
	} else if name, err := syspath.GetDeviceNameFromRolePath(relPath); err == nil {
		handler.WalkDeviceRole(absPath, relPath, name)
	}
}

func callEdgeHandlers(handler systemFileHandler, absPath, relPath string) {
	if syspath.IsEdgeSchemaPath(relPath) {
		handler.WalkEdgeSchema(absPath)
	} else if name, err := syspath.GetEdgeNameFromPath(relPath); err == nil {
		handler.WalkEdge(absPath, relPath, name)
	}
}

func callExternalDatabaseHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetExternalDbNameFromPath(relPath); err == nil {
		handler.WalkExternalDatabase(absPath, relPath, name)
	}
}

func callMessageHistoryStorageHandlers(handler systemFileHandler, absPath, relPath string) {
	if syspath.IsMessageHistoryStorageFile(relPath) {
		handler.WalkMessageHistoryStorage(absPath)
	}
}

func callMessageTypeTriggersHandlers(handler systemFileHandler, absPath, relPath string) {
	if syspath.IsMessageTypeTriggersFile(relPath) {
		handler.WalkMessageTypeTriggers(absPath)
	}
}

func callPluginHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetPluginNameFromPath(relPath); err == nil {
		handler.WalkPlugin(absPath, relPath, name)
	}
}

func callPortalHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetPortalNameFromPath(relPath); err == nil {
		handler.WalkPortal(absPath, relPath, name)
	} else if name, _, err := syspath.GetDatasourceNameFromPath(relPath); err == nil {
		handler.WalkPortalDatasource(absPath, relPath, name)
	} else if name, _, err := syspath.GetInternalResourceNameFromPath(relPath); err == nil {
		handler.WalkPortalInternalResources(absPath, relPath, name)
	} else if name, _, err := syspath.GetWidgetNameFromPath(relPath); err == nil {
		handler.WalkPortalWidget(absPath, relPath, name)
	} else if name, _, err := syspath.GetWidgetParserFromPath(relPath); err == nil {
		handler.WalkPortalWidgetParser(absPath, relPath, name)
	}
}

func callRoleHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetRoleNameFromPath(relPath); err == nil {
		handler.WalkRole(absPath, relPath, name)
	}
}

func callSecretHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetSecretNameFromPath(relPath); err == nil {
		handler.WalkSecret(absPath, relPath, name)
	}
}

func callServiceCacheHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetServiceCacheNameFromPath(relPath); err == nil {
		handler.WalkServiceCache(absPath, relPath, name)
	}
}

func callTimerHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetTimerNameFromPath(relPath); err == nil {
		handler.WalkTimer(absPath, relPath, name)
	}
}

func callTriggerHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetTriggerNameFromPath(relPath); err == nil {
		handler.WalkTrigger(absPath, relPath, name)
	}
}

func callUserHandlers(handler systemFileHandler, absPath, relPath string) {
	if syspath.IsUserSchemaPath(relPath) {
		handler.WalkUserSchema(absPath)
	} else if email, err := syspath.GetUserEmailFromDataPath(relPath); err == nil {
		handler.WalkUser(absPath, relPath, email)
	} else if email, err := syspath.GetUserEmailFromRolePath(relPath); err == nil {
		handler.WalkUserRole(absPath, relPath, email)
	}
}

func callWebhookHandlers(handler systemFileHandler, absPath, relPath string) {
	if name, err := syspath.GetWebhookNameFromPath(relPath); err == nil {
		handler.WalkWebhook(absPath, relPath, name)
	}
}
