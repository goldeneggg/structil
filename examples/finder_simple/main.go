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
		// Find(...string) returns a Finder that fields in NESTED struct are looked up by field name arguments.
		// This example looks up `person.Name` and `person.School` fields.
		Find("Name", "School").
		// Into(...string) returns a Finder that NESTED struct fields are looked up by field name arguments.
		// This example looks up `person.Company.Address` field.
		Into("Company").Find("Address").
		// If multi arguments are assigned for Into method, then execute multi level nesting.
		// This example looks up `person.Company.Group.Name` and `person.Company.Group.Boss` fields.
		Into("Company", "Group").Find("Name", "Boss").
		// ToMap converts from found struct fields to map.
		ToMap()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", m)
	// Output:
	// map[string]interface {}{"Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Name":"Lisa Mary", "School":main.school{Name:"XYZ College", GraduatedYear:2008}}
}
