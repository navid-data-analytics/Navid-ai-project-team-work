package app

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/NYTimes/gziphandler"
	proxyproto "github.com/armon/go-proxyproto"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/vault"
)

// App contains all the app bootstrapping settings
type App struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	logger      log.Logger
	vaultClient vault.Client
	httpsPort   int
	httpPort    int
}

// NewApp returns a new App with the provided context
// The App will shutdown when the context is canceled
func NewApp(ctx context.Context) *App {
	a := &App{
		logger: log.FromContextWithPackageName(ctx, "go-common/app"),
	}
	a.ctx, a.ctxCancel = context.WithCancel(ctx)

	return a
}

// Context returns this apps context
func (a *App) Context() context.Context {
	return a.ctx
}

// VaultClient returns this apps vault client. In case one has not been set up for this client,
// this will default to creating an vault.AutoRenewClient based on ENV variables.
func (a *App) VaultClient() vault.Client {
	if isNil(a.vaultClient) {
		// setup vault client
		vaultOpts, err := vault.OptionsFromEnv()
		if err != nil {
			a.logger.Panic(LogErrFailedToReadVaultOptions, log.Error(err))
		}
		vaultClient, err := vault.NewStandardClient(a.ctx, vaultOpts)
		if err != nil {
			a.logger.Panic(LogErrFailedToSetupVaultClient, log.Error(err))
		}
		vaultAutoRenewClient, err := vault.NewAutoRenewClient(a.ctx, vaultClient)
		if err != nil {
			a.logger.Panic(LogErrFailedToSetupVaultAutoRenewClient, log.Error(err))
		}
		a.vaultClient = vaultAutoRenewClient
	}
	return a.vaultClient
}

// WithVaultClient sets this apps vault client
func (a *App) WithVaultClient(vc vault.Client) *App {
	a.vaultClient = vc
	return a
}

// WithHTTPSPort sets the port to be used for HTTPS server. This overrides the derived setting from EnvHTTPSPort.
func (a *App) WithHTTPSPort(port int) *App {
	a.httpsPort = port
	return a
}

// WithHTTPPort sets the port to be used for HTTP server. This overrides the derived setting from EnvHTTPPort.
func (a *App) WithHTTPPort(port int) *App {
	a.httpPort = port
	return a
}

// ServeHTTPS starts a https server with routing done by the provided mux
func (a *App) ServeHTTPS(mux http.Handler) *App {
	vaultClient := a.VaultClient()
	if isNil(vaultClient) {
		a.logger.Panic(LogErrVaultClientRequired)
	}

	tlsSecret, err := vaultClient.TLSCertSecret(a.ctx)
	if err != nil {
		a.logger.Panic(LogErrTLSCertFetchHTTPS, log.Error(err))
		return a
	}

	if a.httpsPort == 0 {
		port, err := strconv.Atoi(os.Getenv(EnvHTTPSPort))
		if err != nil {
			a.logger.Panic(LogErrFailedToParseEnvHTTPSPort, log.Error(err))
		}
		a.httpsPort = port
	}

	// configured based on the suggestions from https://blog.gopheracademy.com/advent-2016/exposing-go-on-the-internet/
	tlsConfig := &tls.Config{
		// Causes servers to use Go's default ciphersuite preferences,
		// which are tuned to avoid attacks. Does nothing on clients.
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519, // Go 1.8 only
		},
		Certificates: []tls.Certificate{
			*tlsSecret.Certificate,
		},
		NextProtos: []string{
			"h2", //enable HTTP2
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

			// Best disabled, as they don't provide Forward Secrecy,
			// but might be necessary for some clients
			// tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			// tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	srv := http.Server{
		Addr:           ":" + strconv.Itoa(a.httpsPort),
		Handler:        gziphandler.GzipHandler(mux), // enable transparent gzip support
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second, // Go 1.8 only
		MaxHeaderBytes: 16384,
		TLSConfig:      tlsConfig,
		TLSNextProto:   nil, //enable http2
	}

	// start a new net.Listener with tls and haproxy protocol support
	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		a.logger.Panic(LogErrFailedToListenPort, log.Int(LogKeyPort, a.httpsPort), log.Error(err))
	}
	proxyLis := &proxyproto.Listener{Listener: lis}
	tlsLis := tls.NewListener(proxyLis, tlsConfig)

	// start serving with the configured listener
	go func() {
		a.logger.Info(LogMsgServeHTTPS, log.Int(LogKeyPort, a.httpsPort))
		if err := srv.Serve(tlsLis); err != nil {
			a.logger.Info(LogMsgShutdownHTTPS, log.Error(err))
		}
		// shutdown app
		a.Shutdown()
	}()

	// shutdown the server when the context is canceled
	go func() {
		<-a.ctx.Done()
		srv.Close()
	}()

	return a
}

// ServeHTTP starts a http server listening at the given address. This is mainly user to allow an internal status endpoint.
func (a *App) ServeHTTP(mux http.Handler) *App {
	if a.httpPort == 0 {
		port, err := strconv.Atoi(os.Getenv(EnvHTTPPort))
		if err != nil {
			a.logger.Panic(LogErrFailedToParseEnvHTTPPort, log.Error(err))
			return a
		}
		a.httpPort = port
	}

	srv := http.Server{
		Addr:           ":" + strconv.Itoa(a.httpPort),
		Handler:        gziphandler.GzipHandler(mux), // enable transparent gzip support
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second, // Go 1.8 only
		MaxHeaderBytes: 16384,
	}

	// start serving with the configured listener
	go func() {
		a.logger.Info(LogMsgServeHTTP, log.Int(LogKeyPort, a.httpPort))
		if err := srv.ListenAndServe(); err != nil {
			a.logger.Info(LogMsgShutdownHTTP, log.Error(err))
		}
		// shutdown app
		a.Shutdown()
	}()

	// shutdown the server when the context is canceled
	go func() {
		<-a.ctx.Done()
		srv.Close()
	}()

	return a
}

// Shutdown shuts down the app by stopping all running services
func (a *App) Shutdown() {
	a.ctxCancel()
}

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}
