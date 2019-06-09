package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type group struct {
	Name string
	Boss string
}

type company struct {
	Name    string
	Address string
	Period  int
	Group   *group
}

type school struct {
	Name          string
	GraduatedYear int
}

type person struct {
	Name    string
	Age     int
	Company *company
	School  *school
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
		School: &school{
			Name:          "XYZ College",
			GraduatedYear: 2008,
		},
	}

	finder, err := structil.NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.
		Find("Name", "School").
		Into("Company").Find("Address").
		Into("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", m)
}
