__TODO list__

<!-- TOC depthFrom:1 -->

- [`DynamicStruct`](#dynamicstruct)
- [`Decoder`](#decoder)
- [`Getter`](#getter)
- [`Finder`](#finder)
- [Other](#other)

<!-- /TOC -->

## `DynamicStruct`
- [x] add `Definition` method that returns struce definition string of DynamicStruct
- [x] can assign the struct name of DynamicStruct
- [x] add GenericDecoder
- [ ] performance tuning

## `Decoder`
- [ ] support YAML
- [ ] support TOML
- [ ] support XML
- [ ] performance tuning

## `Getter`
- [x] add `Names` method
- [x] current Getter is unsafe because panic is occured. So add safe Getter that returns an error instead of causing a panic ___This is large refactoring___
- [x] performance tuning

## `Finder`
- [ ] may be deprecated `FindTop`
- [ ] add `Get(...names)` method
- [ ] performance tuning (by using `FieldByIndex` ?)

## Other
- [x] performance benchmark comparing on CI workflow with Github Actions
- [x] add usage image
- [ ] SQL type conversion
- [ ] add reference information
  - [JSON\-to\-Go: Convert JSON to Go instantly](https://mholt.github.io/json-to-go/)
  - [golang は ゆるふわに JSON を扱えまぁす\! — KaoriYa](https://www.kaoriya.net/blog/2016/06/25/)
