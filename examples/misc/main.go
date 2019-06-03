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
	exampleEmbedder()
}

func exampleGetter() {
	log.Println("---------- exampleGetter")
	ac, err := structil.NewGetter(hoge)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	name := ac.Get("Name")
	log.Printf("Getter.Get(Name): %s", name)

	name = ac.GetString("NamePtr")
	log.Printf("Getter.GetString(NamePtr): %s", name)

	IsMan := ac.GetBool("IsMan")
	log.Printf("Getter.GetBool(IsMan): %v", IsMan)

	floatVal := ac.GetFloat64("FloatVal")
	log.Printf("Getter.GetFloat64(FloatVal): %v", floatVal)

	// AaPtr
	aaPtr := ac.Get("AaPtr")
	log.Printf("Getter.Get(AaPtr): %v", aaPtr)
	log.Printf("Getter.IsStruct(AaPtr): %v", ac.IsStruct("AaPtr"))
	log.Printf("Getter.IsInterface(AaPtr): %v", ac.IsInterface("AaPtr"))

	aaAc, err := structil.NewGetter(aaPtr)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}

	it := aaAc.Get("Writer")
	log.Printf("AaPtr.Get(Writer): %+v", it)
	log.Printf("AaPtr.Get(Writer).ValueOf().Elem(): %+v", reflect.ValueOf(it).Elem())
	log.Printf("AaPtr.IsStruct(Writer): %v", aaAc.IsStruct("Writer"))
	log.Printf("AaPtr.IsInterface(Writer): %v", aaAc.IsInterface("Writer"))

	// Nil
	rvNil := ac.GetRV("Nil")
	log.Printf("Getter.GetRV(Nil): %v", rvNil)
	aNil := ac.Get("Nil")
	log.Printf("Getter.Get(Nil): %v", aNil)
	log.Printf("Getter.IsStruct(Nil): %v", ac.IsStruct("Nil"))
	log.Printf("Getter.IsInterface(Nil): %v", ac.IsInterface("Nil"))

	aNilAc, err := structil.NewGetter(aNil)
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Getter.Get(Nil).NewGetter: %+v", aNilAc)

	// XArr
	xArr := ac.Get("XArr")
	log.Printf("Getter.Get(XArr): %v", xArr)
	log.Printf("Getter.IsStruct(XArr): %v", ac.IsStruct("XArr"))
	log.Printf("Getter.IsSlice(XArr): %v", ac.IsSlice("XArr"))
	log.Printf("Getter.IsInterface(XArr): %v", ac.IsInterface("XArr"))

	// Map
	fa := func(i int, a structil.Getter) interface{} {
		s1 := a.GetString("Key")
		s2 := a.GetString("Value")
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

func exampleEmbedder() {
	log.Println("---------- exampleEmbedder")
	ac, err := structil.NewGetter(hoge)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	swRes, err := structil.NewEmbedder().
		Seek("AaPtr").Want("Name").
		Seek("AaaPtr").Want("Name").Want("Val").
		From(hoge)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Embedder.From res: %#v", swRes)

	swRes, err = structil.NewEmbedder().
		Seek("AaPtr").Want("Name").
		Seek("AaaPtr").Want("Name").Want("Val").
		FromGetter(ac)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	log.Printf("Embedder.FromGetter res: %#v", swRes)
}
