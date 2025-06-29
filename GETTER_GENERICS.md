# ジェネリクスを活用した Getter の新しいメソッド `GetAs`

`structil` パッケージの `Getter` 型に、Go 1.18 から導入されたジェネリクスを活用した新しいメソッド `GetAs` を導入しました。

## `GetAs` メソッド

`GetAs` は、構造体のフィールドの値を、型パラメータで指定された任意の型として取得するための汎用的なメソッドです。

```go
func GetAs[T any](g *Getter, name string) (T, bool)
```

このメソッドは、従来の型ごとに用意されていた `GetString`, `GetInt`, `GetBool` などのメソッドを置き換えることを目的としています。

### 使用例

```go
package main

import (
	"fmt"
	"github.com/goldeneggg/structil"
)

type MyStruct struct {
	Name string
	Age  int
}

func main() {
	s := MyStruct{Name: "gopher", Age: 10}
	g, _ := structil.NewGetter(s)

	// string 型として Name フィールドの値を取得
	name, ok := structil.GetAs[string](g, "Name")
	if ok {
		fmt.Println(name) // "gopher"
	}

	// int 型として Age フィールドの値を取得
	age, ok := structil.GetAs[int](g, "Age")
	if ok {
		fmt.Println(age) // 10
	}
}
```

### `GetAs` の利点

- **汎用性**: `GetAs` はジェネリクスを使用しているため、任意の型に対応できます。これにより、`Getter` 型に新しい型のためのメソッドを追加する必要がなくなります。
- **コードの簡潔さ**: 型ごとにメソッドを呼び分ける必要がなくなり、コードがよりシンプルになります。
- **タイプセーフ**: コンパイル時に型チェックが行われるため、実行時の型エラーのリスクを低減できます。

### 今後の展望

将来的には、既存の `GetString`, `GetInt` などの型別メソッドは非推奨とし、`GetAs` への移行を推奨していく予定です。これにより、`structil` パッケージの API をより現代的で使いやすいものにしていきます。