package gorkov

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

// Tokenizer allows consuming a stream of tokens.
type Tokenizer interface {
	// Next returns the next token, if possible. If there are no more errors,
	// io.EOF will be returned. Next can also return other errors, so
	// consumers have to check for them.
	Next() (Token, error)
}

// TokenizerFunc can be used to turn a function into a Tokenizer:
type TokenizerFunc func() (Token, error)

// Next calls t().
func (t TokenizerFunc) Next() (Token, error) {
	return t()
}

// ReaderTokenizer turns data from an io.Reader into a stream of tokens. It
// turns newlines ('\n') into End tokens and returns everything else as literal
// tokens. Each literal token either only contains whitespace and punctuation or
// no whitespace and punctuation. Two tokens that follow each other do not
// contain the same type of characters.
//
// Punctuation and whitespace is everything that is a unicode punctuation
// character (category P) or has Unicode's White Space Property. See the
// unicode package for details.
type ReaderTokenizer struct {
	s *bufio.Scanner
	t Tokenizer
}

// NewTokenizer creates a new ReaderTokenizer for the given reader.
func NewTokenizer(r io.Reader) *ReaderTokenizer {
	s := bufio.NewScanner(r)
	s.Split(scanRuneType)
	return &ReaderTokenizer{s: s}
}

// Next returns the next token. See the description of ReaderTokenizer for an
// explanation of which kind of tokens to expect.
func (t *ReaderTokenizer) Next() (Token, error) {
	if t.t == nil {
		t.t = newlineToEnd(newScannerTokenizer(t.s))
	}
	return t.t.Next()
}

// newScannerTokenizer returns a Tokenizer that simply returns a literal token
// for every element the scanner outputs. Errors created by the scanner are
// returned unmodified.
func newScannerTokenizer(s *bufio.Scanner) Tokenizer {
	return TokenizerFunc(func() (Token, error) {
		if !s.Scan() {
			if s.Err() == nil {
				return nil, io.EOF
			}
			return nil, s.Err()
		}
		return Literal(s.Text()), nil
	})
}

const (
	runeNewline = iota
	runePunctuation
	runeLiteral
)

// getRuneType categorises runes into one of the following categories:
// runeNewline, runePunctuation or runeLiteral. For '\n' it will always
// return runeNewline and not runePunctuation. runeLiteral is everything
// that is neither runeNewline nor runePunctuation.
func getRuneType(r rune) int {
	switch {
	case r == '\n':
		return runeNewline
	case unicode.IsPunct(r) || unicode.IsSpace(r):
		return runePunctuation
	default:
		return runeLiteral
	}
}

// scanRuneType is a bufio.SplitFunc. It returns blocks of contiguous runes
// that have the same type as defined by getRuneType.
// TODO: simplify this method, its cyclomatic complexity is a bit high.
//
// nolint: gocyclo
func scanRuneType(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	r, width := utf8.DecodeRune(data)
	if r == utf8.RuneError {
		if !atEOF && !utf8.FullRune(data) {
			// not enough bytes for a full rune, ask for more
			return 0, nil, nil
		}
		return 0, nil, fmt.Errorf("invalid UTF8 starting at %X", data)
	}
	t, pos := getRuneType(r), width
	for pos < len(data) {
		r, width = utf8.DecodeRune(data[pos:])
		if r == utf8.RuneError {
			if !atEOF && !utf8.FullRune(data[pos:]) {
				// not enough bytes for a full rune, ask for more
				return 0, nil, nil
			}
			return 0, nil, fmt.Errorf("invalid UTF8 starting at %X", data)
		}
		if getRuneType(r) != t {
			break
		}
		pos += width
	}
	if pos == len(data) && !atEOF {
		return 0, nil, nil
	}
	return pos, data[:pos], nil
}

// newlineToEnd wraps a Tokenizer. The new Tokenizer returns the same stream
// of tokens as the original except that all newlines are replaced by End
// tokens.
//
// newlineToEnd assumes that the Value of a token is static (it doesn't change
// between calls) and that a token value either contains no newlines or only
// newlines. If a token contains x newlines then x End tokens are returned.
func newlineToEnd(t Tokenizer) Tokenizer {
	newlines := 0
	return TokenizerFunc(func() (Token, error) {
		if newlines > 0 {
			newlines--
			return End, nil
		}
		token, err := t.Next()
		if err != nil {
			return nil, err
		}
		if len(token.Value()) > 0 && getRuneType([]rune(token.Value())[0]) == runeNewline {
			newlines = len(token.Value()) - 1
			return End, nil
		}
		return token, nil
	})
}
