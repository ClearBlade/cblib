package uploadResult

import (
	"errors"
	"fmt"

	cb "github.com/clearblade/Go-SDK"
)

type UploadResult struct {
	*cb.SystemUploadChanges
}

func New(changes *cb.SystemUploadChanges) *UploadResult {
	return &UploadResult{
		SystemUploadChanges: changes,
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
