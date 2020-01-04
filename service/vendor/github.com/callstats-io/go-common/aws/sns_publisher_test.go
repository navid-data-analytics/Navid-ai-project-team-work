package aws_test

import (
	"encoding/json"
	"fmt"

	"github.com/callstats-io/go-common/aws"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/h2non/gock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	successResponse = `<?xml version="1.0" encoding="UTF-8"?>
<PublishResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/">
  <PublishResult>
    <MessageId>fake-123-message-456-id</MessageId>
  </PublishResult>
  <ResponseMetadata>
    <RequestId>fake-123-request-456-id</RequestId>
  </ResponseMetadata>
</PublishResponse>`

	errorResponse = `<?xml version="1.0" encoding="UTF-8"?>
<ErrorResponse xmlns="http://sns.amazonaws.com/doc/2010-03-31/">
  <Error>
    <Code>123</Code>
    <Message>AWS error</Message>
    <Resource>%s</Resource>
  </Error>
  <RequestId>fake-123-request-456-id</RequestId>
</ErrorResponse>`
)

func mockSNSPublishRequest(awsClient *aws.StandardClient, topicArn, msg string, responseStatus int, cb func()) {
	secret, err := awsClient.Credentials(testCtx)
	Expect(err).To(BeNil())
	creds, err := secret.Get()
	Expect(err).To(BeNil())
	Expect(creds.SessionToken).ToNot(BeNil())

	// TODO: figure out the exact parameters so they can be matched
	//qs := strings.Join([]string{
	//	fmt.Sprintf("TopicArn=%s", topicArn),
	//	fmt.Sprintf("Message=%s", url.QueryEscape(msg)),
	//	fmt.Sprintf("X-Amz-Security-Token=%s", creds.SessionToken),
	//	"Action=Publish",
	//}, "&")

	resp := successResponse
	if responseStatus != 200 {
		resp = fmt.Sprintf(errorResponse, topicArn)
	}

	defer gock.DisableNetworking()
	defer gock.Off()

	gock.InterceptClient(awsClient.Options.HTTPClient)

	gock.New(fmt.Sprintf("https://sns.%s.amazonaws.com", awsClient.Options.AWSRegion)).
		Post("/"). // TODO: add + "?"+qs
		Reply(responseStatus).
		XML(resp)

	cb()
}

var _ = Describe("SNSPublisher", func() {
	var fakeTopicArn = "aws:sns:fake:topic:arn"

	BeforeEach(func() {
		var err error
		testOptions, err = aws.OptionsFromEnv()
		Expect(err).To(BeNil())
		testOptions = testOptions.WithMaxRetries(0)
	})

	Context("initialization", func() {
		It("should create an SNSPublisher", func() {
			c := aws.NewStandardClient(testVaultClient, testOptions)
			p := aws.NewSNSPublisher(c, fakeTopicArn)
			Expect(p).ToNot(BeNil())
			Expect(p.TopicARN).To(Equal(fakeTopicArn))
		})
	})

	Context("Publish", func() {
		var (
			snsPublisher *aws.SNSPublisher
			payload      string
		)

		BeforeEach(func() {
			testClient = aws.NewStandardClient(testVaultClient, testOptions)
			Expect(testClient).ToNot(BeNil())
			snsPublisher = aws.NewSNSPublisher(testClient, fakeTopicArn)
			Expect(snsPublisher).ToNot(BeNil())

			msg := map[string]interface{}{
				"organizationName": "Test org",
				"inviteAcceptUrl":  "https://dashboard.localhost/invite/ABBA-123DEADBEEF-654321-FDSA",
				"email":            "test-invitee@callstats.io",
			}
			payloadBytes, err := json.Marshal(msg)
			Expect(err).To(BeNil())
			payload = string(payloadBytes)
		})

		Context("successfully", func() {
			It("should publish a message successfully", func() {
				mockSNSPublishRequest(testClient, fakeTopicArn, payload, 200, func() {
					err := snsPublisher.Publish(testCtx, payload)
					Expect(err).To(BeNil())
				})
			})
		})

		Context("with an error", func() {
			It("should return an error if a message wasn't published successfully", func() {
				mockSNSPublishRequest(testClient, fakeTopicArn, payload, 403, func() {
					err := snsPublisher.Publish(testCtx, payload)
					Expect(err).ToNot(BeNil())
					Expect(err.(awserr.Error).Code()).To(Equal("123"))
					Expect(err.(awserr.Error).Message()).To(Equal("AWS error"))
				})
			})
		})
	})
})
