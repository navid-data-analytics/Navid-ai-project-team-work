package message

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"text/template"

	"github.com/callstats-io/ai-decision/service/src/storage"
)

// TemplateData implements a convenience wrapper for a map[string]interface{} to validate data types
// The wrapping is needed to ensure an error is raised if the datatype is different from expected.
type TemplateData struct {
	values map[string]interface{}
}

// NewTemplateData returns a new *TemplateData initialized with the provided values
func NewTemplateData(values map[string]interface{}) *TemplateData {
	return &TemplateData{
		values: values,
	}
}

// UnmarshalTemplateData unmarshal a set of template values from the provided bytes
func UnmarshalTemplateData(data []byte) (*TemplateData, error) {
	var tmplValues map[string]interface{}
	if err := json.Unmarshal(data, &tmplValues); err != nil {
		return nil, err
	}
	return NewTemplateData(tmplValues), nil
}

// Number returns the value at key as a number.
// Currently JSON unmarshal to interface{} returns always a float64 for numbers.
func (d *TemplateData) Number(key string) (interface{}, error) {
	if v, ok := d.values[key].(float32); ok {
		return v, nil
	}
	if v, ok := d.values[key].(float64); ok {
		return v, nil
	}
	if v, ok := d.values[key].(int); ok {
		return v, nil
	}
	if v, ok := d.values[key].(int32); ok {
		return v, nil
	}
	if v, ok := d.values[key].(int64); ok {
		return v, nil
	}
	return 0, errors.New("invalid number")
}

// String returns the value at key as a string or an error
func (d *TemplateData) String(key string) (string, error) {
	if v, ok := d.values[key].(string); ok {
		return v, nil
	}
	return "", errors.New("invalid string")
}

// Date returns the value at key as date string in day month format or an error
func (d *TemplateData) Date(key string) (string, error) {
	// Cast to float64, as JSON number on the wire is float
	v, ok := d.values[key].(float64)
	if !ok {
		return "", errors.New("invalid timestamp value")
	}
	t := time.Unix(int64(v), 0)
	_, month, day := t.Date()
	return fmt.Sprintf("%v %v", day, month), nil
}

// Template implements a wrapper for text/template with a convenient helper for rendering to string
type Template struct {
	template    *template.Template
	buffer      *bytes.Buffer
	tmplType    string
	tmplVersion int32
}

// NewTemplate initialzes a new template from the given type, version and parsed text template
func NewTemplate(tmpl *storage.MessageTemplate) (*Template, error) {
	parsedTemplate, err := template.New(strconv.Itoa(int(tmpl.ID))).Parse(tmpl.Template)
	if err != nil {
		return nil, err
	}
	return &Template{
		template:    parsedTemplate,
		tmplType:    tmpl.Type,
		tmplVersion: tmpl.Version,
		buffer:      bytes.NewBuffer(make([]byte, 512)),
	}, nil
}

// RenderString returns the rendered value of this template as a string or an error
func (t *Template) RenderString(data *TemplateData) (string, error) {
	t.buffer.Reset()
	if err := t.template.Execute(t.buffer, data); err != nil {
		return "", err
	}
	return t.buffer.String(), nil
}

// Version returns this templates version
func (t *Template) Version() int32 { // func to ensure immutability after creation
	return t.tmplVersion
}

// Type returns this templates type
func (t *Template) Type() string { // func to ensure immutability after creation
	return t.tmplType
}
