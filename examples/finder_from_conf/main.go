package main

import (
	"bytes"
	"fmt"

	"github.com/spf13/viper"

	"github.com/goldeneggg/structil"
)

const (
	DefaultYAML = `
Defs:
  - Company:
    - Group:
      - Boss
    - Period
  - Name
  - Age
`
)

type Definition struct {
	Defs []interface{}
}

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

type hoge struct {
	keys []string
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

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer([]byte(DefaultYAML)))

	var def Definition
	viper.Unmarshal(&def)
	fmt.Printf("%+v\n", def.Defs)

	finder, err := structil.NewFinder(i)
	if err != nil {
		panic(err)
	}

	hoge := &hoge{keys: make([]string, 0, len(def.Defs)+1)}
	for i, d := range def.Defs {
		fmt.Printf("----- %d is string. def: %s\n", i, d)
		hoge.setKeyRecursive(d, "")
	}
	fmt.Printf("hoge.keys: %+v\n", hoge.keys)

	m, err := finder.
		Find("Name").
		Into("Company").Find("Address").
		Into("Company", "Group").Find("Name", "Boss").
		ToMap()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found Map: %#v\n", m)
}

func (h *hoge) setKeyRecursive(d interface{}, prefix string) {
	var res string

	switch t := d.(type) {
	case string:
		res = t
		if prefix != "" {
			res = prefix + "." + res
		}
		h.keys = append(h.keys, res) // set here
	case map[interface{}]interface{}:
		var nk string
		var nd interface{}
		for key, value := range t {
			nk = key.(string)
			if prefix != "" {
				nk = prefix + "." + nk
			}
			nd = value
			break
		}
		h.setKeyRecursive(nd, nk)
	case []interface{}:
		for _, value := range t {
			h.setKeyRecursive(value, prefix)
		}
	default:
		panic(fmt.Sprintf("other. def: %#v, prefix: %s\n", t, prefix))
	}
}
