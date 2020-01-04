package aws

// Environment variable names
const (
	EnvAWSRegion = "AWS_REGION"
	EnvLogLevel  = "LOG_LEVEL"
)

// Defaults
const (
	DefaultAWSRegion = "eu-west-1"
)

// Log messages
const (
	LogAWSCredsVaultFetchFailed  = "Failed to get AWS credentials from Vault"
	LogAWSSessionEstablishFailed = "Failed to establish an AWS session"

	LogAWSClientRequest = "AWS client request"

	LogKeyServiceName = "serviceName"
	LogKeyOperation   = "operation"
	LogKeyParams      = "params"

	LogKeyAWSResponse     = "awsResponse"
	LogKeyAWSErrorCode    = "awsErrorCode"
	LogKeyAWSErrorMessage = "awsErrorMessage"

	LogKeyOriginalPayload = "originalPayload"

	LogSNSMessagePublished = "published a message to SNS"
	LogSNSPublishError     = "error while publishing a message to SNS"
	LogInvalidSNSMessage   = "invalid SNS message"

	LogMessagePublished = "published a message"
	LogPublishError     = "error while publishing a message"
	LogEmptyPayload     = "empty payload"
)
