package go2typings

import (
	"reflect"
)

// its ok to modified (config)
var (
	DocField      = "doc"                     //find document (of struct) from struct field with name "doc"
	DocTag        = "doc"                     //read document (of struct or of field) from tag "doc"
	CustomTypeMap = map[reflect.Kind]string{} //go type <-> ts type
)
