# Generic `GetAs` for `Getter`

This document explains the design and implementation of the generic `GetAs` function
that complements the non-generic `Getter` type in Go 1.19 environments.

## Background

In Go 1.19, methods on non-generic types **cannot** have their own type parameters.  Although Go 1.18 introduced generics,
the language spec only allows type parameters on functions or on generic types.  Attempting to define

```go
func (g *Getter) GetAs[T any](name string) (T, bool) { … }
```

results in a compilation error:

```
method must have no type parameters
```

To work around this in Go 1.19, we provide `GetAs` as a **free function** rather than a method.

## Implementation

```go
// GetAs retrieves the field named "name" from g and attempts to cast its value to T.
// The boolean result reports whether the field exists and the cast to T succeeded.
func GetAs[T any](g *Getter, name string) (T, bool) {
    var zero T
    gf, ok := g.getSafely(name)
    if !ok {
        return zero, false
    }
    res, ok := gf.intf.(T)
    return res, ok
}
```

Here, `g.getSafely` checks presence of the named field and returns its stored `interface{}` value.
The type assertion `gf.intf.(T)` then performs the cast to the generic type parameter.

## Usage Examples

Below snippets show how to use `GetAs` in practice:

```go
g, _ := structil.NewGetter(myStructPtr)

// Get an int field
if age, ok := structil.GetAs[int](g, "Age"); ok {
    fmt.Println("Age is", age)
}

// Get a string slice
if tags, ok := structil.GetAs[[]string](g, "Tags"); ok {
    fmt.Println("Tags", tags)
}

// Attempt nonexistent or mismatched type
if _, ok := structil.GetAs[bool](g, "Name"); !ok {
    fmt.Println("Name is not a bool or does not exist")
}
```

## Future Outlook

When the Go language eventually permits methods with type parameters on non-generic types (e.g., Go 1.21+),
the generic `GetAs` free function can be migrated to a method on `Getter` for more ergonomic calls.
