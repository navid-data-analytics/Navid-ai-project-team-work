package testutil

import (
	"bytes"
	"strings"

	"github.com/callstats-io/go-common/log"
	"github.com/uber-go/zap"
)

// LogBuffer is a fake buffer useful for checking logging output
type LogBuffer struct {
	bytes.Buffer
}

// Sync overridden to a no-op as tests do not need to sync this anywhere
func (b *LogBuffer) Sync() error {
	return nil
}

// Lines returns logged lines
func (b *LogBuffer) Lines() []string {
	output := strings.Split(b.String(), "\n")
	return output[:len(output)-1]
}

// Stripped return strings without last line ending
func (b *LogBuffer) Stripped() string {
	return strings.TrimRight(b.String(), "\n")
}

// Logger returns a new log.Logger using this buffer
func (b *LogBuffer) Logger() log.Logger {
	dyn := zap.DynamicLevel()
	return &log.DynamicLogger{
		DynamicLevel: &dyn,
		Logger:       zap.New(zap.NewJSONEncoder(zap.NoTime()), zap.Output(b)),
	}
}

// NewLogBuffer creates a new logger whose output is saved to buffer
func NewLogBuffer() *LogBuffer {
	return &LogBuffer{}
}
