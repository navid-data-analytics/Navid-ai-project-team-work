package storage

import "time"

// MessageTemplate defines the structure of a message template as stored in postgres
type MessageTemplate struct {
	ID        int32
	Type      string
	Version   int32
	Template  string
	CreatedAt time.Time
}

// Message defines the structure of a message as stored in postgres
type Message struct {
	ID          int32
	AppID       int32
	TemplateID  int32
	Template    *MessageTemplate `pg:",fk:Template"`
	GeneratedAt time.Time
	Data        []byte
}

// AidAnalyticsState defines the structure of a message as stored in postgres
type AidAnalyticsState struct {
	ID      int32
	AppID   int32
	Keyword string
	Data    []byte
	SavedAt time.Time
}
