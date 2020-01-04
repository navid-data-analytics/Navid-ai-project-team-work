package aws_test

import (
	"github.com/callstats-io/go-common/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Context("NewClient", func() {
		It("should initialize with given opts", func() {
			opts := &aws.Options{}
			testClient = aws.NewStandardClient(testVaultClient, opts)
			Expect(testClient).ToNot(BeNil())
			Expect(testClient.Vault).To(Equal(testVaultClient))
			Expect(testClient.Options).To(Equal(opts))
		})
	})

	Context("Session", func() {
		It("should return a new valid session", func() {
			testOptions, err := aws.OptionsFromEnv()
			Expect(err).To(BeNil())
			testClient = aws.NewStandardClient(testVaultClient, testOptions)
			Expect(testClient.Session(testCtx)).ToNot(BeNil())
		})
	})
})
