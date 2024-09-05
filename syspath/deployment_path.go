package syspath

import (
	"fmt"
	"regexp"
)

const (
	deploymentPathRegexStr = `^deployments\/([^\/]+)\.json$`
)

var (
	deploymentPathRegex *regexp.Regexp
)

func init() {
	deploymentPathRegex = regexp.MustCompile(deploymentPathRegexStr)
}

func IsDeploymentPath(path string) bool {
	return topLevelDirectoryIs(path, "deployments")
}

func GetDeploymentNameFromPath(path string) (string, error) {
	matches := deploymentPathRegex.FindStringSubmatch(path)
	if len(matches) != 2 {
		return "", fmt.Errorf("path %q is not a deployment path", path)
	}

	return matches[1], nil
}
