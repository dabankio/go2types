package go2types

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

// default template
const (
	DefaultStructTemplate = `
  //{{.T.PkgPath}}.{{.T.Name}}{{if .Doc}}
  /*** {{.Doc}} */{{end}}
  export interface {{.Name}} {{if.InheritedType}}extends {{.JoinInheritedTypes}}{{end}}{
{{range .Fields}}    {{.MustRender}}{{end}}  }
`

	DefaultEnumTemplate = `
  //{{.T.PkgPath}}.{{.T.Name}}
  export type {{.Name}} = {{.JoinEnumValues}};
`
)

// .
const (
	RegularType = iota
	Enum
)

type (
	// Kind .
	Kind int
	// Struct represent struct should be converted to typescript interface
	Struct struct {
		// Template          string //render template
		RenderFieldIndent string
		Doc               string

		Type          Kind
		ReferenceName string
		Namespace     string
		Name          string
		Fields        []*Field
		InheritedType []string
		Values        []string
		T             reflect.Type
	}
)

// MakeStruct .
func MakeStruct(t reflect.Type, name, namespace string) *Struct {
	if name == "" {
		name = t.Name()
	}
	if namespace == "" {
		namespace = path.Base(t.PkgPath())
	}

	fullName := strings.Title(name)
	ret := &Struct{
		Namespace:     namespace,
		Name:          fullName,
		ReferenceName: namespace + "." + fullName,
		InheritedType: []string{},
		T:             t,
	}

	if t.Kind() == reflect.Struct {
		if docField, ok := t.FieldByName(DocField); ok {
			ret.Doc = docField.Tag.Get(DocTag)
		}
	}
	return ret

}

// JoinInheritedTypes .
func (s *Struct) JoinInheritedTypes() string { return strings.Join(s.InheritedType, ", ") }

// MustRender render and panic on error
func (s *Struct) MustRender() string {
	buffer := bytes.NewBuffer(nil)
	var err error
	if s.Type == Enum {
		err = s.RenderEnum(buffer)
	} else {
		err = s.RenderTo(buffer)
	}
	panicIf(err)
	return string(buffer.Bytes())
}

// RenderTo .
func (s *Struct) RenderTo(w io.Writer) error {
	// TODO template customizable
	tpl, err := template.New("struct_tpl").Parse(DefaultStructTemplate)
	if err != nil {
		return err
	}
	return tpl.Execute(w, s)
}

// JoinEnumValues .
func (s *Struct) JoinEnumValues() string {
	var quoted []string
	for _, v := range s.Values {
		quoted = append(quoted, strconv.Quote(v))
	}
	return strings.Join(quoted, " | ")
}

// RenderEnum .
func (s *Struct) RenderEnum(w io.Writer) (err error) {
	tpl, err := template.New("enum_tpl").Parse(DefaultEnumTemplate)
	if err != nil {
		return err
	}
	return tpl.Execute(w, s)
}
