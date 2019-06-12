# structil

[![GoDoc](https://godoc.org/github.com/goldeneggg/structil?status.png)](https://godoc.org/github.com/goldeneggg/structil)
<br />
<br />
[![Build Status](https://travis-ci.org/goldeneggg/structil.svg?branch=master)](https://travis-ci.org/goldeneggg/structil)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![GolangCI](https://golangci.com/badges/github.com/goldeneggg/gat.svg)](https://golangci.com/r/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

Struct Utilities for runtime and dynamic environment in Go.

## Runtime and Dynamic struct accessor

### `Finder`
We can access usefully nested struct fields using field name string.

[Sample script on playground](https://play.golang.org/p/AcF5c7Prf3z).

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

	// We can use method chain for Find and Into methods.
	//  - FindTop returns a Finder that top level fields in struct are looked up and held named "names".
	//  - Into returns a Finder that nested struct fields are looked up and held named "names".
	//  - Find returns a Finder that fields in struct are looked up and held named "names".
	// And finally, we can call ToMap method for converting from struct to map.
	m, err := finder.
		FindTop("Name", "School").
		Into("Company").Find("Address").
		Into("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", m)
}
```

Result as follows.

```
map[string]interface {}{"Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Name":"Lisa Mary", "School":main.school{Name:"XYZ College", GraduatedYear:2008}}
```

#### with `FinderKeys`
We can create a Finder from the configuration file that have some finding target keys.

We support some file format of configuration file such as `yaml`, `json`, `toml` and more.

Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).

```go
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

  json(i)
  yml(i)
}

func json(i *person) {
  fks, err := structil.NewFinderKeysFromConf("examples/finder_from_conf", "ex_json")
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }
  fmt.Printf("fks.Keys(json): %#v\n", fks.Keys())

  finder, err := structil.NewFinder(i)
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }

  m, err := finder.FromKeys(fks).ToMap()
  fmt.Printf("Found Map(json): %#v, err: %v\n", m, err)
}

func yml(i *person) {
  fks, err := structil.NewFinderKeysFromConf("examples/finder_from_conf", "ex_yml")
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }
  fmt.Printf("fks.Keys(yml): %#v\n", fks.Keys())

  finder, err := structil.NewFinder(i)
  if err != nil {
    fmt.Printf("error: %v\n", err)
    return
  }

  m, err := finder.FromKeys(fks).ToMap()
  fmt.Printf("Found Map(yml): %#v, err: %v\n", m, err)
}
```

Result as follows.

```
fks.Keys(json): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(json): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
fks.Keys(yml): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(yml): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
```

### `Getter`
Use `Getter`

We can access a struct using field name string, like map.

Sample script on playground is https://play.golang.org/p/3CNDJpW3UmN .

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

Result as follows.

```
Name: "Mike Davis", Age: 27, Company: main.company{Name:"Scott inc.", Address:"Osaka", Period:2}
```

#### Collection method for slice of struct
`MapGet` method provides the __Map__ collection function for slice of struct

Sample script on playground is https://play.golang.org/p/98wCWCrs0vf .

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

Result as follows.

```
[]interface {}{"You worked for 3 years since you joined the company Dragons inc.", "You worked for 2 years since you joined the company Swallows inc."}
```

