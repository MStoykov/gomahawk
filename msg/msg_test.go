package msg_test

import (
	"bytes"

	. "github.com/MStoykov/gomahawk/msg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Msg", func() {
	var (
		msg *Msg
		err error
	)

	Context("Parsing an addfiles msg", func() {
		BeforeEach(func() {
			msg, err = ReadMSG(getReaderToFixture("addfiles.msg"))
		})

		It("has a msg", func() {
			Expect(err).ToNot(HaveOccurred())
			Expect(msg).ToNot(BeNil())
		})

		It("has correct DBOP Flag", func() {
			Expect(msg.IsDBOP()).To(BeTrue())
		})

		It("has correct JSON Flag", func() {
			Expect(msg.IsJSON()).To(BeTrue())
		})

		It("payload to be parsable", func() {
			parser := NewCommandParser()
			j, err := parser.ParseCommand(msg)
			Expect(err).ToNot(HaveOccurred())
			Expect(j).ToNot(BeNil())
		})

		It("written to bytes it returns the same bytes", func() {
			buf := new(bytes.Buffer)
			msg.WriteTo(buf)
			Expect(buf.Bytes()).To(Equal(getFixture("addfiles.msg")))
		})
	})

	Context("Parsing a trigger msg", func() {
		BeforeEach(func() {
			msg = NewMsg([]byte(`{"method":"trigger"}`), JSON)
		})

		It("IsTrigger returns true", func() {
			Expect(IsTrigger(msg)).To(BeTrue())
		})
	})
})
