package postgres_test

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/callstats-io/go-common/postgres"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type optionsTestCase struct {
	EnvVariable string
	EnvValue    string
	Error       error
}

var _ = Describe("OptionsFromEnv", func() {
	Context("Success", func() {
		It("should get POSTGRES_CONN_TMPL from env", func() {
			prev := os.Getenv(postgres.EnvConnectionTemplate)
			defer os.Setenv(postgres.EnvConnectionTemplate, prev)
			randstr := strconv.Itoa(rand.Int())
			os.Setenv(postgres.EnvConnectionTemplate, randstr)
			opts, err := postgres.OptionsFromEnv()
			Expect(err).To(BeNil())
			Expect(opts.ConnectionTemplate).To(Equal(randstr))
		})
	})
	Context("Failure", func() {
		testCases := []optionsTestCase{
			optionsTestCase{
				EnvVariable: postgres.EnvConnectionTemplate,
				EnvValue:    "",
				Error:       postgres.ErrEmptyConnectionTemplate,
			},
		}
		for idx := range testCases {
			testCase := &testCases[idx]
			It("should return an error if "+testCase.EnvVariable+" value is \""+testCase.EnvValue+"\"", func() {
				prev := os.Getenv(testCase.EnvVariable)
				defer os.Setenv(testCase.EnvVariable, prev)
				os.Setenv(testCase.EnvVariable, testCase.EnvValue)

				_, err := postgres.OptionsFromEnv()
				Expect(err).To(MatchError(testCase.Error))
			})
		}
	})
})

var _ = Describe("Options", func() {
	Describe("Validate", func() {
		Context("Success", func() {
			It("should return nil if the options are valid", func() {
				opts := &postgres.Options{
					ConnectionTemplate: "postgres://%s:%s@postgres:5432/go_common_test",
				}
				Expect(opts.Validate()).To(BeNil())
			})
		})
		Context("Failure", func() {
			It("should return an error if the options are invalid", func() {
				opts := &postgres.Options{
					ConnectionTemplate: "",
				}
				Expect(opts.Validate()).To(MatchError(postgres.ErrEmptyConnectionTemplate))
			})
		})
	})
})
