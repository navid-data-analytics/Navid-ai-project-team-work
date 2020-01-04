package message_test

import (
	"testing"
	"time"

	"github.com/callstats-io/ai-decision/service/src/message"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/stretchr/testify/require"
)

var tmplIDCounter = int32(0)

func makeTemplate(t *testing.T, mType string, mVersion int32, mTemplate string) *message.Template {
	tmpl, err := message.NewTemplate(&storage.MessageTemplate{ID: mVersion, Type: mType, Version: mVersion, Template: mTemplate, CreatedAt: time.Now()})
	require.Nil(t, err)
	return tmpl
}

func TestTemplateRender(t *testing.T) {
	for _, test := range []struct {
		Description string
		Template    *message.Template
		Data        *message.TemplateData
		ExpMsg      string
		ExpErrMsg   string
	}{
		{
			Description: "valid data string",
			Template:    makeTemplate(t, "asd", 1, `{{.String "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": "def"}),
			ExpMsg:      "def",
		},
		{
			Description: "valid data int",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": 123}),
			ExpMsg:      "123",
		},
		{
			Description: "valid data int32",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": int32(123)}),
			ExpMsg:      "123",
		},
		{
			Description: "valid data int64",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": int64(123)}),
			ExpMsg:      "123",
		},
		{
			Description: "valid data float32",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": float32(123.12)}),
			ExpMsg:      "123.12",
		},
		{
			Description: "valid data float64",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": float64(123.12)}),
			ExpMsg:      "123.12",
		},
		{
			Description: "invalid data string",
			Template:    makeTemplate(t, "asd", 1, `{{.String "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": 123}),
			ExpErrMsg:   "template: 1:1:2: executing \"1\" at <.String>: error calling String: invalid string",
		},
		{
			Description: "invalid data number",
			Template:    makeTemplate(t, "asd", 1, `{{.Number "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": "def"}),
			ExpErrMsg:   "template: 1:1:2: executing \"1\" at <.Number>: error calling Number: invalid number",
		},
		{
			Description: "additional data",
			Template:    makeTemplate(t, "asd", 1, `{{.String "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": "def", "numnum": 456}),
			ExpMsg:      "def",
		},
		{
			Description: "timestamp data",
			Template:    makeTemplate(t, "asd", 1, `{{.Date "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": float64(1531785600.0)}),
			ExpMsg:      "17 July",
		},
		{
			Description: "invalid timestamp data",
			Template:    makeTemplate(t, "asd", 1, `{{.Date "val"}}`),
			Data:        message.NewTemplateData(map[string]interface{}{"val": ""}),
			ExpErrMsg:   "template: 1:1:2: executing \"1\" at <.Date>: error calling Date: invalid timestamp value",
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			rendered, err := test.Template.RenderString(test.Data)
			if test.ExpErrMsg != "" {
				assert.EqualError(err, test.ExpErrMsg)
			} else {
				assert.Nil(err)
				assert.Equal(test.ExpMsg, rendered)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	require.Equal(t, int32(1), makeTemplate(t, "", 1, "").Version())
	require.Equal(t, int32(2), makeTemplate(t, "", 2, "").Version())
}

func TestType(t *testing.T) {
	require.Equal(t, "abc", makeTemplate(t, "abc", 1, "").Type())
	require.Equal(t, "def", makeTemplate(t, "def", 2, "").Type())
}

func TestUnmarshalTemplateData(t *testing.T) {
	t.Run("valid json object", func(t *testing.T) {
		assert := require.New(t)
		val, err := message.UnmarshalTemplateData([]byte(`{"un":"abc"}`))
		assert.Nil(err)
		assert.Equal(message.NewTemplateData(map[string]interface{}{"un": "abc"}), val)
	})
	t.Run("invalid json object", func(t *testing.T) {
		assert := require.New(t)
		val, err := message.UnmarshalTemplateData([]byte(`[{"un":"abc"}]`))
		assert.Nil(val)
		assert.EqualError(err, "json: cannot unmarshal array into Go value of type map[string]interface {}")
	})
}
func TestBadtemplate(t *testing.T) {
	_, err := message.NewTemplate(&storage.MessageTemplate{ID: 1, Type: "abc", Version: 1, Template: "{{}", CreatedAt: time.Now()})
	require.EqualError(t, err, "template: 1:1: unexpected \"}\" in command")
}
