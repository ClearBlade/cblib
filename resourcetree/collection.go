package resourcetree

import (
	"encoding/json"
	"fmt"
	"io"

	ms "github.com/mitchellh/mapstructure"
)

// Collection stores information regarding a collection.
type Collection struct {
	Name    string                   `json:"name" mapstructure:"name"`
	Schema  []ColumnSchema           `json:"schema" mapstructure:"schema"`
	Indexes Indexes                  `json:"indexes" mapstructure:"indexes"`
	Items   []map[string]interface{} `json:"items" mapstructure:"items"`
}

// ColumnSchema stores information regarding a column.
type ColumnSchema struct {
	ColumnName string `json:"ColumnName" mapstructure:"ColumnType"`
	ColumnType string `json:"ColumnType" mapstructure:"ColumnType"`
	PK         bool   `json:"PK" mapstructure:"PK"`
}

// IndexType is one of the index types for a column (primary key, unique, non-unique).
type IndexType string

const (
	// IndexPrimaryKey is a primary key index type.
	IndexPrimaryKey IndexType = "Primary Key"

	// IndexUnique is a unique index type.
	IndexUnique = "Unique Index"

	// IndexNonUnique is a non-unique index type.
	IndexNonUnique = "Nonunique Index"
)

// Index represents an index.
type Index struct {
	Name      string    `json:"name" mapstructure:"name"`
	IndexType IndexType `json:"type" mapstructure:"type"`
}

// Indexes is a slice of Index.
type Indexes struct {
	Total int     `json:"Total" mapstructure:"total"`
	Data  []Index `json:"Data" mapstructure:"Data"`
}

// NewCollectionFromReader returns a new *Collection from the given reader.
func NewCollectionFromReader(r io.Reader) (*Collection, error) {

	coll := Collection{}

	err := json.NewDecoder(r).Decode(&coll)
	if err != nil {
		return nil, fmt.Errorf("collection decode: %s", err)
	}

	return &coll, nil
}

// NewCollectionFromMap returns a new *Collection from the given map.
func NewCollectionFromMap(m map[string]interface{}) (*Collection, error) {

	coll := Collection{}

	err := ms.Decode(m, &coll)
	if err != nil {
		return nil, fmt.Errorf("collection decode: %s", err)
	}

	return &coll, nil
}

// NewIndexesFromReader returns a new *Indexes from the given reader.
func NewIndexesFromReader(r io.Reader) (*Indexes, error) {

	indexes := Indexes{0, nil}

	err := json.NewDecoder(r).Decode(&indexes)
	if err != nil {
		return nil, fmt.Errorf("indexes decode: %s", err)
	}

	indexes.Total = len(indexes.Data)
	return &indexes, nil
}

// NewIndexesFromMap reads a map and returns a new *Indexes struct.
func NewIndexesFromMap(m map[string]interface{}) (*Indexes, error) {

	indexes := Indexes{0, nil}

	err := ms.Decode(m, &indexes)
	if err != nil {
		return nil, fmt.Errorf("indexes decode: %s", err)
	}

	indexes.Total = len(indexes.Data)
	return &indexes, nil
}
