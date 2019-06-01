package main

import (
	"io"
	"log"

	"github.com/goldeneggg/structil"
)

type A struct {
	ID      int64
	Name    string
	NamePtr *string
	IsMan   bool
	AaPtr   *AA
	XArr    []X
	XPtrArr []*X
	StrArr  []string
}

type AA struct {
	Name   string
	Intf   interface{}
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
			Name: "あいう　えおあ",
			Intf: (*io.Writer)(nil),
			AaaPtr: &AAA{
				Name: "かきく　けこか",
				Val:  8,
			},
		},
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

	name, err = ac.GetString("NamePtr")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.GetString(Name): %s", name)

	IsMan, err := ac.GetBool("IsMan")
	if err != nil {
		log.Printf("!!! ERROR: %v", err)
	}
	log.Printf("Accessor.GetBool(IsMan): %v", IsMan)

	hogePtr, err := ac.Get("AaPtr")
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("Accessor.Get(AaPtr): %v", hogePtr)

	log.Printf("Accessor.IsStruct(AaPtr): %v", ac.IsStruct("AaPtr"))
	log.Printf("Accessor.IsSlice(AaPtr): %v", ac.IsSlice("XPtrArr"))

	fIntf := func(i int, a structil.Accessor) interface{} {
		s1, _ := a.GetString("Key")
		s2, _ := a.GetString("Value")
		return s1 + "=" + s2
	}
	results, err := ac.MapStructs("XArr", fIntf)
	if err != nil {
		log.Printf("!!! ERROR: %+v", err)
	}
	log.Printf("results XArr: %v, err: %v", results, err)
	results, err = ac.MapStructs("XPtrArr", fIntf)
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
