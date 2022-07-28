package types

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
	MessageUrl  string
}
