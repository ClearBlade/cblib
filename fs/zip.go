package fs

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"os"

	"github.com/clearblade/cblib/syspath"
)

type SecretPrompter interface {
	PromptForSecret(secretName string) string
}

func GetSystemZipBytes(rootDir string, prompter SecretPrompter, options *ZipOptions) ([]byte, error) {
	path, err := writeSystemZip(rootDir, prompter, options)
	if err != nil {
		return nil, err
	}

	defer os.Remove(path)
	return os.ReadFile(path)
}

func writeSystemZip(rootDir string, prompter SecretPrompter, opts *ZipOptions) (string, error) {
	if rootDir == "" {
		return "", fmt.Errorf("root directory is not set")
	}

	archive, err := os.CreateTemp("", "cb_cli_push_*.zip")
	if err != nil {
		return "", err
	}

	defer archive.Close()
	w := zip.NewWriter(archive)

	if err = walkSystemFiles(rootDir, &zipper{
		prompter: prompter,
		writer:   w,
		opts:     opts,
	}); err != nil {
		return "", fmt.Errorf("could not walk system at %s: %w", rootDir, err)
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	return archive.Name(), nil
}

type zipper struct {
	prompter SecretPrompter
	writer   *zip.Writer
	opts     *ZipOptions
}

// ----------------------
// Special Cases
// ----------------------

/**
 * The collection data and schema are stored in the same file.
 * If the user only wants to push the schema, we need to remove the data before adding to the zip
 */
func (z *zipper) WalkCollection(path, relPath string, collectionName string) {
	if z.opts.shouldPushCollectionSchemaOnly(collectionName) {
		z.copyCollectionSchemaToZip(path, relPath)
	} else if z.opts.shouldPushCollection(collectionName) {
		z.copyFileToZip(path, relPath)
	}
}

/**
 * We don't store external database passwords on disk, so we need to prompt the user for them
 */
func (z *zipper) WalkExternalDatabase(path, relPath string, externalDatabaseName string) {
	if z.opts.shouldPushExternalDatabase(externalDatabaseName) {
		z.copyExternalDatabaseFileToZip(path, relPath)
	}
}

// ----------------------
// Boring Stuff
// ----------------------

func (z *zipper) WalkAdaptor(path, relPath string, adaptorName string) {
	if z.opts.shouldPushAdaptor(adaptorName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkBucketSetMeta(path, relPath string, bucketSetName string) {
	if z.opts.shouldPushBucketSetMeta(bucketSetName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkBucketSetFile(path, relPath string, bucketFile *syspath.FullBucketPath) {
	if z.opts.shouldPushBucketSetFile(bucketFile) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkService(path, relPath string, serviceName string) {
	if z.opts.shouldPushService(serviceName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkLibrary(path, relPath string, libraryName string) {
	if z.opts.shouldPushLibrary(libraryName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkDeployment(path, relPath string, deploymentName string) {
	if z.opts.shouldPushDeployment(deploymentName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkDevice(path, relPath string, deviceName string) {
	if z.opts.shouldPushDevice(deviceName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkDeviceRole(path, relPath string, deviceName string) {
	if z.opts.shouldPushDevice(deviceName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkDeviceSchema(path string) {
	if z.opts.shouldPushDeviceSchema() {
		z.copyFileToZip(path, syspath.DeviceSchemaPath)
	}
}

func (z *zipper) WalkEdge(path, relPath string, edgeName string) {
	if z.opts.shouldPushEdge(edgeName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkEdgeSchema(path string) {
	if z.opts.shouldPushEdgeSchema() {
		z.copyFileToZip(path, syspath.EdgeSchemaPath)
	}
}

func (z *zipper) WalkMessageHistoryStorage(path string) {
	if z.opts.shouldPushMessageHistoryStorage() {
		z.copyFileToZip(path, syspath.MessageHistoryStoragePath)
	}
}

func (z *zipper) WalkMessageTypeTriggers(path string) {
	if z.opts.shouldPushMessageTypeTriggers() {
		z.copyFileToZip(path, syspath.MessageTypeTriggersPath)
	}
}

func (z *zipper) WalkPlugin(path, relPath string, pluginName string) {
	if z.opts.shouldPushPlugin(pluginName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkPortal(path, relPath string, portalName string) {
	if z.opts.shouldPushPortal(portalName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkPortalDatasource(path, relPath string, portalName string) {
	if z.opts.shouldPushPortal(portalName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkPortalInternalResources(path, relPath string, portalName string) {
	if z.opts.shouldPushPortal(portalName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkPortalWidget(path, relPath string, portalName string) {
	if z.opts.shouldPushPortal(portalName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkPortalWidgetParser(path, relPath string, portalName string) {
	if z.opts.shouldPushPortal(portalName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkRole(path, relPath string, roleName string) {
	if z.opts.shouldPushRole(roleName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkSecret(path, relPath string, secretName string) {
	if z.opts.shouldPushSecret(secretName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkServiceCache(path, relPath string, serviceCacheName string) {
	if z.opts.shouldPushServiceCache(serviceCacheName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkTimer(path, relPath string, timerName string) {
	if z.opts.shouldPushTimer(timerName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkTrigger(path, relPath string, triggerName string) {
	if z.opts.shouldPushTrigger(triggerName) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkUser(path, relPath string, email string) {
	if z.opts.shouldPushUser(email) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkUserRole(path, relPath string, email string) {
	if z.opts.shouldPushUser(email) {
		z.copyFileToZip(path, relPath)
	}
}

func (z *zipper) WalkUserSchema(path string) {
	if z.opts.shouldPushUserSchema() {
		z.copyFileToZip(path, syspath.UserSchemaPath)
	}
}

func (z *zipper) WalkWebhook(path, relPath string, webhookName string) {
	if z.opts.shouldPushWebhook(webhookName) {
		z.copyFileToZip(path, relPath)
	}
}

/**
 * Removes the 'items' from a collection file before copying it so that it just contains
 * the schema.
 */
func (z *zipper) copyCollectionSchemaToZip(localPath string, zipPath string) error {
	return z.copyFileToZipWithTransform(localPath, zipPath, func(content []byte) ([]byte, error) {
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}

		delete(data, "items")
		return json.Marshal(data)
	})
}

/**
 * Prompts the user for the password before copying
 */
func (z *zipper) copyExternalDatabaseFileToZip(localPath string, zipPath string) error {
	return z.copyFileToZipWithTransform(localPath, zipPath, func(content []byte) ([]byte, error) {
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}

		name, ok := data["name"].(string)
		if !ok {
			return nil, fmt.Errorf("external database file at %s missing name field", zipPath)
		}

		credentials, ok := data["credentials"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("external database file at %s missing credentials field", zipPath)
		}

		// TODO: Verify that this works setting the passwor dfield
		password := z.prompter.PromptForSecret(fmt.Sprintf("Password for external database '%s'", name))
		credentials["password"] = password
		return json.Marshal(data)
	})
}

func (z *zipper) copyFileToZip(localPath string, zipPath string) error {
	return z.copyFileToZipWithTransform(localPath, zipPath, func(content []byte) ([]byte, error) {
		return content, nil
	})
}

type transformer func([]byte) ([]byte, error)

func (z *zipper) copyFileToZipWithTransform(localPath string, zipPath string, transform transformer) error {
	content, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}

	f, err := z.writer.Create(zipPath)
	if err != nil {
		return err
	}

	newContent, err := transform(content)
	if err != nil {
		return fmt.Errorf("could not transform %s: %w", localPath, err)
	}

	if _, err := f.Write(newContent); err != nil {
		return err
	}

	return nil
}
