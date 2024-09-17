package dryRun

import (
	"fmt"
	"strings"
)

/**
 * Generic printing for "simple" dry run section.
 * These are sections that just have a "create" and "update" list.
 */
type simpleSection struct {
	title   string
	creates []string
	updates []string
}

func newSimpleSection(title string, creates, updates []string) *simpleSection {
	return &simpleSection{
		title:   title,
		creates: creates,
		updates: updates,
	}
}

func (s *simpleSection) Title() string {
	return s.title
}

func (s *simpleSection) HasChanges() bool {
	return len(s.updates)+len(s.creates) > 0
}

func (s *simpleSection) String() string {
	sb := strings.Builder{}

	for _, create := range s.creates {
		sb.WriteString(fmt.Sprintf("Create %q\n", create))
	}

	for _, update := range s.updates {
		sb.WriteString(fmt.Sprintf("Update %q\n", update))
	}

	return sb.String()
}
