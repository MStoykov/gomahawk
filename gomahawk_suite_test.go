package gomahawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGomahawk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomahawk Suite")
}
