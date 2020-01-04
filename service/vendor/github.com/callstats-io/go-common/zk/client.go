package zk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// kafkaBrokerAddress represents Kafka broker connections paramaters
type kafkaBrokerAddress struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// Client represents ZK client object
type Client struct {
	pathPrefix       string
	kafkaBrokersPath string
	connection       *zk.Conn
}

// NewClient creates a new Client object, or returns error if Client creation fails
func NewClient(zkServers []string, sessionTimeout time.Duration, options ...Option) (*Client, error) {
	c := &Client{}

	for _, opt := range options {
		opt(c)
	}

	conn, _, err := zk.Connect(zkServers, sessionTimeout)
	if err != nil {
		return nil, err
	}
	c.connection = conn
	return c, nil
}

// ListChildren returns a list of all children Z-nodes located under "path"
func (client *Client) ListChildren(path string) ([]string, error) {
	children, _, err := client.connection.Children(path)
	return children, err
}

// Get returns value of the Z-node located under "path"
func (client *Client) Get(path string) ([]byte, error) {
	object, _, err := client.connection.Get(path)
	return object, err
}

// Set sets value of the Z-node located under "path"
func (client *Client) Set(path string, data []byte) error {
	_, stat, err := client.connection.Get(path)
	if err != nil && err.Error() != "zk: node does not exist" {
		return err
	}

	if err != nil {
		_, err = client.connection.Create(path, data, 0, zk.WorldACL(zk.PermAll))
		return err
	}

	_, err = client.connection.Set(path, data, stat.Version)
	return err
}

// KafkaBrokers returns list of Kafka brokers stored in Zookeeper identified by zkServers
func (client *Client) KafkaBrokers() ([]string, error) {
	var brokerAddr kafkaBrokerAddress
	var brokerList []string
	if client.kafkaBrokersPath == "" {
		client.kafkaBrokersPath = client.pathPrefix + "/brokers/ids"
	}

	brokers, err := client.ListChildren(client.kafkaBrokersPath)
	if err != nil {
		return nil, err
	}
	for _, brokerName := range brokers {
		brokerPath := client.kafkaBrokersPath + "/" + brokerName
		brokerNode, err := client.Get(brokerPath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(brokerNode, &brokerAddr)
		if err != nil {
			return nil, err
		}

		brokerList = append(brokerList, fmt.Sprintf("%s:%d", brokerAddr.Host, brokerAddr.Port))
	}
	return brokerList, nil
}

// Close the Zookeeper connection
func (client *Client) Close() {
	client.connection.Close()
}
