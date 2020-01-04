package testutil

import "net/http"

// FakeResponseWriter implements http.Writer and tracks calls to it
type FakeResponseWriter struct {
	HeaderData http.Header
	BodyData   []byte
	StatusCode int
}

// Header returns the fake header
func (f *FakeResponseWriter) Header() http.Header {
	if f.HeaderData == nil {
		f.HeaderData = make(http.Header)
	}
	return f.HeaderData
}

func (f *FakeResponseWriter) Write(data []byte) (int, error) {
	f.BodyData = data
	return len(data), nil
}

// WriteHeader stores the written status code to f.Status
func (f *FakeResponseWriter) WriteHeader(status int) {
	f.StatusCode = status
}

// Reset clears current state from the writer
func (f *FakeResponseWriter) Reset() {
	f.StatusCode = 0
	f.BodyData = nil
	f.HeaderData = nil
}
