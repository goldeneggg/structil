package structil

import (
	"fmt"
	"strings"
)

const (
	defaultSep = "."
	rootKey    = "!"
)

type Finder interface {
	Reset() Finder
	Struct(names ...string) Finder
	Find(names ...string) Finder
	ToMap() (map[string]interface{}, error)
	HasError() bool
	Error() string
	GetNameSeparator() string
}

type fImpl struct {
	rootGetter Getter
	gMap       map[string]Getter
	fMap       map[string][]string
	eMap       map[string][]error
	ck         string
	sep        string
}

func NewFinder(i interface{}) (Finder, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return NewFinderWithGetterAndSep(g, defaultSep)
}

func NewFinderWithGetter(g Getter) (Finder, error) {
	return NewFinderWithGetterAndSep(g, defaultSep)
}

func NewFinderWithGetterAndSep(g Getter, sep string) (Finder, error) {
	if sep == "" {
		return nil, fmt.Errorf("sep [%s] is invalid", sep)
	}

	f := &fImpl{rootGetter: g, sep: sep}

	return f.Reset(), nil
}

func (f *fImpl) Reset() Finder {
	gMap := map[string]Getter{}
	gMap[rootKey] = f.rootGetter
	f.gMap = gMap

	fMap := map[string][]string{}
	fMap[rootKey] = []string{}
	f.fMap = fMap

	eMap := map[string][]error{}
	eMap[rootKey] = []error{}
	f.eMap = eMap

	f.ck = rootKey

	return f
}

func (f *fImpl) Struct(names ...string) Finder {
	if f.HasError() {
		return f
	}

	f.ck = rootKey

	var g Getter
	var err error

	for _, name := range names {
		if f.HasError() {
			break
		}

		g, err = NewGetter(f.gMap[f.ck].Get(name))

		if f.ck == rootKey {
			f.ck = name
		} else {
			f.ck = f.ck + f.sep + name
		}

		f.eMap[f.ck] = []error{}
		if err != nil {
			err = fmt.Errorf("Error in name: %s, ck: %s. [%v]", name, f.ck, err)
			f.eMap[f.ck] = append(f.eMap[f.ck], err)
		}
	}

	if !f.HasError() {
		f.gMap[f.ck] = g
		f.fMap[f.ck] = []string{}
	}

	return f
}

func (f *fImpl) Find(names ...string) Finder {
	f.fMap[f.ck] = append(f.fMap[f.ck], names...)

	return f
}

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

			res[key] = getter.Get(name)
		}
	}

	return res, nil
}

func (f *fImpl) HasError() bool {
	for _, errs := range f.eMap {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

func (f *fImpl) Error() string {
	tmp := []string{}

	for _, errs := range f.eMap {
		for _, err := range errs {
			tmp = append(tmp, err.Error())
		}
	}

	// TODO: format prettize
	return strings.Join(tmp, "\n")
}

func (f *fImpl) GetNameSeparator() string {
	return f.sep
}
