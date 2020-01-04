package aws

import (
	"context"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/vault"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Client interface defines common AWS client functions
type Client interface {
	Session(ctx context.Context) (*session.Session, error)
}

// Assert that the StandardClient conforms to the Client interface
var _ = Client(&StandardClient{})

// StandardClient wraps Vault client and AWS session
type StandardClient struct {
	Vault   vault.AWSClient
	Options *Options
}

// NewStandardClient builds and returns a new Client
func NewStandardClient(vc vault.AWSClient, opts *Options) *StandardClient {
	return &StandardClient{
		Vault:   vc,
		Options: opts,
	}
}

// Session returns a new AWS session
func (c *StandardClient) Session(ctx context.Context) (*session.Session, error) {
	creds, err := c.Credentials(ctx)
	if err != nil {
		return nil, err
	}
	awsConf := aws.NewConfig().
		WithCredentials(creds).
		WithRegion(c.Options.AWSRegion).
		WithMaxRetries(c.Options.MaxRetries).
		WithHTTPClient(c.Options.HTTPClient)

	sess, err := session.NewSessionWithOptions(session.Options{Config: *awsConf})
	if err != nil {
		logAWSErr(ctx, LogAWSSessionEstablishFailed, err)
		return nil, err
	}
	if c.Options.DebugLogging {
		setupDebugLogging(ctx, sess)
	}
	return sess, nil
}

// Credentials returns AWS credentials, usually you shouldn't need this
func (c *StandardClient) Credentials(ctx context.Context) (*credentials.Credentials, error) {
	secret, err := c.Vault.AWSSecret(ctx)
	if err != nil {
		logger(ctx).Error(LogAWSCredsVaultFetchFailed, log.Error(err))
		return nil, err
	}
	creds := credentials.NewStaticCredentials(secret.Credentials.AccessKey, secret.Credentials.SecretKey, secret.Credentials.SecurityToken)
	return creds, nil
}

// Logger returns a logger from context with the package name
func logger(ctx context.Context) log.Logger {
	return log.FromContextWithPackageName(ctx, "go-common/aws")
}

func logAWSErr(ctx context.Context, msg string, err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		logger(ctx).
			With(log.String(LogKeyAWSErrorCode, awsErr.Code())).
			With(log.String(LogKeyAWSErrorMessage, awsErr.Message())).
			Error(msg, log.Error(err))
	} else {
		logger(ctx).Error(msg, log.Error(err))
	}
}

// setupDebugLogging adds a request logger for raw AWS client
func setupDebugLogging(ctx context.Context, sess *session.Session) {
	sess.Handlers.Send.PushFront(func(r *request.Request) {
		logger(ctx).
			With(log.String(LogKeyServiceName, r.ClientInfo.ServiceName)).
			With(log.String(LogKeyOperation, r.Operation.Name)).
			With(log.Object(LogKeyParams, r.Params)).
			Debug(LogAWSClientRequest)
	})
}
