package vault_test

import (
	"crypto/tls"
	"errors"
	"time"

	"github.com/callstats-io/go-common/vault"
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func runSharedSecretTestCases(secretFunc func(leaseSec int, parent vault.Secret) vault.Secret) {
	It("should set the created at as time.Now()", func() {
		secret := secretFunc(300, nil)
		Expect(secret.CreateTime().UnixNano()).To(BeNumerically("~", time.Now().UnixNano(), time.Second))
	})
	It("should set the renew from secret lease seconds", func() {
		secret := secretFunc(300, nil)
		// assume renew: lease - 1min - 0-1000ms random offset
		Expect(secret.RenewTime().Sub(time.Now())).To(BeNumerically("~", 240*time.Second, time.Second))
	})
	It("should set the expire from secret lease seconds", func() {
		secret := secretFunc(300, nil)
		// assume expire: lease - 1 second
		Expect(secret.ExpireTime().Sub(time.Now())).To(BeNumerically("~", 299*time.Second, time.Second))
	})
	It("should set the context deadline to expire time", func() {
		secret := secretFunc(300, nil)
		dl, exists := secret.ExpireContext().Deadline()
		Expect(exists).To(BeTrue())
		Expect(dl).To(Equal(secret.ExpireTime()))
	})
	It("should cancel the expire context when parent context expires", func() {
		secret := secretFunc(300, nil)
		secret2 := secretFunc(300, secret)
		secret.Cancel()
		Eventually(func() bool {
			select {
			case <-secret2.ExpireContext().Done():
				return true
			default:
				return false
			}
		}).Should(BeTrue(), "expected parent secret cancel to cancel child secret expire context")
	})
	It("should set the context deadline to renew time", func() {
		secret := secretFunc(300, nil)
		dl, exists := secret.RenewContext().Deadline()
		Expect(exists).To(BeTrue())
		Expect(dl).To(Equal(secret.RenewTime()))
	})
	It("should cancel the renew context when parent context expires", func() {
		secret := secretFunc(300, nil)
		secret2 := secretFunc(300, secret)
		secret.Cancel()
		Eventually(func() bool {
			select {
			case <-secret2.RenewContext().Done():
				return true
			default:
				return false
			}
		}).Should(BeTrue(), "expected parent secret cancel to cancel child secret renew context")
	})
	It("should support leases with long lease times", func() {
		secret := secretFunc(int(24*time.Hour/time.Second), nil)
		dl, _ := secret.RenewContext().Deadline()
		Expect(dl).To(Equal(secret.RenewTime()))
		dl, _ = secret.ExpireContext().Deadline()
		Expect(dl).To(Equal(secret.ExpireTime()))
	})
}

type sharedSecretTestCase struct {
	Name            string
	Key             string
	ErrEmpty        error
	ErrFormat       error
	BaseSecret      func(leaseDuration int) *api.Secret
	TransformSecret func(*api.Secret) (vault.Secret, error)
}

func (t *sharedSecretTestCase) Exec() {
	It("should return error if the "+t.Name+" is not a string", func() {
		secret := t.BaseSecret(300)
		secret.Data[t.Key] = []byte(secret.Data[t.Key].(string))
		_, err := t.TransformSecret(secret)
		Expect(err).To(MatchError(t.ErrFormat))
	})
	It("should return error if the "+t.Name+" is not present in secret", func() {
		secret := t.BaseSecret(300)
		delete(secret.Data, t.Key)
		_, err := t.TransformSecret(secret)
		Expect(err).To(MatchError(t.ErrEmpty))
	})
}

var _ = Describe("StandardSecret", func() {
	makeTestRawSecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			LeaseDuration: leaseDuration,
		}
	}
	makeTestAuthRawSecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			Auth: &api.SecretAuth{
				LeaseDuration: leaseDuration,
			},
		}
	}

	runSharedSecretTestCases(func(leaseDuration int, parent vault.Secret) vault.Secret {
		return vault.NewStandardSecret(makeTestRawSecret(leaseDuration), parent)
	})

	It("should support auth secrets", func() {
		leaseDuration := 24 * time.Hour
		secret := vault.NewStandardSecret(makeTestAuthRawSecret(int(leaseDuration/time.Second)), nil)
		Expect(secret.RenewTime().Sub(time.Now())).To(BeNumerically("~", leaseDuration-time.Minute, time.Second))
		Expect(secret.ExpireTime().Sub(time.Now())).To(BeNumerically("~", leaseDuration-time.Second, time.Second))
	})
	It("should be valid if the secret or its parent has not expired/canceled", func() {
		parent := vault.NewStandardSecret(makeTestRawSecret(int(time.Hour/time.Second)), nil)
		secret := vault.NewStandardSecret(makeTestRawSecret(int(time.Minute/time.Second)), parent) // secret which requires renewal but has not expired
		Expect(secret.RenewTime().Before(time.Now())).To(BeTrue())
		Expect(secret.Valid()).To(BeTrue())

		// should be false if the secret is canceled
		secret.Cancel()
		Expect(secret.Valid()).To(BeFalse())

		// should be false if the secret is canceled
		secret = vault.NewStandardSecret(makeTestRawSecret(int(time.Hour/time.Second)), parent)
		parent.Cancel()
		Expect(secret.Valid()).To(BeFalse())
	})
})

var _ = Describe("UserPassSecret", func() {
	makeTestRawSecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			LeaseDuration: leaseDuration,
			Data: map[string]interface{}{
				vault.SecretDataKeyUsername: "vault",
				vault.SecretDataKeyPassword: "vault",
			},
		}
	}

	runSharedSecretTestCases(func(leaseDuration int, parent vault.Secret) vault.Secret {
		s, err := vault.NewUserPassSecret(vault.NewStandardSecret(makeTestRawSecret(leaseDuration), parent))
		Expect(err).To(BeNil())
		return s
	})

	Context("Success", func() {
		It("should parse credentials from the secret", func() {
			secret := makeTestRawSecret(300)
			testSecret, err := vault.NewUserPassSecret(vault.NewStandardSecret(secret, nil))
			Expect(err).To(BeNil())
			Expect(testSecret.Credentials.User).To(Equal(testSecret.Data[vault.SecretDataKeyUsername]))
			Expect(testSecret.Credentials.Password).To(Equal(testSecret.Data[vault.SecretDataKeyPassword]))
		})
	})
	Context("Failure", func() {
		testCases := []sharedSecretTestCase{
			sharedSecretTestCase{
				Name:       "username",
				Key:        vault.SecretDataKeyUsername,
				ErrEmpty:   vault.ErrEmptySecretUsernameData,
				ErrFormat:  vault.ErrInvalidSecretUsernameFormat,
				BaseSecret: makeTestRawSecret,
				TransformSecret: func(s *api.Secret) (vault.Secret, error) {
					return vault.NewUserPassSecret(vault.NewStandardSecret(s, nil))
				},
			},
			sharedSecretTestCase{
				Name:       "password",
				Key:        vault.SecretDataKeyPassword,
				ErrEmpty:   vault.ErrEmptySecretPasswordData,
				ErrFormat:  vault.ErrInvalidSecretPasswordFormat,
				BaseSecret: makeTestRawSecret,
				TransformSecret: func(s *api.Secret) (vault.Secret, error) {
					return vault.NewUserPassSecret(vault.NewStandardSecret(s, nil))
				},
			},
		}
		for idx := range testCases {
			testCase := &testCases[idx]
			testCase.Exec()
		}
	})
})

var _ = Describe("TLSCertSecret", func() {
	makeTestRawKeySecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			LeaseDuration: leaseDuration,
			Data: map[string]interface{}{
				vault.SecretDataKeyData: testVaultBootstrapConfig.TLSCertKeyData,
			},
		}
	}
	makeTestRawCertSecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			LeaseDuration: leaseDuration,
			Data: map[string]interface{}{
				vault.SecretDataKeyData: testVaultBootstrapConfig.TLSCertData,
			},
		}
	}

	runSharedSecretTestCases(func(leaseDuration int, parent vault.Secret) vault.Secret {
		certSecret := vault.NewStandardSecret(makeTestRawCertSecret(leaseDuration), parent)
		keySecret := vault.NewStandardSecret(makeTestRawKeySecret(leaseDuration), parent)
		s, err := vault.NewTLSCertSecret(certSecret, keySecret)
		Expect(err).To(BeNil())
		return s
	})

	var certSecret, keySecret *vault.StandardSecret
	BeforeEach(func() {
		certSecret = vault.NewStandardSecret(makeTestRawCertSecret(300), nil)
		keySecret = vault.NewStandardSecret(makeTestRawKeySecret(300), nil)
	})

	Context("Success", func() {
		It("should store the raw secret to the authSecret", func() {
			for idx, duration := range []int{240, 360} {
				s := vault.NewStandardSecret(makeTestRawKeySecret(duration), nil)
				secret, err := vault.NewTLSCertSecret(certSecret, s)
				Expect(err).To(BeNil())
				if idx == 0 {
					Expect(secret.Secret).To(Equal(s.Secret))
				} else {
					Expect(secret.Secret).To(Equal(certSecret.Secret))
				}
			}
		})
		It("should parse the certificate from the secret", func() {
			expCert, err := tls.X509KeyPair([]byte(testVaultBootstrapConfig.TLSCertData), []byte(testVaultBootstrapConfig.TLSCertKeyData))
			Expect(err).To(BeNil())
			secret, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(BeNil())
			Expect(secret.Certificate).To(Equal(&expCert))
		})
	})
	Context("Failure", func() {
		It("should return error if the cert data is not present in secret", func() {
			delete(certSecret.Data, vault.SecretDataKeyData)
			_, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(MatchError(vault.ErrEmptySecretTLSCertData))
		})
		It("should return error if the cert key data is not present in secret", func() {
			delete(keySecret.Data, vault.SecretDataKeyData)
			_, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(MatchError(vault.ErrEmptySecretTLSCertKeyData))
		})
		It("should return error if the cert data is not a string", func() {
			certSecret.Data[vault.SecretDataKeyData] = []byte(certSecret.Data[vault.SecretDataKeyData].(string))
			_, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(MatchError(vault.ErrInvalidSecretTLSCertFormat))
		})
		It("should return error if the cert data key is not a string", func() {
			keySecret.Data[vault.SecretDataKeyData] = []byte(keySecret.Data[vault.SecretDataKeyData].(string))
			_, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(MatchError(vault.ErrInvalidSecretTLSCertKeyFormat))
		})
		It("should return error if the cert data is invalid", func() {
			keySecret.Data[vault.SecretDataKeyData] = "abc"
			_, err := vault.NewTLSCertSecret(certSecret, keySecret)
			Expect(err).To(MatchError(errors.New("tls: failed to find any PEM data in key input")))
		})
	})
})

var _ = Describe("AWSSecret", func() {
	makeTestRawAWSSecret := func(leaseDuration int) *api.Secret {
		return &api.Secret{
			LeaseDuration: leaseDuration,
			Data: map[string]interface{}{
				vault.SecretDataKeyAccessKey:     "vault_acc_key",
				vault.SecretDataKeySecretKey:     "vault_sec_key",
				vault.SecretDataKeySecurityToken: "vault_sec_token",
			},
		}
	}

	runSharedSecretTestCases(func(leaseDuration int, parent vault.Secret) vault.Secret {
		s, err := vault.NewAWSSecret(vault.NewStandardSecret(makeTestRawAWSSecret(leaseDuration), parent))
		Expect(err).To(BeNil())
		return s
	})

	Context("Success", func() {
		It("should parse credentials from the secret", func() {
			secret := makeTestRawAWSSecret(300)
			testSecret, err := vault.NewAWSSecret(vault.NewStandardSecret(secret, nil))
			Expect(err).To(BeNil())
			Expect(testSecret.Credentials.AccessKey).To(Equal(testSecret.Data[vault.SecretDataKeyAccessKey]))
			Expect(testSecret.Credentials.SecretKey).To(Equal(testSecret.Data[vault.SecretDataKeySecretKey]))
			Expect(testSecret.Credentials.SecurityToken).To(Equal(testSecret.Data[vault.SecretDataKeySecurityToken]))
		})
	})
	Context("Failure", func() {
		testCases := []sharedSecretTestCase{
			sharedSecretTestCase{
				Name:       "access_key",
				Key:        vault.SecretDataKeyAccessKey,
				ErrEmpty:   vault.ErrEmptySecretAccessKeyData,
				ErrFormat:  vault.ErrInvalidSecretAccessKeyFormat,
				BaseSecret: makeTestRawAWSSecret,
				TransformSecret: func(s *api.Secret) (vault.Secret, error) {
					return vault.NewAWSSecret(vault.NewStandardSecret(s, nil))
				},
			},
			sharedSecretTestCase{
				Name:       "secret_key",
				Key:        vault.SecretDataKeySecretKey,
				ErrEmpty:   vault.ErrEmptySecretSecretKeyData,
				ErrFormat:  vault.ErrInvalidSecretSecretKeyFormat,
				BaseSecret: makeTestRawAWSSecret,
				TransformSecret: func(s *api.Secret) (vault.Secret, error) {
					return vault.NewAWSSecret(vault.NewStandardSecret(s, nil))
				},
			},
			sharedSecretTestCase{
				Name:       "security_token",
				Key:        vault.SecretDataKeySecurityToken,
				ErrEmpty:   vault.ErrEmptySecretSecurityTokenData,
				ErrFormat:  vault.ErrInvalidSecretSecurityTokenFormat,
				BaseSecret: makeTestRawAWSSecret,
				TransformSecret: func(s *api.Secret) (vault.Secret, error) {
					return vault.NewAWSSecret(vault.NewStandardSecret(s, nil))
				},
			},
		}
		for idx := range testCases {
			testCase := &testCases[idx]
			testCase.Exec()
		}
	})
})
