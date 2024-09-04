package cblib

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/clearblade/cblib/syspath"
)

func getSystemZipBytes(options systemPushOptions) ([]byte, error) {
	path, err := writeSystemZip(options)
	if err != nil {
		return nil, err
	}

	defer os.Remove(path)
	return os.ReadFile(path)
}

func writeSystemZip(options systemPushOptions) (string, error) {
	if !RootDirIsSet {
		return "", fmt.Errorf("root directory is not set")
	}

	archive, err := os.CreateTemp("", "cb_cli_push_*.zip")
	if err != nil {
		return "", err
	}

	defer archive.Close()
	w := zip.NewWriter(archive)

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		zipPath, pathErr := filepath.Rel(rootDir, path)
		if pathErr != nil {
			return fmt.Errorf("could not make %s relative to %s: %w", path, rootDir, pathErr)
		}

		// Skip if we don't care about this dir
		if d.IsDir() {
			if syspath.IsClearbladePath(d.Name()) {
				return err
			}

			return filepath.SkipDir
		}

		// Skip if we don't care about this file
		if !options.ShouldPushFile(zipPath) {
			return nil
		}

		if err != nil {
			return err
		}

		// TODO: I'd rather have us pass an interface into here to call the functions
		if options.shouldPushCollectionSchemaFileOnly(zipPath) {
			return copyCollectionSchemaToZip(w, path, zipPath)
		} else {
			return copyFileToZip(w, path, zipPath)
		}
	})

	if err != nil {
		return "", fmt.Errorf("could not walk %s: %w", rootDir, err)
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	return archive.Name(), nil
}

/**
 * Prompts the user for the password before copying
 */
func copyExternalDatabaseFileToZip(w *zip.Writer, localPath string, zipPath string) error {
	return copyFileToZipWithTransform(w, localPath, zipPath, func(content []byte) ([]byte, error) {
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
		password := getOneItem(fmt.Sprintf("Password for external database '%s'", name), true)
		credentials["password"] = password
		return json.Marshal(data)
	})
}

/**
 * Removes the 'items' from a collection file before copying it so that it just contains
 * the schema.
 */
func copyCollectionSchemaToZip(w *zip.Writer, localPath string, zipPath string) error {
	return copyFileToZipWithTransform(w, localPath, zipPath, func(content []byte) ([]byte, error) {
		var data map[string]interface{}
		if err := json.Unmarshal(content, &data); err != nil {
			return nil, err
		}

		delete(data, "items")
		return json.Marshal(data)
	})
}

func copyFileToZip(w *zip.Writer, localPath string, zipPath string) error {
	return copyFileToZipWithTransform(w, localPath, zipPath, func(content []byte) ([]byte, error) {
		return content, nil
	})
}

type transformer func([]byte) ([]byte, error)

func copyFileToZipWithTransform(w *zip.Writer, localPath string, zipPath string, transform transformer) error {
	content, err := os.ReadFile(localPath)
	if err != nil {
		return err
	}

	f, err := w.Create(zipPath)
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
