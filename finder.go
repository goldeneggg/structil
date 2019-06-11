package structil

import (
	"fmt"
	"strings"
)

const (
	defaultSep  = "."
	topLevelKey = "!"
)

// Finder is the interface that builds the nested struct finder.
type Finder interface {
	FindTop(names ...string) Finder
	Find(names ...string) Finder
	Into(names ...string) Finder
	FromKeys(fks *FinderKeys) Finder
	ToMap() (map[string]interface{}, error)
	HasError() bool
	Error() string
	GetNameSeparator() string
	Reset() Finder
}

// FinderImpl is the default Finder implementation.
type FinderImpl struct {
	topLevelGetter Getter
	gMap           map[string]Getter
	fMap           map[string][]string
	eMap           map[string][]error
	ck             string
	sep            string
}

// NewFinder returns a concrete Finder that uses and obtains from i.
// i must be a struct or struct pointer.
func NewFinder(i interface{}) (Finder, error) {
	return NewFinderWithSep(i, defaultSep)
}

// NewFinderWithSep returns a concrete Finder that uses and obtains from i using the separator string.
// i must be a struct or struct pointer.
func NewFinderWithSep(i interface{}, sep string) (Finder, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return NewFinderWithGetterAndSep(g, sep)
}

// NewFinderWithGetter returns a concrete Finder that uses and obtains from g.
// g must be a Getter
func NewFinderWithGetter(g Getter) (Finder, error) {
	return NewFinderWithGetterAndSep(g, defaultSep)
}

// NewFinderWithGetterAndSep returns a concrete Finder that uses and obtains from g using the separator string.
// g must be a Getter
func NewFinderWithGetterAndSep(g Getter, sep string) (Finder, error) {
	if sep == "" {
		return nil, fmt.Errorf("sep [%s] is invalid", sep)
	}

	f := &FinderImpl{topLevelGetter: g, sep: sep}

	return f.Reset(), nil
}

// FindRoot returns a Finder that top level fields in struct are looked up and held named "names".
func (f *FinderImpl) FindTop(names ...string) Finder {
	return f.find(topLevelKey, names...)
}

// Find returns a Finder that fields in struct are looked up and held named "names".
func (f *FinderImpl) Find(names ...string) Finder {
	return f.find(f.ck, names...)
}

func (f *FinderImpl) find(fKey string, names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.fMap[fKey] = make([]string, len(names))
	copy(f.fMap[fKey], names)

	return f
}

// Into returns a Finder that nested struct fields are looked up and held named "names".
func (f *FinderImpl) Into(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.ck = topLevelKey

	var nextGetter Getter
	var ok bool
	var err error
	nextKey := ""

	for _, name := range names {
		if f.HasError() {
			break
		}

		if nextKey == "" {
			nextKey = name
		} else {
			nextKey = nextKey + f.sep + name
		}
		err = nil

		nextGetter, ok = f.gMap[nextKey]
		if !ok {
			if f.gMap[f.ck].Has(name) {
				nextGetter, err = NewGetter(f.gMap[f.ck].Get(name))
			} else {
				err = fmt.Errorf("name %s does not exist", name)
			}
		}

		if err != nil {
			f.addError(nextKey, fmt.Errorf("Error in name: %s, key: %s. [%v]", name, nextKey, err))
		}

		f.gMap[nextKey] = nextGetter
		f.ck = nextKey
	}

	return f
}

func (f *FinderImpl) addError(key string, err error) Finder {
	if _, ok := f.eMap[key]; !ok {
		f.eMap[key] = make([]error, 0, 3)
	}
	f.eMap[key] = append(f.eMap[key], err)

	return f
}

func (f *FinderImpl) FromKeys(fks *FinderKeys) Finder {
	var into, find string
	var ok bool
	m := make(map[string][]string)

	for i := 0; i < fks.Len(); i++ {
		if f.HasError() {
			return f
		}

		into, find = fks.intoAndFindNames(i)
		if _, ok = m[into]; !ok {
			m[into] = make([]string, 0, 10)
		}
		m[into] = append(m[into], find)
	}

	var is []string

	for k, v := range m {
		if k == topLevelKey {
			f.FindTop(v...)
		} else {
			is = strings.Split(k, defaultSep)
			f.Into(is...).Find(v...)
		}
	}

	return f
}

// ToMap returns a map converted from struct.
// Map keys are lookup field names by "Into" method and "Find".
// Map values are lookup field values by "Into" method and "Find".
func (f *FinderImpl) ToMap() (map[string]interface{}, error) {
	if f.HasError() {
		return nil, f
	}

	res := map[string]interface{}{}
	var key string

	for kg, getter := range f.gMap {
		for _, name := range f.fMap[kg] {
			if kg == topLevelKey {
				key = name
			} else {
				key = kg + f.sep + name
			}

			if !getter.Has(name) {
				f.addError(key, fmt.Errorf("field name %s does not exist", name))
				break
			}

			res[key] = getter.Get(name)
		}
	}

	if f.HasError() {
		return nil, f
	}

	return res, nil
}

// HasError tests whether this Finder have any errors.
func (f *FinderImpl) HasError() bool {
	for _, errs := range f.eMap {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

// Error returns error string.
func (f *FinderImpl) Error() string {
	var es []string

	for _, errs := range f.eMap {
		if len(errs) > 0 {
			es = make([]string, len(errs))
			for i, err := range errs {
				es[i] = err.Error()
			}
		}
	}

	// TODO: prettize
	return strings.Join(es, "\n")
}

// GetNameSeparator returns the separator string for nested struct name separating.
// Default is "." (dot).
func (f *FinderImpl) GetNameSeparator() string {
	return f.sep
}

// Reset resets the current build Finder.
func (f *FinderImpl) Reset() Finder {
	gMap := map[string]Getter{}
	gMap[topLevelKey] = f.topLevelGetter
	f.gMap = gMap

	fMap := map[string][]string{}
	f.fMap = fMap

	eMap := map[string][]error{}
	f.eMap = eMap

	f.ck = topLevelKey

	return f
}
