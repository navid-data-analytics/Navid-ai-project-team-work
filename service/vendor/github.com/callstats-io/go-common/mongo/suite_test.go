package mongo_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongo Black Box Suite")
}
