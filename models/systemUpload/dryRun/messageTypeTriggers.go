package dryRun

import (
	"fmt"
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type messageTypeTriggersSection struct {
	run *cb.SystemUploadDryRun
}

func newMessageTypeTriggersSection(run *cb.SystemUploadDryRun) *messageTypeTriggersSection {
	return &messageTypeTriggersSection{run: run}
}

func (l *messageTypeTriggersSection) Title() string {
	return "MESSAGE TYPE TRIGGERS"
}

func (l *messageTypeTriggersSection) HasChanges() bool {
	return len(l.run.MessageTypeTriggers) > 0
}

func (l *messageTypeTriggersSection) String() string {
	sb := strings.Builder{}

	msgTypes := makeMessageTypeToFiltersMap(l.run.MessageTypeTriggers)
	sb.WriteString(msgTypes.String())
	return sb.String()
}

type messageTypeToFilters map[string][]string

func makeMessageTypeToFiltersMap(triggers []*cb.TriggeredMsgType) messageTypeToFilters {
	typeToFilter := map[string][]string{}
	for _, trigger := range triggers {
		if _, ok := typeToFilter[trigger.MessageType]; !ok {
			typeToFilter[trigger.MessageType] = []string{}
		}

		typeToFilter[trigger.MessageType] = append(typeToFilter[trigger.MessageType], trigger.TopicPattern)
	}

	return typeToFilter
}

func (m *messageTypeToFilters) String() string {
	sb := strings.Builder{}

	for messageType, filters := range *m {
		sb.WriteString(fmt.Sprintf("\tMessage Type %q\n", messageType))
		for _, filter := range filters {
			sb.WriteString(fmt.Sprintf("\t\t%s\n", filter))
		}
	}

	return sb.String()
}
