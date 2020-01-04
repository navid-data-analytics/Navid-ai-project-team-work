package aws_test

import (
	"context"
	"testing"

	"github.com/callstats-io/go-common/aws"
	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testOptions     *aws.Options
	testClient      *aws.StandardClient
	testVaultClient *fakeVaultAWSClient

	testCtx, testCtxCancel = context.WithCancel(context.Background())
)

var _ = BeforeSuite(func() {
	secret := &api.Secret{
		LeaseDuration: 3600,
		Data: map[string]interface{}{
			"access_key":     "ASIAJYYYY2AA5K4WIXXX",
			"secret_key":     "HSs0DYYYYYY9W81DXtI0K7X84H+OVZXK5BXXXX",
			"security_token": "AQoDYXdzEEwasAKwQyZUtZaCjVNDiXXXXXXXXgUgBBVUUbSyujLjsw6jYzboOQ89vUVIehUw/9MreAifXFmfdbjTr3g6zc0me9M+dB95DyhetFItX5QThw0lEsVQWSiIeIotGmg7mjT1//e7CJc4LpxbW707loFX1TYD1ilNnblEsIBKGlRNXZ+QJdguY4VkzXxv2urxIH0Sl14xtqsRPboV7eYruSEZlAuP3FLmqFbmA0AFPCT37cLf/vUHinSbvw49C4c9WQLH7CeFPhDub7/rub/QU/lCjjJ43IqIRo9jYgcEvvdRkQSt70zO8moGCc7pFvmL7XGhISegQpEzudErTE/PdhjlGpAKGR3d5qKrHpPYK/k480wk1Ai/t1dTa/8/3jUYTUeIkaJpNBnupQt7qoaXXXXXXXXXX",
		},
	}
	awsSecret, err := vault.NewAWSSecret(vault.NewStandardSecret(secret, nil))
	Expect(err).To(BeNil())
	testVaultClient = &fakeVaultAWSClient{secret: awsSecret}
})

var _ = AfterSuite(func() {
	testCtxCancel()
})

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AWS Black Box Suite")
}

type fakeVaultAWSClient struct {
	secret    *vault.AWSSecret
	secretErr error
}

func (vc *fakeVaultAWSClient) AWSSecret(ctx context.Context) (*vault.AWSSecret, error) {
	if vc.secretErr != nil {
		return nil, vc.secretErr
	}
	return vc.secret, nil
}
