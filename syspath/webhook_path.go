package syspath

import (
	"fmt"
	"regexp"
)

const (
	webhookPathRegexStr = `^webhooks\/([^\/]+)\.json$`
)

var (
	webhookPathRegex *regexp.Regexp
)

func init() {
	webhookPathRegex = regexp.MustCompile(webhookPathRegexStr)
}

func IsWebhookPath(path string) bool {
	return topLevelDirectoryIs(path, "webhooks")
}

func GetWebhookNameFromPath(path string) (string, error) {
	matches := webhookPathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a webhook path", path)
	}

	return matches[1], nil
}
