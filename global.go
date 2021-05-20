package cblib

func setGlobalSystemDotJSONFromSystemMeta(meta *System_meta) {
	setGlobalSystemDotJSON(systemMetaToMap(meta))
}

func setGlobalSystemDotJSON(systemJSON map[string]interface{}) {
	systemDotJSON = systemJSON
}

func setGlobalCBMeta(cbmeta map[string]interface{}) {
	MetaInfo = cbmeta
}
