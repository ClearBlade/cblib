package remote

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertPersistedRemotesEqual(t *testing.T, rootA, rootB string) {
	pathA := path.Join(rootA, ".cb-cli", "remotes")
	dataA, err := ioutil.ReadFile(pathA)
	require.NoError(t, err)

	pathB := path.Join(rootB, ".cb-cli", "remotes")
	dataB, err := ioutil.ReadFile(pathB)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(dataA, dataB), "remotes are not equal: data(%s) != data(%s)", pathA, pathB)
}

func TestSaveToDir(t *testing.T) {
	tempdir := t.TempDir()
	os.MkdirAll(path.Join(tempdir, ".cb-cli"), os.ModePerm)

	var err error

	remotes := NewRemotes()

	// foo

	foo := makeStubRemote("foo")
	remotes.Put(foo)
	err = SaveToDir(tempdir, remotes)
	require.NoError(t, err)

	assertPersistedRemotesEqual(t, tempdir, "./testdata/foo")

	// foo, bar

	fooBar := makeStubRemote("bar")
	remotes.Put(fooBar)
	err = SaveToDir(tempdir, remotes)
	require.NoError(t, err)

	assertPersistedRemotesEqual(t, tempdir, "./testdata/foo-bar")
}

func TestLoadFromDir(t *testing.T) {
	var ok bool

	// empty

	empty, err := LoadFromDir("./testdata/empty")
	require.NoError(t, err)

	assert.Len(t, empty.List(), 0)

	// foo

	foo, err := LoadFromDir("./testdata/foo")
	require.NoError(t, err)

	assert.Len(t, foo.List(), 1)
	_, ok = foo.FindByName("foo")
	assert.True(t, ok)

	// foo, bar

	fooBar, err := LoadFromDir("./testdata/foo-bar")
	require.NoError(t, err)

	assert.Len(t, fooBar.List(), 2)
	_, ok = fooBar.FindByName("foo")
	assert.True(t, ok)
	_, ok = fooBar.FindByName("bar")
	assert.True(t, ok)
}
