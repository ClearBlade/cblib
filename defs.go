package cblib

//
//  These are variables that can be used as
//  flags to a main package using this library, or
//  they can be set directly by unit tests, etc.
//  My, there are a lot of these...
//

const (
	NotExistErrorString  = "Does not exist"
	SpecialNoCBMetaError = "No cbmeta file"
)

var (
	URL                        string
	MsgURL					   string	
	SystemKey                  string
	DevToken                   string
	ShouldImportCollectionRows bool
	ImportRows                 bool
	ExportRows                 bool
	ImportUsers                bool
	ExportUsers                bool
	UserSchema                 bool
	ImportPageSize             int
	ExportPageSize             int
	Email                      string
	Password                   string
	ServiceName                string
	LibraryName                string
	CollectionName             string
	CollectionId               string
	User                       string
	UserId                     string
	RoleName                   string
	TriggerName                string
	TimerName                  string
	DeviceName                 string
	EdgeName                   string
	DashboardName              string
	PluginName                 string
	Message                    bool
	Topic                      string
	Payload                    string
	Help                       bool
	Params                     string
	Push                       bool
	AllServices                bool
	AllLibraries               bool
	AllDevices                 bool
	AllEdges                   bool
	AllDashboards              bool
	AllPlugins                 bool
)

var (
	systemDotJSON map[string]interface{}
	libCode       map[string]interface{}
	svcCode       map[string]interface{}
	rolesInfo     []map[string]interface{}
	MetaInfo      map[string]interface{}
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
	Secret      string
	Description string
	Services    map[string]Service_meta
	PlatformUrl string
	MessageUrl string
}
