package go2types

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

// some const
const (
	DefaultFieldTemplate = `{{.Name}}{{if .IsOptional}}?{{end}}: {{.TsType}}{{if .CanBeNull}} | null{{end}};{{if .Doc}}//{{.Doc}}{{end}}
` //a new line in the end
)

// Field field of struct
type Field struct {
	Template  string //render template
	Anomynous bool
	Omitted   bool //is field ignored, name start with lower case OR json:"-"

	Doc        string
	Name       string `json:"name"`
	TsType     string `json:"type"`
	CanBeNull  bool   `json:"canBeNull"`
	IsOptional bool   `json:"isOptional"` //Ptr type OR json omitempty
	IsDate     bool   `json:"isDate"`
	T          reflect.Type
	// for map[KeyType]ValType
	KeyType string `json:"keyType,omitempty"`
	ValType string `json:"valType,omitempty"`
}

// TagJSON .
type TagJSON struct {
	Exists    bool
	Omitted   bool
	Omitempty bool
	Name      string
}

func (t TagJSON) defaultIfNameEmpty(name string) string {
	if t.Name != "" {
		return t.Name
	}
	return name
}

// ParseField return: parsed field, isAnomynous, isStruct
func ParseField(sf reflect.StructField, go2tsTypes map[reflect.Kind]string) *Field {
	tagJSON := parseTagJSON(sf.Tag.Get("json"))

	typ := sf.Type

	f := Field{
		Anomynous:  sf.Anonymous,
		Doc:        sf.Tag.Get(DocTag),
		T:          typ,
		Omitted:    tagJSON.Omitted || hasLowerCasePrefix(sf.Name),
		Name:       tagJSON.defaultIfNameEmpty(sf.Name),
		IsOptional: tagJSON.Exists && tagJSON.Omitempty,
		CanBeNull:  typ.Kind() == reflect.Ptr && !tagJSON.Omitempty, //TODO consider map, slice...
		IsDate:     isDate(typ),
		// TsType will be setted later
	}

	k := typ.Kind()
	if v, ok := go2tsTypes[k]; ok {
		f.TsType = v
	} else {
		f.TsType = toTypescriptType(typ)
	}

	if !tagJSON.Exists {
		f.TsType = sf.Name
	}
	return &f
}

// MustRender .
func (f *Field) MustRender() string {
	t := f.Template
	if t == "" {
		t = DefaultFieldTemplate
	}
	tpl := template.Must(template.New("field_tpl").Parse(t))
	buffer := bytes.NewBuffer(nil)
	err := tpl.Execute(buffer, f)
	if err != nil {
		panic(fmt.Errorf("template execute error, %v", err))
	}
	return string(buffer.Bytes())
}
