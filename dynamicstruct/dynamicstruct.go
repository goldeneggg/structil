package dynamicstruct

import "reflect"

type ofType int

const (
	smpString = ""
	smpInt    = 0
	smpFloat  = 0.0
	smpBool   = false

	tMap ofType = iota
	tFunc
	tChanBoth
	tChanRecv
	tChanSend
	tStruct
	tSlice
	tPrv
)

var (
	smpMap  map[interface{}]interface{}
	smpIntf interface{}
)

type DynamicStruct interface {
	AddString(name string) DynamicStruct
	AddInt(name string) DynamicStruct
	AddFloat(name string) DynamicStruct
	AddBool(name string) DynamicStruct
	AddMap(name string) DynamicStruct
	AddFunc(name string) DynamicStruct
	AddChanBoth(name string) DynamicStruct
	AddChanRecv(name string) DynamicStruct
	AddChanSend(name string) DynamicStruct
	AddStruct(name string, i interface{}, isPtr bool) DynamicStruct
	AddStructPtr(name string, i interface{}) DynamicStruct
	AddSlice(name string) DynamicStruct
	Build() interface{}
	BuildNonPtr() interface{}
	NumBuiltField() int
	BuiltField(i int) reflect.StructField
	Remove(name string) DynamicStruct
	Exists(name string) bool
}

type DsImpl struct {
	fields     map[string]reflect.Type
	structType reflect.Type
}

func New() DynamicStruct {
	return &DsImpl{fields: map[string]reflect.Type{}}
}

type addParam struct {
	n     string
	i     interface{}
	ot    ofType
	isPtr bool
}

func (ds *DsImpl) AddString(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpString,
		ot:    tPrv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddInt(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpInt,
		ot:    tPrv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddFloat(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpFloat,
		ot:    tPrv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddBool(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpBool,
		ot:    tPrv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddMap(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpMap,
		ot:    tMap,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddFunc(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpIntf,
		ot:    tFunc,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddChanBoth(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpIntf,
		ot:    tChanBoth,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddChanRecv(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpIntf,
		ot:    tChanRecv,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddChanSend(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpIntf,
		ot:    tChanSend,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddStruct(name string, i interface{}, isPtr bool) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     i,
		ot:    tStruct,
		isPtr: isPtr,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) AddStructPtr(name string, i interface{}) DynamicStruct {
	return ds.AddStruct(name, i, true)
}

func (ds *DsImpl) AddSlice(name string) DynamicStruct {
	p := &addParam{
		n:     name,
		i:     smpIntf,
		ot:    tSlice,
		isPtr: false,
	}
	ds.add(p)
	return ds
}

func (ds *DsImpl) add(p *addParam) {
	it := reflect.TypeOf(p.i)
	var typeOf reflect.Type

	switch p.ot {
	case tMap:
		typeOf = reflect.MapOf(it.Key(), it.Elem())
	case tFunc:
		typeOf = reflect.FuncOf([]reflect.Type{it}, []reflect.Type{it}, false)
	case tChanBoth:
		typeOf = reflect.ChanOf(reflect.BothDir, it)
	case tChanRecv:
		typeOf = reflect.ChanOf(reflect.RecvDir, it)
	case tChanSend:
		typeOf = reflect.ChanOf(reflect.SendDir, it)
	case tStruct:
		fs := make([]reflect.StructField, it.NumField())
		for i := 0; i < it.NumField(); i++ {
			fs[i] = it.Field(i)
		}
		typeOf = reflect.StructOf(fs)
	case tSlice:
		typeOf = reflect.SliceOf(it)
	default:
		typeOf = it
	}

	if p.isPtr {
		typeOf = reflect.PtrTo(typeOf)
	}

	ds.fields[p.n] = typeOf
}

func (ds *DsImpl) Build() interface{} {
	return ds.build(true)
}

func (ds *DsImpl) BuildNonPtr() interface{} {
	return ds.build(true)
}

func (ds *DsImpl) build(isPtr bool) interface{} {
	var i int
	fs := make([]reflect.StructField, len(ds.fields))

	for name, typ := range ds.fields {
		fs[i] = reflect.StructField{Name: name, Type: typ}
		i++
	}
	ds.structType = reflect.StructOf(fs)
	n := reflect.New(ds.structType)

	if isPtr {
		return n.Interface()
	} else {
		return reflect.Indirect(n).Interface()
	}
}

func (ds *DsImpl) NumBuiltField() int {
	return ds.structType.NumField()
}

func (ds *DsImpl) BuiltField(i int) reflect.StructField {
	return ds.structType.Field(i)
}

func (ds *DsImpl) Remove(name string) DynamicStruct {
	delete(ds.fields, name)
	return ds
}

func (ds *DsImpl) Exists(name string) bool {
	_, ok := ds.fields[name]
	return ok
}
