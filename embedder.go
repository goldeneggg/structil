package structil

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	defaultNameSep = "."
)

type Embedder interface {
	// TODO: 同一階層で複数の子要素nestを行いたいケースへの対応
	// nameをvariadic(可変) で受ける形式が良さそう
	Seek(name string) Embedder

	Want(name string) Embedder

	// MapConfigure(m map[string]interface{}) Embedder  // TODO: mapでネストを一括設定

	From(st interface{}) (map[string]interface{}, error)

	FromGetter(g Getter) (map[string]interface{}, error)

	GetNameSeparator() string
}

type embImpl struct {
	current string
	seeks   []string
	wants   map[string][]string
	sep     string
}

func NewEmbedder() Embedder {
	return NewEmbedderWithSep(defaultNameSep)
}

func NewEmbedderWithSep(sep string) Embedder {
	return &embImpl{wants: map[string][]string{}, sep: sep}
}

func (e *embImpl) Seek(name string) Embedder {
	e.seeks = append(e.seeks, name)

	e.current = strconv.Itoa(len(e.seeks)-1) + name
	e.wants[e.current] = []string{}
	return e
}

// TODO: 一段もSeekしてなくてもWantできるように
func (e *embImpl) Want(name string) Embedder {
	if _, ok := e.wants[e.current]; ok {
		e.wants[e.current] = append(e.wants[e.current], name)
	}

	return e
}

func (e *embImpl) From(i interface{}) (map[string]interface{}, error) {
	g, err := NewGetter(i)
	if err != nil {
		return nil, err
	}

	return e.FromGetter(g)
}

// TODO: 並列化検討
func (e *embImpl) FromGetter(g Getter) (map[string]interface{}, error) {
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

	for idx, n := range e.seeks {
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
		pvs, ok = e.wants[pk]
		if !ok {
			continue
		}

		if prefix == "" {
			prefix = n
		} else {
			prefix = prefix + e.sep + n
		}
		for _, pv := range pvs {
			resIntf = g.Get(pv)
			res[prefix+e.sep+pv] = resIntf
		}
	}

	return res, nil
}

func (e *embImpl) GetNameSeparator() string {
	return e.sep
}
