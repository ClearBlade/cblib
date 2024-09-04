package syspath

import (
	"fmt"
	"regexp"
)

const (
	portalPathRegexStr             = `^portals\/([^\/]+)\/([^\/]+)\.json$`
	portalDatasourceRegexStr       = `^portals\/([^\/]+)\/config\/datasources\/([^\/]+)\/([^\/]+)\.json$`
	portalInternalResourceRegexStr = `^portals\/([^\/]+)\/config\/internalResources\/([^\/]+)\/([^\/]+)\.(?:js|json)$`
	portalWidgetRegexStr           = `^portals\/([^\/]+)\/config\/widgets\/([^\/]+)\/([^\/]+)\.json$`
	portalWidgetParserRegexStr     = `^portals\/([^\/]+)\/config\/widgets\/([^\/]+)\/parsers\/([^\/]+)\/(.+)$`
)

var (
	portalPathRegex             *regexp.Regexp
	portalDatasourceRegex       *regexp.Regexp
	portalInternalResourceRegex *regexp.Regexp
	portalWidgetRegex           *regexp.Regexp
	portalWidgetParserRegex     *regexp.Regexp
)

func init() {
	portalPathRegex = regexp.MustCompile(portalPathRegexStr)
	portalDatasourceRegex = regexp.MustCompile(portalDatasourceRegexStr)
	portalInternalResourceRegex = regexp.MustCompile(portalInternalResourceRegexStr)
	portalWidgetRegex = regexp.MustCompile(portalWidgetRegexStr)
	portalWidgetParserRegex = regexp.MustCompile(portalWidgetParserRegexStr)
}

func IsPortalPath(path string) bool {
	return topLevelDirectoryIs(path, "portals")
}

func GetPortalNameFromPath(path string) (string, error) {
	matches := portalPathRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 3 {
		return "", fmt.Errorf("path %q is not a portal path", path)
	}
	return matches[1], nil
}

func GetDatasourceNameFromPath(path string) (string, string, error) {
	matches := portalDatasourceRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("path %q is not a portal datasource path", path)
	}
	return matches[1], matches[2], nil
}

func GetInternalResourceNameFromPath(path string) (string, string, error) {
	matches := portalInternalResourceRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("path %q is not a portal internal resource path", path)
	}
	return matches[1], matches[2], nil
}

func GetWidgetNameFromPath(path string) (string, string, error) {
	matches := portalWidgetRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) != 4 {
		return "", "", fmt.Errorf("path %q is not a portal widget path", path)
	}
	return matches[1], matches[2], nil
}

func GetWidgetParserFromPath(path string) (string, string, error) {
	matches := portalWidgetParserRegex.FindStringSubmatch(path)
	if matches == nil || len(matches) < 4 {
		return "", "", fmt.Errorf("path %q is not a portal widget parser path", path)
	}
	return matches[1], matches[2], nil
}
