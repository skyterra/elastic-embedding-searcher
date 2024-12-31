package runner

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEmb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Embedding Suite")
}

var _ = BeforeSuite(func() {

})
