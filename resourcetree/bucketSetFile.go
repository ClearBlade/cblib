package resourcetree

import (
	"fmt"

	ms "github.com/mitchellh/mapstructure"
)

type FileMeta struct {
	PathName     string `json:"path_name" mapstructure:"path_name"`
	RelativeName string `json:"relative_name" mapstructure:"relative_name"`
	BaseName     string `json:"base_name" mapstructure:"base_name"`
	BucketName   string `json:"bucket_name" mapstructure:"bucket_name"`
	Size         int    `json:"size" mapstructure:"size"`
	Permissions  string `json:"permissions" mapstructure:"permissions"`
	LastModified string `json:"last_modified" mapstructure:"last_modified"`
}

func NewFileMetaFromMap(m map[string]interface{}) (*FileMeta, error) {

	fileMeta := FileMeta{}

	err := ms.Decode(m, &fileMeta)
	if err != nil {
		return nil, fmt.Errorf("fileMeta decode: %s", err)
	}

	return &fileMeta, nil
}
