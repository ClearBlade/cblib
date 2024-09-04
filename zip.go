package cblib

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
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

	fileRegex := options.GetFileRegex()
	schemaRegex := options.GetCollectionSchemaRegex()

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		// TODO: Maybe only explore directories that are cblib
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !fileRegex.MatchString(path) {
			return nil
		}

		zipPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return fmt.Errorf("could not make %s relative to %s: %w", path, rootDir, err)
		}

		return copyFileToZip(w, path, zipPath)
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
