package log_test

import (
	"context"

	"github.com/callstats-io/go-common/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/uber-go/zap"
)

var _ = Describe("FromEnv", func() {
	It("should return a new logger", func() {
		firstLogger := log.NewStdoutLogger("PANIC")
		logger := log.FromEnv()
		Expect(logger).ToNot(BeNil())
		Expect(logger).ToNot(Equal(firstLogger))
	})
})

var _ = Describe("NewLogger", func() {
	It("should return a new logger", func() {
		logger := log.NewLogger("PANIC")
		Expect(logger).ToNot(BeNil())
		_, loggerInstance := logger.(log.Logger)
		Expect(loggerInstance).To(BeTrue())
	})
})

var _ = Describe("Contextual logger", func() {
	var ctx context.Context
	var logger log.Logger

	BeforeEach(func() {
		ctx = context.Background()
		logger = log.NewLogger("DEBUG")
	})

	It("should set the logger to context", func() {
		ctx = log.WithLogger(ctx, logger)
		ctxLogger := log.FromContext(ctx)
		// expect logger not to equal the default logger
		Expect(ctxLogger).ToNot(Equal(log.FromContext(context.Background())))
	})

	It("should get a logger from context", func() {
		ctx = log.WithContext(ctx, logger)

		ctxLogger := log.FromContext(ctx)

		Expect(ctxLogger).ToNot(BeNil())
		_, loggerInstance := ctxLogger.(log.Logger)
		Expect(loggerInstance).To(BeTrue())
		Expect(ctxLogger).To(Equal(logger))
	})
})

var _ = Describe("NewStdoutLogger", func() {
	It("should return a new logger", func() {
		testLevels := map[string]zap.Level{
			"":             zap.InfoLevel,
			log.DebugLevel: zap.DebugLevel,
			log.InfoLevel:  zap.InfoLevel,
			log.WarnLevel:  zap.WarnLevel,
			log.ErrorLevel: zap.ErrorLevel,
			log.FatalLevel: zap.FatalLevel,
			log.PanicLevel: zap.PanicLevel,
		}
		for levelStr, level := range testLevels {
			logger := log.NewStdoutLogger(levelStr)
			Expect(logger).ToNot(BeNil())
			Expect(logger.Check(level, levelStr).OK()).To(BeTrue())
		}
	})
})
