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
func main() {
	go2types.CustomTypeMap = map[reflect.Kind]string{
		reflect.TypeOf(user.XTime{}).Kind(): "number", //XTime will be mapped to typescript type number
	}

	w := go2types.NewWorker()
	w.Namespace = "types"
	w.Add(types.T{}, types.User{})
	w.Add(user.Person{})
	w.MustGenerateFile("./types.ts")
}
```


## doc

todo
- doc
