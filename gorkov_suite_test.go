package gorkov_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGorkov(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gorkov Suite")
}
