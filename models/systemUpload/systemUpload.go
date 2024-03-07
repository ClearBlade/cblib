package systemUpload

import (
	"encoding/base64"

	cb "github.com/clearblade/Go-SDK"
	"github.com/clearblade/cblib/internal/types"
)

const emptyZip = "UEsDBAoAAAAAADFIZ1gAAAAAAAAAAAAAAAAFABwAZW1wdHlVVAkAA83W6WXN1ulldXgLAAEE9QEAAAQUAAAAUEsBAh4DCgAAAAAAMUhnWAAAAAAAAAAAAAAAAAUAGAAAAAAAAAAAAKSBAAAAAGVtcHR5VVQFAAPN1ulldXgLAAEE9QEAAAQUAAAAUEsFBgAAAAABAAEASwAAAD8AAAAAAA=="

// TODO: We should do this a bit better
func DoesBackendSupportSystemUpload(systemInfo *types.System_meta, client *cb.DevClient) bool {
	data, err := base64.StdEncoding.DecodeString(emptyZip)
	if err != nil {
		return false
	}

	_, err = client.UploadToSystem(systemInfo.Key, data, true)
	return err == nil
}
