package aws_test

import (
	"encoding/json"

	"github.com/callstats-io/go-common/aws"
	"github.com/callstats-io/go-common/log"
	"github.com/callstats-io/go-common/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LogPublisher", func() {
	var fakeTopicArn = "aws:sns:fake:topic:arn"

	Context("initialization", func() {
		It("should create an SNSPublisher", func() {
			p := aws.NewLogPublisher(fakeTopicArn)
			Expect(p).ToNot(BeNil())
			Expect(p.TopicARN).To(Equal(fakeTopicArn))
		})
	})

	Context("Publish", func() {
		var (
			logBuffer    *testutil.LogBuffer
			logPublisher *aws.LogPublisher
			payload      string
		)

		BeforeEach(func() {
			logBuffer = testutil.NewLogBuffer()
			testCtx = log.WithContext(testCtx, logBuffer.Logger())

			logPublisher = aws.NewLogPublisher(fakeTopicArn)
			Expect(logPublisher).ToNot(BeNil())

			msg := map[string]interface{}{
				"organizationName": "Test org",
				"inviteAcceptUrl":  "https://dashboard.localhost/invite/ABBA-123DEADBEEF-654321-FDSA",
				"email":            "test-invitee@callstats.io",
			}
			payloadBytes, err := json.Marshal(msg)
			Expect(err).To(BeNil())
			payload = string(payloadBytes)
		})

		It("should publish a message successfully", func() {
			err := logPublisher.Publish(testCtx, payload)
			Expect(err).To(BeNil())
			loggedLines := logBuffer.Lines()
			Expect(len(loggedLines)).To(Equal(1))
			Expect(loggedLines[0]).To(ContainSubstring(aws.LogMessagePublished))
			Expect(loggedLines[0]).To(ContainSubstring("https://dashboard.localhost/invite/ABBA-123DEADBEEF-654321-FDSA"))
		})

		It("should return an error when the payload is empty", func() {
			err := logPublisher.Publish(testCtx, "")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal(aws.LogEmptyPayload))
			loggedLines := logBuffer.Lines()
			Expect(len(loggedLines)).To(Equal(1))
			Expect(loggedLines[0]).To(ContainSubstring(aws.LogEmptyPayload))
			Expect(loggedLines[0]).To(ContainSubstring(aws.LogKeyOriginalPayload))
		})
	})
})
