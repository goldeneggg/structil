structil  [![GoDoc](https://godoc.org/github.com/goldeneggg/structil?status.png)](https://godoc.org/github.com/goldeneggg/structil)
==========

[![Build Status](https://travis-ci.org/goldeneggg/structil.svg?branch=master)](https://travis-ci.org/goldeneggg/structil)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

Struct Utilities for runtime and dynamic environment in Go.

__Table of Contents__

<!-- TOC depthFrom:1 -->

- [structil  ![GoDoc](https://godoc.org/github.com/goldeneggg/structil)](#structil-img-src%22httpsgodocorggithubcomgoldenegggstructil%22-alt%22godoc%22)
  - [`Finder`](#finder)
    - [With config file? use `FinderKeys`](#with-config-file-use-finderkeys)
  - [`Getter`](#getter)
    - [`MapGet` method](#mapget-method)
  - [`DynamicStruct`](#dynamicstruct)
    - [JSON unmershal with `DynamicStruct`](#json-unmershal-with-dynamicstruct)

<!-- /TOC -->

## `Finder`
We can access usefully nested struct fields using field name string.

[Sample script on playground](https://play.golang.org/p/AcF5c7Prf3z).

Get `Finder` instance by calling `NewFinder(i interface{})` with an initialized struct.

```go
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
```

Then we can access nested struct by field name. 

```go
// FindTop(...string) returns a Finder that top level fields in struct are looked up by field name arguments.
// This example looks up `person.Name` and `person.School` fields.
finder = finder.FindTop("Name", "School")

// Into(...string) returns a Finder that NESTED struct fields are looked up by field name arguments.
// And Find(...string) returns a Finder that fields in NESTED struct are looked up by field name arguments.
// This example looks up `person.Company.Address` field.
finder = finder.Into("Company").Find("Address")

// If multi arguments are assigned for Into method, then execute multi level nesting.
// This example looks up `person.Company.Group.Name` and `person.Company.Group.Boss`fields.
finder = finder.Into("Company", "Group").Find("Name", "Boss")

// ToMap converts from found struct fields to map.
m, err := finder.ToMap()

fmt.Printf("%#v", m)
```

Result as follows.

```
map[string]interface {}{"Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Name":"Lisa Mary", "School":main.school{Name:"XYZ College", GraduatedYear:2008}}
```

### With config file? use `FinderKeys`
We can create a Finder from the configuration file that have some finding target keys. We support some file formats of configuration file such as `yaml`, `json`, `toml` and more.

Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).

File `examples/finder_from_conf/keys.yml` as follows:

```yml
# "Keys" is required field for top level.
Keys:
  # "- FIELDNAME:" is nest sign for "FIELDNAME"
  - Company:
    - Group:
      - Name
      - Boss
    - Address
    - Period
  # "- FIELDNAME" is finding sign for "FIELDNAME"
  - Name
  - Age
```

Get `FinderKeys` instance by calling `NewFinderKeys(dir, baseName)` with config file dir and filename.

```go
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

finder, err := structil.NewFinder(i)
if err != nil {
	return
}

// Get `FinderKeys` object by calling `NewFinderKeys` with config file dir and baseName
fks, err := structil.NewFinderKeys("examples/finder_from_conf", "keys")
if err != nil {
	fmt.Printf("error: %v\n", err)
	return
}
// And build `Finder` object using `FromKeys` method with `FinderKeys` object
finder = finder.FromKeys(fks)
// This returns same result as follows:
//
// finder = finder.FindTop("Name", "Age").
//   Into("Company").Find("Address", "Period").
//   Into("Company", "Group").Find("Name", "Boss")


// ToMap converts from found struct fields to map.
m, err := finder.ToMap()
fmt.Printf("Found Map(yml): %#v, err: %v\n", m, err)
```

Result as follows.

```
fks.Keys(json): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(json): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
fks.Keys(yml): []string{"Company.Group.Name", "Company.Group.Boss", "Company.Address", "Company.Period", "Name", "Age"}
Found Map(yml): map[string]interface {}{"Age":34, "Company.Address":"Boston", "Company.Group.Boss":"Donald Mac", "Company.Group.Name":"ZZZZZZ Holdings", "Company.Period":11, "Name":"Lisa Mary"}, err: <nil>
```


## `Getter`
We can access a struct using field name string, like map.

[Sample script on playground](https://play.golang.org/p/3CNDJpW3UmN).

```go
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

fmt.Printf("Name: %+v, Age: %+v, Company: %+v\n", getter.String("Name"), getter.Int("Age"), getter.Get("Company"))
```

Result as follows.

```
Name: "Mike Davis", Age: 27, Company: main.company{Name:"Scott inc.", Address:"Osaka", Period:2}
```

### `MapGet` method
`MapGet` method provides the __Map__ collection function for slice of struct

[Sample script on playground](https://play.golang.org/p/98wCWCrs0vf).

```go
// Companies field is slice of struct.
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

// Each of Companies field are applied map function as follows.
fn := func(i int, g structil.Getter) (interface{}, error) {
	return fmt.Sprintf(
		"You worked for %d years since you joined the company %s",
		g.Int("Period"),
		g.String("Name"),
	), nil
}
intfs, err := getter.MapGet("Companies", fn)

fmt.Printf("%#v\n", intfs)
```

Result as follows.

```
[]interface {}{"You worked for 3 years since you joined the company Dragons inc.", "You worked for 2 years since you joined the company Swallows inc."}
```

## `DynamicStruct`
We can create dynamic and runtime struct.

```go
// Add fields using Builder.
// We can use AddXXX method chain.
builder := dynamicstruct.NewBuilder().
	AddString("StringField").
	AddInt("IntField").
	AddFloat("FloatField").
	AddBool("BoolField").
	AddMap("MapField", dynamicstruct.SampleString, dynamicstruct.SampleFloat).
	AddChanBoth("ChanBothField", dynamicstruct.SampleInt).
	AddStructPtr("StructPtrField", hogePtr).
	AddSlice("SliceField", hogePtr)

// Remove removes a field by assigned name.
builder = builder.Remove("FloatField")

// Build generates a DynamicStruct
ds := builder.Build()

// DecodeMap decodes from map to DynamicStruct
input := map[string]interface{}{
	"StringField": "Test String Field",
	"IntField":    12345,
	"BoolField":   true,
}
dec, err := ds.DecodeMap(input)

// Confirm decoded result using Getter
g, err := structil.NewGetter(dec)

fmt.Printf("String: %v, Int: %v, Bool: %v\n", g.String("StringField"), g.Int("IntField"), g.Get("BoolField"))
```

Result as follows.

```
String: Test String Field, Int: 12345, Bool: true
```

### JSON unmershal with `DynamicStruct`

A decoding example from JSON to `DynamicStruct` with `StructTag` using `json.Unmarshal([]byte)` as follows.
This example works correctly not only JSON but also YAML, TOML and more.

```go
builder := dynamicstruct.NewBuilder().
	AddStringWithTag("StringField", `json:"string_field"`).
	AddIntWithTag("IntField", `json:"int_field"`).
	AddFloatWithTag("FloatField", `json:"float_field"`).
	AddBoolWithTag("BoolField", `json:"bool_field"`).
	AddStructPtrWithTag("StructPtrField", hogePtr, `json:"struct_ptr_field"`)

// Get interface of DynamicStruct using Interface() method
ds := builder.Build()
intf := ds.Interface()

// try json unmarshal
input := []byte(`
{
	"string_field":"あいうえお",
	"int_field":9876,
	"float_field":5.67,
	"bool_field":true,
	"struct_ptr_field":{
		"key":"hogekey",
		"value":"hogevalue"
	}
}
`)

err := json.Unmarshal(input, &intf)
if err != nil {
	// error handing
}

g, err := structil.NewGetter(intf)

fmt.Printf("String: %v, Float: %v, StructPtr: %+v\n", g.String("StringField"), g.Float64("FloatField"), g.Get("StructPtrField"))
```

Result as follows.

```
String: あいうえお, Float: 5.67, StructPtr: {Key:hogekey Value:hogevalue}
```
