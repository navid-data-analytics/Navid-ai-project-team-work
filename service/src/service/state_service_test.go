package service_test

import (
	"context"
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

func TestStateSave(t *testing.T) {
	// ensure no leakage between tests
	defer mockStorage.Reset()

	tests := []struct {
		Description string
		ExpErrorMsg string
		Setup       func(req *protos.StateSaveRequest) (*protos.State, error)
	}{
		{
			Description: "valid request",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				mockStorage.Reset()
				savedAt, _ := ptypes.Timestamp(req.GenerationTime)
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: req.Data, SavedAt: savedAt},
				})

				// assume the 'req' to be valid by default and just return the appropriate state from it

				return &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					Data:           req.Data,
					GenerationTime: req.GenerationTime,
				}, nil
			},
		},
		{
			Description: "missing app id",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = app_id: must be a positive integer",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				req.AppId = 0
				return nil, nil
			},
		},
		{
			Description: "missing keyword",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = keyword: cannot be empty",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				req.Keyword = ""
				return nil, nil
			},
		},
		{
			Description: "missing data",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = data: cannot be empty",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				req.Data = nil
				return nil, nil
			},
		},
		{
			Description: "missing generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: cannot be nil",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				req.GenerationTime = nil
				return nil, nil
			},
		},
		{
			Description: "invalid generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: must have positive seconds",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				req.GenerationTime = &timestamp.Timestamp{Nanos: 123}
				return nil, nil
			},
		},
		{
			Description: "state save error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED STATE SAVE TEST ERROR",
			Setup: func(req *protos.StateSaveRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockSaveStateError(errors.New("EXPECTED STATE SAVE TEST ERROR"))

				// assume request to be valid but don't return an expected message as we assume the request errors
				return &protos.State{}, nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			genTime, _ := ptypes.TimestampProto(time.Now())
			// create a valid request, expect Setup to invalidate if needed
			req := &protos.StateSaveRequest{
				AppId:          123,
				Keyword:        fmt.Sprintf("kw-%d", rand.Int()),
				GenerationTime: genTime,
				Data:           []byte(`{"abc": "def"}`),
			}
			expMessage, err := test.Setup(req)
			assert.Nil(err)

			// exec test
			resp, err := testStateClient.Save(context.Background(), req)
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
				assert.Equal(1, mockStorage.SaveStateCalls())
			}
		})
	}
}

func TestStateGet(t *testing.T) {
	// ensure no leakage between tests
	defer mockStorage.Reset()

	tests := []struct {
		Description string
		ExpErrorMsg string
		Setup       func(req *protos.StateGetRequest) (*protos.State, error)
	}{
		{
			Description: "valid request",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				mockStorage.Reset()
				savedAt, _ := ptypes.Timestamp(req.GenerationTime)
				payload := []byte(`{"abc":"def"}`)
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: payload, SavedAt: savedAt},
				})

				// assume the 'req' to be valid by default and just return the appropriate state from it

				return &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					GenerationTime: req.GenerationTime,
					Data:           payload,
				}, nil
			},
		},
		{
			Description: "missing app id",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = app_id: must be a positive integer",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				req.AppId = 0
				return nil, nil
			},
		},
		{
			Description: "missing keyword",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = keyword: cannot be empty",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				req.Keyword = ""
				return nil, nil
			},
		},
		{
			Description: "missing generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: cannot be nil",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				req.GenerationTime = nil
				return nil, nil
			},
		},
		{
			Description: "invalid generation time",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = generation_time: must have positive seconds",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				req.GenerationTime = &timestamp.Timestamp{Nanos: 123}
				return nil, nil
			},
		},
		{
			Description: "not found",
			ExpErrorMsg: "rpc error: code = NotFound desc = not found",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockGetStateError(storage.ErrNotFound)

				// assume request to be valid but don't return an expected message as we assume the request errors
				return &protos.State{}, nil
			},
		},
		{
			Description: "state get error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED STATE GET TEST ERROR",
			Setup: func(req *protos.StateGetRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockGetStateError(errors.New("EXPECTED STATE GET TEST ERROR"))

				// assume request to be valid but don't return an expected message as we assume the request errors
				return &protos.State{}, nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			genTime, _ := ptypes.TimestampProto(time.Now())
			// create a valid request, expect Setup to invalidate if needed
			req := &protos.StateGetRequest{
				AppId:          567,
				Keyword:        fmt.Sprintf("srv-state-get-kw-%d", rand.Int()),
				GenerationTime: genTime,
			}
			expMessage, err := test.Setup(req)
			assert.Nil(err)

			// exec test
			resp, err := testStateClient.Get(context.Background(), req)
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
				assert.Equal(1, mockStorage.GetStateCalls())
			}
		})
	}
}

func TestStateList(t *testing.T) {
	// ensure no leakage between tests
	defer mockStorage.Reset()

	payload := []byte(`{"abc":"def"}`)
	fixedTime := time.Now().Add(-5 * time.Minute)
	fixedProtoTime, _ := ptypes.TimestampProto(fixedTime)

	tests := []struct {
		Description string
		ExpErrorMsg string
		Setup       func(req *protos.StateListRequest) (*protos.State, error)
	}{
		{
			Description: "valid request with just app id",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: payload, SavedAt: fixedTime},
				})
				expState := &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					GenerationTime: fixedProtoTime,
					Data:           payload,
				}

				// reset parts of request
				req.Keyword = ""
				req.GenerationTimeFrom = nil
				req.GenerationTimeTo = nil

				return expState, nil
			},
		},
		{
			Description: "valid request with app id and keyword",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: payload, SavedAt: fixedTime},
				})
				expState := &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					GenerationTime: fixedProtoTime,
					Data:           payload,
				}

				// reset parts of request
				req.GenerationTimeFrom = nil
				req.GenerationTimeTo = nil

				return expState, nil
			},
		},
		{
			Description: "valid request with app id, keyword and from time",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: payload, SavedAt: fixedTime},
				})
				expState := &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					GenerationTime: fixedProtoTime,
					Data:           payload,
				}

				// reset parts of request
				req.GenerationTimeTo = nil

				return expState, nil
			},
		},
		{
			Description: "valid request with app id, keyword, from and to time",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockSavedStates([]*storage.AidAnalyticsState{
					{ID: 1, AppID: req.AppId, Keyword: req.Keyword, Data: payload, SavedAt: fixedTime},
				})
				expState := &protos.State{
					AppId:          req.AppId,
					Keyword:        req.Keyword,
					GenerationTime: fixedProtoTime,
					Data:           payload,
				}

				// assume the 'req' to be valid by default and just return the appropriate state from it

				return expState, nil
			},
		},
		{
			Description: "no states",
			ExpErrorMsg: "rpc error: code = NotFound desc = not found",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockListStatesError(storage.ErrNotFound)
				return nil, nil
			},
		},
		{
			Description: "missing app id",
			ExpErrorMsg: "rpc error: code = InvalidArgument desc = app_id: must be a positive integer",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				req.AppId = 0
				return nil, nil
			},
		},
		{
			Description: "state list error",
			ExpErrorMsg: "rpc error: code = Unavailable desc = EXPECTED STATE LIST TEST ERROR",
			Setup: func(req *protos.StateListRequest) (*protos.State, error) {
				mockStorage.Reset()
				mockStorage.MockListStatesError(errors.New("EXPECTED STATE LIST TEST ERROR"))

				// assume request to be valid but don't return an expected message as we assume the request errors
				return &protos.State{}, nil
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
			req := &protos.StateListRequest{
				AppId:              2001,
				Keyword:            fmt.Sprintf("kw-test-state-list-%d", rand.Int()),
				GenerationTimeFrom: &genTimeFrom,
				GenerationTimeTo:   &genTimeTo,
			}
			expMessage, err := test.Setup(req)
			assert.Nil(err)

			// exec test
			stream, err := testStateClient.List(context.Background(), req)
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
				assert.Equal(1, mockStorage.ListStatesCalls())

				// check no more values
				_, err = stream.Recv()
				assert.EqualError(err, "EOF")
			}
		})
	}
}
