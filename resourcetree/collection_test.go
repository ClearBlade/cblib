package resourcetree

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionFromReaderReturnsCollection(t *testing.T) {

	data := `
	{
		"name": "foobar",
		"schema": [
			{
				"ColumnName": "foo",
				"ColumnType": "int",
				"PK": true
			}
		],
		"indexes": {
			"Total": 1,
			"Data": [ {"name": "foo", "type": "Unique Index"} ]
		},
		"items": []
	}
	`

	coll, err := NewCollectionFromReader(strings.NewReader(data))
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "foobar", coll.Name)
	assert.Len(t, coll.Schema, 1)
	assert.Len(t, coll.Indexes.Data, 1)
	assert.Len(t, coll.Items, 0)
}

func TestIndexesFromNilReturnsEmptyIndexes(t *testing.T) {

	indexes, err := NewIndexesFromMap(nil)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 0, indexes.Total)
	assert.Len(t, indexes.Data, 0)
}

func TestIndexesFromEmptyMapReturnsEmptyIndexes(t *testing.T) {

	data := map[string]interface{}{
		"Data": []interface{}{},
	}

	indexes, err := NewIndexesFromMap(data)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 0, indexes.Total)
	assert.Len(t, indexes.Data, 0)
}

func TestIndexesFromOneIndexReturnsOneIndex(t *testing.T) {

	data := map[string]interface{}{
		"Data": []interface{}{
			map[string]interface{}{"name": "unique", "type": "Unique Index"},
		},
	}

	indexes, err := NewIndexesFromMap(data)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 1, indexes.Total)
	assert.Len(t, indexes.Data, 1)
}

func TestIndexesFromTwoIndexesReturnsTwoIndexes(t *testing.T) {

	data := map[string]interface{}{
		"Data": []interface{}{
			map[string]interface{}{"name": "unique", "type": "Unique Index"},
			map[string]interface{}{"name": "nonunique", "type": "Nonunique Index"},
		},
	}

	indexes, err := NewIndexesFromMap(data)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 2, indexes.Total)
	assert.Len(t, indexes.Data, 2)
}
