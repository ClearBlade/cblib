package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/models/systemUpload"
	"github.com/clearblade/cblib/types"
)

type DryRun struct {
	Errors            []string
	LibrariesToCreate []string
	LibrariesToUpdate []string
	ServicesToCreate  []string
	ServicesToUpdate  []string
}

func New(systemInfo *types.System_meta, client *cb.DevClient, buffer []byte) (*DryRun, error) {
	dryRunResult, err := client.UploadToSystem(systemInfo.Key, buffer, true)
	if err != nil {
		return nil, err
	}

	run := dryRunResult.(map[string]interface{})
	return &DryRun{
		Errors:            systemUpload.ToStringArray(run["errors"]),
		LibrariesToCreate: systemUpload.ToStringArray(run["libraries_to_create"]),
		LibrariesToUpdate: systemUpload.ToStringArray(run["libraries_to_update"]),
		ServicesToCreate:  systemUpload.ToStringArray(run["services_to_create"]),
		ServicesToUpdate:  systemUpload.ToStringArray(run["services_to_update"]),
	}, nil
}

func (d *DryRun) String() string {
	if len(d.Errors) > 0 {
		return fmt.Sprintf("cannot push: %s", strings.Join(d.Errors, "\n"))
	}

	if !d.HasChanges() {
		return ""
	}

	sb := strings.Builder{}
	sb.WriteString("The following changes will be made:\n")
	if d.hasServiceChanges() {
		writeDryRunSection(&sb, "SERVICES", d.ServicesString())
	}

	if d.hasLibraryChanges() {
		writeDryRunSection(&sb, "LIBRARIES", d.LibrariesString())
	}

	return sb.String()
}

func (d *DryRun) HasChanges() bool {
	if len(d.Errors) > 0 {
		return false
	}

	return d.hasServiceChanges() || d.hasLibraryChanges()
}

func (d *DryRun) HasErrors() bool {
	return len(d.Errors) > 0
}

// ----------------------- SERVICES -----------------------
func (d *DryRun) ServicesString() string {
	sb := strings.Builder{}

	for _, service := range d.ServicesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", service))
	}

	for _, service := range d.ServicesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", service))
	}

	return sb.String()
}

func (d *DryRun) hasServiceChanges() bool {
	return len(d.ServicesToCreate)+len(d.ServicesToUpdate) > 0
}

// ----------------------- LIBRARIES -----------------------
func (d *DryRun) LibrariesString() string {
	sb := strings.Builder{}

	for _, library := range d.LibrariesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", library))
	}

	for _, library := range d.LibrariesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", library))
	}

	return sb.String()
}

func (d *DryRun) hasLibraryChanges() bool {
	return len(d.LibrariesToCreate)+len(d.LibrariesToUpdate) > 0
}
