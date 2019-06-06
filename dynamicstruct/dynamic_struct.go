package dynamicstruct

import "reflect"

const (
	smpString = ""
	smpInt    = 0
	smpFloat  = 0.0
	smpBool   = false
)

type ofType int

const (
	tSlice ofType = iota
	tMap
	tChanSend
	tChanRecv
	tChanBoth
	tFunc
	tPrv
)

var (
	smpIntf interface{}
)

type DynamicStruct interface {
	AddString(name string) DynamicStruct
	AddInt(name string) DynamicStruct
	AddFloat(name string) DynamicStruct
	AddBool(name string) DynamicStruct
	AddSlice(name string) DynamicStruct
	AddMap(name string) DynamicStruct
	AddFunc(name string) DynamicStruct
	Remove(name string) DynamicStruct
	Exists(name string) bool
	Build() interface{}
	BuildNonPtr() interface{}
}

type dynImpl struct {
	fields     map[string]reflect.Type
	structType reflect.Type
	st         interface{}
	ptrSt      interface{}
}

func New() DynamicStruct {
	return &dynImpl{fields: map[string]reflect.Type{}}
}

func (ds *dynImpl) AddString(name string) DynamicStruct {
	ds.add(name, smpString, tPrv, false)
	return ds
}

func (ds *dynImpl) AddInt(name string) DynamicStruct {
	ds.add(name, smpInt, tPrv, false)
	return ds
}

func (ds *dynImpl) AddFloat(name string) DynamicStruct {
	ds.add(name, smpFloat, tPrv, false)
	return ds
}

func (ds *dynImpl) AddBool(name string) DynamicStruct {
	ds.add(name, smpBool, tPrv, false)
	return ds
}

func (ds *dynImpl) AddSlice(name string) DynamicStruct {
	ds.add(name, smpIntf, tSlice, false)
	return ds
}

func (ds *dynImpl) AddMap(name string) DynamicStruct {
	ds.add(name, smpIntf, tMap, false)
	return ds
}

func (ds *dynImpl) AddFunc(name string) DynamicStruct {
	ds.add(name, smpIntf, tFunc, false)
	return ds
}

func (ds *dynImpl) add(name string, e interface{}, ot ofType, isPtr bool) {
	et := reflect.TypeOf(e)
	var typeOf reflect.Type

	switch ot {
	case tSlice:
		typeOf = reflect.SliceOf(et)
	case tMap:
		typeOf = reflect.MapOf(et, et)
	case tFunc:
		typeOf = reflect.FuncOf([]reflect.Type{et}, []reflect.Type{et}, false)
	default:
		typeOf = et
	}

	if isPtr {
		typeOf = reflect.PtrTo(typeOf)
	}

	ds.fields[name] = typeOf
}

func (ds *dynImpl) Remove(name string) DynamicStruct {
	delete(ds.fields, name)
	return ds
}

func (ds *dynImpl) Exists(name string) bool {
	_, ok := ds.fields[name]
	return ok
}

func (ds *dynImpl) Build() interface{} {
	return ds.build(true)
}

func (ds *dynImpl) BuildNonPtr() interface{} {
	return ds.build(true)
}

func (ds *dynImpl) build(isPtr bool) interface{} {
	var fs []reflect.StructField

	for name, typ := range ds.fields {
		fs = append(fs, reflect.StructField{Name: name, Type: typ})
	}
	ds.structType = reflect.StructOf(fs)
	n := reflect.New(ds.structType)

	if isPtr {
		return n.Interface()
	} else {
		return reflect.Indirect(n).Interface()
	}
}
