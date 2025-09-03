package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type fileStoresSection struct {
	run *cb.SystemUploadDryRun
}

func newFileStoresSection(run *cb.SystemUploadDryRun) *fileStoresSection {
	return &fileStoresSection{run: run}
}

func (a *fileStoresSection) Title() string {
	return "FILE STORES"
}

func (a *fileStoresSection) HasChanges() bool {
	return (len(a.run.FilestoresToCreate) +
		len(a.run.FilestoresToUpdate) +
		len(a.run.FilestoreFilesToCreate) +
		len(a.run.FilestoreFilesToUpdate)) > 0
}

func (a *fileStoresSection) String() string {
	sb := strings.Builder{}

	for _, bucket := range a.run.FilestoresToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", bucket))
	}

	for _, bucket := range a.run.FilestoresToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", bucket))
	}

	filesToCreate := makeFileStoreToFileMap(a.run.FilestoreFilesToCreate)
	if len(filesToCreate) > 0 {
		sb.WriteString("Files to Create:\n")
		sb.WriteString(filesToCreate.String())
	}

	filesToUpdate := makeFileStoreToFileMap(a.run.FilestoreFilesToUpdate)
	if len(filesToUpdate) > 0 {
		sb.WriteString("Files to Update:\n")
		sb.WriteString(filesToUpdate.String())
	}

	return sb.String()
}

type fileStoreToFileMap map[string][]string

func makeFileStoreToFileMap(files []cb.FileStoreFileUpdate) fileStoreToFileMap {
	fileStoreToFile := map[string][]string{}
	for _, file := range files {
		if _, ok := fileStoreToFile[file.FileStoreName]; !ok {
			fileStoreToFile[file.FileStoreName] = []string{}
		}

		fileStoreToFile[file.FileStoreName] = append(fileStoreToFile[file.FileStoreName], file.Path)
	}

	return fileStoreToFile
}

func (f *fileStoreToFileMap) String() string {
	sb := strings.Builder{}

	for fileStore, files := range *f {
		sb.WriteString(fmt.Sprintf("\tFile Store %q\n", fileStore))
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("\t\t%s\n", file))
		}
	}

	return sb.String()
}
