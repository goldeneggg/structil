package structil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	defaultNameSep = "."
	rootKey        = "!"
)

type Finder interface {
	// TODO: 同一階層で複数の子要素nestを行いたいケースへの対応
	// nameをvariadic(可変) で受ける形式が良さそう
	//Struct(name string) Finder
	Struct(names ...string) Finder
	Find(names ...string) Finder
	// MapConfigure(m map[string]interface{}) Finder  // TODO: mapでネストを一括設定
	From(st interface{}) (map[string]interface{}, error)
	FromGetter(g Getter) (map[string]interface{}, error)
	HasError() bool
	Error() string
	GetNameSeparator() string
}

type fImpl struct {
	getters map[string]Getter
	finds   map[string][]string
	errs    map[string][]error
	ck      string
	sep     string
}

func NewFinder(i interface{}) (Finder, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return NewFinderWithGetterAndSep(g, defaultNameSep)
}

func NewFinderWithGetter(g Getter) (Finder, error) {
	return NewFinderWithGetterAndSep(g, defaultNameSep)
}

func NewFinderWithGetterAndSep(g Getter, sep string) (Finder, error) {
	getters := map[string]Getter{}
	getters[rootKey] = g

	finds := map[string][]string{}
	finds[rootKey] = []string{}

	errs := map[string][]error{}
	errs[rootKey] = []error{}

	return &fImpl{getters: getters, finds: finds, errs: errs, ck: rootKey, sep: sep}, nil
}

func (f *fImpl) Struct(names ...string) Finder {
	if f.HasError() {
		return f
	}

	var g Getter
	var err error

	f.ck = rootKey

	for _, name := range names {
		if errs, _ := f.errs[f.ck]; len(errs) > 0 {
			break
		}

		g, err = NewGetter(f.getters[f.ck].Get(name))

		f.ck = strings.Join(names, f.sep)

		f.errs[f.ck] = []error{}
		if err != nil {
			f.errs[f.ck] = append(f.errs[f.ck], err)
		}

		f.getters[f.ck] = g
		f.finds[f.ck] = []string{}
	}

	return f
}

func (f *fImpl) Find(names ...string) Finder {
	f.finds[f.ck] = append(f.finds[f.ck], names...)

	return f
}

func (f *fImpl) From(i interface{}) (map[string]interface{}, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return f.FromGetter(g)
}

func (f *fImpl) FromGetter(g Getter) (res map[string]interface{}, err error) {
	var ck string
	var ng Getter

	cGetters := map[string]Getter{}
	res = map[string]interface{}{}

	for sName, targets := range f.finds {
		ck = ""

		for _, n := range strings.Split(sName, f.sep) {
			ng, err = NewGetter(g.Get(n))
			if err != nil {
				return
			}

			cGetters
		}

	}

	return
}

func (f *fImpl) HasError() bool {
	for _, errs := range f.errs {
		if len(errs) > 0 {
			return true
		}
	}

	return false
}

func (f *fImpl) Error() string {
	tmp := []string{}

	for _, errs := range f.errs {
		for _, err := range errs {
			tmp = append(tmp, err.Error())
		}
	}

	return strings.Join(tmp, "\n")
}

func (f *fImpl) GetNameSeparator() string {
	return f.sep
}

// TODO: 並列化検討
func (f *fImpl) FromGetterOld(g Getter) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	var prefix string
	var resIntf interface{}

	var nested interface{}
	var typ reflect.Type
	var kind reflect.Kind
	var pk string
	var pvs []string
	var ok bool

	var err error

	for idx, n := range f.structs {
		// get nested struct
		if !g.IsStruct(n) {
			return nil, fmt.Errorf("struct named %s does not exist in %+v", n, g)
		}

		nested = g.Get(n)

		typ = reflect.TypeOf(nested)
		kind = typ.Kind()
		if kind != reflect.Struct {
			return nil, fmt.Errorf("name %s is not struct: %+v", n, nested)
		}

		g, err = NewGetter(nested)
		if err != nil {
			return nil, err
		}

		// retrieve items from nested struct
		pk = strconv.Itoa(idx) + n
		pvs, ok = f.finds[pk]
		if !ok {
			continue
		}

		if prefix == "" {
			prefix = n
		} else {
			prefix = prefix + f.sep + n
		}
		for _, pv := range pvs {
			resIntf = g.Get(pv)
			res[prefix+f.sep+pv] = resIntf
		}
	}

	return res, nil
}

func (f *fImpl) StructOld(name string) Finder {
	f.structs = append(f.structs, name)

	f.current = strconv.Itoa(len(f.structs)-1) + name
	f.finds[f.current] = []string{}
	return f
}

// TODO: 一段もStructしてなくてもFindできるように
func (f *fImpl) FindOld(name string) Finder {
	if _, ok := f.finds[f.current]; ok {
		f.finds[f.current] = append(f.finds[f.current], name)
	}

	return f
}

// TODO: 並列化検討
func (f *fImpl) FromGetterOld(g Getter) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	var prefix string
	var resIntf interface{}

	var nested interface{}
	var typ reflect.Type
	var kind reflect.Kind
	var pk string
	var pvs []string
	var ok bool

	var err error

	for idx, n := range f.structs {
		// get nested struct
		if !g.IsStruct(n) {
			return nil, fmt.Errorf("struct named %s does not exist in %+v", n, g)
		}

		nested = g.Get(n)

		typ = reflect.TypeOf(nested)
		kind = typ.Kind()
		if kind != reflect.Struct {
			return nil, fmt.Errorf("name %s is not struct: %+v", n, nested)
		}

		g, err = NewGetter(nested)
		if err != nil {
			return nil, err
		}

		// retrieve items from nested struct
		pk = strconv.Itoa(idx) + n
		pvs, ok = f.finds[pk]
		if !ok {
			continue
		}

		if prefix == "" {
			prefix = n
		} else {
			prefix = prefix + f.sep + n
		}
		for _, pv := range pvs {
			resIntf = g.Get(pv)
			res[prefix+f.sep+pv] = resIntf
		}
	}

	return res, nil
}

func (f *fImpl) StructOld(name string) Finder {
	f.structs = append(f.structs, name)

	f.current = strconv.Itoa(len(f.structs)-1) + name
	f.finds[f.current] = []string{}
	return f
}

// TODO: 一段もStructしてなくてもFindできるように
func (f *fImpl) FindOld(name string) Finder {
	if _, ok := f.finds[f.current]; ok {
		f.finds[f.current] = append(f.finds[f.current], name)
	}

	return f
}
