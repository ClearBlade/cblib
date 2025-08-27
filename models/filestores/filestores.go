package filestores

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
)

var (
	FileStoresFilesDir string
)

func PullFile(client *cb.DevClient, systemKey, fileStore, fileName string) error {
	contents, err := client.ReadFilestoreFile(systemKey, fileStore, fileName)
	if err != nil {
		return err
	}

	return writeFileStoreFile(fileStore, fileName, contents)
}

func PullFiles(client *cb.DevClient, systemKey, fileStore string) error {
	fmt.Printf("Pulling files for file store %s\n", fileStore)
	opts := cb.ListOptions{
		Depth: -1,
		Limit: 100,
	}

	for {
		files, err := client.ListFilestoreFiles(systemKey, fileStore, &opts)
		if err != nil {
			return err
		}

		if err := writeFileStoreFiles(client, systemKey, fileStore, files.Files); err != nil {
			return err
		}

		opts.ContinuationToken = files.ContinuationToken
		if opts.ContinuationToken == "" {
			return nil
		}
	}
}

func PullFilesForAllFileStores(client *cb.DevClient, systemKey string) error {
	fileStores, err := client.GetFilestores(systemKey)
	if err != nil {
		return err
	}

	for _, fileStore := range fileStores {
		if err := PullFiles(client, systemKey, fileStore.Name); err != nil {
			return err
		}
	}

	return nil
}

func writeFileStoreFiles(client *cb.DevClient, systemKey, fileStore string, files []cb.FileMeta) error {
	for _, file := range files {
		if file.IsDir {
			continue
		}

		if err := PullFile(client, systemKey, fileStore, file.FullPath); err != nil {
			return err
		}
	}

	return nil
}

func writeFileStoreFile(fileStore string, path string, contents []byte) error {
	filePath := filepath.Join(FileStoresFilesDir, fileStore, path)
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, 0777); err != nil {
		return err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err := io.Copy(file, bytes.NewReader(contents)); err != nil {
		return err
	}

	return nil
}
