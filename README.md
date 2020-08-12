structil [![PkgGoDev](https://pkg.go.dev/badge/github.com/goldeneggg/structil)](https://pkg.go.dev/github.com/goldeneggg/structil)
==========

[![Workflow Status](https://github.com/goldeneggg/structil/workflows/CI/badge.svg)](https://github.com/goldeneggg/structil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/goldeneggg/structil)](https://goreportcard.com/report/github.com/goldeneggg/structil)
[![Codecov](https://codecov.io/github/goldeneggg/structil/coverage.svg?branch=master)](https://codecov.io/github/goldeneggg/structil?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-lightgrey.svg)](https://github.com/goldeneggg/structil/blob/master/LICENSE)

struct + util = __structil__, for runtime and dynamic environment in Go.


## Why?

I'd like to ...

- conveniently dive into the specific field in nested struct
- simply verify if a field with the specified name and type exists in object
- conveniently handle known ___and unknown___ formatted JSON/YAML
- etc

```
*** JSON and YAML format is known or unknown ***


JSON →↓        →→ (known case) struct  →→→→→→→→→↓→→ (use struct directly)
      ↓        ↑                                ↓
      ↓→→ map →→→ (unknown case) DynamicStruct →→→ Getter, Finder
      ↑
YAML →↑
```

Please see [my medium post](https://medium.com/@s0k0mata/dynamic-and-runtime-struct-utilities-in-go-go-golang-reflection-25c154335185) as well.


## Examples
are [here](/examples).


## Benchmark
See [this file](https://github.com/goldeneggg/structil/blob/bench-latest/BENCHMARK_LATEST.txt)

It's the latest benchmark result that is executed on GitHub Actions runner instance.
