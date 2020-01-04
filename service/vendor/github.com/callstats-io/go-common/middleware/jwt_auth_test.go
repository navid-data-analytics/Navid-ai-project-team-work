package middleware_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/callstats-io/go-common/auth"
	"github.com/callstats-io/go-common/middleware"
	"github.com/callstats-io/go-common/response"
	"github.com/callstats-io/go-common/testutil"
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndpointJWTAuth", func() {
	var testValidSignMethod = jwt.SigningMethodHS256
	expectAuthErrorResponse := func(w *testutil.FakeResponseWriter, err error) {
		Expect(w.StatusCode).To(Equal(http.StatusUnauthorized))
		Expect(w.BodyData).To(Equal(response.MarshalError(err)))
	}

	dummyRespData := []byte("DUMMYRESP")
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(dummyRespData)
	}

	var testSignSecret []byte
	var claims *auth.EndpointClaims
	var writer *testutil.FakeResponseWriter
	BeforeEach(func() {
		testSignSecret = randomSignSecret()
		claims = randomClaims()
		writer = &testutil.FakeResponseWriter{}
	})
	Context("Success", func() {
		It("should call the handler with a correctly signed valid JWT", func() {
			token := createJWT(testSignSecret, claims, testValidSignMethod)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(dummyRespData))
		})
	})
	Context("Failure", func() {
		It("should fail if authorization header is not present", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			expectAuthErrorResponse(writer, middleware.ErrMissingAuthHeader)
		})
		It("should fail if authorization header does not start with Bearer", func() {
			token := createJWT(testSignSecret, claims, testValidSignMethod)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderAuthorization, token)
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			expectAuthErrorResponse(writer, middleware.ErrMissingAuthBearerPrefix)
		})
		It("should fail if authorization header does not have a JWT token", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderAuthorization, "Bearer ")
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			expectAuthErrorResponse(writer, auth.ErrEmptyAuthToken)
		})
		It("should fail if the JWT could not be parsed", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderAuthorization, "Bearer abc")
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			expectAuthErrorResponse(writer, auth.ErrInvalidAuthToken)
		})
		It("should fail if the JWT token has expired", func() {
			claims.ExpiresAt = time.Now().Add(-time.Hour).UTC().Unix()
			token := createJWT(testSignSecret, claims, testValidSignMethod)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
			middleware.EndpointJWTAuth(testSignSecret, testValidSignMethod)(dummyHandler)(writer, req)
			expectAuthErrorResponse(writer, auth.ErrAuthExpiredToken)
		})
	})
})

var _ = Describe("EndpointRequireJWTScopes", func() {
	dummyRespData := []byte("DUMMYRESP")
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(dummyRespData)
	}

	var requiredScopes []string
	var claims *auth.EndpointClaims
	var writer *testutil.FakeResponseWriter
	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
		requiredScopes = []string{
			strconv.Itoa(rand.Int()),
		}
		claims = randomClaims()
		writer = &testutil.FakeResponseWriter{}
	})
	Context("Success", func() {
		It("should call the handler when endpoint claims has the correct scope", func() {
			claims.Scope = append(claims.Scope, requiredScopes...)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopes(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(dummyRespData))
		})
	})
	Context("Failure", func() {
		It("should fail if request context doesn't have the claims", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())

			middleware.EndpointRequireJWTScopes(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusInternalServerError))
			Expect(writer.BodyData).To(Equal(response.RespInternalServerError))
		})
		It("should fail if endpoint claims has doesn't have correct scope", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopes(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(writer.BodyData).To(Equal(middleware.RespInvalidScope))
		})
		It("should fail if claims scope is empty", func() {
			claims.Scope = []string{}
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopes(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(writer.BodyData).To(Equal(middleware.RespInvalidScope))
		})
		It("should fail claims scope is nil", func() {
			claims.Scope = nil
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopes(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusUnauthorized))
			Expect(writer.BodyData).To(Equal(middleware.RespInvalidScope))
		})
	})
})

var _ = Describe("EndpointRequireJWTScopesForHTTP1", func() {
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DUMMYRESP"))
	}
	var requiredScopes []string
	var claims *auth.EndpointClaims
	var writer *testutil.FakeResponseWriter
	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
		requiredScopes = []string{
			strconv.Itoa(rand.Int()),
		}
		claims = randomClaims()
		writer = &testutil.FakeResponseWriter{}
	})

	Context("Success", func() {
		It("should call handler if HTTP major is 2", func() {
			req, err := http.NewRequest("POST", "", nil)
			Expect(err).To(BeNil())
			writer := &testutil.FakeResponseWriter{}

			called := false
			req.ProtoMajor = 2
			middleware.RequireHTTP2()(func(w http.ResponseWriter, r *http.Request) {
				called = true
			})(writer, req)
			Expect(called).To(BeTrue())
		})
		It("should call handler if HTTP major is 1 and the claims have the requiredScopes", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopesForHTTP1(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
			Expect(writer.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		})
	})

	Context("Failure", func() {
		It("should fail if request context doesn't have the claims", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())

			middleware.EndpointRequireJWTScopesForHTTP1(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusInternalServerError))
			Expect(writer.BodyData).To(Equal(response.RespInternalServerError))
		})
		It("should fail if endpoint claims has doesn't have correct scope and HTTP major is 1", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopesForHTTP1(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
			Expect(writer.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		})
		It("should fail if claims scope is empty and HTTP major is 1", func() {
			claims.Scope = []string{}
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopesForHTTP1(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
			Expect(writer.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		})
		It("should fail claims scope is nil and HTTP major is 1", func() {
			claims.Scope = nil
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointRequireJWTScopesForHTTP1(requiredScopes)(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusHTTPVersionNotSupported))
			Expect(writer.BodyData).To(Equal(response.RespInvalidHTTPVersion))
		})
	})
})

var _ = Describe("EndpointJWTCORSMiddleware", func() {
	dummyRespData := []byte("DUMMYRESP")
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(dummyRespData)
	}
	testOriginURL := "http://testorigin"

	var claims *auth.EndpointClaims
	var writer *testutil.FakeResponseWriter
	BeforeEach(func() {
		rand.Seed(time.Now().UnixNano())
		claims = randomClaims()
		writer = &testutil.FakeResponseWriter{}
	})
	Context("Success", func() {
		It("should call the handler when endpoint claims has the correct originURLs", func() {
			claims.OriginURLs = append(claims.OriginURLs, testOriginURL)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderOrigin, testOriginURL)
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointJWTCORS()(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(dummyRespData))
			Expect(writer.HeaderData[middleware.HeaderAccessControlAllowOrigin]).To(Equal([]string{testOriginURL}))
		})
		It("should succeed if request context doesn't have claims", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderOrigin, testOriginURL)

			middleware.EndpointJWTCORS()(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(dummyRespData))
			Expect(writer.HeaderData[middleware.HeaderAccessControlAllowOrigin]).To(Equal([]string{testOriginURL}))
		})
		It("should succeed if claims doesn't have originURLs limitations", func() {
			claims.OriginURLs = []string{}
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderOrigin, testOriginURL)
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointJWTCORS()(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusOK))
			Expect(writer.BodyData).To(Equal(dummyRespData))
			Expect(writer.HeaderData[middleware.HeaderAccessControlAllowOrigin]).To(Equal([]string{testOriginURL}))
		})
		It("should succeed if origin is not present", func() {
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			claimsReq, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			claimsReq = claimsReq.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			for _, r := range []*http.Request{req, claimsReq} {
				middleware.EndpointJWTCORS()(dummyHandler)(writer, r)
				Expect(writer.StatusCode).To(Equal(http.StatusOK))
				Expect(writer.BodyData).To(Equal(dummyRespData))
				// Expect the header to always have been added to the response, even if technically the request is invalid for CORS
				Expect(writer.HeaderData[middleware.HeaderAccessControlAllowOrigin]).To(Equal([]string{""}))
			}
		})
	})
	Context("Failure", func() {
		It("should fail if origin is not in claims OriginURLs", func() {
			claims.OriginURLs = append(claims.OriginURLs, testOriginURL)
			req, err := http.NewRequest(http.MethodPost, "", nil)
			Expect(err).To(BeNil())
			req.Header.Set(middleware.HeaderOrigin, testOriginURL+"abc")
			req = req.WithContext(auth.WithEndpointClaims(req.Context(), claims))

			middleware.EndpointJWTCORS()(dummyHandler)(writer, req)
			Expect(writer.StatusCode).To(Equal(http.StatusPreconditionFailed))
			Expect(writer.BodyData).To(Equal(middleware.RespInvalidOriginURL))
			Expect(writer.HeaderData[middleware.HeaderAccessControlAllowOrigin]).To(Equal([]string{""}))
		})
	})
})
