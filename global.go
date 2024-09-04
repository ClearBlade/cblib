package cblib

import "github.com/clearblade/cblib/types"

func setGlobalSystemDotJSONFromSystemMeta(meta *types.System_meta) {
	setGlobalSystemDotJSON(systemMetaToMap(meta))
}

func setGlobalSystemDotJSON(systemJSON map[string]interface{}) {
	systemDotJSON = systemJSON
}

func setGlobalCBMeta(cbmeta map[string]interface{}) {
	MetaInfo = cbmeta
}
