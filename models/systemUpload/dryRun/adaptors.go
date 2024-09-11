package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type adaptorsSection struct {
	run *cb.SystemUploadDryRun
}

func newAdaptorsSection(run *cb.SystemUploadDryRun) *adaptorsSection {
	return &adaptorsSection{run: run}
}

func (a *adaptorsSection) Title() string {
	return "ADAPTORS"
}

func (a *adaptorsSection) HasChanges() bool {
	return (len(a.run.AdaptorFilesToCreate) +
		len(a.run.AdaptorFilesToUpdate) +
		len(a.run.AdaptorsToCreate) +
		len(a.run.AdaptorsToUpdate)) > 0
}

func (a *adaptorsSection) String() string {
	sb := strings.Builder{}

	for _, adaptor := range a.run.AdaptorsToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", adaptor))
	}

	for _, adaptor := range a.run.AdaptorsToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", adaptor))
	}

	filesToCreate := makeAdaptorToFileMap(a.run.AdaptorFilesToCreate)
	if len(filesToCreate) > 0 {
		sb.WriteString("Files to Create:\n")
		sb.WriteString(filesToCreate.String())
	}

	filesToUpdate := makeAdaptorToFileMap(a.run.AdaptorFilesToUpdate)
	if len(filesToUpdate) > 0 {
		sb.WriteString("Files to Update:\n")
		sb.WriteString(filesToUpdate.String())
	}

	return sb.String()
}

type adaptorToFileMap map[string][]string

func makeAdaptorToFileMap(files []cb.AdaptorFileUpdate) adaptorToFileMap {
	adaptorToFile := map[string][]string{}
	for _, file := range files {
		if _, ok := adaptorToFile[file.AdaptorName]; !ok {
			adaptorToFile[file.AdaptorName] = []string{}
		}

		adaptorToFile[file.AdaptorName] = append(adaptorToFile[file.AdaptorName], file.FileName)
	}

	return adaptorToFile
}

func (a *adaptorToFileMap) String() string {
	sb := strings.Builder{}

	for adaptor, files := range *a {
		sb.WriteString(fmt.Sprintf("Adaptor %q\n", adaptor))
		for _, file := range files {
			sb.WriteString(fmt.Sprintf("\t%s\n", file))
		}
	}

	return sb.String()
}
