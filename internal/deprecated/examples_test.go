package deprecated

import (
	"fmt"
)

func ExampleGetter() {
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
		Name: "Tony",
		Age:  25,
		Company: &Company{
			Name:    "Tiger inc.",
			Address: "Tokyo",
			Period:  3,
		},
	}

	// i must be a struct or struct pointer
	getter, err := NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"num of fields=%d\nfield names=%v\n'Name'=%s\n'Age'=%d\n'Company'=%+v",
		getter.NumField(),     // get num of fields
		getter.Names(),        // get field names
		getter.String("Name"), // get as string
		getter.Int("Age"),     // get as int
		getter.Get("Company"), // get as interface{}
	)
	// Output:
	// num of fields=3
	// field names=[Name Age Company]
	// 'Name'=Tony
	// 'Age'=25
	// 'Company'={Name:Tiger inc. Address:Tokyo Period:3}
}

func ExampleGetter_MapGet() {
	type Company struct {
		Name    string
		Address string
		Period  int
	}

	type Person struct {
		Name      string
		Age       int
		Companies []*Company
	}

	i := &Person{
		Name: "Tony",
		Age:  25,
		Companies: []*Company{
			{
				Name:    "Tiger inc.",
				Address: "Tokyo",
				Period:  3,
			},
			{
				Name:    "Dragon inc.",
				Address: "Osaka",
				Period:  4,
			},
		},
	}

	getter, err := NewGetter(i)
	if err != nil {
		panic(err)
	}

	// Each of "Companies" field are applied map function as follows.
	fn := func(i int, g *Getter) (interface{}, error) {
		return fmt.Sprintf(
			"You worked for %d years since you joined the company %s",
			g.Int("Period"),
			g.String("Name"),
		), nil
	}

	// 1st argeument must be a field name of array or slice field.
	// function assigned 2nd argument is applied each "Companies" element.
	intfs, err := getter.MapGet("Companies", fn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", intfs)
	// Output:
	// []interface {}{"You worked for 3 years since you joined the company Tiger inc.", "You worked for 4 years since you joined the company Dragon inc."}
}

func ExampleFinder() {
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

	// 2nd argument is the separator string for nested field names separating
	finder, err := NewFinderWithSep(i, ">")
	// Note:
	// If "NewFinder(i)" is called instead of "NewFinderWithSep", default separator "." is automatically used.
	// finder, err := NewFinder(i)

	if err != nil {
		panic(err)
	}

	// Finder provides method chain mechanism
	m, err := finder.
		// Find(...string) returns a Finder that fields in NESTED struct are looked up by field name arguments.
		Find("School").
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

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Company>Address":"New York", "Company>Group>Boss":"Donald", "Company>Group>Name":"YYY Group Holdings", "School":deprecated.School{Name:"ABC College", GraduatedYear:1995}}
}

func ExampleFinder_FromKeys_yml() {
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

	// testdata/finder_from_conf/ex_yml.yml as follows:
	//
	// Keys:
	//   - Company:
	//     - Group:
	//       - Name
	//       - Boss
	//     - Address
	//     - Period
	//   - Name
	//   - Age

	// Get `FinderKeys` object by calling `NewFinderKeys` with config file dir and baseName
	// This config file path is "testdata/finder_from_conf/ex_json.json"
	fks, err := NewFinderKeys("../../testdata/finder_from_conf", "ex_yml")
	if err != nil {
		panic(err)
	}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	// And build `Finder` object using `FromKeys` method with `FinderKeys` object
	// This returns the same result as follows:
	//
	// finder = finder.Find("Name", "Age").
	//   Into("Company").Find("Address", "Period").
	//   Into("Company", "Group").Find("Name", "Boss")
	m, err := finder.FromKeys(fks).ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":45, "Company.Address":"New York", "Company.Group.Boss":"Donald", "Company.Group.Name":"YYY Group Holdings", "Company.Period":20, "Name":"Joe Davis"}
}

func ExampleFinder_FromKeys_json() {
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

	// testdata/finder_from_conf/ex_json.json as follows:
	//
	// {
	//   "Keys":[
	//     {
	//       "Company":[
	//         {
	//           "Group":[
	//             "Name",
	//             "Boss"
	//           ]
	//         },
	//         "Address",
	//         "Period"
	//       ]
	//     },
	//     "Name",
	//     "Age"
	//   ]
	// }
	fks, err := NewFinderKeys("../../testdata/finder_from_conf", "ex_json")
	if err != nil {
		panic(err)
	}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.FromKeys(fks).ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":45, "Company.Address":"New York", "Company.Group.Boss":"Donald", "Company.Group.Name":"YYY Group Holdings", "Company.Period":20, "Name":"Joe Davis"}
}
