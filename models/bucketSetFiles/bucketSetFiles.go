package bucketSetFiles

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/types"
	"github.com/clearblade/cblib/resourcetree"
)

var (
	BucketSetFilesDir string
)

func PullFilesForAllBucketSets(systemInfo *types.System_meta, client *cb.DevClient) error {
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
		err = PullFiles(systemInfo, client, bucketSet.Name, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func PullFile(systemInfo *types.System_meta, client *cb.DevClient, bucketSetName string, boxName string, fileName string) error {
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

func PullFiles(systemInfo *types.System_meta, client *cb.DevClient, bucketSetName string, boxName string) error {
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
	bucketSetFileDirectory := filepath.Join(BucketSetFilesDir, bucketSetName, box, relativeDirectory)

	if err := os.MkdirAll(bucketSetFileDirectory, 0777); err != nil {
		return err
	}

	data, err := base64.StdEncoding.DecodeString(fileContents)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(bucketSetFileDirectory, fileName), data, 0666); err != nil {
		return err
	}

	return nil
}

func readBucketSetFile(bucketSetName string, boxName string, fileName string) (string, error) {
	file, err := ioutil.ReadFile(filepath.Join(BucketSetFilesDir, bucketSetName, boxName, fileName))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(file), nil
}

func PushFile(systemInfo *types.System_meta, client *cb.DevClient, bucketSetName string, boxName string, fileName string) error {
	fileContents, err := readBucketSetFile(bucketSetName, boxName, fileName)
	if err != nil {
		return err
	}
	_, err = client.CreateBucketSetFile(systemInfo.Key, bucketSetName, boxName, fileName, fileContents)

	return err
}

type pushFilesOptions struct {
	logDoesNotExistMessage bool
}

func pushFilesForBox(systemInfo *types.System_meta, client *cb.DevClient, bucketSetName string, boxName string, options pushFilesOptions) error {
	boxDirectory := path.Join(BucketSetFilesDir, bucketSetName, boxName)
	if _, err := os.Stat(boxDirectory); os.IsNotExist(err) {
		if options.logDoesNotExistMessage {
			fmt.Printf("Box '%s' does not exist on local filesystem", boxName)
		}
		return nil
	}
	filepath.WalkDir(boxDirectory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			fileName, err := filepath.Rel(boxDirectory, path)
			if err != nil {
				return err
			}
			err = PushFile(systemInfo, client, bucketSetName, boxName, fileName)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func PushFiles(systemInfo *types.System_meta, client *cb.DevClient, bucketSetName string, boxName string) error {

	if boxName != "" {
		return pushFilesForBox(systemInfo, client, bucketSetName, boxName, pushFilesOptions{logDoesNotExistMessage: true})
	} else {
		// no box specified, push files that are in every box
		err := pushFilesForBox(systemInfo, client, bucketSetName, "inbox", pushFilesOptions{logDoesNotExistMessage: false})
		if err != nil {
			return err
		}
		err = pushFilesForBox(systemInfo, client, bucketSetName, "outbox", pushFilesOptions{logDoesNotExistMessage: false})
		if err != nil {
			return err
		}
		err = pushFilesForBox(systemInfo, client, bucketSetName, "sandbox", pushFilesOptions{logDoesNotExistMessage: false})
		if err != nil {
			return err
		}
	}

	return nil
}
