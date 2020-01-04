package zk_test

import (
	"os"
	"strings"
	"time"

	"github.com/callstats-io/go-common/zk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	zkAddresses := strings.Split(os.Getenv("ZOOKEEPER_URLS"), ",")
	zkSessionTimeout := 10 * time.Second
	zkPathPrefix := "/test"

	Context("NewZkClient", func() {
		Context("Success", func() {
			It("should return a new client with valid connection", func() {
				client, err := zk.NewClient(zkAddresses, zkSessionTimeout, zk.OptionPathPrefix(zkPathPrefix))
				Expect(err).To(BeNil())
				Expect(client).ToNot(BeNil())
				client.Close()
			})
		})
		Context("Failure", func() {
			It("should return an error if servers are invalid", func() {
				_, err := zk.NewClient([]string{"abc%2F"}, zkSessionTimeout)
				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("Reading ZK data", func() {
		var client *zk.Client

		BeforeEach(func() {
			client, _ = zk.NewClient(zkAddresses, zkSessionTimeout)
			client.Set("/data", []byte("test-data"))
		})
		AfterEach(func() {
			client.Close()
		})

		It("should read existing data from Zookeeper", func() {
			data, err := client.Get("/data")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("test-data")))
		})

		It("should fail to read not existing data from Zookeeper", func() {
			_, err := client.Get("/aaa")
			Expect(err).ToNot(BeNil())
		})
	})

	Context("Setting ZK data", func() {
		var client *zk.Client

		BeforeEach(func() {
			client, _ = zk.NewClient(zkAddresses, zkSessionTimeout)
		})
		AfterEach(func() {
			client.Close()
		})

		It("should set data in Zookeeper", func() {
			client.Set("/data", []byte("test-data"))
			data, err := client.Get("/data")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("test-data")))
		})

		It("should update existing data in Zookeeper", func() {
			client.Set("/data", []byte("test-data"))
			data, err := client.Get("/data")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("test-data")))

			client.Set("/data", []byte("new-test-data"))
			data, err = client.Get("/data")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("new-test-data")))
		})

		It("should list ZK nodes", func() {
			client.Set("/data", nil)
			client.Set("/data/test", []byte("test-data"))
			data, err := client.Get("/data/test")
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("test-data")))

			nodes, err := client.ListChildren("/data")
			Expect(err).To(BeNil())
			Expect(nodes).To(Equal([]string{"test"}))
		})
	})
	Context("KafkaBrokers", func() {
		It("should return list of brokers", func() {
			client, err := zk.NewClient(zkAddresses, zkSessionTimeout, zk.OptionPathPrefix(zkPathPrefix))
			Expect(err).To(BeNil())
			defer client.Close()
			client.Set(zkPathPrefix, nil)
			client.Set(zkPathPrefix+"/brokers", nil)
			client.Set(zkPathPrefix+"/brokers/ids", nil)
			client.Set(zkPathPrefix+"/brokers/ids/0", []byte("{\"host\":\"kafka\",\"port\":9092}"))

			brokerList, err := client.KafkaBrokers()
			Expect(err).To(BeNil())
			Expect(brokerList).ToNot(BeNil())
			Expect(brokerList).To(HaveLen(1))
		})

		It("should return list of brokers without prefix", func() {
			client, err := zk.NewClient(zkAddresses, zkSessionTimeout)
			Expect(err).To(BeNil())
			defer client.Close()
			client.Set("/brokers", nil)
			client.Set("/brokers/ids", nil)
			client.Set("/brokers/ids/0", []byte("{\"host\":\"kafka\",\"port\":9092}"))

			brokerList, err := client.KafkaBrokers()
			Expect(err).To(BeNil())
			Expect(brokerList).ToNot(BeNil())
			Expect(brokerList).To(HaveLen(1))
		})

		It("should return an error if path does not exist", func() {
			client, err := zk.NewClient(zkAddresses, zkSessionTimeout, zk.OptionPathPrefix("/testnonexistent"))
			Expect(err).To(BeNil())
			defer client.Close()

			_, err = client.KafkaBrokers()
			Expect(err).ToNot(BeNil())
		})
	})
})
