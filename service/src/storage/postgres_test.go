package storage_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/callstats-io/ai-decision/service/src/storage"
	"github.com/callstats-io/go-common/postgres"
	"github.com/callstats-io/go-common/testutil"
	"github.com/stretchr/testify/require"
)

// TODO(SH): Full combinatiorial test cases, these are now best effort coverage

type badConnectionClient struct{}

func (c *badConnectionClient) DB(ctx context.Context) (*postgres.DB, error) {
	return nil, errors.New("FAKE BAD CONNECTION")
}

func (c *badConnectionClient) Close() {}

func TestFetchMessageTemplates(t *testing.T) {
	const (
		tmplType1     = "test_fa_tmpls_1"
		tmplType2     = "test_fa_tmpls_2"
		tmplTypeEmpty = "test_fa_tmpls_empty"
	)

	filterTemplates := func(templates []*storage.MessageTemplate, mType string, maxVersion int32) []*storage.MessageTemplate {
		ret := []*storage.MessageTemplate{}
		for _, t := range templates {
			if t.Type == mType && (maxVersion == 0 || t.Version <= maxVersion) {
				tmpl := t
				ret = append(ret, tmpl)
			}
		}
		return ret
	}
	testTemplates := []*storage.MessageTemplate{
		{Type: tmplType1, Version: 1, Template: "{.String \"val1\"}"},
		{Type: tmplType1, Version: 2, Template: "{.String \"val1\"} {.String \"val2\"}"},
		{Type: tmplType2, Version: 1, Template: "{.String \"val\"}"},
	}

	db, err := testPostgresClient.DB(testCtx)
	require.Nil(t, err)
	_, err = db.Model(&testTemplates).Returning("*").Insert()
	require.Nil(t, err)

	for _, test := range []struct {
		Description     string
		TemplateType    string
		TemplateVersion int32
		ExpErrMsg       string
		ExpTemplates    []*storage.MessageTemplate
		Storage         *storage.Postgres
	}{
		{
			Description:  "all templates of type 1",
			TemplateType: tmplType1,
			ExpTemplates: filterTemplates(testTemplates, tmplType1, 0),
			Storage:      storage.NewPostgres(testPostgresClient),
		},
		{
			Description:  "all templates of type 2",
			TemplateType: tmplType2,
			ExpTemplates: filterTemplates(testTemplates, tmplType2, 0),
			Storage:      storage.NewPostgres(testPostgresClient),
		},
		{
			Description:     "subset of templates",
			TemplateType:    tmplType1,
			TemplateVersion: 1,
			ExpTemplates:    filterTemplates(testTemplates, tmplType1, 1),
			Storage:         storage.NewPostgres(testPostgresClient),
		},
		{
			Description:  "no templates",
			TemplateType: tmplTypeEmpty,
			ExpErrMsg:    storage.ErrNotFound.Error(),
			Storage:      storage.NewPostgres(testPostgresClient),
		},
		{
			Description:  "fail if unable to connect",
			TemplateType: tmplType1,
			ExpErrMsg:    "failed to connect to database",
			Storage:      storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description:  "fail if query error",
			TemplateType: tmplType1,
			ExpErrMsg:    "pg: database is closed",
			Storage:      storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				tmpls, err := test.Storage.FetchMessageTemplates(ctx, test.TemplateType, test.TemplateVersion)
				if test.ExpErrMsg != "" {
					assert.EqualError(err, test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.Equal(test.ExpTemplates, tmpls)
				}
			}))
		})
	}
}
func TestCreateMessageTemplate(t *testing.T) {
	for _, test := range []struct {
		Description string
		Template    storage.MessageTemplate
		ExpErrMsg   string
		Storage     *storage.Postgres
	}{
		{
			Description: "valid message",
			Template:    storage.MessageTemplate{Type: "test_create_template_1", Version: 1, Template: "{.String \"val1\"}"},
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail type",
			Template:    storage.MessageTemplate{Version: 1, Template: "{.String \"val1\"}"},
			ExpErrMsg:   "null value in column \"type\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing version",
			Template:    storage.MessageTemplate{Type: "test_create_template_1", Template: "{.String \"val1\"}"},
			ExpErrMsg:   "null value in column \"version\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing template",
			Template:    storage.MessageTemplate{Type: "test_create_template_1", Version: 1},
			ExpErrMsg:   "null value in column \"template\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "fail if unable to connect",
			Template:    storage.MessageTemplate{Type: "test_create_template_2", Version: 1, Template: "{.String \"val1\"}"},
			ExpErrMsg:   "failed to connect to database",
			Storage:     storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description: "fail if query error",
			Template:    storage.MessageTemplate{Type: "test_create_template_3", Version: 1, Template: "{.String \"val1\"}"},
			ExpErrMsg:   "pg: database is closed",
			Storage:     storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				err := test.Storage.CreateMessageTemplate(ctx, &test.Template)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.NotEqual(0, test.Template.ID) // expect an ID to have been set
					storedTemplate := &storage.MessageTemplate{ID: test.Template.ID}
					assert.Nil(testPostgresDB.Select(storedTemplate))
					assert.Equal(&test.Template, storedTemplate)
				}
			}))
		})
	}
}
func TestCreateMessage(t *testing.T) {
	testTemplate := &storage.MessageTemplate{Type: "test_fa_create_tmpls_1", Version: 1, Template: "{.String \"val1\"}"}
	_, err := testPostgresDB.Model(testTemplate).Returning("*").Insert()
	require.Nil(t, err)

	validMessage := storage.Message{AppID: 123, TemplateID: testTemplate.ID, GeneratedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)}

	for _, test := range []struct {
		Description string
		Message     storage.Message
		ExpErrMsg   string
		Storage     *storage.Postgres
	}{
		{
			Description: "valid message",
			Message:     validMessage,
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing app id",
			Message:     storage.Message{TemplateID: testTemplate.ID, GeneratedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:   "null value in column \"app_id\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing template id",
			Message:     storage.Message{AppID: 123, GeneratedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:   "null value in column \"template_id\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail template does not exist",
			Message:     storage.Message{AppID: 123, TemplateID: -1, GeneratedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:   "insert or update on table \"messages\" violates foreign key constraint \"messages_template_id_fkey\"",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing generation time",
			Message:     storage.Message{AppID: 123, TemplateID: testTemplate.ID, Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:   "null value in column \"generated_at\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "postgres fail missing data",
			Message:     storage.Message{AppID: 123, TemplateID: testTemplate.ID, GeneratedAt: time.Now()},
			ExpErrMsg:   "null value in column \"data\" violates not-null constraint",
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "fail if unable to connect",
			Message:     validMessage,
			ExpErrMsg:   "failed to connect to database",
			Storage:     storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description: "fail if query error",
			Message:     validMessage,
			ExpErrMsg:   "pg: database is closed",
			Storage:     storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				err := test.Storage.CreateMessage(ctx, &test.Message)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.NotEqual(0, test.Message.ID) // expect an ID to have been set
					storedMessage := &storage.Message{ID: test.Message.ID}
					assert.Nil(testPostgresDB.Select(storedMessage))
					assert.Equal(&test.Message, storedMessage)
				}
			}))
		})
	}
}
func TestListMessages(t *testing.T) {
	const (
		app1            = int32(1567)
		app2            = int32(1678)
		type1           = "type-tlm-1567-1"
		type2           = "type-tlm-1567-2"
		typeNonExistent = "type-tlm-doesnotexist"
	)

	time1 := time.Now().Add(-5 * time.Minute)
	time2 := time.Now().Add(-3 * time.Minute)
	time3 := time.Now().Add(-time.Minute)
	timeNonExistentAfter := time.Now()
	timeNonExistentBefore := time.Now().Add(-10 * time.Minute)
	template := `{{.String "val1" }}`
	payload := []byte(`{"val1":"abc"}`)

	// note that changing these arrays is likely to break some of the existing tests as they lookup right from here
	// to validate requests
	createdMessageTemplates := []*storage.MessageTemplate{
		{Template: template, Type: type1, Version: 1},
		{Template: template, Type: type1, Version: 2},
		{Template: template, Type: type2, Version: 2}, // no version one for message type 2 to check maxVersion criteria
	}
	_, err := testPostgresDB.Model(&createdMessageTemplates).Returning("*").Insert()
	require.Nil(t, err)

	// set up messages so that app 1 has 1 message for each template and app2 one message for one template
	// to allow tests to check against app, time, type and version queries
	createdMessages := []*storage.Message{
		{AppID: app1, Template: createdMessageTemplates[0], TemplateID: createdMessageTemplates[0].ID, GeneratedAt: time1, Data: payload},
		{AppID: app1, Template: createdMessageTemplates[1], TemplateID: createdMessageTemplates[1].ID, GeneratedAt: time2, Data: payload},
		{AppID: app1, Template: createdMessageTemplates[2], TemplateID: createdMessageTemplates[2].ID, GeneratedAt: time3, Data: payload},
		{AppID: app2, Template: createdMessageTemplates[0], TemplateID: createdMessageTemplates[0].ID, GeneratedAt: time1, Data: payload},
	}
	_, err = testPostgresDB.Model(&createdMessages).Returning("*").Insert()
	require.Nil(t, err)

	for _, test := range []struct {
		Description string
		AppID       int32
		Type        string
		MinVersion  int32
		MaxVersion  int32
		From        time.Time
		To          time.Time
		ExpMessages []*storage.Message
		ExpErrMsg   string
		Storage     *storage.Postgres
	}{
		{
			Description: "list all app1 messages",
			AppID:       app1,
			ExpMessages: createdMessages[:len(createdMessages)-1], // skip all app2
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "list all app2 messages",
			AppID:       app2,
			ExpMessages: createdMessages[len(createdMessages)-1:], // skip all app1
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id and type",
			AppID:       app1,
			Type:        type1,
			ExpMessages: createdMessages[0:2], // assume only one record for app1 + type1
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id and time",
			AppID:       app1,
			From:        time1,
			To:          time2,
			ExpMessages: createdMessages[0:2], // assume app1 records all within time range
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id, type and time",
			AppID:       app1,
			Type:        type1,
			From:        time2,
			To:          time3,
			ExpMessages: createdMessages[1:2], // assume only a single record from time 1 to time 2 with type1
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records with type",
			AppID:       app1,
			Type:        typeNonExistent,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records with version above specified",
			AppID:       app1,
			Type:        type2,
			MinVersion:  int32(len(createdMessageTemplates)),
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records with version below specified",
			AppID:       app1,
			Type:        type2,
			MaxVersion:  1,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records after time",
			AppID:       app1,
			Type:        type1,
			From:        timeNonExistentAfter,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records before time",
			AppID:       app1,
			Type:        type1,
			To:          timeNonExistentBefore,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "fail if unable to connect",
			ExpMessages: []*storage.Message{},
			ExpErrMsg:   "failed to connect to database",
			Storage:     storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description: "fail if query error",
			ExpMessages: []*storage.Message{},
			ExpErrMsg:   "pg: database is closed",
			Storage:     storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				from, to := &test.From, &test.To
				if test.From.IsZero() {
					from = nil
				}
				if test.To.IsZero() {
					to = nil
				}
				messages, err := test.Storage.ListMessages(ctx, test.AppID, test.Type, test.MinVersion, test.MaxVersion, from, to)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.Equal(test.ExpMessages, messages)
					for _, m := range messages {
						assert.NotNil(m.Template) // verify template was preloaded correctly
					}
				}
			}))
		})
	}
}
func TestCreateAidAnalyticsState(t *testing.T) {
	validAidAnalyticsState := storage.AidAnalyticsState{AppID: 123, Keyword: fmt.Sprintf("kw-tss-%d", rand.Int()), SavedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)}
	duplicateAidAnalyticsState := validAidAnalyticsState
	duplicateAidAnalyticsState.Data = []byte(`{"val1":"abcd"}`)

	for _, test := range []struct {
		Description       string
		AidAnalyticsState storage.AidAnalyticsState
		ExpErrMsg         string
		Storage           *storage.Postgres
	}{
		{
			Description:       "valid message",
			AidAnalyticsState: validAidAnalyticsState,
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "valid duplicate app_id-time-keyword",
			AidAnalyticsState: duplicateAidAnalyticsState,
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "postgres fail missing app id",
			AidAnalyticsState: storage.AidAnalyticsState{Keyword: fmt.Sprintf("kw-tss-%d", rand.Int()), SavedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:         "null value in column \"app_id\" violates not-null constraint",
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "postgres fail missing keyword",
			AidAnalyticsState: storage.AidAnalyticsState{AppID: 123, SavedAt: time.Now(), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:         "null value in column \"keyword\" violates not-null constraint",
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "postgres fail missing generation time",
			AidAnalyticsState: storage.AidAnalyticsState{AppID: 123, Keyword: fmt.Sprintf("kw-tss-%d", rand.Int()), Data: []byte(`{"val1":"abc"}`)},
			ExpErrMsg:         "null value in column \"saved_at\" violates not-null constraint",
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "postgres fail missing data",
			AidAnalyticsState: storage.AidAnalyticsState{AppID: 123, Keyword: fmt.Sprintf("kw-tss-%d", rand.Int()), SavedAt: time.Now()},
			ExpErrMsg:         "null value in column \"data\" violates not-null constraint",
			Storage:           storage.NewPostgres(testPostgresClient),
		},
		{
			Description:       "fail if unable to connect",
			AidAnalyticsState: validAidAnalyticsState,
			ExpErrMsg:         "failed to connect to database",
			Storage:           storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description:       "fail if query error",
			AidAnalyticsState: validAidAnalyticsState,
			ExpErrMsg:         "pg: database is closed",
			Storage:           storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				err := test.Storage.SaveState(ctx, &test.AidAnalyticsState)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.NotEqual(0, test.AidAnalyticsState.ID) // expect an ID to have been set
					storedState := &storage.AidAnalyticsState{ID: test.AidAnalyticsState.ID}
					assert.Nil(testPostgresDB.Select(storedState))
					assert.Equal(&test.AidAnalyticsState, storedState)
				}
			}))
		})
	}
}
func TestGetAnalyticsState(t *testing.T) {
	const (
		app1       = int32(1000)
		app2       = int32(1001)
		kwBothApps = "kw-tgs-1000-1"
		kwApp2Only = "kw-tgs-1000-2"
	)

	timeExpected := time.Now().Add(-5 * time.Minute)
	timeNonExistentAfter := time.Now()
	timeNonExistentBefore := time.Now().Add(-10 * time.Minute)
	payload := []byte(`{"val1":"abc"}`)

	// note that changing this array is likely to break some of the existing tests as they lookup right from here
	// to validate requests
	createdStates := []*storage.AidAnalyticsState{
		{AppID: app1, Keyword: kwBothApps, SavedAt: timeExpected, Data: payload},
		{AppID: app2, Keyword: kwBothApps, SavedAt: timeExpected, Data: payload},
		{AppID: app2, Keyword: kwApp2Only, SavedAt: timeExpected, Data: payload},
	}

	_, err := testPostgresDB.Model(&createdStates).Returning("*").Insert()
	require.Nil(t, err)

	for _, test := range []struct {
		Description string
		QueryState  *storage.AidAnalyticsState
		ExpState    *storage.AidAnalyticsState
		ExpErrMsg   string
		Storage     *storage.Postgres
	}{
		{
			Description: "state exists",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwBothApps, SavedAt: timeExpected},
			ExpState:    createdStates[0],
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "state with keyword does not exist",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwApp2Only, SavedAt: timeExpected},
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "state with time after does not exist",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwApp2Only, SavedAt: timeNonExistentAfter},
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "state with time before does not exist",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwApp2Only, SavedAt: timeNonExistentBefore},
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "fail if unable to connect",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwBothApps, SavedAt: timeExpected},
			ExpErrMsg:   "failed to connect to database",
			Storage:     storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description: "fail if query error",
			QueryState:  &storage.AidAnalyticsState{AppID: app1, Keyword: kwBothApps, SavedAt: timeExpected},
			ExpErrMsg:   "pg: database is closed",
			Storage:     storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				err := test.Storage.GetState(ctx, test.QueryState)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.Equal(test.ExpState, test.QueryState)
				}
			}))
		})
	}
}
func TestListAidAnalyticsStates(t *testing.T) {
	const (
		app1          = int32(567)
		app2          = int32(678)
		kw1           = "kw-tls-567-1"
		kw2           = "kw-tls-567-2"
		kw3           = "kw-tls-678-1"
		kwNonExistent = "kw-tls-doesnotexist"
	)

	time1 := time.Now().Add(-5 * time.Minute)
	time2 := time.Now().Add(-time.Minute)
	timeNonExistentAfter := time.Now()
	timeNonExistentBefore := time.Now().Add(-10 * time.Minute)
	payload := []byte(`{"val1":"abc"}`)

	// note that changing this array is likely to break some of the existing tests as they lookup right from here
	// to validate requests
	createdStates := []*storage.AidAnalyticsState{
		{AppID: app1, Keyword: kw1, SavedAt: time1, Data: payload},
		{AppID: app1, Keyword: kw2, SavedAt: time1, Data: payload},
		{AppID: app1, Keyword: kw2, SavedAt: time2, Data: payload},
		{AppID: app2, Keyword: kw3, SavedAt: time1, Data: payload},
	}

	_, err := testPostgresDB.Model(&createdStates).Returning("*").Insert()
	require.Nil(t, err)

	for _, test := range []struct {
		Description string
		AppID       int32
		Keyword     string
		From        time.Time
		To          time.Time
		ExpStates   []*storage.AidAnalyticsState
		ExpErrMsg   string
		Storage     *storage.Postgres
	}{
		{
			Description: "list all app1 states",
			AppID:       app1,
			ExpStates:   createdStates[:len(createdStates)-1], // skip all app2
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "list all app2 states",
			AppID:       app2,
			ExpStates:   createdStates[len(createdStates)-1:], // skip all app1
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id and keyword",
			AppID:       app1,
			Keyword:     kw1,
			ExpStates:   createdStates[0:1], // assume only one record for app1 + kw1
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id and time",
			AppID:       app1,
			From:        time1,
			To:          time2,
			ExpStates:   createdStates[0:3], // assume app1 records all within time range
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "filter by app id, keyword and time",
			AppID:       app1,
			Keyword:     kw2,
			From:        time1,
			To:          time2,
			ExpStates:   createdStates[1:3], // assume only two records for app1 + kw2
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records with keyword",
			AppID:       app1,
			Keyword:     kwNonExistent,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records after time",
			AppID:       app1,
			From:        timeNonExistentAfter,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "no records before time",
			AppID:       app1,
			To:          timeNonExistentBefore,
			ExpErrMsg:   storage.ErrNotFound.Error(),
			Storage:     storage.NewPostgres(testPostgresClient),
		},
		{
			Description: "fail if unable to connect",
			ExpStates:   []*storage.AidAnalyticsState{},
			ExpErrMsg:   "failed to connect to database",
			Storage:     storage.NewPostgres(&badConnectionClient{}),
		},
		{
			Description: "fail if query error",
			ExpStates:   []*storage.AidAnalyticsState{},
			ExpErrMsg:   "pg: database is closed",
			Storage:     storage.NewPostgres(testPostgresClosedConnClient),
		},
	} {
		t.Run(test.Description, func(t *testing.T) {
			assert := require.New(t)

			assert.Nil(testutil.WithDeadlineContext(time.Second, func(ctx context.Context) {
				from, to := &test.From, &test.To
				if test.From.IsZero() {
					from = nil
				}
				if test.To.IsZero() {
					to = nil
				}
				states, err := test.Storage.ListStates(ctx, test.AppID, test.Keyword, from, to)
				if test.ExpErrMsg != "" {
					assert.NotNil(err)
					assert.Contains(err.Error(), test.ExpErrMsg)
				} else {
					assert.Nil(err)
					assert.Equal(test.ExpStates, states)
				}
			}))
		})
	}
}
