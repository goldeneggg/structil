package structil_test

import (
	"fmt"

	. "github.com/goldeneggg/structil"
)

func ExampleFinder_ToMap_simpleFind() {
	type Person struct {
		Name string
		Age  int
	}

	i := &Person{"Scott Tiger", 25}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.Find("Name", "Age").ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":25, "Name":"Scott Tiger"}
}

func ExampleFinder_ToMap_singleNestStruct() {
	type Company struct {
		Name    string
		Address string
		Period  int
	}

	type Person struct {
		Name string
		Age  int
		*Company
	}

	i := &Person{
		Name: "Mark Hunt",
		Age:  25,
		Company: &Company{
			Name:    "Tiger inc.",
			Address: "Tokyo",
			Period:  3,
		},
	}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.
		Find("Name", "Age").
		Struct("Company").Find("Period").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":25, "Company.Period":3, "Name":"Mark Hunt"}
}

func ExampleFinder_ToMap_multiNestStruct() {
	type Group struct {
		Name string
		Boss string
	}

	type Company struct {
		Name    string
		Address string
		Period  int
		*Group
	}

	type School struct {
		Name          string
		GraduatedYear int
	}

	type Person struct {
		Name string
		Age  int
		*Company
		*School
	}

	i := &Person{
		Name: "Joe Davis",
		Age:  45,
		Company: &Company{
			Name:    "XXX Cars inc.",
			Address: "New York",
			Period:  20,
			Group: &Group{
				Name: "YYY Group Holdings",
				Boss: "Donald",
			},
		},
		School: &School{
			Name:          "ABC College",
			GraduatedYear: 1995,
		},
	}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.
		Find("School").
		Struct("Company").Find("Address").
		Struct("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Company.Address":"New York", "Company.Group.Boss":"Donald", "Company.Group.Name":"YYY Group Holdings", "School":structil_test.School{Name:"ABC College", GraduatedYear:1995}}
}
