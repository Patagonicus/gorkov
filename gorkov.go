package gorkov

// Token is an one element of a markov chain. Usually this is a word or some
// whitespace.
type Token interface {
	// Type provides a string identifying the type of the token. Type must
	// always return the same string and that string must be not be used
	// by any other type of token used in one Gorkov instance. It may not
	// contain any null bytes.
	//
	// Tokens generated by this package will only use type namse consisting
	// of a single ASCII letter or digit (a-z, A-Z and 0-9).
	Type() string

	// Identifier returns a string identifying this particular token.
	// Identifier must always return the same string for one token and
	// that string may not contain any null bytes.
	//
	// For two tokens a and b the following must hold:
	//   a.Type() == b.Type() && a.Identifier() == b.Identifier()
	// is true iff a and b are considered equal.
	Identifier() string

	// Value returns the string that is used when generating a text using
	// this token. This is usually a static string, but can also be
	// dynamically generated.
	Value() string
}

const (
	// LiteralType is the type used for literal tokens.
	LiteralType = "l"
)

type token struct {
	t, value string
}

func (t token) Type() string {
	return t.t
}

func (t token) Identifier() string {
	return t.value
}

func (t token) Value() string {
	return t.value
}

// NewToken creates a new token with a static value. The value is also used as
// the identifier. It can be used for static tokens, such as literal words.
func NewToken(t, value string) Token {
	return token{
		t:     t,
		value: value,
	}
}

// Literal creates a new token for a literal. This is a convenience function
// and is equal to calling NewToken(LiteralType, value).
func Literal(value string) Token {
	return NewToken(LiteralType, value)
}
