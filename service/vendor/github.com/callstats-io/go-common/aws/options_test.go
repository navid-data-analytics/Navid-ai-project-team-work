package aws_test

import (
	"net/http"
	"os"

	"github.com/callstats-io/go-common/aws"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options", func() {
	Context("from environment", func() {
		It("should read the AWS region", func() {
			prev := os.Getenv(aws.EnvAWSRegion)
			defer os.Setenv(aws.EnvAWSRegion, prev)

			os.Setenv(aws.EnvAWSRegion, "test-region-1")
			opts, err := aws.OptionsFromEnv()
			Expect(err).To(BeNil())
			Expect(opts).ToNot(BeNil())
			Expect(opts.AWSRegion).To(Equal("test-region-1"))
		})
		It("should use default AWS region", func() {
			prev := os.Getenv(aws.EnvAWSRegion)
			defer os.Setenv(aws.EnvAWSRegion, prev)

			os.Unsetenv(aws.EnvAWSRegion)
			opts, err := aws.OptionsFromEnv()
			Expect(err).To(BeNil())
			Expect(opts).ToNot(BeNil())
			Expect(opts.AWSRegion).To(Equal(aws.DefaultAWSRegion))
		})

		It("should have debug logging off by default", func() {
			prev := os.Getenv(aws.EnvLogLevel)
			defer os.Setenv(aws.EnvLogLevel, prev)

			os.Unsetenv(aws.EnvLogLevel)
			opts, err := aws.OptionsFromEnv()
			Expect(err).To(BeNil())
			Expect(opts).ToNot(BeNil())
			Expect(opts.DebugLogging).To(Equal(false))
		})
		It("should set the debug logging on", func() {
			prev := os.Getenv(aws.EnvLogLevel)
			defer os.Setenv(aws.EnvLogLevel, prev)

			os.Setenv(aws.EnvLogLevel, "DEBUG")
			opts, err := aws.OptionsFromEnv()
			Expect(err).To(BeNil())
			Expect(opts).ToNot(BeNil())
			Expect(opts.DebugLogging).To(Equal(true))
		})
	})

	Context("using builder", func() {
		It("should set MaxRetries", func() {
			opts := &aws.Options{MaxRetries: 5}
			opts = opts.WithMaxRetries(2)
			Expect(opts.MaxRetries).To(Equal(2))
		})
		It("should set HTTP client", func() {
			opts := &aws.Options{HTTPClient: http.DefaultClient}
			httpClient := &http.Client{Timeout: 1}
			opts = opts.WithHTTPClient(httpClient)
			Expect(opts.HTTPClient).To(Equal(httpClient))
			Expect(opts.HTTPClient).ToNot(Equal(http.DefaultClient))
		})
		It("should set DebugLogging on", func() {
			opts := &aws.Options{DebugLogging: false}
			opts = opts.WithDebugLogging()
			Expect(opts.DebugLogging).To(Equal(true))
		})
	})
})
