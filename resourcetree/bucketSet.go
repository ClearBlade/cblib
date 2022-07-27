package resourcetree

import (
	"fmt"

	ms "github.com/mitchellh/mapstructure"
)

type BucketSet struct {
	EdgeStorage     string `json:"edge_storage" mapstructure:"edge_storage"`
	Name            string `json:"name" mapstructure:"name"`
	PlatformStorage string `json:"platform_storage" mapstructure:"platform_storage"`
	SystemKey       string `json:"system_key" mapstructure:"system_key"`
}

func NewBucketSetFromMap(m map[string]interface{}) (*BucketSet, error) {

	bucketSet := BucketSet{}

	err := ms.Decode(m, &bucketSet)
	if err != nil {
		return nil, fmt.Errorf("bucketSet decode: %s", err)
	}

	return &bucketSet, nil
}
