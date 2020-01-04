package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/callstats-io/ai-decision/service/gen/protos"
	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/require"
)

func TestMessageCreate(t *testing.T) {
	// ensure no leakage between tests
	defer mockStorage.Reset()

	const (
		tmplType  = "growth"
		value1Key = "value1"
		value2Key = "value2"
	)

	sharedTemplates := []*storage.MessageTemplate{
		&storage.MessageTemplate{ID: 1, Type: tmplType, Version: 1, Template: "tmpl with {{.Number \"" + value1Key + "\" }}", CreatedAt: time.Now().Add(-2 * time.Hour)},
		&storage.MessageTemplate{ID: 2, Type: tmplType, Version: 2, Template: "tmpl with {{.Number \"" + value1Key + "\" }} and {{.String \"" + value2Key + "\" }}", CreatedAt: time.Now().Add(-time.Hour)},
	}

	tests := []struct {
		Description string
		ExpErrorMsg string
		Setup       func(req *protos.MessageCreateRequest) (*protos.Message, error)
	}{
		{
			Description: "valid request with latest template version",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				version := sharedTemplates[len(sharedTemplates)-1].Version
				data, _ := json.Marshal(map[string]interface{}{
					value1Key: 123,
					value2Key: "awesomeness",
				})
				generatedAt, _ := ptypes.Timestamp(req.GenerationTime)
				mockStorage.Reset()
				mockStorage.MockSavedMessageTemplates(sharedTemplates)
				mockStorage.MockSavedMessages([]*storage.Message{
					{
						AppID:       req.AppId,
						TemplateID:  sharedTemplates[len(sharedTemplates)-1].ID,
						Template:    sharedTemplates[len(sharedTemplates)-1],
						Data:        data,
						GeneratedAt: generatedAt,
					},
				})

				req.Version = version
				req.Data = data
				return &protos.Message{
					AppId:          req.AppId,
					Type:           req.Type,
					Version:        req.Version,
					Data:           req.Data,
					GenerationTime: req.GenerationTime,
					Message:        "tmpl with 123 and awesomeness",
				}, nil
			},
		},
		{
			Description: "valid request with earlier template version",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				version := sharedTemplates[0].Version
				data, _ := json.Marshal(map[string]interface{}{
					value1Key: 123,
					value2Key: "awesomeness",
				})
				generatedAt, _ := ptypes.Timestamp(req.GenerationTime)
				mockStorage.Reset()
				mockStorage.MockSavedMessageTemplates(sharedTemplates)
				mockStorage.MockSavedMessages([]*storage.Message{
					{
						AppID:       req.AppId,
						TemplateID:  sharedTemplates[len(sharedTemplates)-1].ID,
						Template:    sharedTemplates[len(sharedTemplates)-1],
						Data:        data,
						GeneratedAt: generatedAt,
					},
				})

				req.Version = version
				req.Data = data
				return &protos.Message{
					AppId:          req.AppId,
					Type:           req.Type,
					Version:        req.Version,
					Data:           req.Data,
					GenerationTime: req.GenerationTime,
					Message:        "tmpl with 123",
				}, nil
			},
		},
		{
			Description: "missing template data",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = template: 2:1:38: executing \"2\" at <.String>: error calling String: invalid string",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessageTemplates(sharedTemplates)

				req.Version = int32(len(sharedTemplates))
				req.Data, _ = json.Marshal(map[string]interface{}{
					value1Key: 123,
				})

				return nil, nil
			},
		},
		{
			Description: "invalid template data",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = template: 2:1:12: executing \"2\" at <.Number>: error calling Number: invalid number",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessageTemplates(sharedTemplates)

				req.Version = int32(len(sharedTemplates))
				req.Data, _ = json.Marshal(map[string]interface{}{
					value1Key: "123",
					value2Key: "awesomeness",
				})

				return nil, nil
			},
		},
		{
			Description: "missing app id",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = app_id: must be a positive integer",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.AppId = 0
				return nil, nil
			},
		},
		{
			Description: "missing template type",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = type: cannot be empty",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.Type = ""
				return nil, nil
			},
		},
		{
			Description: "missing template version",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = version: must be a positive integer",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.Version = 0
				return nil, nil
			},
		},
		{
			Description: "missing template data",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = data: cannot be empty",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.Data = nil
				return nil, nil
			},
		},
		{
			Description: "invalid template data json",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = data: invalid character 'a' looking for beginning of value",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.Data = []byte("abc")
				return nil, nil
			},
		},
		{
			Description: "missing generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: cannot be nil",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.GenerationTime = nil
				return nil, nil
			},
		},
		{
			Description: "invalid generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: must have positive seconds",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				req.GenerationTime = &timestamp.Timestamp{Nanos: 123}
				return nil, nil
			},
		},
		{
			Description: "fetch message templates error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED FETCH TEMPLATES TEST ERROR",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockFetchMessageTemplatesError(errors.New("EXPECTED FETCH TEMPLATES TEST ERROR"))
				return nil, nil
			},
		},
		{
			Description: "create message error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED CREATE MESSAGE TEST ERROR",
			Setup: func(req *protos.MessageCreateRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessageTemplates(sharedTemplates)
				mockStorage.MockCreateMessageError(errors.New("EXPECTED CREATE MESSAGE TEST ERROR"))

				// set valid data
				req.Version = int32(len(sharedTemplates))
				req.Data, _ = json.Marshal(map[string]interface{}{
					value1Key: 123,
					value2Key: "awesomeness",
				})

				return nil, nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			genTime, _ := ptypes.TimestampProto(time.Now())
			// create a request
			req := &protos.MessageCreateRequest{
				AppId:          123,
				Type:           tmplType,
				Version:        1,
				GenerationTime: genTime,
				Data:           []byte(`{"abc": "def"}`), // json object with random garbage, expect Setup to overwrite where required
			}
			expMessage, err := test.Setup(req)
			assert.Nil(err)

			// exec test
			resp, err := testMessageClient.Create(context.Background(), req)
			if test.ExpErrorMsg != "" {
				assert.EqualError(err, test.ExpErrorMsg)
			} else {
				assert.Nil(err)
				if expMessage != nil {
					// Timestamp deep equality fails with the assertion library so validate them manually and reset to nil
					if expMessage.GenerationTime != nil {
						assert.NotNil(resp.GenerationTime)
						assert.Equal(expMessage.GenerationTime.Seconds, resp.GenerationTime.Seconds)
						assert.Equal(expMessage.GenerationTime.Nanos, resp.GenerationTime.Nanos)
						resp.GenerationTime = nil
						expMessage.GenerationTime = nil
					} else {
						assert.Nil(resp.GenerationTime)
					}
				}
				assert.Equal(expMessage, resp)
				assert.Equal(1, mockStorage.FetchMessageTemplatesCalls())
				assert.Equal(1, mockStorage.CreateMessageCalls())
			}
		})
	}
}

func TestMessageList(t *testing.T) {
	// ensure no leakage between tests
	defer mockStorage.Reset()
	payload := []byte(`{"abc":"def"}`)
	fixedTime := time.Now().Add(-5 * time.Minute)
	fixedProtoTime, _ := ptypes.TimestampProto(fixedTime)
	fixedTemplate := &storage.MessageTemplate{ID: 1, Type: "t-msg-list-type-1", Version: 1, Template: `{{.String "abc"}}`, CreatedAt: fixedTime}

	tests := []struct {
		Description string
		ExpErrorMsg string
		Setup       func(req *protos.MessageListRequest) (*protos.Message, error)
	}{
		{
			Description: "valid request with just app id",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// reset parts of request
				req.Type = ""
				req.GenerationTimeFrom = nil
				req.GenerationTimeTo = nil
				req.MinVersion = 0
				req.MaxVersion = 0

				return expMessage, nil
			},
		},
		{
			Description: "valid request with app id and type",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// reset parts of request
				req.GenerationTimeFrom = nil
				req.GenerationTimeTo = nil
				req.MinVersion = 0
				req.MaxVersion = 0

				return expMessage, nil
			},
		},
		{
			Description: "valid request with app id, type and from time",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// reset parts of request
				req.GenerationTimeTo = nil
				req.MinVersion = 0
				req.MaxVersion = 0

				return expMessage, nil
			},
		},
		{
			Description: "valid request with app id, type, from time and to time",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// reset parts of request
				req.MinVersion = 0
				req.MaxVersion = 0

				return expMessage, nil
			},
		},
		{
			Description: "valid request with app id, type, from time, to time and minVersion",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// reset parts of request
				req.MaxVersion = 0

				return expMessage, nil
			},
		},
		{
			Description: "valid request with app id, type, from time, to time, minVersiona and maxVersion",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockSavedMessages([]*storage.Message{
					{ID: 1, AppID: req.AppId, Template: fixedTemplate, TemplateID: fixedTemplate.ID, Data: payload, GeneratedAt: fixedTime},
				})
				expMessage := &protos.Message{
					AppId:          req.AppId,
					Type:           fixedTemplate.Type,
					Version:        fixedTemplate.Version,
					GenerationTime: fixedProtoTime,
					Data:           payload,
					Message:        "def",
				}

				// assume the 'req' to be valid by default and just return the appropriate state from it

				return expMessage, nil
			},
		},
		{
			Description: "no messages",
			ExpErrorMsg: "rpc error: code = NotFound desc = not found",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockListMessagesError(storage.ErrNotFound)
				return nil, nil
			},
		},
		{
			Description: "missing app id",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = app_id: must be a positive integer",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				req.AppId = 0
				return nil, nil
			},
		},
		{
			Description: "message list error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED MESSAGE LIST TEST ERROR",
			Setup: func(req *protos.MessageListRequest) (*protos.Message, error) {
				mockStorage.Reset()
				mockStorage.MockListMessagesError(errors.New("EXPECTED MESSAGE LIST TEST ERROR"))

				// assume request to be valid but don't return an expected message as we assume the request errors
				return &protos.Message{}, nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			genTimeFrom := *fixedProtoTime
			genTimeFrom.Seconds -= 60
			genTimeTo := *fixedProtoTime
			genTimeTo.Seconds += 60
			// create a valid request, expect Setup to invalidate if needed
			req := &protos.MessageListRequest{
				AppId:              2001,
				Type:               fmt.Sprintf("kw-test-state-list-%d", rand.Int()),
				GenerationTimeFrom: &genTimeFrom,
				GenerationTimeTo:   &genTimeTo,
			}
			expMessage, err := test.Setup(req)
			assert.Nil(err)

			// exec test
			stream, err := testMessageClient.List(context.Background(), req)
			assert.Nil(err)
			resp, err := stream.Recv()
			if test.ExpErrorMsg != "" {
				assert.EqualError(err, test.ExpErrorMsg)
			} else {
				assert.Nil(err)
				if expMessage != nil {
					// Timestamp deep equality fails with the assertion library so validate them manually and reset to nil
					if expMessage.GenerationTime != nil {
						assert.NotNil(resp.GenerationTime)
						assert.Equal(expMessage.GenerationTime.Seconds, resp.GenerationTime.Seconds)
						assert.Equal(expMessage.GenerationTime.Nanos, resp.GenerationTime.Nanos)
						resp.GenerationTime = nil
						expMessage.GenerationTime = nil
					} else {
						assert.Nil(resp.GenerationTime)
					}
				}
				assert.Equal(expMessage, resp)
				assert.Equal(1, mockStorage.ListMessagesCalls())

				// check no more values
				_, err = stream.Recv()
				assert.EqualError(err, "EOF")
			}
		})
	}
}
