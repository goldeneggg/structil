__TODO list__

<!-- TOC depthFrom:1 -->

- [`Finder`](#finder)
- [`Getter`](#getter)
- [`DynamicStruct`](#dynamicstruct)
- [`GenericDecoder`](#genericdecoder)
- [Other](#other)

<!-- /TOC -->

## `Finder`
- [ ] may be deprecated `FindTop`
- [ ] add `Get(...names)` method
- [ ] performance tuning (by using `FieldByIndex` ?)

## `Getter`
- [x] add `Names` method
- [ ] current Getter is unsafe because panic is occured. So add safe Getter that returns an error instead of causing a panic ___This is large refactoring___
- [ ] performance tuning

## `DynamicStruct`
- [x] add `Definition` method that returns struce definition string of DynamicStruct
- [x] can assign the struct name of DynamicStruct
- [x] add GenericDecoder
- [ ] performance tuning

## `GenericDecoder`
- [ ]  support YAML
- [ ]  performance tuning

## Other
- [x] performance benchmark comparing on CI workflow with Github Actions
- [x] add usage image
- [ ] add reference information
  - [JSON\-to\-Go: Convert JSON to Go instantly](https://mholt.github.io/json-to-go/)
  - [golang は ゆるふわに JSON を扱えまぁす\! — KaoriYa](https://www.kaoriya.net/blog/2016/06/25/)
  - 
