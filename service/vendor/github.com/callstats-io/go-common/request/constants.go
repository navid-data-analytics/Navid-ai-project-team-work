package request

type ctxKey int

const (
	ctxKeyRequestID ctxKey = iota
)

// Constants
const (
	CsioClientIPHeader = "X-Csio-Client-Ip" // go uses canonical header names so make the constant canonical
)

// Log messages
const (
	LogErrFailedToCreateRequestID = "Failed to create unique request id"
)
