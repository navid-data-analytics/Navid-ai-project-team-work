package vault

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/callstats-io/go-common/log"
	"github.com/hashicorp/vault/api"
)

// Client is the interface a vault client that can do any supported operation implements
type Client interface {
	VaultHTTPClient() *http.Client
	Options() *Options
	SecretReader
	SecretWriter
	AuthClient
	MongoClient
	PostgresClient
	TLSCertClient
	AWSClient
}

// SecretReader is the interface client which are able to read arbitrary vault secrets implement
type SecretReader interface {
	Read(ctx context.Context, path string) (*StandardSecret, error)
}

// SecretWriter is the interface client which are able to write arbitrary vault secrets implement
type SecretWriter interface {
	Write(ctx context.Context, path string, data map[string]interface{}) (*StandardSecret, error)
}

// AuthClient is the interface which clients that authenticate implement
// DEPRECATED, use SecretReader interface and raw secret data instead
type AuthClient interface {
	Authenticate(ctx context.Context) (*StandardSecret, error)
}

// MongoClient is the interface clients which are able to give out mongo secrets implement
// DEPRECATED, use SecretReader interface and raw secret data instead
type MongoClient interface {
	MongoSecret(ctx context.Context) (*UserPassSecret, error)
}

// PostgresClient is the interface clients which are able to give out Postgres secrets implement
// DEPRECATED, use SecretReader interface and raw secret data instead
type PostgresClient interface {
	PostgresSecret(ctx context.Context) (*UserPassSecret, error)
}

// TLSCertClient is the interface clients which are able to give out tls certificate secrets implement
// DEPRECATED, use SecretReader interface and raw secret data instead
type TLSCertClient interface {
	TLSCertSecret(ctx context.Context) (*TLSCertSecret, error)
}

// AWSClient is the interface clients which are able to give out aws credentials secrets implement
// DEPRECATED, use SecretReader interface and raw secret data instead
type AWSClient interface {
	AWSSecret(ctx context.Context) (*AWSSecret, error)
}

// Assert StandardClient conforms to Client interface
var _ = Client(&StandardClient{})

// StandardClient implements the vault client abstraction
type StandardClient struct {
	vaultHTTPClient *http.Client
	options         *Options
	vaultClient     *api.Client
	vaultRevoker    *api.Client // vault client used for delayed auth revokes
	lock            sync.Mutex
	logger          log.Logger
	authSecret      *StandardSecret
	mongoFetcher    secretHandler
	postgresFetcher secretHandler
	tlsKeyFetcher   secretHandler
	tlsCertFetcher  secretHandler
	awsFetcher      secretHandler
	secretHandlers  map[string]*secretHandler
}

// NewStandardClient returns a new vault client
func NewStandardClient(clientCtx context.Context, opts *Options) (*StandardClient, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	config := api.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, err
	}
	vaultClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	vaultRevokeClient, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	logger := log.FromContextWithPackageName(clientCtx, "go-common/vault")
	client := &StandardClient{
		vaultHTTPClient: config.HttpClient,
		options:         opts,
		vaultClient:     vaultClient,
		vaultRevoker:    vaultRevokeClient,
		logger:          logger,
		secretHandlers:  make(map[string]*secretHandler),
		mongoFetcher: secretHandler{
			client:          vaultClient,
			logger:          logger,
			successMsg:      LogMongoCredsReadSuccess,
			failureMsg:      LogMongoCredsReadFailure,
			renewSuccessMsg: LogMongoCredsRenewSuccess,
			renewFailureMsg: LogMongoCredsRenewFailure,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewUserPassSecret(NewStandardSecret(s, authSecret))
			},
		},
		postgresFetcher: secretHandler{
			client:          vaultClient,
			logger:          logger,
			successMsg:      LogPostgresCredsReadSuccess,
			failureMsg:      LogPostgresCredsReadFailure,
			renewSuccessMsg: LogPostgresCredsRenewSuccess,
			renewFailureMsg: LogPostgresCredsRenewFailure,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewUserPassSecret(NewStandardSecret(s, authSecret))
			},
		},
		tlsKeyFetcher: secretHandler{
			client:          vaultClient,
			logger:          logger,
			successMsg:      LogTLSCertKeyReadSuccess,
			failureMsg:      LogTLSCertKeyReadFailure,
			renewSuccessMsg: LogTLSCertKeyRenewSuccess,
			renewFailureMsg: LogTLSCertKeyRenewFailure,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewStandardSecret(s, authSecret), nil
			},
		},
		tlsCertFetcher: secretHandler{
			client:          vaultClient,
			logger:          logger,
			successMsg:      LogTLSCertReadSuccess,
			failureMsg:      LogTLSCertReadFailure,
			renewSuccessMsg: LogTLSCertRenewSuccess,
			renewFailureMsg: LogTLSCertRenewFailure,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewStandardSecret(s, authSecret), nil
			},
		},
		awsFetcher: secretHandler{
			client:          vaultClient,
			logger:          logger,
			successMsg:      LogAWSCredsReadSuccess,
			failureMsg:      LogAWSCredsReadFailure,
			renewSuccessMsg: LogAWSCredsRenewSuccess,
			renewFailureMsg: LogAWSCredsRenewFailure,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewAWSSecret(NewStandardSecret(s, authSecret))
			},
		},
	}
	return client, nil
}

// Options returns the current options
func (sc *StandardClient) Options() *Options {
	return sc.options
}

// VaultHTTPClient returns the HTTP client vault uses and is exposed for tests, should not be used otherwise
func (sc *StandardClient) VaultHTTPClient() *http.Client {
	return sc.vaultHTTPClient
}

// Authenticate gets a new auth token using approle authentication from vault.
// It will first try to renew the existing token (if any) and if that fails try to (re)authenticate.
// If that fails an error is returned.
func (sc *StandardClient) Authenticate(ctx context.Context) (*StandardSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	secret, err := sc.vaultAuthSecret(ctx, sc.authSecret)
	if err != nil {
		return nil, err
	}

	if sc.authSecret != nil {
		// cancel the old auth secret to require all dependent secrets to refresh as they might now be invalid
		sc.authSecret.Cancel()
	}

	// store the auth secret + set authentication token
	sc.authSecret = secret

	sc.logger.Info(
		LogAuthSuccessMsg,
		log.Int(LogKeyCurrentSecretID, int(sc.authSecret.ID())),
		log.Int(LogKeyLeaseDuration, sc.authSecret.Auth.LeaseDuration),
		log.String(LogKeyAccessor, sc.authSecret.Auth.Accessor),
		log.Time(LogKeyCreateTime, sc.authSecret.CreateTime()),
		log.Time(LogKeyRenewTime, sc.authSecret.RenewTime()),
		log.Time(LogKeyExpireTime, sc.authSecret.ExpireTime()),
	)

	return sc.authSecret, nil
}

func (sc *StandardClient) vaultAuthSecret(ctx context.Context, oldAuthSecret *StandardSecret) (*StandardSecret, error) {
	// if already authenticated, try to renew auth token lease for an hour
	if oldAuthSecret != nil && oldAuthSecret.Auth.Renewable {
		s, err := sc.vaultClient.Auth().Token().RenewSelf(3600)
		if err != nil {
			sc.logger.Error(LogErrFailedRenewAuth, log.Error(err))
		} else {
			secret := NewStandardSecret(s, nil)
			if secret.ExpireTime().After(time.Now().Add(MinLeaseExpireDuration)) {
				sc.updateFetcherAuthSecrets(secret)
				return secret, nil
			}
		}
	}

	// otherwise reauthenticate completely
	creds := &AppRoleCredentials{}
	if err := creds.ReadEnvironment(); err != nil {
		sc.logger.Error(LogErrFailedReadCredentials, log.Error(err))
		return nil, err
	}

	secret, err := sc.vaultClient.Logical().Write(sc.options.AppRoleAuthPath, creds.Map())
	if err != nil {
		sc.logger.Error(LogErrFailedAuth, log.Error(err))
		return nil, err
	}

	// if this client was authenticated, try to revoke previous token and all the old secrets
	// this is done after a 30 second delay to allow any existing work depending on existing secrets to finish
	// CRITICAL: delaying the revoke means that all existing secrets may still be renewed.
	// This causes invalid secret state when the secrets are revoked if all existing secrets are not removed on auth secret change.
	if sc.vaultClient.Token() != "" {
		sc.vaultRevoker.SetToken(oldAuthSecret.Auth.ClientToken)
		time.AfterFunc(30*time.Second, func() {
			if err := sc.vaultRevoker.Auth().Token().RevokeSelf(""); err != nil {
				sc.logger.Error("Failed to revoke token and tree of old secrets", log.Error(err))
			} else {
				sc.logger.Info("Revoked token and tree of old secrets", log.String(LogKeyAccessor, oldAuthSecret.Auth.Accessor))
			}
		})
	}

	// set the new secrets client token as the acive token
	sc.vaultClient.SetToken(secret.Auth.ClientToken)

	authSecret := NewStandardSecret(secret, nil)
	sc.updateFetcherAuthSecrets(authSecret)
	return authSecret, nil
}

func (sc *StandardClient) updateFetcherAuthSecrets(authSecret *StandardSecret) {
	// reset all secret handler authentications
	for _, handler := range sc.secretHandlers {
		handler.SetAuthSecret(authSecret)
	}
	sc.mongoFetcher.SetAuthSecret(authSecret)
	sc.postgresFetcher.SetAuthSecret(authSecret)
	sc.tlsKeyFetcher.SetAuthSecret(authSecret)
	sc.tlsCertFetcher.SetAuthSecret(authSecret)
	sc.awsFetcher.SetAuthSecret(authSecret)
}

// Read reads a secret from vault
func (sc *StandardClient) Read(ctx context.Context, path string) (*StandardSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	handler, exists := sc.secretHandlers[path]
	if !exists {
		handler = &secretHandler{
			client:          sc.vaultClient,
			logger:          sc.logger,
			authSecret:      sc.authSecret,
			successMsg:      "Successfully read secret at path " + path,
			failureMsg:      "Failed to read secret at path " + path,
			renewSuccessMsg: "Successfully renewed secret at path " + path,
			renewFailureMsg: "Failed to renewed secret at path " + path,
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewStandardSecret(s, authSecret), nil
			},
		}
		sc.secretHandlers[path] = handler
	}

	if err := handler.Renew(ctx); err != nil {
		if err := handler.Read(ctx, path); err != nil {
			return nil, err
		}
	}
	return handler.secret.(*StandardSecret), nil
}

// Write writes a secret to vault
func (sc *StandardClient) Write(ctx context.Context, path string, data map[string]interface{}) (*StandardSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	handler, exists := sc.secretHandlers[path]
	if !exists {
		handler = &secretHandler{
			client:          sc.vaultClient,
			logger:          sc.logger,
			authSecret:      sc.authSecret,
			successMsg:      "Successfully wrote secret at path " + path,
			failureMsg:      "Failed to write secret at path " + path,
			renewSuccessMsg: "Successfully renewed secret at path " + path, // unused but unsure if required, code is bit hackish
			renewFailureMsg: "Failed to renew secret at path " + path,      // unused but unsure if required, code is bit hackish
			transformVaultSecret: func(s *api.Secret, authSecret *StandardSecret) (Secret, error) {
				return NewStandardSecret(s, authSecret), nil
			},
		}
		sc.secretHandlers[path] = handler
	}

	if err := handler.Write(ctx, path, data); err != nil {
		return nil, err
	}
	return handler.secret.(*StandardSecret), nil
}

// MongoSecret returns a new secret containing an username and password which can be used to connect to mongo.
// It returns an error if the client has not been set up to work with mongo connections or a valid secret could not be fetched from vault
func (sc *StandardClient) MongoSecret(ctx context.Context) (*UserPassSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if !sc.options.EnableMongo {
		return nil, ErrMongoDisabled
	}

	if err := sc.mongoFetcher.Renew(ctx); err != nil {
		if err := sc.mongoFetcher.Read(ctx, sc.options.MongoCredsReadPath); err != nil {
			return nil, err
		}
	}
	return sc.mongoFetcher.secret.(*UserPassSecret), nil
}

// PostgresSecret returns a new secret containing an username and password which can be used to connect to postgres.
// It returns an error if the client has not been set up to work with postgres connections or a valid secret could not be fetched from vault
func (sc *StandardClient) PostgresSecret(ctx context.Context) (*UserPassSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if !sc.options.EnablePostgres {
		return nil, ErrPostgresDisabled
	}

	if err := sc.postgresFetcher.Renew(ctx); err != nil {
		if err := sc.postgresFetcher.Read(ctx, sc.options.PostgresCredsReadPath); err != nil {
			return nil, err
		}
	}
	return sc.postgresFetcher.secret.(*UserPassSecret), nil
}

// TLSCertSecret returns a new secret containing an TLS certificate.
// It returns an error if the client has not been set up to work with postgres connections or a valid secret could not be fetched from vault
func (sc *StandardClient) TLSCertSecret(ctx context.Context) (*TLSCertSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if !sc.options.EnableTLSCert {
		return nil, ErrTLSCertDisabled
	}

	if err := sc.tlsCertFetcher.Renew(ctx); err != nil {
		if err := sc.tlsCertFetcher.Read(ctx, sc.options.TLSCertReadPath); err != nil {
			return nil, err
		}
	}

	if err := sc.tlsKeyFetcher.Renew(ctx); err != nil {
		if err := sc.tlsKeyFetcher.Read(ctx, sc.options.TLSCertKeyReadPath); err != nil {
			return nil, err
		}
	}

	certSecret := sc.tlsCertFetcher.secret
	keySecret := sc.tlsKeyFetcher.secret

	return NewTLSCertSecret(certSecret.(*StandardSecret), keySecret.(*StandardSecret))
}

// AWSSecret returns a new secret containing a aws STS secret token which can be used to connect to aws.
// It returns an error if the client has not been set up to work with aws or a valid secret could not be fetched from vault
func (sc *StandardClient) AWSSecret(ctx context.Context) (*AWSSecret, error) {
	sc.lock.Lock()
	defer sc.lock.Unlock()

	if !sc.options.EnableAWS {
		return nil, ErrAWSDisabled
	}

	if err := sc.awsFetcher.Renew(ctx); err != nil {
		if err := sc.awsFetcher.Read(ctx, sc.options.AWSCredsReadPath); err != nil {
			return nil, err
		}
	}
	return sc.awsFetcher.secret.(*AWSSecret), nil
}

type secretHandler struct {
	client               *api.Client
	logger               log.Logger
	successMsg           string
	failureMsg           string
	renewSuccessMsg      string
	renewFailureMsg      string
	vaultSecret          *api.Secret
	secret               Secret
	authSecret           *StandardSecret
	transformVaultSecret func(*api.Secret, *StandardSecret) (Secret, error)
}

func (s *secretHandler) SetAuthSecret(authSecret *StandardSecret) {
	if s.authSecret != nil && s.authSecret.VaultID() != authSecret.VaultID() {
		// set the new auth secret and reset all secret state to prevent renewals of old secrets
		s.vaultSecret = nil
		s.secret = nil
	}
	s.authSecret = authSecret
}

func (s *secretHandler) Write(ctx context.Context, path string, data map[string]interface{}) error {
	return s.handleSecretOperation(ctx, func() (*api.Secret, error) { return s.client.Logical().Write(path, data) })
}

func (s *secretHandler) Read(ctx context.Context, path string) error {
	return s.handleSecretOperation(ctx, func() (*api.Secret, error) { return s.client.Logical().Read(path) })
}

func (s *secretHandler) handleSecretOperation(ctx context.Context, cb func() (*api.Secret, error)) error {
	if s.authSecret == nil {
		return ErrUnauthenticated
	}

	secret, err := cb()
	if err != nil {
		s.logger.Error(s.failureMsg, log.Error(err))
		return err
	}

	// In rare cases Vault returns nil, nil (no error + path not found) so we need to verify it here.
	// We consider not found to be an error which could happen due to e.g. invalid configuration in the service.
	if secret == nil {
		s.logger.Error(s.failureMsg, log.Error(ErrSecretNotFound))
		return ErrSecretNotFound
	}

	secretTransformed, err := s.transformVaultSecret(secret, s.authSecret)
	if err != nil {
		s.logger.Error(s.failureMsg, log.Error(err))
		return err
	}

	s.vaultSecret = secret
	s.secret = secretTransformed
	s.loggerWithSecrets(ctx, s.vaultSecret, s.secret).Info(s.successMsg)
	return nil
}

func (s *secretHandler) Renew(ctx context.Context) error {
	if s.authSecret == nil {
		return ErrUnauthenticated
	}

	if s.vaultSecret == nil {
		return ErrVaultSecretNotSet
	}

	if !s.vaultSecret.Renewable {
		return ErrVaultSecretNotRenewable
	}

	// By default try to renew for one hour, actual lease might be shorter
	secret, err := s.client.Sys().Renew(s.vaultSecret.LeaseID, 3600)
	if err != nil {
		s.logger.Error(s.renewFailureMsg, log.Error(ErrVaultSecretRenewFailed))
		return ErrVaultSecretRenewFailed
	}

	// In rare cases Vault returns nil, nil (no error + path not found) so we need to verify it here.
	// We consider not found to be an error which could happen due to e.g. invalid configuration in the service.
	if secret == nil {
		s.logger.Error(s.failureMsg, log.Error(ErrSecretNotFound))
		return ErrSecretNotFound
	}

	// update lease id and duration to the ones in the renewed secret
	s.vaultSecret.RequestID = secret.RequestID
	s.vaultSecret.LeaseID = secret.LeaseID
	s.vaultSecret.LeaseDuration = secret.LeaseDuration

	// update config state
	newSecret, err := s.transformVaultSecret(s.vaultSecret, s.authSecret)
	if err != nil {
		return err
	}

	// if the renewed secret has less than 10 minutes left, refetch the secrets.
	// This should only happen when we're close to the max lease time
	if newSecret.ExpireTime().Before(time.Now().Add(MinLeaseExpireDuration)) {
		s.loggerWithSecrets(ctx, s.vaultSecret, newSecret).Info(ErrVaultMaxLeaseExceeded.Error())
		return ErrVaultMaxLeaseExceeded
	}
	s.loggerWithSecrets(ctx, s.vaultSecret, newSecret).Info(s.renewSuccessMsg)
	s.secret = newSecret
	return nil
}

func (s *secretHandler) loggerWithSecrets(ctx context.Context, vaultSecret *api.Secret, secret Secret) log.Logger {
	// assume that both secret will either be available or nil
	if vaultSecret == nil || secret == nil {
		return s.logger
	}

	return s.logger.With(
		log.Int(LogKeyCurrentSecretID, int(secret.ID())),
		log.String(LogKeyLeaseID, vaultSecret.LeaseID),
		log.Int(LogKeyLeaseDuration, vaultSecret.LeaseDuration),
		log.Time(LogKeyCreateTime, secret.CreateTime()),
		log.Time(LogKeyRenewTime, secret.RenewTime()),
		log.Time(LogKeyExpireTime, secret.ExpireTime()),
	)
}
