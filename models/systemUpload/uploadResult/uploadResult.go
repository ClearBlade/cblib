package uploadResult

import (
	"errors"
	"fmt"

	"github.com/clearblade/cblib/models/systemUpload"
)

type UploadResult struct {
	Errors []string
}

func New(body interface{}) *UploadResult {
	return &UploadResult{
		Errors: systemUpload.ToStringArray(body.(map[string]interface{})["errors"]),
	}
}

func (r *UploadResult) Error() error {
	if len(r.Errors) == 0 {
		return nil
	}

	errs := make([]error, len(r.Errors))
	for i, err := range r.Errors {
		errs[i] = errors.New(err)
	}

	return fmt.Errorf("encountered the following errors while pushing: %w", errors.Join(errs...))
}
