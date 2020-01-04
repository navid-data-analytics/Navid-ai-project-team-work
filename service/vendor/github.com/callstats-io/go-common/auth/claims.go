package auth

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/callstats-io/go-common/log"
	jwt "github.com/dgrijalva/jwt-go"
)

var validSigningMethods = []jwt.SigningMethod{
	jwt.SigningMethodHS256,
	jwt.SigningMethodES256,
}

// EndpointClaims contains the information user EndpointClaims to be
type EndpointClaims struct {
	jwt.StandardClaims
	CollectSDP         bool     `json:"collectSDP"`
	SubmissionInterval string   `json:"submissionInterval"`
	UserID             string   `json:"userID"`
	RawUserID          string   `json:"rawUserID"`
	EndpointType       string   `json:"endpointType"`
	AppID              int      `json:"appID"`
	OriginURLs         []string `json:"origin"`
	Scope              []string `json:"scope"`
}

// ParseAndVerify parses the auth token to auth claims
func (c *EndpointClaims) ParseAndVerify(ctx context.Context, secret []byte, method jwt.SigningMethod, token string) error {
	if token == "" {
		return ErrEmptyAuthToken
	}

	_, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		if t.Method != method {
			return nil, ErrInvalidSigningMethod
		}
		return secret, nil
	})

	// if the token has expired, return an explicit error
	if c.ExpiresAt != 0 && c.ExpiresAt < time.Now().UTC().Unix() {
		log.FromContext(ctx).Error(LogMsgErrExpiredAuthToken, log.Object(LogKeyAuthToken, c))
		return ErrAuthExpiredToken
	}

	if err != nil {
		log.FromContext(ctx).Error(LogMsgErrInvalidAuthToken, log.Error(err))
		// The JWT library wraps all errors to its own internal error implementation (why?) so we need to
		// assert and check for inner error match here to figure out if this was signing method failure or not
		if vErr, ok := err.(*jwt.ValidationError); ok && vErr.Inner == ErrInvalidSigningMethod {
			return ErrInvalidSigningMethod
		}
		return ErrInvalidAuthToken
	}

	return nil
}

// ValidateAppID validates the auth claims against actually submitted appID
func (c *EndpointClaims) ValidateAppID(appID int) error {
	if appID != c.AppID {
		return ErrInvalidAuthAppID
	}

	return nil
}

// ValidateUserID validates the auth claims against actually submitted userID
func (c *EndpointClaims) ValidateUserID(userID string) error {
	if userID != c.UserID {
		return ErrInvalidAuthUserID
	}
	return nil
}

// ValidateRawUserID validates the auth claims against non-encoded submitted userID
func (c *EndpointClaims) ValidateRawUserID(rawUserID string) error {
	if rawUserID != c.RawUserID {
		return ErrInvalidAuthUserID
	}
	return nil
}

// ValidateAuthScopes validates that all the required scopes are found from the auth token
func (c *EndpointClaims) ValidateAuthScopes(scopes []string) error {
	if !hasAll(c.Scope, scopes) {
		return ErrInvalidScope
	}

	return nil
}

// ValidateOrigin validates that origin matches one of the allowed originURLs
func (c *EndpointClaims) ValidateOrigin(origin string) error {
	if len(c.OriginURLs) == 0 {
		return nil
	}

	for _, url := range c.OriginURLs {
		// drop all suffix slahes. Node collector has this so we preserve compatibility for multiple slashes.
		for strings.HasSuffix(url, "/") {
			url = strings.TrimSuffix(url, "/")
		}

		if strings.Contains(url, "*.") {
			// handle wildcard URLs, follows exactly the same process as node collector:
			// https://github.com/callstats-io/cs-js-common/blob/4c1ed465790d94fd5ba8e446323c086b18557315/index.js#L36
			replaced := url
			replaced = strings.Replace(replaced, "*.", "(([a-z0-9][a-z0-9_-]*[a-z0-9]|[a-z0-9]).)+", -1)
			replaced = strings.Replace(replaced, "/", "\\/", -1)
			replaced = strings.Replace(replaced, ".", "\\.", -1)
			matchRegex, err := regexp.Compile("^" + replaced + "$")
			if err != nil {
				return ErrParseOriginURL
			}

			if matchRegex.MatchString(origin) {
				return nil
			}
		} else if origin == url {
			return nil
		}
	}

	return ErrInvalidOriginURL
}

// check that the first set contains all items in the second set
func hasAll(first, second []string) bool {
	for _, ss := range second {
		if !has(first, ss) {
			return false
		}
	}

	return true
}

// check that the set contains the value
func has(haystack []string, needle string) bool {
	for _, v := range haystack {
		if needle == v {
			return true
		}
	}

	return false
}
