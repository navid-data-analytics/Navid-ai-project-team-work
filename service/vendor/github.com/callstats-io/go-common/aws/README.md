### Using with Vault

First [bootstrap dev Vault](../tools) with AWS, then [setup a Vault auto-renew client](../vault).

Add to `docker-compose.yaml`:

```yaml
- VAULT_ENABLE_AWS=true
- VAULT_AWS_CREDS_PATH=<env>/aws/sts/<role name>
- VAULT_AWS_REGION=eu-west-1
- ORG_INVITE_SNS_TOPIC_ARN="arn:aws:sns:eu-west-1:123412341234:not-a-real-test-topic-name"
```

Match
- `<env>` with your running environment
- `<role name>` with an arbitrary role name configured to Vault

An example to setup the AWS client:

```golang
package main

import (
  "context"
  "encoding/json"
  "os"

  "github.com/callstats-io/go-common/aws"
  "github.com/callstats-io/go-common/log"
  "github.com/callstats-io/go-common/vault"
)

func main() {
  ctx := context.Background()
  logger, _ := log.FromEnv()

  // set up vault auto-renew client, details in Vault README
  var avc *vault.AutoRenewClient
  avc = setupAvc()

  awsOpts, err := aws.OptionsFromEnv()
  if err != nil {
    logger.Panic("Invalid AWS env options", log.Error(err))
  }

  topic := os.Getenv("ORG_INVITE_SNS_TOPIC_ARN")

  awsClient := aws.NewStandardClient(avc, awsOpts)
  orgInvitePublisher := aws.NewSNSPublisher(awsClient, topic)

  msg := map[string]interface{}{
    "appId": "123456798",
    "inviteUrl": "https://dashboard.callstats.io/invite/ABBA-123DEADBEEF-654321-FDSA",
    "emailAddress": "invitee@callstats.io",
  }
  payload, _ := json.Marshal(msg)

  err := orgInvitePublisher.Publish(ctx, payload)
  if err != nil {
    logger.Error("SNS publish error", log.Error(err))
    return
  }
  logger.Info("published an SNS message")
}
```
