package structil

import (
	"fmt"
)

func ExampleGetter_String() {
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

func ExampleGetter_Int() {
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

func ExampleGetter_Get_struct() {
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
func ExampleGetter_Get_structSliceField() {
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

func ExampleGetter_MapGet_joinElements() {
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
