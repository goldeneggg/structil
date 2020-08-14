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

	// Get `FinderKeys` object by calling `NewFinderKeys` with config file dir and baseName
	// This config file path is "examples/finder_from_conf/ex_json.json"
	fks, err := structil.NewFinderKeys(".", "ex_json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("fks.Keys(json): %#v\n", fks.Keys())
	// Output:
	// fks.Keys(json): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}

	finder, err := structil.NewFinder(i)
	if err != nil {
		panic(err)
	}

	// And build `Finder` object using `FromKeys` method with `FinderKeys` object
	finder = finder.FromKeys(fks)
	// This returns same result as follows:
	//
	// finder = finder.FindTop("Name", "Age").
	//   Into("Company").Find("Address", "Period").
	//   Into("Company", "Group").Find("Name", "Boss")

	// ToMap converts from found struct fields to map.
	m, err := finder.ToMap()

	fmt.Printf("Found Map(json): %#v, err: %v\n", m, err)
	// Output:
	// Found Map(json): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>

	// YAML example as follows
	fks, err = structil.NewFinderKeys(".", "ex_yml")
	if err != nil {
		panic(err)
	}

	fmt.Printf("fks.Keys(yml): %#v\n", fks.Keys())
	// Output:
	// fks.Keys(yml): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}

	m, err = finder.Reset().FromKeys(fks).ToMap()

	fmt.Printf("Found Map(yml): %#v, err: %v\n", m, err)
	// Output:
	// Found Map(yml): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
}
