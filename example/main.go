package main

import (
	"example/types"
	"example/user"
	"local/go2types"

	"reflect"
)

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
