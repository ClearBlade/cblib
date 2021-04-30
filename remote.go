package cblib

import "github.com/clearblade/cblib/internal/remote"

// useRemote makes the given remote active, which implies updating the system.json
// file (system key, secret), as well as cbmeta (credentials).
func useRemote(remote *remote.Remote) {
}

func remoteTransformSystemJSON(data map[string]interface{}, remote *remote.Remote) error {
	return nil
}

func remoteTransformCBMeta(data map[string]interface{}, remote *remote.Remote) error {
	return nil
}

func remoteUpdateSystemJSON(path string, remote *remote.Remote) error {
	return nil
}

func remoteUpdateCBMeta(path string, remote *remote.Remote) error {
	return nil
}
