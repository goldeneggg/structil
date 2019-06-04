package structil

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	defaultNameSep = "."
)

type Finder interface {
	// TODO: 同一階層で複数の子要素nestを行いたいケースへの対応
	// nameをvariadic(可変) で受ける形式が良さそう
	Struct(name string) Finder
	Find(name string) Finder
	// MapConfigure(m map[string]interface{}) Finder  // TODO: mapでネストを一括設定
	From(st interface{}) (map[string]interface{}, error)
	FromGetter(g Getter) (map[string]interface{}, error)
	GetNameSeparator() string
}

type fImpl struct {
	current string
	structs []string
	finds   map[string][]string
	sep     string
}

func NewFinder() Finder {
	return NewFinderWithSep(defaultNameSep)
}

func NewFinderWithSep(sep string) Finder {
	return &fImpl{finds: map[string][]string{}, sep: sep}
}

func (f *fImpl) Struct(name string) Finder {
	f.structs = append(f.structs, name)

	f.current = strconv.Itoa(len(f.structs)-1) + name
	f.finds[f.current] = []string{}
	return f
}

// TODO: 一段もStructしてなくてもFindできるように
func (f *fImpl) Find(name string) Finder {
	if _, ok := f.finds[f.current]; ok {
		f.finds[f.current] = append(f.finds[f.current], name)
	}

	return f
}

func (f *fImpl) From(i interface{}) (map[string]interface{}, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return f.FromGetter(g)
}

// TODO: 並列化検討
func (f *fImpl) FromGetter(g Getter) (map[string]interface{}, error) {
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

func (e *fImpl) GetNameSeparator() string {
	return e.sep
}
