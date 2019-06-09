package structil

import (
	"fmt"
	"strings"
)

const (
	defaultSep = "."
	rootKey    = "!"
)

// Finder is the interface that builds the nested struct finder.
type Finder interface {
	Struct(names ...string) Finder
	Find(names ...string) Finder
	ToMap() (map[string]interface{}, error)
	HasError() bool
	Error() string
	GetNameSeparator() string
	Reset() Finder
}

type fImpl struct {
	rootGetter Getter
	gMap       map[string]Getter
	fMap       map[string][]string
	eMap       map[string][]error
	ck         string
	sep        string
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

	f := &fImpl{rootGetter: g, sep: sep}

	return f.Reset(), nil
}

// Reset resets the current build Finder.
func (f *fImpl) Reset() Finder {
	gMap := map[string]Getter{}
	gMap[rootKey] = f.rootGetter
	f.gMap = gMap

	fMap := map[string][]string{}
	f.fMap = fMap

	eMap := map[string][]error{}
	f.eMap = eMap

	f.ck = rootKey

	return f
}

// Struct returns a Finder that nested struct fields are looked up and held named "names".
func (f *fImpl) Struct(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.ck = rootKey

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

func (f *fImpl) addError(key string, err error) Finder {
	if _, ok := f.eMap[key]; !ok {
		f.eMap[key] = make([]error, 0, 3)
	}
	f.eMap[key] = append(f.eMap[key], err)

	return f
}

// Find returns a Finder that fields in struct are looked up and held named "names".
func (f *fImpl) Find(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.fMap[f.ck] = make([]string, len(names))
	copy(f.fMap[f.ck], names)

	return f
}

// ToMap returns a map converted from struct.
// Map keys are lookup field names by "Struct" method and "Find".
// Map values are lookup field values by "Struct" method and "Find".
func (f *fImpl) ToMap() (map[string]interface{}, error) {
	if f.HasError() {
		return nil, f
	}

	res := map[string]interface{}{}
	var key string

	for kg, getter := range f.gMap {
		for _, name := range f.fMap[kg] {
			if kg == rootKey {
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
func (f *fImpl) HasError() bool {
	for _, errs := range f.eMap {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

// Error returns error string.
func (f *fImpl) Error() string {
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
func (f *fImpl) GetNameSeparator() string {
	return f.sep
}
