package matchers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/Patagonicus/gorkov"
	. "github.com/Patagonicus/gorkov/internal/matchers"
)

var _ = Describe("MatchToken", func() {
	var (
		foo    gorkov.Token
		bar    gorkov.Token
		foobar gorkov.Token
		start  gorkov.Token
		end    gorkov.Token
	)

	BeforeEach(func() {
		foo = gorkov.Literal("foo")
		bar = gorkov.Literal("bar")
		foobar = gorkov.Literal("foobar")
		start = gorkov.NewToken("s", "")
		end = gorkov.End
	})

	DescribeTable("Calling Match",
		func(expected, actual *gorkov.Token) {
			ok, err := MatchToken(*expected).Match(*actual)
			Expect(err).NotTo(HaveOccurred())
			Expect(ok).To(Equal(gorkov.TokensEqual(*expected, *actual)))
		},
		Entry("foo == foo", &foo, &foo),
		Entry("foo != bar", &foo, &bar),
		Entry("foo != foobar", &foo, &foobar),
		Entry("foo != start", &foo, &start),
		Entry("foo != end", &foo, &end),
		Entry("bar != foo", &bar, &foo),
		Entry("bar == bar", &bar, &bar),
		Entry("bar != foobar", &bar, &foobar),
		Entry("bar != start", &bar, &start),
		Entry("bar != end", &bar, &end),
		Entry("foobar != foo", &foobar, &foo),
		Entry("foobar != bar", &foobar, &bar),
		Entry("foobar == foobar", &foobar, &foobar),
		Entry("foobar != start", &foobar, &start),
		Entry("foobar != end", &foobar, &end),
		Entry("start != foo", &start, &foo),
		Entry("start != bar", &start, &bar),
		Entry("start != foobar", &start, &foobar),
		Entry("start == start", &start, &start),
		Entry("start != end", &start, &end),
		Entry("end != foo", &end, &foo),
		Entry("end != bar", &end, &bar),
		Entry("end != foobar", &end, &foobar),
		Entry("end != start", &end, &start),
		Entry("end == end", &end, &end),
	)

	Context("Calling Match(nil)", func() {
		It("should return an error", func() {
			_, err := MatchToken(foo).Match(nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Calling Match() with a non-Token", func() {
		It("should return an error", func() {
			_, err := MatchToken(foo).Match(42)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Using nil as expected", func() {
		Specify("Match() should return an error", func() {
			_, err := MatchToken(nil).Match(foo)
			Expect(err).To(HaveOccurred())
		})
	})

	// some sanity checks (making sure that there are no panics)
	Context("Calling FailureMessage", func() {
		It("should not panic", func() {
			MatchToken(foo).FailureMessage(gorkov.Literal("not matching"))
			MatchToken(bar).FailureMessage(gorkov.Literal("not matching"))
			MatchToken(foobar).FailureMessage(gorkov.Literal("not matching"))
			MatchToken(start).FailureMessage(gorkov.Literal("not matching"))
			MatchToken(end).FailureMessage(gorkov.Literal("not matching"))
		})
	})
	Context("Calling NegatedFailureMessage", func() {
		It("should not panic", func() {
			MatchToken(foo).NegatedFailureMessage(gorkov.Literal("not matching"))
			MatchToken(bar).NegatedFailureMessage(gorkov.Literal("not matching"))
			MatchToken(foobar).NegatedFailureMessage(gorkov.Literal("not matching"))
			MatchToken(start).NegatedFailureMessage(gorkov.Literal("not matching"))
			MatchToken(end).NegatedFailureMessage(gorkov.Literal("not matching"))
		})
	})
})
