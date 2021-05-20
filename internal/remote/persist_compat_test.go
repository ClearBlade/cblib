package remote

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromLegacy(t *testing.T) {
	legacy, err := loadLegacyRemote("./testdata/legacy")
	require.NoError(t, err)

	remotes := NewRemotes()
	err = remotes.Put(legacy)
	require.NoError(t, err)

	tempdir := t.TempDir()
	os.MkdirAll(path.Join(tempdir, ".cb-cli"), os.ModePerm)
	err = SaveToDir(tempdir, remotes)
	require.NoError(t, err)

	assertPersistedRemotesEqual(t, tempdir, "./testdata/legacy-remotes")
}
