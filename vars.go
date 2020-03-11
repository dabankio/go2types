package go2types

import (
	"reflect"
	"time"
)

// its ok to modified (config)
var (
	DocField = "doc"           //find document (of struct) from struct field with name "doc"
	DocTags  = []string{"doc"} //read document (of struct or of field) from tag "doc"

	IgnoreTypes   = []reflect.Type{} //thoese types will not render
	CustomTypeMap = map[reflect.Type]string{
		reflect.TypeOf(time.Time{}):  "string",
		reflect.TypeOf(&time.Time{}): "string",
	} //go type <-> ts type
	DefaultFileDoc = `// DO NOT EDIT.
// Generated by go2types (https://github.com/dabankio/go2types)
// tslint:disable
// eslint-disable`
)
