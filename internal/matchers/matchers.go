package matchers

import (
	"fmt"

	"github.com/Patagonicus/gorkov"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

type tokenMatcher struct {
	expected gorkov.Token
}

// MatchToken compares the Type and Identifier of the two tokens. It is
// an error for expected to be nil or for actual to be nil or not a Token.
func MatchToken(expected gorkov.Token) types.GomegaMatcher {
	return tokenMatcher{expected}
}

func (m tokenMatcher) Match(actual interface{}) (bool, error) {
	if m.expected == nil {
		return false, fmt.Errorf("expected must not be nil")
	}

	actualT, ok := actual.(gorkov.Token)
	if !ok {
		return false, fmt.Errorf("expected Token but got %#v of type %T", actual, actual)
	}
	return m.expected.Type() == actualT.Type() && m.expected.Identifier() == actualT.Identifier(), nil
}

func (m tokenMatcher) FailureMessage(actual interface{}) string {
	return format.Message(actual, "to match", m.expected)
}

func (m tokenMatcher) NegatedFailureMessage(actual interface{}) string {
	return format.Message(actual, "not to match", m.expected)
}
