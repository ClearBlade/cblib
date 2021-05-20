package remote

import (
	"encoding/json"
	"os"
	"path"
)

func makeSystemJSONPath(rootDir string) string {
	return path.Join(rootDir, "system.json")
}

// systemJSONData is used for unmarshaling legacy system.json files.
type systemJSONData struct {
	SystemKeyCamel    string `json:"systemKey" yaml:"systemKey"`
	SystemKeySnake    string `json:"system_key" yaml:"system_key"`
	SystemSecretCamel string `json:"systemSecret" yaml:"systemSecret"`
	SystemSecretSnake string `json:"system_secret" yaml:"system_secret"`
}

func (s *systemJSONData) SystemKey() string {
	if len(s.SystemKeyCamel) > 0 {
		return s.SystemKeyCamel
	}
	return s.SystemKeySnake
}

func (s *systemJSONData) SystemSecret() string {
	if len(s.SystemSecretCamel) > 0 {
		return s.SystemSecretCamel
	}
	return s.SystemSecretSnake
}

func makeCBMetaPath(rootDir string) string {
	return path.Join(rootDir, hiddenDir, "cbmeta")
}

// cbmetaData is used for unmarshaling cbmeta files.
type cbmetaData struct {
	PlatformURLCamel  string `json:"platformURL" yaml:"platformURL"`
	PlatformURLSnake  string `json:"platform_url" yaml:"platform_url"`
	MessagingURLCamel string `json:"messagingURL" yaml:"messagingURL"`
	MessagingURLSnake string `json:"messaging_url" yaml:"messaging_url"`
	Token_            string `json:"token" yaml:"token"`
}

func (c *cbmetaData) PlatformURL() string {
	if len(c.PlatformURLCamel) > 0 {
		return c.PlatformURLCamel
	}
	return c.PlatformURLSnake
}

func (c *cbmetaData) MessagingURL() string {
	if len(c.MessagingURLCamel) > 0 {
		return c.MessagingURLCamel
	}
	return c.MessagingURLSnake
}

func (c *cbmetaData) Token() string {
	return c.Token_
}

// loadLegacyRemote loads a single remote from legacy folder structure.
// It will load system info from the system.json file, and credentials
// from the cbmeta file.
func loadLegacyRemote(rootDir string) (*Remote, error) {
	systemJSONPath := makeSystemJSONPath(rootDir)
	systemJSONFile, err := os.Open(systemJSONPath)
	if err != nil {
		return nil, err
	}

	systemJSON := systemJSONData{}
	err = json.NewDecoder(systemJSONFile).Decode(&systemJSON)
	if err != nil {
		return nil, err
	}

	cbmetaPath := makeCBMetaPath(rootDir)
	cbmetaFile, err := os.Open(cbmetaPath)
	if err != nil {
		return nil, err
	}

	cbmeta := cbmetaData{}
	err = json.NewDecoder(cbmetaFile).Decode(&cbmeta)
	if err != nil {
		return nil, err
	}

	remote := Remote{
		Name:         "legacy",
		PlatformURL:  cbmeta.PlatformURL(),
		MessagingURL: cbmeta.MessagingURL(),
		SystemKey:    systemJSON.SystemKey(),
		SystemSecret: systemJSON.SystemSecret(),
		Token:        cbmeta.Token(),
	}

	return &remote, nil
}

// LoadFromDirLegacy loads legacy remotes from the given directory root.
func LoadFromDirLegacy(rootDir string) (*Remotes, error) {
	legacy, err := loadLegacyRemote(rootDir)
	if err != nil {
		return nil, err
	}

	remotes := NewRemotes()
	err = remotes.Put(legacy)
	if err != nil {
		return nil, err
	}

	return remotes, nil
}

// LoadFromDirOrLegacy loads the remotes from the given directory root. If there's
// no remotes, it tries to infer remotes from the existing project.
func LoadFromDirOrLegacy(rootDir string) (*Remotes, error) {
	remotes, err := LoadFromDir(".")
	if err != nil {
		return nil, err
	}

	if remotes.Len() == 0 {
		return LoadFromDirLegacy(rootDir)
	}

	return remotes, nil
}
