package cblib

var (
	URL                        string
	ShouldImportCollectionRows bool
	ImportPageSize             int
)

type Role_meta struct {
	Name        string
	Description string
	Permission  []map[string]interface{}
}

type Column struct {
	ColumnName string
	ColumnType string
}

type Collection_meta struct {
	Name          string
	Collection_id string
	Columns       []Column
}

type User_meta struct {
	Columns []Column
}

type Service_meta struct {
	Name    string
	Version int
	Hash    string
	Params  []string
}

type System_meta struct {
	Name        string
	Key         string
	Description string
	Services    map[string]Service_meta
	PlatformUrl string
}
