package gorkov

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
