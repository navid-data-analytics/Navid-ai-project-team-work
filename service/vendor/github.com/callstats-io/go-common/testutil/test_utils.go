package testutil

import "testing"

// MustBeNil verifies that the error is nil. If not it calls t.Errorf with message and t.FailNow to stop execution.
func MustBeNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected error to be nil, got: %s", err)
		t.FailNow()
	}
}
