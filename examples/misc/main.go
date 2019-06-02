package main

import (
	"io"
	"log"
	"os"
	"reflect"

	"github.com/goldeneggg/structil"
)

type A struct {
	ID      int64
	Name    string
	NamePtr *string
	IsMan   bool
	AaPtr   *AA
	Nil     *AA
	XArr    []X
	XPtrArr []*X
	StrArr  []string
}

type AA struct {
	Name   string
	Writer io.Writer
	AaaPtr *AAA
}

type AAA struct {
	Name string
	Val  int
}

type X struct {
	Key   string
	Value string
}

var (
	name = "ほげ　ふがお"

	hoge = &A{
		ID:      1,
		Name:    name,
		NamePtr: &name,
		IsMan:   true,
		AaPtr: &AA{
			Name:   "あいう　えおあ",
			Writer: os.Stdout,
			AaaPtr: &AAA{
				Name: "かきく　けこか",
				Val:  8,
			},
		},
		Nil: nil,
		XArr: []X{
			{
				Key:   "key1",
				Value: "value1",
			},
			{
				Key:   "key2",
				Value: "value2",
			},
		},
		XPtrArr: []*X{
			{
				Key:   "key100",
				Value: "value100",
			},
			{
				Key:   "key200",
				Value: "value200",
			},
		},
		StrArr: []string{"key1", "value1", "key2", "value2"},
	}
)

func main() {
	// logger, _ := zap.NewDevelopment()
	// logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))

	exampleAccessor()
	exampleRetriever()
	exampleUtil()
}

func exampleAccessor() {
	log.Println("---------- exampleAccessor")
	ac, err := structil.NewAccessor(hoge)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	name, err := ac.Get("Name")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.Get(Name): %s", name)

	hoge.Name = "あほ　ぼけお"
	name, err = ac.Get("Name")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.Get(Name) AFTER: %s", name)

	name, err = ac.GetString("NamePtr")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.GetString(NamePtr): %s", name)

	IsMan, err := ac.GetBool("IsMan")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.GetBool(IsMan): %v", IsMan)

	// AaPtr
	aaPtr, err := ac.Get("AaPtr")
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("Accessor.Get(AaPtr): %v", aaPtr)
	log.Printf("Accessor.IsStruct(AaPtr): %v", ac.IsStruct("AaPtr"))
	log.Printf("Accessor.IsInterface(AaPtr): %v", ac.IsInterface("AaPtr"))

	aaAc, err := structil.NewAccessor(aaPtr)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	it, err := aaAc.Get("Writer")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("AaPtr.Get(Writer): %+v", it)
	log.Printf("AaPtr.Get(Writer).ValueOf().Elem(): %+v", reflect.ValueOf(it).Elem())
	log.Printf("AaPtr.IsStruct(Writer): %v", aaAc.IsStruct("Writer"))
	log.Printf("AaPtr.IsInterface(Writer): %v", aaAc.IsInterface("Writer"))

	// Nil
	rvNil, err := ac.GetRV("Nil")
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("Accessor.GetRV(Nil): %v", rvNil)
	aNil, err := ac.Get("Nil")
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("Accessor.Get(Nil): %v", aNil)
	log.Printf("Accessor.IsStruct(Nil): %v", ac.IsStruct("Nil"))
	log.Printf("Accessor.IsInterface(Nil): %v", ac.IsInterface("Nil"))

	aNilAc, err := structil.NewAccessor(aNil)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.Get(Nil).NewAccessor: %+v", aNilAc)

	// XArr
	xArr, err := ac.Get("XArr")
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("Accessor.Get(XArr): %v", xArr)
	log.Printf("Accessor.IsStruct(XArr): %v", ac.IsStruct("XArr"))
	log.Printf("Accessor.IsSlice(XArr): %v", ac.IsSlice("XArr"))
	log.Printf("Accessor.IsInterface(XArr): %v", ac.IsInterface("XArr"))

	// Map
	fa := func(i int, a structil.Accessor) interface{} {
		s1, _ := a.GetString("Key")
		s2, _ := a.GetString("Value")
		return s1 + "=" + s2
	}

	results, err := ac.MapStructs("XArr", fa)
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("results XArr: %v, err: %v", results, err)

	results, err = ac.MapStructs("XPtrArr", fa)
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("results XPtrArr: %v, err: %v", results, err)
}

func exampleRetriever() {
	log.Println("---------- exampleRetriever")
	ac, err := structil.NewAccessor(hoge)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	swRes, err := structil.NewRetriever().
		Nest("AaPtr").Want("Name").
		Nest("AaaPtr").Want("Name").Want("Val").
		From(hoge)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Retriever.From res: %#v", swRes)

	swRes, err = structil.NewRetriever().
		Nest("AaPtr").Want("Name").
		Nest("AaaPtr").Want("Name").Want("Val").
		FromAccessor(ac)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Retriever.FromAccessor res: %#v", swRes)
}

func exampleUtil() {
	log.Println("---------- exampleUtil")
	var v reflect.Value

	v = structil.IElemOf(hoge.ID)
	log.Printf("IElemOf int: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.Name)
	log.Printf("IElemOf string: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.NamePtr)
	log.Printf("IElemOf string ptr: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.IsMan)
	log.Printf("IElemOf bool: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.AaPtr)
	log.Printf("IElemOf struct ptr: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.AaPtr.Writer)
	log.Printf("IElemOf interface: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
	v = structil.IElemOf(hoge.Nil)
	log.Printf("IElemOf struct ptr nil: %+v", v)
	v = structil.IElemOf(hoge.XPtrArr)
	log.Printf("IElemOf struct slice ptr: %+v, Type: %+v, Kind: %+v", v, v.Type(), v.Kind())
}
