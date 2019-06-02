package structil

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	defaultNameSep = "."
)

type Retriever interface {
	Nest(name string) Retriever
	Want(name string) Retriever
	From(st interface{}) (map[string]interface{}, error)
	FromAccessor(ac Accessor) (map[string]interface{}, error)
	GetNameSeparator() string
}

type collectorImpl struct {
	current string
	nests   []string
	wants   map[string][]string
	sep     string
}

func NewRetriever() Retriever {
	return NewRetrieverWithSep(defaultNameSep)
}

func NewRetrieverWithSep(sep string) Retriever {
	return &collectorImpl{wants: map[string][]string{}, sep: sep}
}

func (sf *collectorImpl) Nest(name string) Retriever {
	sf.nests = append(sf.nests, name)

	sf.current = strconv.Itoa(len(sf.nests)-1) + name
	sf.wants[sf.current] = []string{}
	return sf
}

func (sf *collectorImpl) Want(name string) Retriever {
	if _, ok := sf.wants[sf.current]; ok {
		sf.wants[sf.current] = append(sf.wants[sf.current], name)
	}

	return sf
}

func (sf *collectorImpl) From(st interface{}) (map[string]interface{}, error) {
	ac, err := NewAccessor(st)
	if err != nil {
		return nil, err
	}

	return sf.FromAccessor(ac)
}

func (sf *collectorImpl) FromAccessor(ac Accessor) (map[string]interface{}, error) {
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

	for idx, n := range sf.nests {
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

		ac, err = NewAccessor(nested)
		if err != nil {
			return nil, err
		}

		// retrieve items from nested struct
		pk = strconv.Itoa(idx) + n
		pvs, ok = sf.wants[pk]
		if !ok {
			continue
		}

		if prefix == "" {
			prefix = n
		} else {
			prefix = prefix + sf.sep + n
		}
		for _, pv := range pvs {
			resIntf = ac.Get(pv)
			res[prefix+sf.sep+pv] = resIntf
		}
	}

	return res, nil
}

func (sf *collectorImpl) GetNameSeparator() string {
	return sf.sep
}
