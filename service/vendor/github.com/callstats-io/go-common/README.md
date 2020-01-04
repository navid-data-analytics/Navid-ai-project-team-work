# go-common
Golang common lib setups used by multiple internal services. This README contains highlevel purpose documentation for each of the packages. For more detailed documentation see each of the packages.

## Components

### [Auth](/auth)
JWT Based auth token validations. Includes middlewares for parsing and verifying auth tokens from Authorization-header, validation of appID to appID url match etc.

### [Logger](/log)
An abstraction of [zap Logger](https://github.com/uber-go/zap) with configurations from ENV variables

Additionally supports getting a specified or a default logger from context.

### [Metrics](/metrics)
Preconfigured [go-metrics](https://github.com/rcrowley/go-metrics) timers for tracking e.g. time it takes for a http request, a mongo query etc.

### [Middleware](/middleware)
Tools to create middleware chains out of golang http handlers. The middlewares are executed in order with the expectation of the middleware either writing a response or calling the next handler.
Each middleware should be done as a function that takes any required configuration as parameter, returns a function that takes the next http.HandlerFunc as argument and returns a http.HandlerFunc compatible function. 

In code:
```golang
func BeforeAfterMiddleware(before, after string) Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            logger := log.FromContextWithPackageName(r.Context(), "go-common/middleware/beforeafter")
            logger.Info(before)
            next(w,r)
            logger.Info(after)
        }
    }
}
```

Example real middleware:
```golang
func RequestMetrics() Middleware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h(w, r)
			metrics.RequestTimer().Update(time.Now().Sub(start) / time.Millisecond)
		}
	}
}
```

### [Response](/response)
Utility functions for common HTTP responses. Naming should roughly follow the HTTP semantic name, e.g. a response method of HTTP STATUS 400 Bad Request should be called response.BadRequest.
Currently every response also adds content type of application/json and the content length to every response.

Also contains a simple response structure we've used previously

```golang
type Response struct {
	Status string `json:"status"` // HTTP status code
	Msg    string `json:"msg"`    // Description
}
```

Finally includes a MarshalError method that creates an marshaled response of
```golang
var err error
Response{
    Status: "error",
    Msg: err.Error(),
}
```


### [Server](/server)
Utility methods for running a HTTP(S) servers

### [Tools](/tools)
Cross-project tools to eg. bootstrap Vault and MongoDB.

### [Vaultbootstrap](/vaultbootstrap)
Utilities for common dev bootstrapping of vault. E.g. setting policies, mounting rules etc. for local development.

### [Vault](/vault)
Abstractions on top of the official vault client for credentials with automatic refresh on expire/renew.

#### StandardClient
Client which contains the basic operations for fetching secrets from vault. The methods should always fetch a new secret from vault.

#### AutoRenewClient
Client which wraps another vault client and fetches a new secret automatically when a secret expires. 
It returns the cached secret while the secret has not expired or a newer secret has not been fetched.

#### [Postgres](/postgres)
Postgres contains a client which can be used to get a valid postgres connection based on the vault provided credentials. 
Generally you'll want to use the AutoRenewClient from vault as the vault client as the postgres standard client verifies 
the secret on every database connection request and recycles connections when the secrets change.
This process will be very slow if the client has to fetch a secret every time from vault.
This is a separate package so that all code which uses the vault client does not need to import postgres packages.

#### [Kafka](/kafka)
Kafka producer abstraction library

#### [Zookeeper](/zk) 
Zookeeper client abstraction library

#### [Mongo](/mongo)
Mongo contains a client which can be used to get a valid mongo session based on the vault provided credentials. 
Generally you'll want to use the AutoRenewClient from vault as the vault client as the mongo standard client verifies 
the secret on every mongo session request and recycles connections when the secrets change. 
This process will be very slow if the client has to fetch a secret every time from vault.
This is a separate package so that all code which uses the vault client does not need to import mongo packages.
