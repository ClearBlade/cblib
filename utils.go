package cblib

import (
	cb "github.com/clearblade/Go-SDK"
	"strings"
)

func setupAddrs(paddr string) {
	cb.CB_ADDR = paddr
	preIdx := strings.Index(paddr, "://")
	baseAddress := paddr[preIdx+3:]
	postIdx := strings.Index(baseAddress, ":")
	if postIdx != -1 {
		baseAddress = baseAddress[:postIdx]
	}
	cb.CB_MSG_ADDR = baseAddress + ":1883"
}
