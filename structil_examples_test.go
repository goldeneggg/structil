package structil

import (
	"fmt"
)

func ExampleGetterImpl_String() {
	type Person struct {
		Name string
		Age  int
	}

	i := &Person{
		Name: "Tony",
		Age:  25,
	}

	getter, err := NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", getter.String("Name"))
	// Output:
	// Tony
}

func ExampleGetterImpl_Int() {
	type Person struct {
		Name string
		Age  int
	}

	i := &Person{
		Name: "Tony",
		Age:  25,
	}

	getter, err := NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", getter.Int("Age"))
	// Output:
	// 25
}

func ExampleGetterImpl_Get_struct() {
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

	getter, err := NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", getter.Get("Company"))
	// Output:
	// {Name:Tiger inc. Address:Tokyo Period:3}
}

/*
func ExampleGetterImpl_Get_structSliceField() {
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

	fmt.Printf("%#v", getter.Get("Companies"))
	// Output:
	// {Name:Tiger inc. Address:Tokyo Period:3}
}
*/

func ExampleGetterImpl_MapGet_joinElements() {
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

	fn := func(i int, g Getter) (interface{}, error) {
		return fmt.Sprintf(
			"You worked for %d years since you joined the company %s",
			g.Int("Period"),
			g.String("Name"),
		), nil
	}

	intfs, err := getter.MapGet("Companies", fn)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", intfs)
	// Output:
	// []interface {}{"You worked for 3 years since you joined the company Tiger inc.", "You worked for 4 years since you joined the company Dragon inc."}
}

func ExampleFinderImpl_ToMap_simpleFind() {
	type Person struct {
		Name string
		Age  int
	}

	i := &Person{"Scott Tiger", 25}

	finder, err := NewFinder(i)
	if err != nil {
		panic(err)
	}

	m, err := finder.FindTop("Name", "Age").ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":25, "Name":"Scott Tiger"}
}

func ExampleFinderImpl_ToMap_singleNestInto() {
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
		FindTop("Name", "Age").
		Into("Company").Find("Period").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Age":25, "Company.Period":3, "Name":"Mark Hunt"}
}

func ExampleFinderImpl_ToMap_multiNestInto() {
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
		FindTop("School").
		Into("Company").Find("Address").
		Into("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
	// Output:
	// map[string]interface {}{"Company.Address":"New York", "Company.Group.Boss":"Donald", "Company.Group.Name":"YYY Group Holdings", "School":structil.School{Name:"ABC College", GraduatedYear:1995}}
}

func ExampleFinderImpl_FromKeys_yml() {
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

	// examples/finder_from_conf/ex_yml.yml as follows:
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

	fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_yml")
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

func ExampleFinderImpl_FromKeys_json() {
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

	// examples/finder_from_conf/ex_json.json as follows:
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
	fks, err := NewFinderKeysFromConf("examples/finder_from_conf", "ex_json")
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
