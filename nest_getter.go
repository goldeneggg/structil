package structil

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	defaultNameSep = "."
)

type NestGetter interface {
	Nest(name string) NestGetter
	Want(name string) NestGetter
	From(st interface{}) (map[string]interface{}, error)
	FromGetter(ac Getter) (map[string]interface{}, error)
	GetNameSeparator() string
}

type ngImpl struct {
	current string
	nests   []string
	wants   map[string][]string
	sep     string
}

func NewNestGetter() NestGetter {
	return NewNestGetterWithSep(defaultNameSep)
}

func NewNestGetterWithSep(sep string) NestGetter {
	return &ngImpl{wants: map[string][]string{}, sep: sep}
}

func (ng *ngImpl) Nest(name string) NestGetter {
	ng.nests = append(ng.nests, name)

	ng.current = strconv.Itoa(len(ng.nests)-1) + name
	ng.wants[ng.current] = []string{}
	return ng
}

func (ng *ngImpl) Want(name string) NestGetter {
	if _, ok := ng.wants[ng.current]; ok {
		ng.wants[ng.current] = append(ng.wants[ng.current], name)
	}

	return ng
}

func (ng *ngImpl) From(st interface{}) (map[string]interface{}, error) {
	g, err := NewGetter(st)
	if err != nil {
		return nil, err
	}

	return ng.FromGetter(g)
}

func (ng *ngImpl) FromGetter(ac Getter) (map[string]interface{}, error) {
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

	for idx, n := range ng.nests {
		// get nested struct
		if !ac.IsStruct(n) {
			return nil, fmt.Errorf("struct named %s does not exist in %+v", n, ac)
		}

		nested = ac.Get(n)

		typ = reflect.TypeOf(nested)
		kind = typ.Kind()
		if kind != reflect.Struct {
			return nil, fmt.Errorf("name %s is not struct: %+v", n, nested)
		}

		ac, err = NewGetter(nested)
		if err != nil {
			return nil, err
		}

		// retrieve items from nested struct
		pk = strconv.Itoa(idx) + n
		pvs, ok = ng.wants[pk]
		if !ok {
			continue
		}

		if prefix == "" {
			prefix = n
		} else {
			prefix = prefix + ng.sep + n
		}
		for _, pv := range pvs {
			resIntf = ac.Get(pv)
			res[prefix+ng.sep+pv] = resIntf
		}
	}

	return res, nil
}

func (ng *ngImpl) GetNameSeparator() string {
	return ng.sep
}
