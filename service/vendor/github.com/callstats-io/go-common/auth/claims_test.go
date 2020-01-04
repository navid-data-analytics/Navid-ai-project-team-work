package auth_test

import (
	"context"
	"time"

	"github.com/callstats-io/go-common/auth"
	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// NOTE SH: these tests expect validCommonHeader to have been used in setting up all the events that are validated
// This is to get to the correct localID easily without using reflect
var _ = Describe("Claims", func() {
	var testSignMethod = auth.SigningMethodHS256
	var testSignSecret []byte
	var claims, parsedClaims *auth.EndpointClaims
	BeforeEach(func() {
		testSignSecret = randomSignSecret()
		claims = randomClaims()
		parsedClaims = &auth.EndpointClaims{}
	})
	Context("ParseAndVerify", func() {
		It("should not return error for a valid JWT token", func() {
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(BeNil())
			Expect(parsedClaims.AppID).To(Equal(claims.AppID))
			Expect(parsedClaims.UserID).To(Equal(claims.UserID))
			Expect(parsedClaims.ExpiresAt).To(Equal(claims.ExpiresAt))
		})
		It("should fail if the token is empty", func() {
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, "")).To(MatchError(auth.ErrEmptyAuthToken))
		})
		It("should fail if the parsing fails", func() {
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, "abcdef")).To(MatchError(auth.ErrInvalidAuthToken))
		})
		It("should fail if the signature isn't valid", func() {
			token := createJWT([]byte("abcdef"), claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(MatchError(auth.ErrInvalidAuthToken))
		})
		It("should fail if the token has expired", func() {
			claims.ExpiresAt = time.Now().Add(-time.Hour).Unix()
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(MatchError(auth.ErrAuthExpiredToken))
		})
		It("should fail if the token has unsupported signing method", func() {
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS512)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(MatchError(auth.ErrInvalidSigningMethod))
		})
	})

	Context("ValidateAppID", func() {
		BeforeEach(func() {
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(BeNil())
		})
		It("should succeed if the appID matches claims appID", func() {
			Expect(parsedClaims.ValidateAppID(claims.AppID)).To(BeNil())
		})
		It("should fail if the appID doesn't match claims appID", func() {
			Expect(parsedClaims.ValidateAppID(claims.AppID + 1)).To(MatchError(auth.ErrInvalidAuthAppID))
		})
	})

	Context("ValidateUserID", func() {
		BeforeEach(func() {
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(BeNil())
		})
		It("should succeed if the userID matches claims userID", func() {
			Expect(claims.ValidateUserID(claims.UserID)).To(BeNil())
		})
		It("should fail if the userID doesn't match claims userID", func() {
			Expect(claims.ValidateUserID(claims.UserID + "abc")).To(MatchError(auth.ErrInvalidAuthUserID))
		})
	})

	Context("ValidateUserIDWithRawUserIDInToken", func() {
		claimsWithRaw := randomClaimsWithRawUserID()
		BeforeEach(func() {
			token := createJWTWithRawUserID(testSignSecret, claimsWithRaw, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(BeNil())
		})
		It("should succeed if the userID matches claims userID", func() {
			Expect(claimsWithRaw.ValidateRawUserID(claimsWithRaw.RawUserID)).To(BeNil())
		})
		It("should fail if the userID doesn't match claims userID", func() {
			Expect(claimsWithRaw.ValidateRawUserID(claimsWithRaw.RawUserID + "abc")).To(MatchError(auth.ErrInvalidAuthUserID))
		})
	})

	Context("ValidateAuthScopes", func() {
		BeforeEach(func() {
			token := createJWT(testSignSecret, claims, jwt.SigningMethodHS256)
			Expect(parsedClaims.ParseAndVerify(context.Background(), testSignSecret, testSignMethod, token)).To(BeNil())
		})
		It("should succeed if the required scopes matches claims scopes", func() {
			Expect(claims.ValidateAuthScopes(claims.Scope)).To(BeNil())
		})
		It("should fail if the claim is missing required scopes", func() {
			Expect(claims.ValidateAuthScopes(append(claims.Scope, "non-existing"))).To(MatchError(auth.ErrInvalidScope))
		})
	})

	Context("ValidateOrigin", func() {
		for _, testCase := range []struct {
			Description string
			TestOrigin  string
			OriginURLs  []string
			ExpError    error
		}{
			// copy of test cases for node collector:
			// https://github.com/callstats-io/cs-js-common/blob/4c1ed465790d94fd5ba8e446323c086b18557315/test/index.js#L76
			{
				Description: "matches \"https://foo.bar.com\" to \"https://foo.bar.com\"",
				TestOrigin:  "https://foo.bar.com",
				OriginURLs:  []string{"https://foo.bar.com"},
				ExpError:    nil,
			},

			{
				Description: "matches \"https://foo.bar.com\" to \"https://foo.bar.com/\"",
				TestOrigin:  "https://foo.bar.com",
				OriginURLs:  []string{"https://foo.bar.com/"},
				ExpError:    nil,
			},

			{
				Description: "matches \"https://foo.bar.com\" to \"https://foo.bar.com///\"",
				TestOrigin:  "https://foo.bar.com",
				OriginURLs:  []string{"https://foo.bar.com///"},
				ExpError:    nil,
			},

			{
				Description: "does not match \"https://foo2.bar.com\" to \"https://foo.bar.com\"",
				TestOrigin:  "https://foo2.bar.com",
				OriginURLs:  []string{"https://foo.bar.com"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "matches \"https://foo.bar.com\" to [\"https://bar.com\", \"https://foo.bar.com\"]",
				TestOrigin:  "https://foo.bar.com",
				OriginURLs:  []string{"https://bar.com", "https://foo.bar.com"},
				ExpError:    nil,
			},

			{
				Description: "does not match \"https://foo2.bar.com\" to [\"https://foo.bar.com\"]",
				TestOrigin:  "https://foo2.bar.com",
				OriginURLs:  []string{"https://foo.bar.com"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "does not match \"https://foo.bar.company\" to \"https://foo.bar.com\"",
				TestOrigin:  "https://foo2.bar.company",
				OriginURLs:  []string{"https://foo.bar.com"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "does not match \"https://foo.bar.com\" to \"https://foo.bar.company\"",
				TestOrigin:  "https://foo2.bar.com",
				OriginURLs:  []string{"https://foo.bar.company"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "does not match \"https://foo.bar.com\" to \"https://foo.bar.company\"",
				TestOrigin:  "https://foo2.bar.com",
				OriginURLs:  []string{"https://foo.bar.company"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "matches \"https://foo2.bar.com\" to \"https://*.bar.com\"",
				TestOrigin:  "https://foo2.bar.com",
				OriginURLs:  []string{"https://*.bar.com"},
				ExpError:    nil,
			},

			{
				Description: "matches \"https://foo2.2.bar.com\" to \"https://*.bar.com\"",
				TestOrigin:  "https://foo2.2.bar.com",
				OriginURLs:  []string{"https://*.bar.com"},
				ExpError:    nil,
			},

			{
				Description: "does not match \"https://foo2.2-.bar.com\" to \"https://*.bar.com\"",
				TestOrigin:  "https://foo2.2-.bar.com",
				OriginURLs:  []string{"https://*.bar.com"},
				ExpError:    auth.ErrInvalidOriginURL,
			},

			{
				Description: "does not match \"https://foo.bar.company\" to \"https://*.bar.com\"",
				TestOrigin:  "https://foo.bar.company",
				OriginURLs:  []string{"https://*.bar.com"},
				ExpError:    auth.ErrInvalidOriginURL,
			},
		} {
			testCase := testCase

			It(testCase.Description, func() {
				claims := &auth.EndpointClaims{OriginURLs: testCase.OriginURLs}
				if testCase.ExpError == nil {
					Expect(claims.ValidateOrigin(testCase.TestOrigin)).To(BeNil())
				} else {
					Expect(claims.ValidateOrigin(testCase.TestOrigin)).To(MatchError(testCase.ExpError))
				}
			})

		}
	})
})
