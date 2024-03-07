package systemUpload

import (
	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/types"
)

func DoesBackendSupportSystemUpload(systemInfo *types.System_meta, client *cb.DevClient) bool {
	_, err := client.GetSystemUploadVersion(systemInfo.Key)
	return err == nil
}
