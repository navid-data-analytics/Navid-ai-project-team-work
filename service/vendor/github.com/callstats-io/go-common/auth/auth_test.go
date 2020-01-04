package auth_test

import (
	"fmt"

	"github.com/callstats-io/go-common/auth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// NOTE SH: these tests expect validCommonHeader to have been used in setting up all the events that are validated
// This is to get to the correct localID easily without using reflect
var _ = Describe("SigningMethodFromString", func() {
	testCases := map[string]auth.SigningMethod{
		"HS256": auth.SigningMethodHS256,
		"HS384": auth.SigningMethodHS384,
		"HS512": auth.SigningMethodHS512,
		"ES256": auth.SigningMethodES256,
		"ES384": auth.SigningMethodES384,
		"ES512": auth.SigningMethodES512,
	}
	for s, m := range testCases {
		It(fmt.Sprintf("should return correct method for %s", s), func() {
			Expect(auth.SigningMethodFromString(s)).To(Equal(m))
		})
	}
	It("should panic if the method is unknown", func() {
		Expect(func() { auth.SigningMethodFromString("HS") }).To(Panic())
	})
})
