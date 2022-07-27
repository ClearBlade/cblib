package cblib

import (
	"io/ioutil"
	"os"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/resourcetree"
)

func pullFilesForAllBucketSets(systemInfo *System_meta, client *cb.DevClient) error {
	theBucketSets, err := client.GetBucketSets(systemInfo.Key)
	if err != nil {
		return err
	}

	for i := 0; i < len(theBucketSets); i++ {
		bucketSet, err := resourcetree.NewBucketSetFromMap(theBucketSets[i].(map[string]interface{}))
		if err != nil {
			return err
		}
		// empty string for boxName signifies all boxes
		err = pullFiles(systemInfo, client, bucketSet.Name, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func pullFile(systemInfo *System_meta, client *cb.DevClient, bucketSetName string, boxName string, fileName string) error {
	fileMetaMap, err := client.GetBucketSetFile(systemInfo.Key, bucketSetName, boxName, fileName)
	if err != nil {
		return err
	}

	fileMeta, err := resourcetree.NewFileMetaFromMap(fileMetaMap)
	if err != nil {
		return err
	}

	fileContents, err := client.ReadBucketSetFile(systemInfo.Key, bucketSetName, fileMeta.BucketName, fileName)
	if err != nil {
		return err
	}

	return writeBucketSetFile(bucketSetName, fileMeta, fileContents)
}

func pullFiles(systemInfo *System_meta, client *cb.DevClient, bucketSetName string, boxName string) error {
	fileMetaDict, err := client.GetBucketSetFiles(systemInfo.Key, bucketSetName, boxName)
	if err != nil {
		return err
	}

	for _, v := range fileMetaDict {
		fileMeta, err := resourcetree.NewFileMetaFromMap(v.(map[string]interface{}))
		if err != nil {
			return err
		}

		fileContents, err := client.ReadBucketSetFile(systemInfo.Key, bucketSetName, fileMeta.BucketName, fileMeta.RelativeName)
		if err != nil {
			return err
		}

		err = writeBucketSetFile(bucketSetName, fileMeta, fileContents)
		if err != nil {
			return err
		}

	}

	return nil
}

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
