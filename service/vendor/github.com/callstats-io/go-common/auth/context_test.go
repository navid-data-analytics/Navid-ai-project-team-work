package auth_test

import (
	"context"

	"github.com/callstats-io/go-common/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndpointClaimsFromContext", func() {
	It("should return nil if claims has not been stored in context", func() {
		Expect(auth.EndpointClaimsFromContext(context.Background())).To(BeNil())
	})
	It("should return the claims if stored in context with WithEndpointClaims", func() {
		claims := randomClaims()
		ctx := auth.WithEndpointClaims(context.Background(), claims)
		Expect(auth.EndpointClaimsFromContext(ctx)).To(Equal(claims))
	})
})
