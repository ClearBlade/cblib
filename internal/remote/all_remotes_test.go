package remote

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeStubRemote(name string) *Remote {
	return &Remote{
		Name:         fmt.Sprintf("%s", name),
		PlatformURL:  fmt.Sprintf("https://%s.remote", name),
		MessagingURL: fmt.Sprintf("%s.remote:1883", name),
		SystemKey:    fmt.Sprintf("%s-system-key", name),
		SystemSecret: fmt.Sprintf("%s-system-secret", name),
		Token:        fmt.Sprintf("%s-token", name),
	}
}

func TestPutNewRemote(t *testing.T) {
	remotes := NewRemotes()

	foo := makeStubRemote("foo")
	err := remotes.Put(foo)
	require.NoError(t, err)

	assert.Equal(t, []*Remote{foo}, remotes.List())

	bar := makeStubRemote("bar")
	err = remotes.Put(bar)
	require.NoError(t, err)

	assert.Equal(t, []*Remote{bar, foo}, remotes.List())
}

func TestPutExistingRemoteUpdatesRemote(t *testing.T) {
	remotes := NewRemotes()

	foo := makeStubRemote("foo")
	err := remotes.Put(foo)
	require.NoError(t, err)

	foo.Token = "overridden-token"
	err = remotes.Put(foo)
	require.NoError(t, err)

	assert.Equal(t, []*Remote{foo}, remotes.List())
}

func TestRemoveExistingRemote(t *testing.T) {
	remotes := NewRemotes()

	foo := makeStubRemote("foo")
	err := remotes.Put(foo)
	require.NoError(t, err)

	assert.Equal(t, []*Remote{foo}, remotes.List())

	err = remotes.Remove(foo)
	require.NoError(t, err)

	assert.Equal(t, []*Remote{}, remotes.List())
}

func TestRemoveNonExistingRemoteReturnsError(t *testing.T) {
	remotes := NewRemotes()

	foo := makeStubRemote("foo")
	err := remotes.Remove(foo)

	assert.Error(t, err)
}

func TestFindByName(t *testing.T) {
	remotes := NewRemotes()

	r, ok := remotes.FindByName("foo")
	assert.Nil(t, r)
	assert.False(t, ok)

	foo := makeStubRemote("foo")
	err := remotes.Put(foo)
	require.NoError(t, err)

	r, ok = remotes.FindByName("foo")
	assert.NotNil(t, r)
	assert.True(t, ok)
}

func TestCurrent(t *testing.T) {
	foo := makeStubRemote("foo")
	bar := makeStubRemote("bar")

	var err error

	t.Run("Set to first remote", func(t *testing.T) {
		remotes := NewRemotes()

		err = remotes.Put(foo)
		require.NoError(t, err)

		err = remotes.Put(bar)
		require.NoError(t, err)

		curr, ok := remotes.Current()
		assert.True(t, ok)
		assert.Equal(t, foo, curr)
	})

	t.Run("Set to empty when only remote removed", func(t *testing.T) {
		remotes := NewRemotes()

		err = remotes.Put(foo)
		require.NoError(t, err)

		err = remotes.Remove(foo)
		require.NoError(t, err)

		curr, ok := remotes.Current()
		assert.False(t, ok)
		assert.Nil(t, curr)
	})

	t.Run("Set to another remote when removed", func(t *testing.T) {
		remotes := NewRemotes()

		err = remotes.Put(foo)
		require.NoError(t, err)

		err = remotes.Put(bar)
		require.NoError(t, err)

		err = remotes.Remove(foo)
		require.NoError(t, err)

		curr, ok := remotes.Current()
		assert.True(t, ok)
		assert.Equal(t, bar, curr)
	})
}
