package cblib

import (
	// "io/ioutil"
	// "os"
	// "path/filepath"

	"io/ioutil"
	"os"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/resourcetree"
)

func pullFile(systemInfo *System_meta, client *cb.DevClient, bucketSetName string, boxName string, fileName string) error {
	fileMetaMap, err := client.GetBucketSetFile(systemInfo.Key, bucketSetName, boxName, fileName)
	if err != nil {
		return err
	}

	fileContents, err := client.ReadBucketSetFile(systemInfo.Key, bucketSetName, boxName, fileName)
	if err != nil {
		return err
	}

	fileMeta, err := resourcetree.NewFileMetaFromMap(fileMetaMap)
	if err != nil {
		return err
	}

	return writeBucketSetFile(bucketSetName, fileMeta, fileContents)
}

// func pullAllFiles(systemInfo *System_meta, client *cb.DevClient, bucketSetName string) error {
// 	if files, err := client.GetBucketSetFiles(systemInfo.Key, bucketSetName, ""); err != nil {
// 		return err
// 	} else {
// 		return writeBucketSetFiles(bucketSetName, files)
// 	}
// }

// func writeBucketSetFiles(bucketSetName string, files map[string]interface{}) error {
// 	myBucketSetDir := filepath.Join(bucketSetFilesDir, bucketSetName)
// 	if err := os.MkdirAll(myBucketSetDir, 0777); err != nil {
// 		return err
// 	}

// 	for k, v := range files {

// 	}

// }

func writeBucketSetFile(bucketSetName string, fileMeta *resourcetree.FileMeta, fileContents string) error {
	box := fileMeta.BucketName
	fileName := fileMeta.BaseName

	relativeDirectory, _ := filepath.Split(fileMeta.RelativeName)
	bucketSetFileDirectory := filepath.Join(bucketSetFilesDir, bucketSetName, box, relativeDirectory)

	if err := os.MkdirAll(bucketSetFileDirectory, 0777); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(bucketSetFileDirectory, fileName), []byte(fileContents), 0666); err != nil {
		return err
	}

	return nil
}

func whitelistFileMeta(fileMeta *resourcetree.FileMeta) {

}
