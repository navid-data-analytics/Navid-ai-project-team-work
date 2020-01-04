package storage

import (
	"context"
	"errors"
	"time"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/postgres"
)

// Postgres defines a postgres backed storage
type Postgres struct {
	pgClient postgres.Client
}

// NewPostgres returns a new postgres backed storage with the specified client
func NewPostgres(pgClient postgres.Client) *Postgres {
	return &Postgres{
		pgClient: pgClient,
	}
}

// FetchMessageTemplates returns all message templates matching to a given type up to the specified version.
// If maxVersion is zero, all versions are returned.
func (s *Postgres) FetchMessageTemplates(ctx context.Context, mType string, maxVersion int32) ([]*MessageTemplate, error) {
	db, err := s.db(ctx)
	if err != nil {
		return nil, err
	}

	templates := []*MessageTemplate{}
	query := db.Model(&templates).Where("type = ?", mType)
	if maxVersion > 0 {
		query = query.Where("version <= ?", maxVersion)
	}
	if err := query.Select(); err != nil {
		return nil, err
	}

	if len(templates) == 0 {
		// assume at least one template has to always be found
		return nil, ErrNotFound
	}

	return templates, nil
}

// CreateMessage adds a new message to postgres. The message validation is expected to be performed before calling this function.
func (s *Postgres) CreateMessage(ctx context.Context, msg *Message) error {
	db, err := s.db(ctx)
	if err != nil {
		return err
	}
	if _, err := db.Model(msg).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}

// ListMessages fetches all message by app id.
// If message type is provided, all messages must additionally have the type of template
// If minVersion and/or maxVersion are provided, all messages must additionally be within the specified range (0 = beginning/end)
// If from and/or to are provided, all messages must additionally be within the specified range (nil = beginning/end)
func (s *Postgres) ListMessages(ctx context.Context, appID int32, mType string, minVersion, maxVersion int32, from, to *time.Time) ([]*Message, error) {
	db, err := s.db(ctx)
	if err != nil {
		return nil, err
	}

	var messages []*Message
	query := db.Model(&messages).
		Column("message.*", "Template").
		Relation("Template").
		Where("app_id = ?", appID)
	if mType != "" {
		query = query.Where("type = ?", mType)
	}
	if from != nil {
		query = query.Where("generated_at >= ?", from)
	}
	if to != nil {
		query = query.Where("generated_at <= ?", to)
	}
	if minVersion != 0 {
		query = query.Where("version >= ?", minVersion)
	}
	if maxVersion != 0 {
		query = query.Where("version <= ?", maxVersion)
	}
	if err := query.Select(); err != nil {
		return nil, err
	}
	if len(messages) == 0 {
		return nil, ErrNotFound
	}
	return messages, nil
}

// CreateMessageTemplate adds a new message template to postgres. The message template validation is expected to be performed before calling this function.
func (s *Postgres) CreateMessageTemplate(ctx context.Context, tmpl *MessageTemplate) error {
	db, err := s.db(ctx)
	if err != nil {
		return err
	}
	if _, err := db.Model(tmpl).Returning("*").Insert(); err != nil {
		return err
	}
	return nil
}

// SaveState saves the provided state to postgres.
// The message validation is expected to be performed before calling this function.
// If a conflicting state existed, it is overridden by the new state.
func (s *Postgres) SaveState(ctx context.Context, state *AidAnalyticsState) error {
	db, err := s.db(ctx)
	if err != nil {
		return err
	}
	query := db.Model(state).
		OnConflict("ON CONSTRAINT aid_analytics_states_keyword_idx DO UPDATE").
		Set("keyword = EXCLUDED.keyword, data = EXCLUDED.data").
		Returning("*")
	if _, err := query.Insert(); err != nil {
		return err
	}
	return nil
}

// GetState returns a state by app id, keyword and timestamp.
// The message validation is expected to be performed before calling this function.
func (s *Postgres) GetState(ctx context.Context, state *AidAnalyticsState) error {
	db, err := s.db(ctx)
	if err != nil {
		return err
	}

	if err := db.Model(state).Where("app_id = ?app_id AND keyword = ?keyword AND saved_at = ?saved_at").First(); err != nil {
		if err == postgres.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	return nil
}

// ListStates fetches all state by app id.
// If keyword is provided, all states must additionally match the keyword
// If from and/or to are provided, all states must additionally be within the specified range (nil = beginning/end)
func (s *Postgres) ListStates(ctx context.Context, appID int32, keyword string, from, to *time.Time) ([]*AidAnalyticsState, error) {
	db, err := s.db(ctx)
	if err != nil {
		return nil, err
	}

	var states []*AidAnalyticsState
	query := db.Model(&states).Where("app_id = ?", appID)
	if keyword != "" {
		query = query.Where("keyword = ?", keyword)
	}
	if from != nil {
		query = query.Where("saved_at >= ?", from)
	}
	if to != nil {
		query = query.Where("saved_at <= ?", to)
	}
	if err := query.Select(); err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nil, ErrNotFound
	}
	return states, nil
}

func (s *Postgres) db(ctx context.Context) (*postgres.DB, error) {
	db, err := s.pgClient.DB(ctx)
	if err != nil {
		log.FromContext(ctx).Error("failed to get db connection", log.Error(err))
		return nil, errors.New("failed to connect to database")
	}
	return db, nil
}
