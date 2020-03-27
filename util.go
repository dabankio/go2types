package go2types

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"path"
	"reflect"
	"strings"

	"golang.org/x/tools/go/packages"
)

func parseTagJSON(tag string) TagJSON {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 {
		return TagJSON{Exists: false}
	}
	if parts[0] == "-" {
		return TagJSON{Exists: true, Omitted: true}
	}
	return TagJSON{
		Exists:    true,
		Omitted:   false,
		Omitempty: len(parts) > 1 && parts[1] == "omitempty",
		Name:      parts[0],
	}
}

func panicIf(e error) {
	if e != nil {
		panic(e)
	}
}

func isNumber(k reflect.Kind) bool {
	switch k {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func indirect(t reflect.Type) reflect.Type {
	k := t.Kind()
	for k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}
	return t
}

func isDate(t reflect.Type) bool {
	return t.Name() == "Time" && t.PkgPath() == "time"
}

func inArray(val int, array []*Struct) bool {
	return len(array) > val && val > 0
}

func isEnum(t reflect.Type) bool {
	return t.PkgPath() != ""
}

type constantValueDoc struct {
	name string
	constant.Value
	doc string
}

func getEnumValues(pkgName, typename string) ([]constantValueDoc, error) {
	res, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedName | packages.LoadSyntax | packages.NeedSyntax,
	}, pkgName)
	if err != nil {
		return nil, err
	}
	enums := []constantValueDoc{}
	if len(res) > 1 {
		return nil, errors.New("more than one result package")
	}

	docMap := map[string]string{} //doc map, <enumName, doc>
	for _, f := range res[0].Syntax {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						// fmt.Println("[dbg] try get doc", spec.Name, toJSON(spec.Doc))
					case *ast.ValueSpec:
						if spec.Doc.Text() != "" {
							// fmt.Println("[dbg] try get doc", typename, spec.Names, spec.Doc.Text())
							for _, name := range spec.Names {
								docMap[name.Name] = strings.TrimSpace(spec.Doc.Text())
							}
						}
					}
				}
			case *ast.FuncDecl:
				// fmt.Println("[dbg] try get doc:", toJSON(decl.Doc))
			}
		}
	}

	pkg := res[0].Types.Scope()
	// fmt.Println("dbg pkg.Names", pkg.Names())
	// fmt.Println("[dbg] info.Types ", info.Types)
	for _, name := range pkg.Names() {
		v := pkg.Lookup(name)
		// It has format similar to this "type.T".
		baseTypename := path.Base(v.Type().String())
		if v != nil && baseTypename == typename {
			switch t := v.(type) {
			case *types.Const:
				enums = append(enums, constantValueDoc{name: name, Value: t.Val(), doc: docMap[name]})
			}
		}
	}
	return enums, nil
}

var jsonRawMessageType = reflect.TypeOf((*json.RawMessage)(nil)).Elem()

// Call this func on each step of type processing.
// This func returns type string representation
func toTypescriptType(t reflect.Type) string {
	k := t.Kind()
	switch {
	case k == reflect.Ptr:
		t = indirect(t)
		return toTypescriptType(t)
	case k == reflect.Struct:
		if isDate(t) {
			return "string"
		}
		return t.Name()
	case isNumber(k) && isEnum(t):
		// TODO validate t.Name() is ok
		return t.Name()
	case isNumber(k):
		return "number"
	case k == reflect.String && isEnum(t):
		// TODO validate t.Name() is ok
		return t.Name()
	case k == reflect.String:
		return "string"
	case k == reflect.Bool:
		return "boolean"
	case k == reflect.Slice || k == reflect.Array:
		return fmt.Sprintf("Array<%s>", toTypescriptType(t.Elem()))
	case k == reflect.Interface || t == jsonRawMessageType:
		return "any"
	case k == reflect.Map:
		return "<TODO_MAP_TYPE_NAME>" //TODO
		// KeyType, ValType := toTypescriptType(t.Key(), getTypeName), toTypescriptType(t.Elem(), getTypeName)
		// return fmt.Sprintf("Record<%s, %s>", KeyType, ValType)
	}
	return t.String()
}

// x enum
type xenum struct {
	Name, Value, Doc string
}

func getEnumStringValues(t reflect.Type, pkgNames ...string) []xenum {
	if len(pkgNames) == 0 {
		pkgNames = []string{t.PkgPath()}
	}

	enumStrValues := []xenum{}
	for _, pkg := range pkgNames {
		values, err := getEnumValues(pkg, t.String())
		if err != nil {
			panic(err)
		}
		for _, xv := range values {
			v := xv.Value
			reflectValue := reflect.New(t).Elem()
			newVal := constant.Val(v)
			switch t.Kind() {
			case reflect.String:
				reflectValue.SetString(constant.StringVal(v))
			case reflect.Int:
				value, ok := constant.Int64Val(v)
				if !ok {
					panic("failed to convert")
				}
				reflectValue.SetInt(value)
			default:
				fmt.Println(reflect.TypeOf(newVal), newVal, reflectValue, v.Kind(), t)
				panic("unknown type")
			}
			strVal := fmt.Sprintf("%v", reflectValue)

			x := xenum{strings.Trim(strVal, "\""), strVal, xv.doc}
			for _, num := range "0123456789" {
				if strings.HasPrefix(strVal, string(num)) {
					x.Name = xv.name
					break
				}
			}

			enumStrValues = append(enumStrValues, x)
		}
	}
	return enumStrValues
}

func hasLowerCasePrefix(s string) bool {
	if s == "" {
		return false
	}
	return s[0] >= 'a' && s[0] <= 'z' // return strings.ToLower(s[:1]) == s[:1]
}

func structFieldTags(sf reflect.StructField) string {
	var tags []string
	for _, tag := range DocTags {
		if t := sf.Tag.Get(tag); t != "" {
			tags = append(tags, fmt.Sprintf("%s:%v", tag, t))
		}
	}
	return strings.Join(tags, ", ")
}
