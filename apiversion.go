package apiversion

import (
	"errors"
	"net/http"
	"reflect"
	"sort"
	"time"
)

var (
	ErrInvalidVersion = errors.New("Invalid Version")

	ErrVersionNotSupplied = errors.New("Version not supplied")

	ErrVersionDeprecated = errors.New("Version Deprecated")
)

const (
	defaultLayout = "1994-20-08"
	defaultHeader = "version"
	defaultQuery  = "v"
)

type VersionManager struct {
	Layout   string
	Header   string
	Query    string
	versions []*Version
}

func (vm *VersionManager) Versions() []string {
	versions := make([]string, len(vm.versions))
	for i := range versions {
		versions[i] = vm.versions[i].String()
	}

	return versions
}
func (vm *VersionManager) layout() string {
	if vm.Layout != "" {
		return vm.Layout
	}

	return defaultLayout
}

func (vm *VersionManager) header() string {
	if vm.Header != "" {
		return vm.Header
	}

	return defaultHeader
}

func (vm *VersionManager) query() string {
	if vm.Query != "" {
		return vm.Query
	}

	return defaultQuery
}

func (vm *VersionManager) Add(v *Version) error {
	var err error
	v.layout = vm.Layout
	v.date, err = time.Parse(vm.Layout, v.Date)
	if err != nil {
		return err
	}
	vm.versions = append(vm.versions, v)
	sort.Sort(sort.Reverse(versions(vm.versions)))
	return nil
}

func (vm *VersionManager) LatestVersion() *Version {
	if len(vm.versions) == 0 {
		return nil
	}
	return vm.versions[0]
}

func (vm *VersionManager) getVersionByTime(t time.Time) (*Version, error) {
	for _, v := range vm.versions {
		if v.date.Equal(t) {
			return v, nil
		}
	}
	return nil, ErrInvalidVersion
}

func (vm *VersionManager) Parse(r *http.Request) (*Version, error) {
	h := r.Header.Get(vm.header())
	q := r.URL.Query().Get(vm.query())
	if h == "" && q == "" {
		return nil, ErrVersionNotSupplied
	}
	hDate, qDate := time.Time{}, time.Time{}
	var err error
	if h != "" {
		hDate, err = time.Parse(vm.layout(), h)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}
	if q != "" {
		qDate, err = time.Parse(vm.layout(), q)
		if err != nil {
			return nil, ErrInvalidVersion
		}
	}
	t := hDate
	if hDate.Before(qDate) {
		t = qDate
	}
	v, err := vm.getVersionByTime(t)
	if v.Deprecated {
		return nil, ErrVersionDeprecated
	}
	return v, nil
}

func (vm *VersionManager) Apply(version *Version, obj Versionable) (map[string]interface{}, error) {
	data := obj.Data()
	for _, ver := range vm.versions {
		if version.date.After(ver.date) || version.date.Equal(ver.date) {
			break
		}
		for _, c := range ver.Changes {
			typ := reflect.TypeOf(obj).Elem().Name()

			// If there is an action for this obj type
			// execute the action.
			a, ok := c.Actions[typ]
			if ok {
				data = a(data)
			}
		}
	}
	return data, nil
}
