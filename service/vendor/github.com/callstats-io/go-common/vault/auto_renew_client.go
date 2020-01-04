package vault

import (
	"context"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

const (
	// set max backoff to a minute
	maxBackoff = int(time.Minute / time.Millisecond)
)

// Assert AutoRenewClient conforms to Client interface
var _ = Client(&AutoRenewClient{})

// AutoRenewClient internally caches any of the requested
type AutoRenewClient struct {
	ctx            context.Context
	client         Client
	lock           sync.Mutex
	secretHandlers map[string]*autoRenewSecretHandler

	// secret auto renew operations
	authSecretHandler     *autoRenewSecretHandler
	postgresSecretHandler *autoRenewSecretHandler
	mongoSecretHandler    *autoRenewSecretHandler
	tlsCertSecretHandler  *autoRenewSecretHandler
	awsSecretHandler      *autoRenewSecretHandler
}

// NewAutoRenewClient returns a new caching client that caches secrets until they require renewal or expire
func NewAutoRenewClient(clientCtx context.Context, c Client) (*AutoRenewClient, error) {
	arc := &AutoRenewClient{
		ctx:            clientCtx,
		client:         c,
		secretHandlers: make(map[string]*autoRenewSecretHandler),
	}

	callCtx, callCtxCancel := context.WithCancel(context.Background())
	defer callCtxCancel()

	{
		arc.authSecretHandler = &autoRenewSecretHandler{
			NewSecret: func(ctx context.Context) (Secret, error) {
				return c.Authenticate(ctx)
			},
		}
		if _, err := arc.authSecretHandler.Fetch(callCtx); err != nil {
			return nil, err
		}
		arc.authSecretHandler.RenewAutomatically(clientCtx)
	}

	if c.Options().EnableMongo {
		arc.mongoSecretHandler = &autoRenewSecretHandler{
			NewSecret: func(ctx context.Context) (Secret, error) {
				return c.MongoSecret(ctx)
			},
		}
		if _, err := arc.mongoSecretHandler.Fetch(callCtx); err != nil {
			return nil, err
		}
		arc.mongoSecretHandler.RenewAutomatically(clientCtx)
	}

	if c.Options().EnablePostgres {
		arc.postgresSecretHandler = &autoRenewSecretHandler{
			NewSecret: func(ctx context.Context) (Secret, error) {
				return c.PostgresSecret(ctx)
			},
		}
		if _, err := arc.postgresSecretHandler.Fetch(callCtx); err != nil {
			return nil, err
		}
		arc.postgresSecretHandler.RenewAutomatically(clientCtx)
	}

	if c.Options().EnableTLSCert {
		arc.tlsCertSecretHandler = &autoRenewSecretHandler{
			NewSecret: func(ctx context.Context) (Secret, error) {
				return c.TLSCertSecret(ctx)
			},
		}
		if _, err := arc.tlsCertSecretHandler.Fetch(callCtx); err != nil {
			return nil, err
		}
		arc.tlsCertSecretHandler.RenewAutomatically(clientCtx)
	}

	if c.Options().EnableAWS {
		arc.awsSecretHandler = &autoRenewSecretHandler{
			NewSecret: func(ctx context.Context) (Secret, error) {
				return c.AWSSecret(ctx)
			},
		}
		if _, err := arc.awsSecretHandler.Fetch(callCtx); err != nil {
			return nil, err
		}
		arc.awsSecretHandler.RenewAutomatically(clientCtx)
	}

	return arc, nil
}

// Options returns this clients underlying clients Options
func (arc *AutoRenewClient) Options() *Options {
	return arc.client.Options()
}

// VaultHTTPClient returns this clients underlying clients HTTP client vault uses and is exposed for tests, should not be used otherwise
func (arc *AutoRenewClient) VaultHTTPClient() *http.Client {
	return arc.client.VaultHTTPClient()
}

// Authenticate returns the most recent auth secret returned by the underlying client.
// The secrets are automatically refetched whenever they require renewal or expire.
func (arc *AutoRenewClient) Authenticate(ctx context.Context) (*StandardSecret, error) {
	secret, err := arc.authSecretHandler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*StandardSecret), nil
}

// Read returns the most recent secret for the path in the underlying client.
// The secret is automatically refetched whenever it requires renewal or expires.
func (arc *AutoRenewClient) Read(ctx context.Context, path string) (*StandardSecret, error) {
	handler, err := arc.handlerForPath(ctx, path)
	if err != nil {
		return nil, err
	}

	secret, err := handler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*StandardSecret), nil
}

// Write returns the vault write result. Write is never cached.
func (arc *AutoRenewClient) Write(ctx context.Context, path string, data map[string]interface{}) (*StandardSecret, error) {
	return arc.client.Write(ctx, path, data)
}

func (arc *AutoRenewClient) handlerForPath(ctx context.Context, path string) (*autoRenewSecretHandler, error) {
	arc.lock.Lock()
	defer arc.lock.Unlock()
	handler, exists := arc.secretHandlers[path]

	if !exists {
		handler = &autoRenewSecretHandler{
			NewSecret: func(callCtx context.Context) (Secret, error) {
				return arc.client.Read(callCtx, path)
			},
		}
		// fetch secret once using reads context
		if _, err := handler.Fetch(ctx); err != nil {
			return nil, err
		}
		handler.RenewAutomatically(arc.ctx)
		arc.secretHandlers[path] = handler
	}
	return handler, nil
}

// MongoSecret returns the most recent MongoSecret returned by the underlying client.
// The secrets are automatically refetched whenever they require renewal or expire.
func (arc *AutoRenewClient) MongoSecret(ctx context.Context) (*UserPassSecret, error) {
	if !arc.client.Options().EnableMongo {
		return nil, ErrMongoDisabled
	}
	secret, err := arc.mongoSecretHandler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*UserPassSecret), nil
}

// PostgresSecret returns the most recent PostgresSecret returned by the underlying client.
// The secrets are automatically refetched whenever they require renewal or expire.
func (arc *AutoRenewClient) PostgresSecret(ctx context.Context) (*UserPassSecret, error) {
	if !arc.client.Options().EnablePostgres {
		return nil, ErrPostgresDisabled
	}
	secret, err := arc.postgresSecretHandler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*UserPassSecret), nil
}

// TLSCertSecret returns the most recent TLSCertSecret returned by the underlying client.
// The secrets are automatically refetched whenever they require renewal or expire.
func (arc *AutoRenewClient) TLSCertSecret(ctx context.Context) (*TLSCertSecret, error) {
	if !arc.client.Options().EnableTLSCert {
		return nil, ErrTLSCertDisabled
	}
	secret, err := arc.tlsCertSecretHandler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*TLSCertSecret), nil
}

// AWSSecret returns the most recent AWSSecret returned by the underlying client.
// The secrets are automatically refetched whenever they require renewal or expire.
func (arc *AutoRenewClient) AWSSecret(ctx context.Context) (*AWSSecret, error) {
	if !arc.client.Options().EnableAWS {
		return nil, ErrAWSDisabled
	}
	secret, err := arc.awsSecretHandler.Secret()
	if err != nil {
		return nil, err
	}
	return secret.(*AWSSecret), nil
}

// autoRenewSecretHandler is the interface all secret operations that should be handled by the AutoRenewClient should implement
type autoRenewSecretHandler struct {
	lock      sync.Mutex
	secret    Secret
	fetchErr  error
	NewSecret func(context.Context) (Secret, error)
}

func (h *autoRenewSecretHandler) Fetch(ctx context.Context) (Secret, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	secret, err := h.NewSecret(ctx)

	// reflect is required to check that the value of an interface is nil
	if err != nil && !nilSecret(h.secret) && h.secret.Valid() {
		// do not assign to old secret so active clients can still use it as it is still valid
		return secret, err
	}

	h.secret, h.fetchErr = secret, err
	return h.secret, h.fetchErr
}

func (h *autoRenewSecretHandler) Secret() (Secret, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.secret, h.fetchErr
}

func (h *autoRenewSecretHandler) RenewAutomatically(clientCtx context.Context) {
	go func() {
		backoff := 0
		for {
			h.lock.Lock()
			secret := h.secret
			h.lock.Unlock()

			if !nilSecret(secret) {
				// wait for a change
				select {
				case <-clientCtx.Done():
					return // shutdown
				case <-secret.RenewContext().Done():
					// waits for renew and fetches new secret when it happens on lines below
				case <-secret.ExpireContext().Done():
					//  waits for expire and fetches new secret when it happens on lines below
				}
			}

			// try to fetch
			// if there was an error, retry with increasing backoff
			if err := h.fetchSecret(); err != nil {
				backoff += 50
				if backoff > maxBackoff {
					// always retry at least once roughly every minute
					backoff = maxBackoff
				}
				backoff += rand.Intn(20)
				time.Sleep(time.Duration(backoff) * time.Millisecond)
			} else {
				backoff = 0
			}
		}
	}()
}

func (h *autoRenewSecretHandler) fetchSecret() error {
	callCtx, callCtxCancel := context.WithCancel(context.Background())
	_, err := h.Fetch(callCtx)
	callCtxCancel()
	return err
}
