package request_test

import (
	"net/http"

	"github.com/callstats-io/go-common/request"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IP", func() {
	Context("Success", func() {
		Context("IPv4", func() {
			It("should return the IP only with IPv4", func() {
				Expect(request.IPFromAddr("127.0.0.1:1234")).To(Equal("127.0.0.1"))
			})
			It("should return the value itself for IP onlyP", func() {
				Expect(request.IPFromAddr("127.0.0.1")).To(Equal("127.0.0.1"))
			})
		})
		Context("IPv6", func() {
			It("should return the IP only", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromAddr("[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234")).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
			It("should return the IP only with zone", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromAddr("[2001:db8:85a3:8d3:1319:8a2e:370:7348%11]:1234")).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
			It("should return the value itself for IP only", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromAddr("[2001:db8:85a3:8d3:1319:8a2e:370:7348]")).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
		})
	})
})

var _ = Describe("IPFromRequest", func() {
	Context("Success", func() {
		makeRequestWithIP := func(ip string) *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set(request.CsioClientIPHeader, ip)
			return req
		}
		makeRequestWithoutHeaderIP := func(ip string) *http.Request {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = ip
			return req
		}
		Context("IPv4", func() {
			It("should return the IP only with IPv4", func() {
				Expect(request.IPFromRequest(makeRequestWithIP("127.0.0.1:1234"))).To(Equal("127.0.0.1"))
				Expect(request.IPFromRequest(makeRequestWithoutHeaderIP("127.0.0.1:1234"))).To(Equal("127.0.0.1"))
			})
			It("should return the value itself for IP only", func() {
				Expect(request.IPFromRequest(makeRequestWithIP("127.0.0.1"))).To(Equal("127.0.0.1"))
				Expect(request.IPFromRequest(makeRequestWithoutHeaderIP("127.0.0.1"))).To(Equal("127.0.0.1"))
			})
		})
		Context("IPv6", func() {
			It("should return the IP only", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromRequest(makeRequestWithIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
				Expect(request.IPFromRequest(makeRequestWithoutHeaderIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1234"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
			It("should return the IP only with zone", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromRequest(makeRequestWithIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348%11]:1234"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
				Expect(request.IPFromRequest(makeRequestWithoutHeaderIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348%11]:1234"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
			It("should return the value itself for IP only", func() {
				// add brackets to reflect IPv6 URL IP + port
				Expect(request.IPFromRequest(makeRequestWithIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348]"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
				Expect(request.IPFromRequest(makeRequestWithoutHeaderIP("[2001:db8:85a3:8d3:1319:8a2e:370:7348]"))).To(Equal("2001:db8:85a3:8d3:1319:8a2e:370:7348"))
			})
		})
	})
})
