package systemUpload

import (
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/types"
)

func DoesBackendSupportSystemUpload(systemInfo *types.System_meta, client *cb.DevClient) bool {
	_, err := client.GetSystemUploadVersion(systemInfo.Key)
	return err == nil
}

func ToStringArray(val interface{}) []string {
	result := make([]string, 0)

	arr := val.([]interface{})
	for _, item := range arr {
		result = append(result, item.(string))
	}

	return result
}
