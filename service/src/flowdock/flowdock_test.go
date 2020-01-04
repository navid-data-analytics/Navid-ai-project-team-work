package flowdock_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/callstats-io/ai-decision/service/src/flowdock"
)

func mustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	return cli, s.Close
}

func TestMessageSendSuccess(t *testing.T) {
	flow_token := "secretflowtoken"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		mustBeNil(err)
		/*    "encoding/json"
		if valid := json.Valid(body); !valid {
		  panic("Not valid json")
		}*/ // go 1.9
		if con := bytes.Contains(body, []byte(`"flow_token": "`+flow_token+`"`)); !con {
			panic("Flow token missing or incorrect")
		}
		w.Write([]byte("{}"))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := flowdock.NewClient(flow_token)
	cli.HTTPClient = httpClient

	err := cli.SendAiNotificationMessage(1, "type", "msg")
	mustBeNil(err)
}

func TestMessageSendFailed(t *testing.T) {
	failure_message := "{failure}"
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		mustBeNil(err)
		w.Write([]byte(failure_message))
	})
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	cli := flowdock.NewClient("secretflowtoken")
	cli.HTTPClient = httpClient

	err := cli.SendAiNotificationMessage(1, "type", "msg")
	if err.Error() != failure_message {
		panic("wrong error handling")
	}
}
