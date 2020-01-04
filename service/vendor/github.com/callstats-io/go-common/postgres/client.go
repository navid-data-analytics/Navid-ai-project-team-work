package postgres

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-pg/pg"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/vault"
)

// Client is the interface postgres clients implement
type Client interface {
	DB(ctx context.Context) (*DB, error)
	Close()
}

// ParseURL parses connection string passed in as argument as Options
func ParseURL(raw string) (*pg.Options, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}

	if parsed.Host == "" {
		return nil, ErrInvalidAddr
	}
	if parsed.Path == "" {
		return nil, ErrInvalidDB
	}

	// ignore forward slash
	db := parsed.Path
	if strings.HasPrefix(db, "/") {
		db = db[1:]
	}

	opts := &pg.Options{
		Addr:     parsed.Host,
		Database: db,
	}

	if parsed.User != nil {
		password, _ := parsed.User.Password()
		opts.User = parsed.User.Username()
		opts.Password = password
	}

	return opts, nil
}

var _ = Client(&StandardClient{})

// StandardClient contains the postgres abstraction used by other services
type StandardClient struct {
	ctx          context.Context
	options      *Options
	vaultClient  vault.PostgresClient
	lock         sync.Mutex
	closed       bool
	cachedSecret *vault.UserPassSecret
	cachedDB     *DB
}

// NewStandardClient returns a new postgres StandardClient which internally uses the vault
func NewStandardClient(ctx context.Context, vc vault.PostgresClient, opts *Options) (*StandardClient, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	client := &StandardClient{
		ctx:         ctx,
		options:     opts,
		vaultClient: vc,
	}

	return client, nil
}

// DB returns a new DB bound to the ctx lifetime. The caller is expected to close the DB when they're done with it.
func (pgc *StandardClient) DB(ctx context.Context) (*DB, error) {
	pgc.lock.Lock()
	defer pgc.lock.Unlock()

	if pgc.closed {
		return nil, ErrClosed
	}

	// get a new or cached secret
	secret, err := pgc.vaultClient.PostgresSecret(ctx)
	if err != nil {
		return nil, err
	}

	// if the secret equals to the cached one, return a new copy of the current DB (retains postgres auth info)
	if pgc.cachedSecret != nil && pgc.cachedDB != nil && secret.ID() == pgc.cachedSecret.ID() {
		return pgc.cachedDB, nil
	}

	// get a new "root" DB with the new secret
	opts, err := ParseURL(fmt.Sprintf(pgc.options.ConnectionTemplate, secret.Credentials.User, secret.Credentials.Password))
	if err != nil {
		return nil, err
	}

	// if there is an old connection, close it with delay to allow any existing calls to complete before closing
	if pgc.cachedDB != nil {
		pgc.closeExistingConnection()
		pgc.cachedDB = nil
	}

	logger := log.FromContextWithPackageName(ctx, "go-common/postgres").With(log.Int(LogKeyCurrentSecretID, int(secret.ID())))
	if pgc.cachedSecret != nil {
		logger = logger.With(log.Int(LogKeyPrevSecretID, int(pgc.cachedSecret.ID())))
	}
	logger.Info(LogMsgNewSession)

	pgc.cachedSecret = secret
	pgc.cachedDB = asAliasedDB(pg.Connect(opts))

	//return a new copy of the current DB (retains auth info) to prevent the cached DB from being accidentally closed
	return pgc.cachedDB, nil
}

// Close closes this StandardClient.
func (pgc *StandardClient) Close() {
	pgc.lock.Lock()
	defer pgc.lock.Unlock()

	if pgc.closed {
		return
	}
	pgc.closed = true

	if pgc.cachedDB != nil {
		pgc.cachedDB.Close()
	}
}

func (pgc *StandardClient) closeExistingConnection() {
	oldDB := pgc.cachedDB
	closeCtx, closeCtxCancel := context.WithTimeout(pgc.ctx, 30*time.Second)
	go func() {
		// delay close to allow for active operations to complete
		// close happens when either the client context is done or the timer expires
		<-closeCtx.Done()
		oldDB.Close()

		// cancel to make linter happy, this is a no-op
		closeCtxCancel()
	}()
}

var _ = Client(&StaticClient{})

// StaticClient contains the postgres abstraction used by other services via static postgres connection URL
type StaticClient struct {
	lock     sync.Mutex
	cachedDB *DB
	closed   bool
}

// NewStaticClient returns a new postgres StaticClient which uses a static connection url to postgres
func NewStaticClient(ctx context.Context, o *Options) (*StaticClient, error) {
	opts, err := ParseURL(o.ConnectionTemplate)
	if err != nil {
		return nil, err
	}

	client := &StaticClient{
		cachedDB: asAliasedDB(pg.Connect(opts)),
	}

	// close automatically on context cancel, allows simple shutdowns of everything based on single context
	go func() {
		<-ctx.Done()
		client.Close()
	}()

	return client, nil
}

// DB returns the static db connection created via the connection url on create
func (s *StaticClient) DB(_ context.Context) (*DB, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.closed {
		return nil, ErrClosed
	}

	return s.cachedDB, nil
}

// Close closes this client immediately
func (s *StaticClient) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.closed = true
	s.cachedDB.Close()
}
