package dryRun

import (
	"fmt"
	"strings"
)

func writeDryRunSection(sb *strings.Builder, title, contents string) {
	sb.WriteString(fmt.Sprintf("-- %s --\n", title))
	sb.WriteString(contents)
	sb.WriteString("\n\n")
}
