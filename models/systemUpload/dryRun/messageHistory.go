package dryRun

import (
	"strings"

	cb "github.com/clearblade/Go-SDK"
)

type messageHistorySection struct {
	run *cb.SystemUploadDryRun
}

func newMessageHistorySection(run *cb.SystemUploadDryRun) *messageHistorySection {
	return &messageHistorySection{run: run}
}

func (l *messageHistorySection) Title() string {
	return "MESSAGE HISTORY"
}

func (l *messageHistorySection) HasChanges() bool {
	return len(l.run.MessageHistoryStorageTopics) > 0
}

func (l *messageHistorySection) String() string {
	sb := strings.Builder{}

	if len(l.run.MessageHistoryStorageTopics) > 0 {
		sb.WriteString("Message History Storage Topics: ")
		writeList(&sb, l.run.MessageHistoryStorageTopics)
	}

	return sb.String()
}
