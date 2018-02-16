package gorkov_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/Patagonicus/gorkov"
)

var _ = Describe("Prefix", func() {
	assertPrefixEquals := func(actual, expected Prefix) {
		ExpectWithOffset(1, len(actual)).To(Equal(len(expected)))
		for i := range actual {
			ExpectWithOffset(1, TokensEqual(actual[i], expected[i])).To(BeTrue())
		}
	}

	Describe("Creating a new prefix", func() {
		for n := 1; n < 10; n++ {
			Context(fmt.Sprintf("Of length %d", n), func() {
				length := n
				prefix := NewPrefix(length)

				It(fmt.Sprintf("Should have length %d", length), func() {
					Expect(len(prefix)).To(Equal(length))
				})

				It("Should only contain Start tokens", func() {
					for _, t := range prefix {
						Expect(TokensEqual(t, Start)).To(BeTrue(), "a new prefix should only contain Start tokens")
					}
				})
			})
		}

		Context("Of length 0", func() {
			It("Should panic", func() {
				Expect(func() { NewPrefix(0) }).To(Panic())
			})
		})

		Context("Of length -1", func() {
			It("Should panic", func() {
				Expect(func() { NewPrefix(-1) }).To(Panic())
			})
		})
	})

	Describe("Shifting tokens", func() {
		var (
			fooToken Token
			barToken Token
		)

		BeforeEach(func() {
			fooToken = Literal("foo")
			barToken = Literal("bar")
		})

		Context("For a prefix of length 1", func() {
			It("should replace the token", func() {
				prefix := NewPrefix(1)

				prefix.Shift(fooToken)
				assertPrefixEquals(prefix, Prefix{fooToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{barToken})
			})
		})
		Context("for a prefix of length 2", func() {
			It("should drop the first and add the new token at the end", func() {
				prefix := NewPrefix(2)

				prefix.Shift(fooToken)
				assertPrefixEquals(prefix, Prefix{Start, fooToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{fooToken, barToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{barToken, barToken})
			})
		})
		Context("for a prefix of length 5", func() {
			It("should drop the first and add the new token at the end", func() {
				prefix := NewPrefix(5)

				prefix.Shift(fooToken)
				assertPrefixEquals(prefix, Prefix{Start, Start, Start, Start, fooToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{Start, Start, Start, fooToken, barToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{Start, Start, fooToken, barToken, barToken})

				prefix.Shift(fooToken)
				assertPrefixEquals(prefix, Prefix{Start, fooToken, barToken, barToken, fooToken})

				prefix.Shift(barToken)
				assertPrefixEquals(prefix, Prefix{fooToken, barToken, barToken, fooToken, barToken})

				prefix.Shift(fooToken)
				assertPrefixEquals(prefix, Prefix{barToken, barToken, fooToken, barToken, fooToken})
			})
		})
	})

	Describe("Generating a key", func() {
		assertKeyContains := func(p Prefix, s string) {
			ExpectWithOffset(1, string(p.Key())).To(ContainSubstring(s))
		}

		var (
			fooToken    Token
			barToken    Token
			foobarToken Token
		)

		BeforeEach(func() {
			fooToken = Literal("foo")
			barToken = Literal("bar")
			foobarToken = Literal("foobar")
		})

		Context("With a length of 1", func() {
			var (
				start  Prefix
				foo    Prefix
				bar    Prefix
				foobar Prefix
			)

			BeforeEach(func() {
				start = Prefix{Start}
				foo = Prefix{fooToken}
				bar = Prefix{barToken}
				foobar = Prefix{foobarToken}
			})

			It("should contain the type of the token", func() {
				assertKeyContains(start, Start.Type())
				assertKeyContains(foo, fooToken.Type())
				assertKeyContains(bar, barToken.Type())
				assertKeyContains(foobar, foobarToken.Type())
			})

			It("should contain the identifier of the token", func() {
				assertKeyContains(start, Start.Identifier())
				assertKeyContains(foo, fooToken.Identifier())
				assertKeyContains(bar, barToken.Identifier())
				assertKeyContains(foobar, foobarToken.Identifier())
			})

			Specify("different prefixes have different keys", func() {
				Expect(start.Key()).NotTo(Equal(foo.Key()))
				Expect(start.Key()).NotTo(Equal(bar.Key()))
				Expect(start.Key()).NotTo(Equal(foobar.Key()))

				Expect(foo.Key()).NotTo(Equal(bar.Key()))
				Expect(foo.Key()).NotTo(Equal(foobar.Key()))

				Expect(bar.Key()).NotTo(Equal(foobar.Key()))
			})

			Specify("equal prefixes have the same key", func() {
				Expect(Prefix{Start}.Key()).To(Equal(start.Key()))
				Expect(Prefix{fooToken}.Key()).To(Equal(foo.Key()))
				Expect(Prefix{barToken}.Key()).To(Equal(bar.Key()))
				Expect(Prefix{foobarToken}.Key()).To(Equal(foobar.Key()))
			})
		})

		Context("With a length of 2", func() {
			var (
				startStart Prefix
				startFoo   Prefix
				fooBar     Prefix
				fooFoo     Prefix
				foobarEnd  Prefix
			)

			BeforeEach(func() {
				startStart = Prefix{Start, Start}
				startFoo = Prefix{Start, fooToken}
				fooBar = Prefix{fooToken, barToken}
				fooFoo = Prefix{fooToken, fooToken}
				foobarEnd = Prefix{foobarToken, End}
			})

			It("should contain the type of the first token", func() {
				assertKeyContains(startStart, Start.Type())
				assertKeyContains(startFoo, Start.Type())
				assertKeyContains(fooBar, fooToken.Type())
				assertKeyContains(fooFoo, fooToken.Type())
				assertKeyContains(foobarEnd, foobarToken.Type())
			})

			It("should contain the type of the second token", func() {
				assertKeyContains(startStart, Start.Type())
				assertKeyContains(startFoo, fooToken.Type())
				assertKeyContains(fooBar, barToken.Type())
				assertKeyContains(fooFoo, fooToken.Type())
				assertKeyContains(foobarEnd, End.Type())
			})

			It("should contain the identifier of the first token", func() {
				assertKeyContains(startStart, Start.Identifier())
				assertKeyContains(startFoo, Start.Identifier())
				assertKeyContains(fooBar, fooToken.Identifier())
				assertKeyContains(fooFoo, fooToken.Identifier())
				assertKeyContains(foobarEnd, foobarToken.Identifier())
			})

			It("should contain the identifier of the second token", func() {
				assertKeyContains(startStart, Start.Identifier())
				assertKeyContains(startFoo, fooToken.Identifier())
				assertKeyContains(fooBar, barToken.Identifier())
				assertKeyContains(fooFoo, fooToken.Identifier())
				assertKeyContains(foobarEnd, End.Identifier())
			})

			Specify("different prefixes have different keys", func() {
				Expect(startStart.Key()).NotTo(Equal(startFoo.Key()))
				Expect(startStart.Key()).NotTo(Equal(fooBar.Key()))
				Expect(startStart.Key()).NotTo(Equal(fooFoo.Key()))
				Expect(startStart.Key()).NotTo(Equal(foobarEnd.Key()))

				Expect(startFoo.Key()).NotTo(Equal(fooBar.Key()))
				Expect(startFoo.Key()).NotTo(Equal(fooFoo.Key()))
				Expect(startFoo.Key()).NotTo(Equal(foobarEnd.Key()))

				Expect(fooBar.Key()).NotTo(Equal(fooFoo.Key()))
				Expect(fooBar.Key()).NotTo(Equal(foobarEnd.Key()))

				Expect(fooFoo.Key()).NotTo(Equal(foobarEnd.Key()))
			})

			Specify("equal prefixes have the same key", func() {
				Expect(Prefix{Start, Start}.Key()).To(Equal(startStart.Key()))
				Expect(Prefix{Start, fooToken}.Key()).To(Equal(startFoo.Key()))
				Expect(Prefix{fooToken, barToken}.Key()).To(Equal(fooBar.Key()))
				Expect(Prefix{fooToken, fooToken}.Key()).To(Equal(fooFoo.Key()))
				Expect(Prefix{foobarToken, End}.Key()).To(Equal(foobarEnd.Key()))
			})
		})
	})
})
