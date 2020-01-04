package middleware_test

import (
	"net/http"
	"strings"

	"github.com/callstats-io/go-common/middleware"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Middleware", func() {
	var calls []string
	middlewareA := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "A1")
			h(w, r)
			calls = append(calls, "A2")
		}
	}
	middlewareB := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "B1")
			h(w, r)
			calls = append(calls, "B2")
		}
	}
	middlewareC := func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			calls = append(calls, "C1")
			h(w, r)
			calls = append(calls, "C2")
		}
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "Handler")
	}

	BeforeEach(func() {
		calls = make([]string, 0, 4)
	})

	It("should call the registered chain in order", func() {
		c := middleware.NewConfigurator()
		c.Use(middlewareA)
		c.Use(middlewareB)
		c.Use(middlewareC)
		c.ToHandler(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|C1|Handler|C2|B2|A2"))
	})

	It("should call the registered chain in order when added by Append", func() {
		middleware.NewConfigurator().
			Append(middlewareA).
			Append(middlewareB).
			Append(middlewareC).
			Sink(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|C1|Handler|C2|B2|A2"))
	})

	It("should call the registered chain in reverse order when added by Prepend", func() {
		middleware.NewConfigurator().
			Prepend(middlewareA).
			Prepend(middlewareB).
			Prepend(middlewareC).
			Sink(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("C1|B1|A1|Handler|A2|B2|C2"))
	})

	It("should call the registered chain in order when added by both Append and Prepend", func() {
		middleware.NewConfigurator().
			Append(middlewareA).
			Append(middlewareB).
			Prepend(middlewareC).
			Sink(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("C1|A1|B1|Handler|B2|A2|C2"))
	})

	It("should stop processing on first middleware return", func() {
		middlewareQuit := func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, "Quit")
				return
			}
		}
		c := middleware.NewConfigurator()
		c.Use(middlewareA)
		c.Use(middlewareB)
		c.Use(middlewareQuit)
		c.Use(middlewareC)
		c.ToHandler(handler)(nil, nil)
		// Expect C and handler not to be called
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|Quit|B2|A2"))
	})

	It("should give an immutable chain on ToHandler", func() {
		c := middleware.NewConfigurator()
		c.Use(middlewareA)
		c.Use(middlewareB)
		chainA := c.ToHandler(handler)
		c.Use(middlewareC)
		chainB := c.ToHandler(handler)

		// Expect no C middleware calls
		chainA(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|Handler|B2|A2"))

		// Expect second chain with C middleware calls
		calls = make([]string, 0, 4)
		chainB(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|C1|Handler|C2|B2|A2"))
	})

	It("should give a full copy on Clone", func() {
		c := middleware.NewConfigurator()
		c.Use(middlewareA)
		c.Use(middlewareB)
		mOther := c.Clone()
		c.Use(middlewareC)

		// Expect C middleware calls
		c.ToHandler(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|C1|Handler|C2|B2|A2"))

		// Expect second chain without C calls
		calls = make([]string, 0, 3)
		mOther.ToHandler(handler)(nil, nil)
		Expect(strings.Join(calls, "|")).To(Equal("A1|B1|Handler|B2|A2"))
	})
})
