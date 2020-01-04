package http_test

import (
	"context"
	"encoding/json"
	"errors"

	"net/http/httptest"
	"testing"

	"github.com/callstats-io/ai-decision/service/src/config"
	"github.com/callstats-io/ai-decision/service/src/http"
	"github.com/stretchr/testify/require"
)

func TestStatus(t *testing.T) {
	for _, test := range []struct {
		Name      string
		StatusErr error
		Status    string
	}{
		{Name: "up", StatusErr: nil, Status: http.StatusUp},
		{Name: "down", StatusErr: errors.New("EXP ERR"), Status: http.StatusDown},
	} {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			assert := require.New(t)

			w := httptest.NewRecorder()
			called := false
			http.NewInternalRequestRouter(nil, func(ctx context.Context) error {
				called = true
				return test.StatusErr
			}).ServeHTTP(w, httptest.NewRequest("GET", "/status", nil))
			assert.True(called)
			actual := http.Status{}
			assert.Nil(json.Unmarshal(w.Body.Bytes(), &actual))
			expected := http.Status{
				Status:    test.Status,
				BuildTime: config.ImageBuildTime,
				Version:   config.ServiceVersion,
				StartTime: actual.StartTime,
			}
			assert.Equal(expected, actual)
			assert.NotEmpty(actual.StartTime)
		})
	}
}
