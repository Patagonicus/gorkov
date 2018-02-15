package gorkov_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/Patagonicus/gorkov"
)

var _ = Describe("Token", func() {
	// variables with a 2 contain an equal but no identical object to the one
	// without the 2 in the name.
	var (
		foo         Token
		foo2        Token
		bar         Token
		bar2        Token
		whitespace  Token
		whitespace2 Token
		newline     Token
		newline2    Token
		customType  Token
		customType2 Token
	)

	BeforeEach(func() {
		foo = Literal("foo")
		foo2 = Literal("foo")
		bar = Literal("bar")
		bar2 = Literal("bar")
		whitespace = Literal(" ")
		whitespace2 = Literal(" ")
		newline = Literal("\n")
		newline2 = Literal("\n")
		customType = NewToken("custom", "token using a custom type")
		customType2 = NewToken("custom", "token using a custom type")
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

	// nolint: dupl
	Describe("Getting the identifier", func() {
		It("should return the value", func() {
			Expect(foo.Identifier()).To(Equal("foo"))
			Expect(bar.Identifier()).To(Equal("bar"))
			Expect(whitespace.Identifier()).To(Equal(" "))
			Expect(newline.Identifier()).To(Equal("\n"))
			Expect(customType.Identifier()).To(Equal("token using a custom type"))
		})
	})

	// nolint: dupl
	Describe("Getting the value", func() {
		It("should return the value", func() {
			Expect(foo.Value()).To(Equal("foo"))
			Expect(bar.Value()).To(Equal("bar"))
			Expect(whitespace.Value()).To(Equal(" "))
			Expect(newline.Value()).To(Equal("\n"))
			Expect(customType.Value()).To(Equal("token using a custom type"))
		})
	})

	var entries []TableEntry
	names := []string{"foo", "bar", "whitespace", "newline", "customType"}
	tokens := []*Token{&foo, &bar, &whitespace, &newline, &customType}
	tokens2 := []*Token{&foo2, &bar2, &whitespace2, &newline2, &customType2}
	for i, a := range tokens {
		for j, b := range tokens {
			entries = append(entries, createTokensEqualEntry(a, b, names[i], names[j], i == j))
		}
		for j, b := range tokens2 {
			entries = append(entries, createTokensEqualEntry(a, b, names[i], names[j]+"2", i == j))
			entries = append(entries, createTokensEqualEntry(b, a, names[j]+"2", names[i], i == j))
		}
	}
	DescribeTable("TokensEqual",
		// We have to use pointers to Token because BeforeEach is run after the
		// entries for this table are created.
		func(a, b *Token, expected bool) {
			Expect(TokensEqual(*a, *b)).To(Equal(expected))
		},
		entries...,
	)
})

func createTokensEqualEntry(a, b *Token, nameA, nameB string, equal bool) TableEntry {
	var comp string
	if equal {
		comp = "=="
	} else {
		comp = "!="
	}
	return Entry(
		fmt.Sprintf("%s %s %s", nameA, comp, nameB),
		a,
		b,
		equal,
	)
}
