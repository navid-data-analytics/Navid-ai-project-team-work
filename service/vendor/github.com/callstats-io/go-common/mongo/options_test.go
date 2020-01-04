package mongo_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/callstats-io/go-common/mongo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OptionsFromEnv", func() {
	It("should get MONGO_CONN_TMPL from env", func() {
		prev := os.Getenv(mongo.EnvConnectionTemplate)
		randstr := strconv.Itoa(rand.Int())
		os.Setenv(mongo.EnvConnectionTemplate, randstr)
		opts, err := mongo.OptionsFromEnv()
		Expect(err).To(BeNil())
		Expect(opts.ConnectionTemplate).To(Equal(randstr))
		os.Setenv(mongo.EnvConnectionTemplate, prev)
	})
	It("should get MONGO_DIAL_TIMEOUT from env", func() {
		prev := os.Getenv(mongo.EnvDialTimeout)
		os.Setenv(mongo.EnvDialTimeout, "abc")
		_, timeErr := time.ParseDuration("abc")
		opts, err := mongo.OptionsFromEnv()
		Expect(err).To(MatchError(fmt.Errorf("Invalid value for %s, error: %s", mongo.EnvDialTimeout, timeErr)))

		os.Setenv(mongo.EnvDialTimeout, "")
		opts, err = mongo.OptionsFromEnv()
		Expect(err).To(BeNil())
		Expect(opts.DialTimeout).To(Equal(5 * time.Second))

		os.Setenv(mongo.EnvDialTimeout, "1s")
		opts, err = mongo.OptionsFromEnv()
		Expect(err).To(BeNil())
		Expect(opts.DialTimeout).To(Equal(time.Second))
		os.Setenv(mongo.EnvDialTimeout, prev)
	})
	It("should return error if the options are invalid", func() {
		prev := os.Getenv(mongo.EnvConnectionTemplate)
		os.Unsetenv(mongo.EnvConnectionTemplate)
		_, err := mongo.OptionsFromEnv()
		Expect(err).To(MatchError(mongo.ErrEmptyConnectionTemplate))
		os.Setenv(mongo.EnvConnectionTemplate, prev)
	})
})

var _ = Describe("Options", func() {
	Context("Validate", func() {
		It("should return nil if the options are valid", func() {
			opts := &mongo.Options{
				ConnectionTemplate: "mongodb://%s:%s@mongo:27017,mongo2:27018/test",
			}
			Expect(opts.Validate()).To(BeNil())
		})
	})
})
