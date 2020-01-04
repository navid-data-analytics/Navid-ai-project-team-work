package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/vault"
)

// Client is the interface all models expect for operations
type Client interface {
	DB(ctx context.Context) (*sql.DB, error)
}

// VaultSQLClient implements a client that uses vault to fetch the sql secrets
type VaultSQLClient struct {
	vaultReader     vault.SecretReader
	ctx             context.Context
	ctxCancel       context.CancelFunc
	logger          log.Logger
	lock            sync.Mutex
	secret          *vault.StandardSecret
	conn            *sql.DB
	sqlConnTemplate string
	sqlDriver       string
	vaultSecretPath string
	pingInterval    time.Duration
}

// assert implements Client
var _ = Client(&VaultSQLClient{})

// NewVaultSQLClient returns a new SQL database client which creates and closes connections based on vault credentials
// If no options are provided, the configurable fields default to:
// - sql driver: "postgres" (option: OptionSQLDriver)
// - sql connection url: ENV["POSTGRES_CONN_TMPL"] (option: OptionSQLConnectionTemplate)
// - credentials vault path: ENV["ENV"]/postgresql/ENV["SERVICE_NAME"]/creds/ENV["SERVICE_NAME"] (option: OptionCredentialsVaultPath)
func NewVaultSQLClient(ctx context.Context, vaultReader vault.SecretReader, options ...Option) (*VaultSQLClient, error) {
	wrappedCtx, wrappedCtxCancel := context.WithCancel(ctx)

	// default to postgres, assume client has loaded a library exposing the driver, all overridable by options
	client := &VaultSQLClient{
		ctx:             wrappedCtx,
		ctxCancel:       wrappedCtxCancel,
		logger:          log.FromContextWithPackageName(wrappedCtx, "go-common/database/sql"),
		vaultReader:     vaultReader,
		sqlDriver:       "postgres",
		sqlConnTemplate: os.Getenv("POSTGRES_CONN_TMPL"),
		vaultSecretPath: os.Getenv("ENV") + "/postgresql/" + os.Getenv("SERVICE_NAME") + "/creds/" + os.Getenv("SERVICE_NAME"),
		pingInterval:    30 * time.Second,
	}

	// configure client
	for _, opt := range options {
		opt(client)
	}

	return client, nil
}

// DB returns a *sql.DB based on vault credentials.
// It will block until either a valid connection could be established or the context expires.
// If the context expires, the context error is returned.
// If the client is closed, ErrClosed is returned
func (c *VaultSQLClient) DB(ctx context.Context) (*sql.DB, error) {
	select {
	case <-c.ctx.Done():
		return nil, ErrClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return c.connection()
	}
}

// Close closes this client and any open connections.
// It is equivalent to cancelling the context the client was created with.
func (c *VaultSQLClient) Close() {
	c.ctxCancel()
}

func (c *VaultSQLClient) connection() (*sql.DB, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.secret == nil || c.conn == nil {
		return c.newConnection()
	}

	select {
	case <-c.secret.ExpireContext().Done():
		return c.newConnection()
	case <-c.secret.RenewContext().Done():
		return c.newConnection()
	default:
		// do nothing, pass through to returning existing connection
		// do not return here since select is not deterministic (potentially possible to return conn after ctx cancel)
	}
	return c.conn, nil
}

func (c *VaultSQLClient) newConnection() (*sql.DB, error) {
	secret, err := c.vaultReader.Read(c.ctx, c.vaultSecretPath)
	if err != nil {
		return nil, err
	}

	// create new connection
	conn, err := sql.Open(c.sqlDriver, fmt.Sprintf(c.sqlConnTemplate, secret.Data["username"], secret.Data["password"]))
	if err != nil {
		return nil, err
	}

	// verify connection works
	if err := conn.Ping(); err != nil {
		return nil, err
	}

	if c.conn != nil {
		// delay closing to allow any currently active operations to finish
		oldConn := c.conn
		time.AfterFunc(30*time.Second, func() {
			oldConn.Close()
		})
		c.conn = nil
	}

	c.conn = conn
	c.secret = secret

	return conn, nil
}
