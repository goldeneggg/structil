package main

import (
	"io"
	"log"
	"os"
	"reflect"

	"github.com/goldeneggg/structil"
)

type A struct {
	ID       int64
	Name     string
	NamePtr  *string
	IsMan    bool
	FloatVal float64
	AaPtr    *AA
	Nil      *AA
	XArr     []X
	XPtrArr  []*X
	StrArr   []string
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
		ID:       1,
		Name:     name,
		NamePtr:  &name,
		IsMan:    true,
		FloatVal: 3.14,
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
	exampleGetter()
	exampleFinder()
}

func exampleGetter() {
	log.Println("---------- exampleGetter")
	g, err := structil.NewGetter(hoge)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	name := g.Get("Name")
	log.Printf("Getter.Get(Name): %s", name)

	name = g.GetString("NamePtr")
	log.Printf("Getter.GetString(NamePtr): %s", name)

	intVal := g.GetInt64("ID")
	log.Printf("Getter.GetInt64(ID): %v", intVal)

	floatVal := g.GetFloat64("FloatVal")
	log.Printf("Getter.GetFloat64(FloatVal): %v", floatVal)

	IsMan := g.GetBool("IsMan")
	log.Printf("Getter.GetBool(IsMan): %v", IsMan)

	// AaPtr
	aaPtr := g.Get("AaPtr")
	log.Printf("Getter.Get(AaPtr): %v", aaPtr)
	log.Printf("Getter.IsStruct(AaPtr): %v", g.IsStruct("AaPtr"))

	aaAc, err := structil.NewGetter(aaPtr)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	it := aaAc.Get("Writer")
	log.Printf("AaPtr.Get(Writer): %+v", it)
	log.Printf("AaPtr.Get(Writer).ValueOf().Elem(): %+v", reflect.ValueOf(it).Elem())
	log.Printf("AaPtr.IsStruct(Writer): %v", aaAc.IsStruct("Writer"))

	// Nil
	aNil := g.Get("Nil")
	log.Printf("Getter.Get(Nil): %v", aNil)
	log.Printf("Getter.IsStruct(Nil): %v", g.IsStruct("Nil"))

	aNilAc, err := structil.NewGetter(aNil)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Getter.Get(Nil).NewGetter: %+v", aNilAc)

	// XArr
	xArr := g.Get("XArr")
	log.Printf("Getter.Get(XArr): %v", xArr)
	log.Printf("Getter.IsStruct(XArr): %v", g.IsStruct("XArr"))
	log.Printf("Getter.IsSlice(XArr): %v", g.IsSlice("XArr"))

	// Map
	fa := func(i int, a structil.Getter) interface{} {
		s1 := a.GetString("Key")
		s2 := a.GetString("Value")
		return s1 + "=" + s2
	}

	results, err := g.MapGet("XArr", fa)
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("results XArr: %v, err: %v", results, err)

	results, err = g.MapGet("XPtrArr", fa)
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("results XPtrArr: %v, err: %v", results, err)

	g.DumpRVs()
}

func exampleFinder() {
	log.Println("---------- exampleFinder")

	finder, err := structil.NewFinder(hoge)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Finder: %#v", finder)

	swRes, err := finder.
		Struct("AaPtr").Find("Name").
		Struct("AaPtr", "AaaPtr").Find("Name", "Val").
		ToMap()
	log.Printf("Finder.ToMap res: %+v, err: %v", swRes, err)
}
