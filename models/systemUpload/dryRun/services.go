package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type servicesSection struct {
	run *cb.SystemUploadDryRun
}

func newServicesSection(run *cb.SystemUploadDryRun) *servicesSection {
	return &servicesSection{run: run}
}

func (s *servicesSection) Title() string {
	return "SERVICES"
}

func (s *servicesSection) HasChanges() bool {
	return len(s.run.ServicesToCreate)+len(s.run.ServicesToUpdate) > 0
}

func (s *servicesSection) String() string {
	sb := strings.Builder{}

	for _, service := range s.run.ServicesToCreate {
		sb.WriteString(fmt.Sprintf("Create %q\n", service))
	}

	for _, service := range s.run.ServicesToUpdate {
		sb.WriteString(fmt.Sprintf("Update %q\n", service))
	}

	return sb.String()
}
