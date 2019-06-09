package main

import (
	"fmt"

	"github.com/goldeneggg/structil"
)

type company struct {
	Name    string
	Address string
	Period  int
}

type person struct {
	Name      string
	Age       int
	Companies []*company
}

func main() {
	i := &person{
		Name: "John",
		Age:  28,
		Companies: []*company{
			{
				Name:    "Dragons inc.",
				Address: "Nagoya",
				Period:  3,
			},
			{
				Name:    "Swallows inc.",
				Address: "Tokyo",
				Period:  2,
			},
		},
	}

	getter, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fn := func(i int, g structil.Getter) (interface{}, error) {
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
}
