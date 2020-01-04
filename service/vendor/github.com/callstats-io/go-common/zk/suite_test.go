package zk_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zookeeper Black Box Suite")
}
