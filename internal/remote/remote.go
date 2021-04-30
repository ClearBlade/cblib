// package remote helps you manage multiple remote endponds to ClearBlade instances.
//
// Most manipulation of remotes happen in-memory, with the option to persist the
// whole data-structure to disk if needed.
package remote

import (
	"fmt"
	"regexp"
)

type Remote struct {
	Name         string
	PlatformURL  string
	MessagingURL string
	SystemKey    string
	SystemSecret string
	Token        string
}

func validateRemoteName(name string) error {
	re := regexp.MustCompile(`^[a-z]+[a-z0-9_-]*[a-z0-9]$`)

	if !re.MatchString(name) {
		return fmt.Errorf("remote name must match regex: %s", re)
	}

	return nil
}
