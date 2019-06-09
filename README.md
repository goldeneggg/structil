# structil

[![Build Status](https://travis-ci.org/goldeneggg/structil.svg?branch=master)](https://travis-ci.org/goldeneggg/structil)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![GolangCI](https://golangci.com/badges/github.com/goldeneggg/gat.svg)](https://golangci.com/r/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)

Struct Utilities for runtime and dynamic environment in Go.

## Usage

### `Getter`
Use `Getter`

We can access a struct using field name string, like map.

```go
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
```
```
Name: Mike Davis, Age: 27, Company: {Name:Scott inc. Address:Osaka Period:2}
```

#### Collection method for slice of struct
`MapGet` method provides the __Map__ collection function for slice of struct

```go
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

	fmt.Printf("%#v\n", intfs)
}
```
```
[]interface {}{"You worked for 3 years since you joined the company Dragons inc.", "You worked for 2 years since you joined the company Swallows inc."}
```

### `Finder`
Use `Finder`

We can access usefully nested struct fields using field name string

```go
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
		Struct("Company").Find("Address").
		Struct("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
}
```
```
map[string]interface {}{"Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Name":"Lisa Mary", "School":main.school{Name:"XYZ College", GraduatedYear:2008}}
```
