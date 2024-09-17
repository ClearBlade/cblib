package dryRun

import (
	"fmt"
	"strings"
)

type dryRunSection interface {
	HasChanges() bool
	Title() string
	fmt.Stringer
}

func writeDryRunSection(sb *strings.Builder, section dryRunSection) {
	sb.WriteString(fmt.Sprintf("-- %s --\n", section.Title()))
	sb.WriteString(section.String())
	sb.WriteString("\n\n")
}

func writeList(sb *strings.Builder, list []string) {
	for i, item := range list {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(item)
	}
	sb.WriteString("\n")
}
