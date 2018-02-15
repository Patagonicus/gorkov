package gorkov_test

import (
	"fmt"
	"io"
	"strings"
	"testing/iotest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "github.com/Patagonicus/gorkov"
)

var _ = Describe("ReaderTokenizer", func() {
	DescribeTable("reading without any io errors",
		func(input string, expected []Token) {
			tokenizer := NewTokenizer(strings.NewReader(input))
			for _, t := range expected {
				token, err := tokenizer.Next()
				Expect(err).NotTo(HaveOccurred())
				Expect(TokensEqual(t, token)).To(BeTrue(), "%q should equal %q", t, token)
			}
			_, err := tokenizer.Next()
			Expect(err).To(Equal(io.EOF))
		},
		Entry(
			"empty input",
			"",
			nil,
		),
		Entry(
			"single word",
			"foo",
			makeTokens("foo"),
		),
		Entry(
			"only punctuation",
			",,. ",
			makeTokens(",,. "),
		),
		Entry(
			"onld a newline",
			"\n",
			[]Token{End},
		),
		Entry(
			"short text",
			"foo bar, baz\n",
			makeTokens("foo", " ", "bar", ", ", "baz", End),
		),
		Entry(
			"long text",
			"foobar baz, foo: foo foo bar baz\n",
			makeTokens("foobar", " ", "baz", ", ", "foo", ": ", "foo", " ", "foo", " ", "bar", " ", "baz", End),
		),
		Entry(
			"multiple newlines",
			"foo\nbar baz\nbarfoo\n",
			makeTokens("foo", End, "bar", " ", "baz", End, "barfoo", End),
		),
		Entry(
			"multiple newlines in a row",
			"foo\n\nbar\n\n\nfoo bar\n",
			makeTokens("foo", End, End, "bar", End, End, End, "foo", " ", "bar", End),
		),
	)
	Describe("reading with an io error", func() {
		It("should return the error", func() {
			tokenizer := NewTokenizer(iotest.TimeoutReader(iotest.OneByteReader(strings.NewReader("foo bar baz\n"))))
			for {
				_, err := tokenizer.Next()
				if err == io.EOF {
					Fail("got EOF before io error")
				} else if err != nil {
					return
				}
			}
		})
	})
	Describe("reading multibyte runes one byte at a time", func() {
		It("should return the correct token", func() {
			token, err := NewTokenizer(iotest.OneByteReader(strings.NewReader("☹"))).Next()
			Expect(err).NotTo(HaveOccurred())
			Expect(token).To(Equal(Literal("☹")))
		})
	})
	Describe("reading invalid UTF8", func() {
		Context("as the first rune", func() {
			It("should return an error", func() {
				_, err := NewTokenizer(strings.NewReader("☹"[:1])).Next()
				Expect(err).To(HaveOccurred())
				Expect(err).NotTo(Equal(io.EOF))
			})
		})
		Context("in the middle of the stream", func() {
			It("should return an error", func() {
				_, err := NewTokenizer(iotest.OneByteReader(strings.NewReader("abc" + "☹"[:1] + "def"))).Next()
				Expect(err).To(HaveOccurred())
				Expect(err).NotTo(Equal(io.EOF))
			})
		})
		Context("at the end of the stream", func() {
			It("should return an error", func() {
				_, err := NewTokenizer(iotest.OneByteReader(strings.NewReader("abc" + "☹"[:1]))).Next()
				Expect(err).To(HaveOccurred())
				Expect(err).NotTo(Equal(io.EOF))
			})
		})
	})
})

// makeTokens takes a list of Tokens and strings and turns them into a slice
// of Tokens. Tokens are simply copied, strings are run through Literal(). Any
// value that cannot be cast to Token or string will result in a panic.
func makeTokens(tokens ...interface{}) []Token {
	result := make([]Token, len(tokens))
	for i, v := range tokens {
		switch val := v.(type) {
		case Token:
			result[i] = val
		case string:
			result[i] = Literal(val)
		default:
			panic(fmt.Errorf("invalid input of type %T: %#v", val, val))
		}
	}
	return result
}
