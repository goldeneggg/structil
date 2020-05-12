<!-- TOC depthFrom:1 -->

- [`Finder`](#finder)
  - [With config file? use `FinderKeys`](#with-config-file-use-finderkeys)
- [`Getter`](#getter)
  - [`MapGet` method](#mapget-method)
- [`DynamicStruct`](#dynamicstruct)
  - [JSON unmershal with `DynamicStruct`](#json-unmershal-with-dynamicstruct)

<!-- /TOC -->


## `Finder`
We can access usefully nested struct fields using field name string.

See [example code](/examples/finder_simple/main.go)


### With config file? use `FinderKeys`
We can create a Finder from the configuration file that have some finding target keys. We support some file formats of configuration file such as `yaml`, `json`, `toml` and more.

See [example code](/examples/finder_simple/main.go)

___Thanks for the awesome configuration management library [spf13/viper](https://github.com/spf13/viper).___


## `Getter`
We can access a struct using field name string, like map.

See [example code](/examples/getter_simple/main.go)


### `MapGet` method
`MapGet` method provides the __Map__ collection function for slice of struct

See [example code](/examples/getter_map/main.go)


## `DynamicStruct`
We can create dynamic and runtime struct.

See [example code](/examples/dystruct_simple/main.go)


### JSON unmershal with `DynamicStruct`

A decoding example from JSON to `DynamicStruct` with `StructTag` using `json.Unmarshal([]byte)` as follows.
This example works correctly not only JSON but also YAML, TOML and more.

See [example code](/examples/dystruct_json_unmarshel/main.go)
