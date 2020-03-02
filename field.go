package go2types

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"
)

// some const
const (
	DefaultFieldTemplate = `{{.Name}}{{if .IsOptional}}?{{end}}: {{.TsType}}{{if .CanBeNull}} | null{{end}};{{if .Doc}}//{{.Doc}}{{end}}`
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

	typ, kind := sf.Type, sf.Type.Kind()
	f := Field{
		Anomynous:  sf.Anonymous,
		T:          typ,
		Omitted:    tagJSON.Omitted || hasLowerCasePrefix(sf.Name),
		Name:       tagJSON.defaultIfNameEmpty(sf.Name),
		IsOptional: tagJSON.Exists && tagJSON.Omitempty,
		CanBeNull:  !tagJSON.Omitempty && (kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map),
		IsDate:     isDate(typ),
	}

	if v, ok := go2tsTypes[kind]; ok {
		f.TsType = v
	} else {
		f.TsType = toTypescriptType(typ)
	}

	if !tagJSON.Exists {
		f.TsType = sf.Name
	}

	for _, t := range DocTags {
		if v := sf.Tag.Get(t); v != "" {
			f.Doc = fmt.Sprintf("%s:%s, %s", t, v, f.Doc)
		}
	}
	return &f
}

// MustRender .
func (f *Field) MustRender() string {
	t := f.Template
	if t == "" {
		t = DefaultFieldTemplate
	}
	buffer := bytes.NewBuffer(nil)
	err := template.Must(template.New("field_tpl").Parse(t)).Execute(buffer, f)
	if err != nil {
		panic(fmt.Errorf("template execute error, %v", err))
	}
	return string(buffer.Bytes())
}
