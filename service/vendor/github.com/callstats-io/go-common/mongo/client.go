package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	mgo "github.com/globalsign/mgo"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/vault"
)

// Errors
var (
	ErrInvalidSessionSecret = errors.New("Failed to get a valid vault mongo secret")
	ErrClosed               = errors.New("StandardClient has been closed")
)

// Client is the interface mongo clients implement
type Client interface {
	SessionClone(ctx context.Context) (*Session, error)
	Session(ctx context.Context) (*Session, error)
	Close()
}

// Assert StandardClient conforms to Client interface
var _ = Client(&StandardClient{})

// StandardClient contains the mongo abstraction used by other services
type StandardClient struct {
	options       *Options
	vaultClient   vault.MongoClient
	lock          sync.Mutex
	closed        bool
	cachedSecret  *vault.UserPassSecret
	cachedSession *mgo.Session
}

// NewStandardClient returns a new mongo StandardClient which internally uses the vault
func NewStandardClient(vc vault.MongoClient, opts *Options) (*StandardClient, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	client := &StandardClient{
		options:     opts,
		vaultClient: vc,
	}

	return client, nil
}

// SessionClone returns a new session bound to the ctx lifetime. The caller is expected to close the session when they're done with it.
// Internally uses Session.Clone() from mgo.
func (mc *StandardClient) SessionClone(ctx context.Context) (*Session, error) {
	return mc.session(ctx, func(s *mgo.Session) *mgo.Session { return s.Clone() })
}

// Session returns a new session bound to the ctx lifetime. The caller is expected to close the session when they're done with it.
// Internally uses Session.Copy() from mgo.
func (mc *StandardClient) Session(ctx context.Context) (*Session, error) {
	return mc.session(ctx, func(s *mgo.Session) *mgo.Session { return s.Copy() })
}

func (mc *StandardClient) session(ctx context.Context, sessionFunc func(*mgo.Session) *mgo.Session) (*Session, error) {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	if mc.closed {
		return nil, ErrClosed
	}

	// get a new or cached secret
	secret, err := mc.vaultClient.MongoSecret(ctx)
	if err != nil {
		return nil, err
	}

	// if the secret equals to the cached one, return a new copy of the current session (retains mongo auth info)
	if mc.cachedSecret != nil && mc.cachedSession != nil && secret.ID() == mc.cachedSecret.ID() {
		return asAliasedSession(sessionFunc(mc.cachedSession)), nil
	}

	// otherwise close the cached session as the secret changed it might now be invalid
	if mc.cachedSession != nil {
		// delay closing to allow any currently active operations to finish
		oldSession := mc.cachedSession
		time.AfterFunc(30*time.Second, func() {
			oldSession.Close()
		})
		mc.cachedSession = nil
	}

	// get a new "root" session with the new secret
	session, err := mgo.DialWithTimeout(fmt.Sprintf(mc.options.ConnectionTemplate, secret.Credentials.User, secret.Credentials.Password), mc.options.DialTimeout)
	if err != nil {
		return nil, err
	}

	logger := log.FromContextWithPackageName(ctx, "go-common/mongo").With(log.Int(LogKeyCurrentSecretID, int(secret.ID())))
	if mc.cachedSecret != nil {
		logger = logger.With(log.Int(LogKeyPrevSecretID, int(mc.cachedSecret.ID())))
	}
	logger.Info(LogMsgNewSession)

	mc.cachedSecret = secret
	mc.cachedSession = session

	//return a new copy of the current session (retains auth info) to prevent the cached session from being accidentally closed
	return asAliasedSession(sessionFunc(session)), nil
}

// Status checks if the client has an alive session in a context
func (mc *StandardClient) Status(ctx context.Context) (M, error) {
	session, err := mc.Session(ctx)
	if err != nil {
		return nil, err
	}
	return session.Status()
}

// Close closes this StandardClient.
func (mc *StandardClient) Close() {
	mc.lock.Lock()
	defer mc.lock.Unlock()

	if mc.closed {
		return
	}
	mc.closed = true

	if mc.cachedSession != nil {
		mc.cachedSession.Close()
	}
}

var _ = Client(&StaticClient{})

// StaticClient contains the mongo abstraction used by other services via static mongo connection URL
type StaticClient struct {
	lock          sync.Mutex
	cachedSession *mgo.Session
	closed        bool
}

// NewStaticClient returns a new mongo StaticClient which uses a static connection url to mongodb
func NewStaticClient(ctx context.Context, o *Options) (*StaticClient, error) {
	ses, err := mgo.DialWithTimeout(o.ConnectionTemplate, o.DialTimeout)
	if err != nil {
		return nil, err
	}

	client := &StaticClient{
		cachedSession: ses,
	}

	// close automatically on context cancel, allows simple shutdowns of everything based on single context
	go func() {
		<-ctx.Done()
		client.Close()
	}()

	return client, nil
}

// SessionClone returns a new session bound to the ctx lifetime. The caller is expected to close the session when they're done with it.
// Internally uses Session.Clone() from mgo.
func (s *StaticClient) SessionClone(ctx context.Context) (*Session, error) {
	return s.session(ctx, func(s *mgo.Session) *mgo.Session { return s.Clone() })
}

// Session returns a new session bound to the ctx lifetime. The caller is expected to close the session when they're done with it.
// Internally uses Session.Copy() from mgo.
func (s *StaticClient) Session(ctx context.Context) (*Session, error) {
	return s.session(ctx, func(s *mgo.Session) *mgo.Session { return s.Copy() })
}

func (s *StaticClient) session(_ context.Context, sessionFunc func(*mgo.Session) *mgo.Session) (*Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.closed {
		return nil, ErrClosed
	}

	return asAliasedSession(sessionFunc(s.cachedSession)), nil
}

// Close closes this client immediately
func (s *StaticClient) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.closed = true
	s.cachedSession.Close()
}
