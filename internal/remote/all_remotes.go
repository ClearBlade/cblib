package remote

import (
	"fmt"
	"sort"
)

const (
	errRemoteNotPresent = "remote not present"
)

// remotesData is a private structure that contains the "state" or "data" of
// the remotes.
type remotesData struct {
	Remotes map[string]*Remote `json:"remotes" yaml:"remotes"`
	Current string             `json:"current" yaml:"current"`
}

// makeRemotesData returns a new remotesData instance.
func makeRemotesData() remotesData {
	return remotesData{
		make(map[string]*Remote),
		"",
	}
}

// Remotes for managing multiple remotes.
type Remotes struct {
	data remotesData
}

// NewRemotes returns a new remote.Remotes instance.
func NewRemotes() *Remotes {
	return &Remotes{
		makeRemotesData(),
	}
}

// updateCurrentIfNeeded makes sure the `current` invariant is valid.
func (rs *Remotes) updateCurrentIfNeeded() {
	if rs.HasByName(rs.data.Current) {
		return
	}

	remoteList := rs.List()

	if len(remoteList) == 0 {
		rs.data.Current = ""
	} else {
		rs.data.Current = remoteList[0].Name
	}
}

// Put adds the given remote to the mapping of remotes.
func (rs *Remotes) Put(remote *Remote) error {
	err := validateRemoteName(remote.Name)
	if err != nil {
		return err
	}
	rs.data.Remotes[remote.Name] = remote
	rs.updateCurrentIfNeeded()
	return nil
}

// Remove removes the given remote from the mapping.
func (rs *Remotes) Remove(remote *Remote) error {
	_, ok := rs.data.Remotes[remote.Name]
	if !ok {
		return fmt.Errorf("%s: %s", errRemoteNotPresent, remote.Name)
	}
	delete(rs.data.Remotes, remote.Name)
	rs.updateCurrentIfNeeded()
	return nil
}

// List lists all the remotes sorted by name.
func (rs *Remotes) List() []*Remote {
	result := make([]*Remote, 0, len(rs.data.Remotes))

	for _, r := range rs.data.Remotes {
		result = append(result, r)
	}

	sort.Slice(result, func(i int, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// HasByName returns whenever the given remote exists by name.
func (rs *Remotes) HasByName(name string) bool {
	_, ok := rs.data.Remotes[name]
	return ok
}

// FindByName finds a remote given the name.
func (rs *Remotes) FindByName(name string) (*Remote, bool) {
	r, ok := rs.data.Remotes[name]
	return r, ok
}

// Current returns the current remote.
func (rs *Remotes) Current() (*Remote, bool) {
	return rs.FindByName(rs.data.Current)
}

// SetCurrent sets the current remote.
// Returns error if the given remote does not exists.
func (rs *Remotes) SetCurrent(remote *Remote) error {
	if !rs.HasByName(remote.Name) {
		return fmt.Errorf("%s: %s", errRemoteNotPresent, remote.Name)
	}
	rs.data.Current = remote.Name
	return nil
}
