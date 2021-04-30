package remote

import (
	"encoding/json"
	"os"
	"path"

	"github.com/clearblade/cblib/internal/fsutil"
)

const (
	hiddenDir = ".cb-cli"
)

func makePersistPath(rootDir string) string {
	return path.Join(rootDir, hiddenDir, "remotes")
}

// SaveToDir writes the given remotes to the given directory, overwriting any
// other remotes.
func SaveToDir(rootDir string, remotes *Remotes) error {
	err := fsutil.EnsureDirectory(rootDir)
	if err != nil {
		return err
	}

	persistPath := makePersistPath(rootDir)

	f, err := os.Create(persistPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(remotes.data)
	if err != nil {
		return err
	}

	return nil
}

// LoadFromDir loads the remotes from the given directory root. If there's no
// remotes, it returns empty remotes.
func LoadFromDir(rootDir string) (*Remotes, error) {
	err := fsutil.EnsureDirectory(rootDir)
	if err != nil {
		return nil, err
	}

	remotes := NewRemotes()

	persistPath := makePersistPath(rootDir)

	f, err := os.Open(persistPath)
	if os.IsNotExist(err) {
		return remotes, nil
	} else if err != nil {
		return nil, err
	}

	err = json.NewDecoder(f).Decode(&remotes.data)
	if err != nil {
		return nil, err
	}

	return remotes, nil
}
