package auth_test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/callstats-io/go-common/auth"
	jwt "github.com/dgrijalva/jwt-go"
)

func createJWT(secret []byte, c *auth.EndpointClaims, sm jwt.SigningMethod) string {
	claims := jwt.MapClaims{
		"appID":  c.AppID,
		"userID": c.UserID,
		"exp":    c.ExpiresAt,
	}
	token := jwt.NewWithClaims(sm, &claims)
	signed, _ := token.SignedString(secret)
	return signed
}

func randomSignSecret() []byte {
	rand.Seed(time.Now().UnixNano())
	s := make([]byte, 32, 32)
	rand.Read(s)
	return s
}

func randomClaims() *auth.EndpointClaims {
	rand.Seed(time.Now().UnixNano())
	return &auth.EndpointClaims{
		AppID:  rand.Int(),
		UserID: strconv.Itoa(rand.Int()),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Scope: []string{
			strconv.Itoa(rand.Int()),
		},
	}
}

func createJWTWithRawUserID(secret []byte, c *auth.EndpointClaims, sm jwt.SigningMethod) string {
	claims := jwt.MapClaims{
		"appID":     c.AppID,
		"userID":    c.UserID,
		"rawUserID": c.RawUserID,
		"exp":       c.ExpiresAt,
	}
	token := jwt.NewWithClaims(sm, &claims)
	signed, _ := token.SignedString(secret)
	return signed
}

func randomClaimsWithRawUserID() *auth.EndpointClaims {
	rand.Seed(time.Now().UnixNano())
	randID := rand.Int()
	return &auth.EndpointClaims{
		AppID:     rand.Int(),
		UserID:    "some%2Fuser%2Fid" + strconv.Itoa(randID),
		RawUserID: "some/user/id" + strconv.Itoa(randID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		Scope: []string{
			strconv.Itoa(rand.Int()),
		},
	}
}
