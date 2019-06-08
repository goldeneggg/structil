package structil_test

import (
	"fmt"

	. "github.com/goldeneggg/structil"
)

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
		Name: "Mark Hunt",
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
			"You worked for %d yeards since you joined the company %s",
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
	// []interface {}{"You worked for 3 yeards since you joined the company Tiger inc.", "You worked for 4 yeards since you joined the company Dragon inc."}
}
