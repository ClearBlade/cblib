package systemUpload

import (
	"fmt"
	"strconv"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/types"
)

func DoesBackendSupportSystemUploadForCode(systemInfo *types.System_meta, client *cb.DevClient) bool {
	version, err := GetSystemUploadVersion(systemInfo, client)
	if err != nil {
		return false
	}

	return version >= 1
}

func GetSystemUploadVersion(systemInfo *types.System_meta, client *cb.DevClient) (int, error) {
	resp, err := client.GetSystemUploadVersion(systemInfo.Key)
	if err != nil {
		return 0, err
	}

	respAsMap, ok := resp.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("expected system upload response to be map, got: %T", resp)
	}

	version, ok := respAsMap["version"].(string)
	if !ok {
		return 0, fmt.Errorf("expected version to be string, got: %T", respAsMap["version"])
	}

	numericVersion, err := strconv.Atoi(version)
	if err != nil {
		return 0, fmt.Errorf("could not parse version %q as int: %w", version, err)
	}

	return numericVersion, nil
}

func ToStringArray(val interface{}) []string {
	result := make([]string, 0)

	arr := val.([]interface{})
	for _, item := range arr {
		result = append(result, item.(string))
	}

	return result
}
