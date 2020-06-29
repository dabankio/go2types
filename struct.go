package go2types

import (
	"bytes"
	"io"
	"path"
	"reflect"
	"strings"
	"text/template"
)

// default template
var (
	DefaultStructTemplate = `{{if .NotIgnored}}
{{if .Doc}}{{.Indent}}/*** {{.Doc}} */{{end}}
{{.Indent}}export interface {{.Name}} {{if.InheritedType}}extends {{.JoinInheritedTypes}}{{end}}{
{{range .Fields}}{{.MustRender}}
{{end}}{{.Indent}}}
{{end}}`

	DefaultEnumTemplate = `
{{if .Doc}}{{.Indent}}/**
{{.Doc}}
*/{{end}}
{{.Indent}}export enum {{.Name}} {
{{$e := .}}{{range .Values}}{{$e.Indent}}  {{.Name}} = '{{.Value}}',{{if .Doc}} // {{.Doc}}{{end}}
{{end}}{{.Indent}}}`
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
		Indent string //whole block indent
		Doc    string

		Type          Kind
		ReferenceName string
		Namespace     string
		Name          string
		Fields        []*Field
		InheritedType []string
		Values        []xenum
		T             reflect.Type
	}
)

// NotIgnored struct not ignored would be rendered
func (ctx *Struct) NotIgnored() bool {
	for _, typ := range IgnoreTypes {
		if ctx.T == typ {
			return false
		}
	}
	return true
}

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
		ret.Doc = getDoc(t.PkgPath(), t.Name(), docTypeType)
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

// RenderEnum .
func (s *Struct) RenderEnum(w io.Writer) (err error) {
	tpl, err := template.New("enum_tpl").Parse(DefaultEnumTemplate)
	if err != nil {
		return err
	}
	return tpl.Execute(w, s)
}
