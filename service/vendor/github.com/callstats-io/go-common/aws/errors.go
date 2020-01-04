package aws

import "errors"

// Error messages
var (
	ErrEmptyAWSAccessKey   = errors.New("AWS access key is empty")
	ErrEmptyAWSSecretKey   = errors.New("AWS secret key is empty")
	ErrEmptyAWSSecretToken = errors.New("AWS secret token is empty")

	ErrAWSSessionEstablishFailed = errors.New(LogAWSSessionEstablishFailed)
)

// SNS error messages
var (
	ErrEmptyPayload      = errors.New(LogEmptyPayload)
	ErrInvalidSNSMessage = errors.New(LogInvalidSNSMessage)
)
