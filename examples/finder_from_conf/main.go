package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type person struct {
	Name    string
	Age     int
	Company *company
	Schools []*school
}

type school struct {
	Name          string
	GraduatedYear int
}

type company struct {
	Name    string
	Address string
	Period  int
	Group   *group
}

type group struct {
	Name string
	Boss string
}

type hoge struct {
	keys []string
}

func main() {
	i := &person{
		Name: "Lisa Mary",
		Age:  34,
		Company: &company{
			Name:    "ZZZ Air inc.",
			Address: "Boston",
			Period:  11,
			Group: &group{
				Name: "ZZZZZZ Holdings",
				Boss: "Donald Mac",
			},
		},
		Schools: []*school{
			{
				Name:          "STU High School",
				GraduatedYear: 2005,
			},
			{
				Name:          "XYZ College",
				GraduatedYear: 2008,
			},
		},
	}

	fks, err := structil.NewFinderKeys("examples/finder_from_conf", "ex_json")
	if err != nil {
		panic(err)
	}
	fmt.Printf("fks.Keys(json): %#v\n", fks.Keys())

	finder, err := structil.NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.FromKeys(fks).ToMap()
	fmt.Printf("Found Map(json): %#v, err: %v\n", m, err)

	fks, err = structil.NewFinderKeys("examples/finder_from_conf", "ex_yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("fks.Keys(yml): %#v\n", fks.Keys())

	m, err = finder.Reset().FromKeys(fks).ToMap()
	fmt.Printf("Found Map(yml): %#v, err: %v\n", m, err)
}
