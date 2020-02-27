package main

import (
	"github.com/dabankio/go2types"
	"github.com/dabankio/go2types/example/types"
	"github.com/dabankio/go2types/example/user"
	"reflect"
)

func main() {
	go2types.CustomTypeMap = map[reflect.Kind]string{
		reflect.TypeOf(user.XTime{}).Kind(): "number",
	}

	w := go2types.NewWorker()
	w.Namespace = "types"
	w.Add(types.T{}, types.User{})
	w.Add(user.Person{})
	w.MustGenerateFile("./types.ts")
}
