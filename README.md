# Golang template engine with Js syntax !!

# Dependency
```
github.com/CoinSummer/go-template@v0.2.0
```
# Features
1. Custom operator and functions
2. Interpolate from env
3. Expressions only, easy to use

# Usage

## Hello world
```go
package main

import (
	"fmt"

	gt "github.com/CoinSummer/go-template"
)

func main() {
	tp := gt.NewTemplate("Everyone knows {$a.b + $a.b} == { round($a.b + $a.b, 1)}", nil)
	res, err := tp.Render( `{"a": {"b": 1}}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(res) // Everyone knows 2 == 2
}
```

## With config
```go
package main

import (
	"fmt"

	gt "github.com/CoinSummer/go-template"
)

func main() {
	tp := gt.NewTemplateWithConfig("Everyone knows {$a.b + $a.b} == { round($a.b + $a.b, 1)}", nil, &gt.TemplateConfig{
		TimeOffset: 8,
		TimeFormat: time.RFC3339,
    })
	res, err := tp.Render( `{"a": {"b": 1}}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(res) // Everyone knows 2 == 2
}
```

## Custom operator
```go
package main

import (
	"fmt"

	gt "github.com/CoinSummer/go-template"
	"github.com/shopspring/decimal"
)

func main() {
	engine := gt.NewTemplateEngine()
	engine.OperatorsMgr.RegisterFunc("%", func(arg1, arg2 interface{}) (interface{}, error) {
		a, ok := arg1.(decimal.Decimal)
		if !ok {
			return nil, fmt.Errorf("%% with NaN: %v", a)
		}
		b, ok := arg2.(decimal.Decimal)
		if !ok {
			return nil, fmt.Errorf("%% to NaN: %v", b)
		}
		return a.Mod(b), nil
	})

	tp := gt.NewTemplate("100 % 3 = {100 % 3}", engine)
	res, err := tp.Render(``)
	if err != nil {
		panic(err)
	}
	fmt.Println(res) 
}
```


