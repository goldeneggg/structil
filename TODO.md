__TODO list__

<!-- TOC depthFrom:1 -->

- [`Finder`](#finder)
- [`Getter`](#getter)
- [`DynamicStruct`](#dynamicstruct)
- [Other](#other)

<!-- /TOC -->

## `Finder`
- [ ] performance tuning (by using `FieldByIndex` ?)
- [ ] may be deprecated `FindTop`
- [ ] add `Get(...names)` method

## `Getter`
- [x] add `Names` method
- [ ] current Getter is unsafe because panic is occured. So add safe Getter that returns an error instead of causing a panic ___This is large refactoring___


## `DynamicStruct`
- [ ]  performance tuning

## Other
- [ ] performance benchmark comparing on CI workflow with Github Actions
- [x] add usage image
- [ ] add reference information
  - [JSON\-to\-Go: Convert JSON to Go instantly](https://mholt.github.io/json-to-go/)
  - [golang は ゆるふわに JSON を扱えまぁす\! — KaoriYa](https://www.kaoriya.net/blog/2016/06/25/)
  - 
