package request_test

import (
	"context"
	"strings"

	"github.com/callstats-io/go-common/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IDFromContext & WithID", func() {
	Context("Success", func() {
		It("should return a new ID", func() {
			ctx := context.Background()
			id := request.IDFromContext(ctx)
			Expect(id).ToNot(BeEmpty())
			// naively test for uuid-like (sanity check)
			Expect(strings.Count(id, "-")).To(Equal(4))
		})
		It("should return an existing ID", func() {
			ctx := context.Background()
			requestID := "abcdef"
			ctx = request.WithID(ctx, requestID)
			Expect(request.IDFromContext(ctx)).To(Equal(requestID))
		})
	})
})
