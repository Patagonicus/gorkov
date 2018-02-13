package gorkov_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/Patagonicus/gorkov"
)

var _ = Describe("Token", func() {
	var (
		foo        Token
		bar        Token
		whitespace Token
		newline    Token
		customType Token
	)

	BeforeEach(func() {
		foo = Literal("foo")
		bar = Literal("bar")
		whitespace = Literal(" ")
		newline = Literal("\n")
		customType = NewToken("custom", "token using a custom type")
	})

	Describe("Getting the type", func() {
		Context("With literal tokens", func() {
			It("should return LiteralType", func() {
				Expect(foo.Type()).To(Equal(LiteralType))
				Expect(bar.Type()).To(Equal(LiteralType))
				Expect(whitespace.Type()).To(Equal(LiteralType))
				Expect(newline.Type()).To(Equal(LiteralType))
			})
		})
		Context("With custom type", func() {
			It("should return the custom type", func() {
				Expect(customType.Type()).To(Equal("custom"))
			})
		})
	})

	Describe("Getting the identifier", func() {
		It("should return the value", func() {
			Expect(foo.Identifier()).To(Equal("foo"))
			Expect(bar.Identifier()).To(Equal("bar"))
			Expect(whitespace.Identifier()).To(Equal(" "))
			Expect(newline.Identifier()).To(Equal("\n"))
			Expect(customType.Identifier()).To(Equal("token using a custom type"))
		})
	})

	Describe("Getting the value", func() {
		It("should return the value", func() {
			Expect(foo.Value()).To(Equal("foo"))
			Expect(bar.Value()).To(Equal("bar"))
			Expect(whitespace.Value()).To(Equal(" "))
			Expect(newline.Value()).To(Equal("\n"))
			Expect(customType.Value()).To(Equal("token using a custom type"))
		})
	})
})
