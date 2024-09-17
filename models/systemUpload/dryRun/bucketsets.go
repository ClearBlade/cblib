package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type bucketSetsSection struct {
	run *cb.SystemUploadDryRun
}

func newBucketSetsSection(run *cb.SystemUploadDryRun) *bucketSetsSection {
	return &bucketSetsSection{run: run}
}

func (a *bucketSetsSection) Title() string {
	return "BUCKET SETS"
}

func (a *bucketSetsSection) HasChanges() bool {
	return (len(a.run.BucketsToCreate) +
		len(a.run.BucketsToUpdate) +
		len(a.run.BucketFilesToCreate) +
		len(a.run.BucketFilesToUpdate)) > 0
}

func (a *bucketSetsSection) String() string {
	sb := strings.Builder{}

	for _, bucket := range a.run.BucketsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", bucket))
	}

	for _, bucket := range a.run.BucketFilesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", bucket))
	}

	filesToCreate := makeBucketToFileMap(a.run.BucketFilesToCreate)
	if len(filesToCreate) > 0 {
		sb.WriteString("Files to Create:\n")
		sb.WriteString(filesToCreate.String())
	}

	filesToUpdate := makeBucketToFileMap(a.run.BucketFilesToUpdate)
	if len(filesToUpdate) > 0 {
		sb.WriteString("Files to Update:\n")
		sb.WriteString(filesToUpdate.String())
	}

	return sb.String()
}

type bucketToFileMap map[string][]string

func makeBucketToFileMap(files []cb.BucketFileUpdate) bucketToFileMap {
	bucketToFile := map[string][]string{}
	for _, file := range files {
		if _, ok := bucketToFile[file.BucketName]; !ok {
			bucketToFile[file.BucketName] = []string{}
		}

		filename := fmt.Sprintf("%s/%s", file.BucketBox, file.RelativePath)
		bucketToFile[file.BucketName] = append(bucketToFile[file.BucketName], filename)
	}

	return bucketToFile
}

func (a *bucketToFileMap) String() string {
	sb := strings.Builder{}

	for adaptor, files := range *a {
		sb.WriteString(fmt.Sprintf("\tBucket %q\n", adaptor))
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("\t\t%s\n", file))
		}
	}

	return sb.String()
}
