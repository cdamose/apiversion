package apiversion

import "time"

// Versionable represents an object that is subject to versioning.
// An object must implement Versionable in order to be affected.
type Versionable interface {
	// Data returns the object as it stands in the latest version.
	Data() map[string]interface{}
}

// Action represents an action to take on a object in order to make
// it compatible. An action takes an interface as input and returns
// an updated interface.

type Action func(map[string]interface{}) map[string]interface{}

type Change struct {
	Description string
	Actions     map[string]Action
}

type version struct {
	Date       string
	Changes    []*Change
	Deprecated bool
	date       time.Time
	layout     string
}

func (v *version) String() string {
	return v.date.Format(v.layout)
}

type versions []*version

func (vs versions) len() int           { return len(vs) }
func (vs versions) swap(i, j int)      { vs[i], vs[j] = vs[j], vs[i] }
func (vs versions) Less(i, j int) bool { return vs[i].date.Before(vs[j].date) }
