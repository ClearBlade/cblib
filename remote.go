package cblib

import "github.com/clearblade/cblib/internal/remote"

// useRemoteByMerging makes the given remote active, which implies updating the
// system.json file (system key, secret), as well as cbmeta (credentials). Ideally,
// we would use the remote directly, but there's a lot of code scattered around that
// depends on the aforementioned files.
func useRemoteByMerging(systemJSON, cbmeta map[string]interface{}, remote *remote.Remote) error {
	var err error

	err = remoteTransformSystemJSON(systemJSON, remote)
	if err != nil {
		return err
	}
	err = remoteTransformCBMeta(cbmeta, remote)
	if err != nil {
		return err
	}

	setGlobalSystemDotJSON(systemJSON)
	setGlobalCBMeta(cbmeta)

	err = storeSystemDotJSON(systemJSON)
	if err != nil {
		return err
	}
	err = storeCBMeta(cbmeta)
	if err != nil {
		return err
	}

	return nil
}

// useRemoteByMergingFromFlobals is simular to useRemoteByMerging, but obtains
// the metadata from the global state.
func useRemoteByMergingFromGlobals(remote *remote.Remote) error {
	systemMeta, err := getSysMeta()
	if err != nil {
		return err
	}
	systemJSON := systemMetaToMap(systemMeta)

	cbmeta, err := getCbMeta()
	if err != nil {
		return err
	}

	return useRemoteByMerging(systemJSON, cbmeta, remote)
}

func remoteTransformSystemJSON(data map[string]interface{}, remote *remote.Remote) error {
	data["platform_url"] = remote.PlatformURL
	data["messaging_url"] = remote.MessagingURL
	data["system_key"] = remote.SystemKey
	data["system_secret"] = remote.SystemSecret
	return nil
}

func remoteTransformCBMeta(data map[string]interface{}, remote *remote.Remote) error {
	data["developer_email"] = ""
	data["platform_url"] = remote.PlatformURL
	data["token"] = remote.Token
	return nil
}
