package mocks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/callstats-io/ai-decision/service/src/storage"
)

// Storage implements a mock storage with identical methods as the actual storages
type Storage struct {
	mockCallCounts           map[string]int
	mockedErrors             map[string]error
	mockedMessageTemplates   []*storage.MessageTemplate
	mockedMessages           []*storage.Message
	mockedAidAnalyticsStates []*storage.AidAnalyticsState
}

// NewMockedStorage returns a new initilized storage mock
func NewMockedStorage() *Storage {
	return &Storage{
		mockCallCounts:           map[string]int{},
		mockedErrors:             map[string]error{},
		mockedMessageTemplates:   []*storage.MessageTemplate{},
		mockedMessages:           []*storage.Message{},
		mockedAidAnalyticsStates: []*storage.AidAnalyticsState{},
	}
}

// Reset clears the storage
func (s *Storage) Reset() {
	s.mockedErrors = map[string]error{}
	s.mockCallCounts = map[string]int{}
	s.mockedMessageTemplates = []*storage.MessageTemplate{}
	s.mockedMessages = []*storage.Message{}
	s.mockedAidAnalyticsStates = []*storage.AidAnalyticsState{}
}

// FetchMessageTemplatesCalls returns the number of FetchMessageTemplates calls
func (s *Storage) FetchMessageTemplatesCalls() int {
	return s.calls("FetchMessageTemplates")
}

// CreateMessageCalls returns the number of CreateMessage calls
func (s *Storage) CreateMessageCalls() int {
	return s.calls("CreateMessage")
}

// CreateMessageTemplateCalls returns the number of CreateMessageTemplate calls
func (s *Storage) CreateMessageTemplateCalls() int {
	return s.calls("CreateMessageTemplate")
}

// SaveStateCalls returns the number of SaveState calls
func (s *Storage) SaveStateCalls() int {
	return s.calls("SaveState")
}

// GetStateCalls returns the number of GetState calls
func (s *Storage) GetStateCalls() int {
	return s.calls("GetState")
}

// ListStatesCalls returns the number of ListStates calls
func (s *Storage) ListStatesCalls() int {
	return s.calls("ListStates")
}

// ListMessagesCalls returns the number of ListMessages calls
func (s *Storage) ListMessagesCalls() int {
	return s.calls("ListMessages")
}

// MockFetchMessageTemplatesError sets the FetchMessageTemplates mocked error
func (s *Storage) MockFetchMessageTemplatesError(err error) {
	s.mockError("FetchMessageTemplates", err)
}

// MockCreateMessageError sets the CreateMessage mocked error
func (s *Storage) MockCreateMessageError(err error) {
	s.mockError("CreateMessage", err)
}

// MockListMessagesError sets the ListMessages mocked error
func (s *Storage) MockListMessagesError(err error) {
	s.mockError("ListMessages", err)
}

// MockCreateMessageTemplateError sets the CreateMessageTemplate mocked error
func (s *Storage) MockCreateMessageTemplateError(err error) {
	s.mockError("CreateMessageTemplate", err)
}

// MockSaveStateError sets the SaveState mocked error
func (s *Storage) MockSaveStateError(err error) {
	s.mockError("SaveState", err)
}

// MockGetStateError sets the GetState mocked error
func (s *Storage) MockGetStateError(err error) {
	s.mockError("GetState", err)
}

// MockListStatesError sets the ListStates mocked error
func (s *Storage) MockListStatesError(err error) {
	s.mockError("ListStates", err)
}

// MockSavedMessageTemplates sets the message templates stored in mock
func (s *Storage) MockSavedMessageTemplates(states []*storage.MessageTemplate) {
	s.mockedMessageTemplates = states
}

// MockSavedMessages sets the mocked message templates to be returned by calls to Messages
func (s *Storage) MockSavedMessages(msgs []*storage.Message) {
	s.mockedMessages = msgs
}

// MockSavedStates sets the mocked message templates to be returned by calls to SavedStates
func (s *Storage) MockSavedStates(states []*storage.AidAnalyticsState) {
	s.mockedAidAnalyticsStates = states
}

// FetchMessageTemplates returns all mocked message templates for a given type up to max version
func (s *Storage) FetchMessageTemplates(ctx context.Context, mType string, maxVersion int32) ([]*storage.MessageTemplate, error) {
	s.called("FetchMessageTemplates")
	if err := s.mockedErrors["FetchMessageTemplates"]; err != nil {
		return nil, err
	}

	templates := []*storage.MessageTemplate{}
	for _, t := range s.mockedMessageTemplates {
		if t.Type == mType && t.Version <= maxVersion {
			t := t // copy pointer to ensure no leak
			templates = append(templates, t)
		}
	}
	return templates, nil
}

// CreateMessage returns an error if mocked
func (s *Storage) CreateMessage(ctx context.Context, msg *storage.Message) error {
	s.called("CreateMessage")
	if err := s.mockedErrors["CreateMessage"]; err != nil {
		return err
	}
	s.copy(s.mockedMessages[0], msg)
	return nil
}

// ListMessages returns an error if mocked
func (s *Storage) ListMessages(ctx context.Context, appID int32, keyword string, minVersion, maxVersion int32, from, to *time.Time) ([]*storage.Message, error) {
	s.called("ListMessages")
	if err := s.mockedErrors["ListMessages"]; err != nil {
		return nil, err
	}
	return s.mockedMessages, nil
}

// CreateMessageTemplate returns an error if mocked
func (s *Storage) CreateMessageTemplate(ctx context.Context, tmpl *storage.MessageTemplate) error {
	s.called("CreateMessageTemplate")
	if err := s.mockedErrors["CreateMessageTemplate"]; err != nil {
		return err
	}

	s.copy(s.mockedMessageTemplates[0], tmpl)
	return nil
}

// SaveState returns an error if mocked
func (s *Storage) SaveState(ctx context.Context, state *storage.AidAnalyticsState) error {
	s.called("SaveState")
	if err := s.mockedErrors["SaveState"]; err != nil {
		return err
	}
	s.copy(s.mockedAidAnalyticsStates[0], state)
	return nil
}

// GetState returns an error if mocked
func (s *Storage) GetState(ctx context.Context, state *storage.AidAnalyticsState) error {
	s.called("GetState")
	if err := s.mockedErrors["GetState"]; err != nil {
		return err
	}
	s.copy(s.mockedAidAnalyticsStates[0], state)
	return nil
}

// ListStates returns an error if mocked
func (s *Storage) ListStates(ctx context.Context, appID int32, keyword string, from, to *time.Time) ([]*storage.AidAnalyticsState, error) {
	s.called("ListStates")
	if err := s.mockedErrors["ListStates"]; err != nil {
		return nil, err
	}
	return s.mockedAidAnalyticsStates, nil
}

// calls returns the number of calls made to the given method since last reset
func (s *Storage) calls(method string) int {
	return s.mockCallCounts[method]
}

func (s *Storage) mockError(method string, err error) {
	s.mockedErrors[method] = err
}

func (s *Storage) called(method string) {
	s.mockCallCounts[method] = s.mockCallCounts["method"] + 1
}

func (s *Storage) copy(src, dst interface{}) {
	data, _ := json.Marshal(src)
	json.Unmarshal(data, dst)
}
