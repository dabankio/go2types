# Golang structs to typescript typings convertor

## TODO
- map type
- more customable namespace
- tests
- type embed (inherit)
- not generate some type

## Example

[example/main.go](example/main.go)

## How to setup

- create go file with the code bellow
- run this code with `go run`

```golang
package main

import (
  "github.com/dabankio/go2types"
   // you can use your own
  "github.com/dabankio/go2types/example/types"
)

type Root struct {
	User types.User
	T    types.T
}

func main() {
	s := go2types.New()
	s.Add(types.T{})
	s.Add(types.User{})

	err := s.GenerateFile("./test.ts")
	if err != nil {
		panic(err)
	}
}
```

# Custom tags

we support custom tag `ts` it has the following syntax

```
type M struct {
	Username string `json:"Username2" ts:"string,optional"`
}
```

tsTag type

```
tsTag[0] = "string"|"date"|"-"
tsTag[1] = "optional"|"no-null"|"null"
```

see field.go for more info
