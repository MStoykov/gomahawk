package gomahawk

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GomahawkServer", func() {
	Context("With a brand new GomahawkServer", func() {
		It("Will start :)", func() {
			tg := NewFakeGomahawk()
			gs, err := NewGomahawkServer(tg)
			Expect(err).ToNot(HaveOccurred())
			Expect(gs.Start()).ToNot(HaveOccurred())
		})
	})
})
