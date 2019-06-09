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
	Name    string
	Age     int
	Company *company
}

func main() {
	i := &person{
		Name: "Mike Davis",
		Age:  27,
		Company: &company{
			Name:    "Scott inc.",
			Address: "Osaka",
			Period:  2,
		},
	}

	getter, err := structil.NewGetter(i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Name: %+v, Age: %+v, Company: %+v\n", getter.String("Name"), getter.Int("Age"), getter.Get("Company"))
}
