package request_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ===== TEST SETUP =====
func TestAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Request Black Box Suite")
}
