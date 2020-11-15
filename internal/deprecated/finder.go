package deprecated

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultSep  = "."
	topLevelKey = "!"
)

// // Finder is the interface that builds the nested struct finder.
// type Finder interface {
// 	FindTop(names ...string) Finder
// 	Find(names ...string) Finder
// 	Into(names ...string) Finder
// 	FromKeys(fks *FinderKeys) Finder
// 	ToMap() (map[string]interface{}, error)
// 	HasError() bool
// 	Error() string
// 	GetNameSeparator() string
// 	Reset() Finder
// }

// Finder is the struct that builds the nested struct finder.
type Finder struct {
	topLevelGetter *Getter
	gMap           map[string]*Getter
	fMap           map[string][]string
	eMap           map[string][]error
	ck             string
	sep            string
}

// NewFinder returns a concrete Finder that uses and obtains from i.
// i must be a struct or struct pointer.
func NewFinder(i interface{}) (*Finder, error) {
	return NewFinderWithSep(i, defaultSep)
}

// NewFinderWithSep returns a concrete Finder that uses and obtains from i using the separator string.
// i must be a struct or struct pointer.
func NewFinderWithSep(i interface{}, sep string) (*Finder, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return NewFinderWithGetterAndSep(g, sep)
}

// NewFinderWithGetter returns a concrete Finder that uses and obtains from g.
// g must be a Getter
func NewFinderWithGetter(g *Getter) (*Finder, error) {
	return NewFinderWithGetterAndSep(g, defaultSep)
}

// NewFinderWithGetterAndSep returns a concrete Finder that uses and obtains from g using the separator string.
// g must be a Getter
func NewFinderWithGetterAndSep(g *Getter, sep string) (*Finder, error) {
	if sep == "" {
		return nil, fmt.Errorf("sep [%s] is invalid", sep)
	}

	f := &Finder{topLevelGetter: g, sep: sep}

	return f.Reset(), nil
}

// FindTop returns a Finder that top level fields in struct are looked up and held named names.
// Deprecated: planning to remove this method.
func (f *Finder) FindTop(names ...string) *Finder {
	return f.find(topLevelKey, names...)
}

// Find returns a Finder that fields in struct are looked up and held named names.
func (f *Finder) Find(names ...string) *Finder {
	return f.find(f.ck, names...)
}

func (f *Finder) find(fKey string, names ...string) *Finder {
	if f.HasError() {
		return f
	}

	f.fMap[fKey] = make([]string, len(names))
	copy(f.fMap[fKey], names)

	return f
}

// Into returns a Finder that nested struct fields are looked up and held named names.
func (f *Finder) Into(names ...string) *Finder {
	if f.HasError() {
		return f
	}

	f.ck = topLevelKey

	var nextGetter *Getter
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

func (f *Finder) addError(key string, err error) *Finder {
	if _, ok := f.eMap[key]; !ok {
		f.eMap[key] = make([]error, 0, 3)
	}
	f.eMap[key] = append(f.eMap[key], err)

	return f
}

// FromKeys returns a Finder that looked up by FinderKeys generated from configuration file.
func (f *Finder) FromKeys(fks *FinderKeys) *Finder {
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
func (f *Finder) ToMap() (map[string]interface{}, error) {
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
func (f *Finder) HasError() bool {
	for _, errs := range f.eMap {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

// Error returns error string.
func (f *Finder) Error() string {
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
func (f *Finder) GetNameSeparator() string {
	return f.sep
}

// Reset resets the current build Finder.
func (f *Finder) Reset() *Finder {
	gMap := map[string]*Getter{}
	gMap[topLevelKey] = f.topLevelGetter
	f.gMap = gMap

	fMap := map[string][]string{}
	f.fMap = fMap

	eMap := map[string][]error{}
	f.eMap = eMap

	f.ck = topLevelKey

	return f
}

// FinderKeys is the struct that have keys for Finder.
type FinderKeys struct {
	keys []string
}

type confKeys struct {
	Keys []interface{}
}

// NewFinderKeys returns a FinderKeys object
// that is created from configuration file indicated by dir and name file.
func NewFinderKeys(dir string, baseName string) (*FinderKeys, error) {
	viper.SetConfigName(baseName)
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var ck confKeys
	if err := viper.Unmarshal(&ck); err != nil {
		return nil, err
	}

	return newFinderKeysFromConf(ck)
}

func newFinderKeysFromConf(ck confKeys) (*FinderKeys, error) {
	if len(ck.Keys) == 0 {
		return nil, fmt.Errorf("failed to parse or no keys exist in file")
	}

	fks := &FinderKeys{keys: make([]string, 0, len(ck.Keys)+1)}

	var err error
	for _, ckk := range ck.Keys {
		err = fks.addRecursive(ckk, "")
		if err != nil {
			return nil, err
		}
	}

	return fks, nil
}

func (fks *FinderKeys) addRecursive(key interface{}, prefix string) error {
	var res string

	switch t := key.(type) {
	case string:
		res = t
		if prefix != "" {
			res = prefix + defaultSep + res
		}
		fks.keys = append(fks.keys, res) // set here
	case map[string]interface{}:
		var nk string
		var nd interface{}
		for key, value := range t {
			nk = key
			if prefix != "" {
				nk = prefix + "." + nk
			}
			nd = value
			break
		}
		err := fks.addRecursive(nd, nk)
		if err != nil {
			return err
		}
	case map[interface{}]interface{}:
		var nk string
		var nd interface{}
		for key, value := range t {
			nk = key.(string)
			if prefix != "" {
				nk = prefix + "." + nk
			}
			nd = value
			break
		}
		err := fks.addRecursive(nd, nk)
		if err != nil {
			return err
		}
	case []interface{}:
		var err error
		for _, value := range t {
			err = fks.addRecursive(value, prefix)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported type: %#v, prefix: %s", t, prefix)
	}

	return nil
}

// Len returns length of FinderKeys
func (fks *FinderKeys) Len() int {
	return len(fks.keys)
}

// Keys returns keys
func (fks *FinderKeys) Keys() []string {
	return fks.keys
}

func (fks *FinderKeys) intoAndFindNames(i int) (string, string) {
	s := strings.Split(fks.keys[i], defaultSep)
	if len(s) == 1 {
		return topLevelKey, s[0]
	}

	return strings.Join(s[0:len(s)-1], defaultSep), s[len(s)-1]
}
