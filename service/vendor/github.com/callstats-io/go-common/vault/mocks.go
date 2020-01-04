package vault

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/hashicorp/vault/api"
)

// MockReader creates a new mocked vault reader
type MockReader struct {
	readCalls int
	secret    *StandardSecret
	err       error
	release   chan bool
}

// Read returns either the mocked error or, if nil, the mocked secret and tracks the number of calls made
func (r *MockReader) Read(ctx context.Context, path string) (*StandardSecret, error) {
	r.readCalls++
	if r.release != nil {
		<-r.release
	}

	if r.err != nil {
		return nil, r.err
	}
	return r.secret, nil
}

// Reset resets the internal state of this mocked reader
func (r *MockReader) Reset() {
	r.secret.Cancel()
	r.readCalls = 0
	r.secret = nil
	r.err = nil

	// if in block mode close the release channel to release all blocked operations
	if r.release != nil {
		close(r.release)
		r.release = nil
	}
}

// MockSecret mocks a new standard secret with the provided data
func (r *MockReader) MockSecret(data map[string]interface{}) {
	r.secret = NewStandardSecret(&api.Secret{
		LeaseID:       fmt.Sprintf("ABCDEF%d", rand.Intn(100)),
		LeaseDuration: rand.Intn(600) + 600,
		Data:          data,
	}, nil)
}

// MockError sets the mock reader to return the provided error
func (r *MockReader) MockError(err error) {
	r.err = err
}

// BlockReads sets the reader into a mode where it won't return the secret and/or error until Release is called.
func (r *MockReader) BlockReads() {
	if r.release != nil {
		return
	}
	r.release = make(chan bool)
}

// ReleaseOnce releases the blocked channel once. It is a no-op if the reader is not in blocked mode
func (r *MockReader) ReleaseOnce() {
	go func() {
		defer func() {
			recover() // ensure write on closed channel doesn't break
		}()
		if r.release != nil {
			r.release <- true
		}
	}()
}
