package app

// environment variables
const (
	EnvHTTPSPort = "HTTPS_PORT"
	EnvHTTPPort  = "HTTP_PORT"
)

// Log message keys
const (
	LogKeyPort = "port"
)

// Log messages
const (
	LogErrFailedToReadVaultOptions          = "Failed to read vault options"
	LogErrFailedToSetupVaultClient          = "Failed to setup vault base client"
	LogErrFailedToSetupVaultAutoRenewClient = "Failed to setup vault auto renew client"
	LogErrVaultClientRequired               = "Vault client cannot be nil with ServeHTTPS"
	LogErrTLSCertFetchHTTPS                 = "Failed to fetch TLS secret with HTTPS"
	LogErrFailedToParseEnvHTTPSPort         = "Failed to parse env HTTPS_PORT to a number"
	LogErrFailedToParseEnvHTTPPort          = "Failed to parse env HTTP_PORT to a number"
	LogErrFailedToListenPort                = "Failed to listen"
	LogMsgServeHTTPS                        = "Starting HTTPS server"
	LogMsgShutdownHTTPS                     = "Shutting down HTTPS server"
	LogMsgServeHTTP                         = "Starting HTTP server"
	LogMsgShutdownHTTP                      = "Shutting down HTTP server"
)
