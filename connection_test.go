package gomahawk

import (
	"io"

	msg "github.com/MStoykov/gomahawk/msg"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeConn struct {
	input        *io.PipeReader
	inputWriter  *io.PipeWriter
	output       *io.PipeWriter
	outputReader *io.PipeReader
}

func newFakeConn() *FakeConn {
	f := new(FakeConn)
	f.input, f.inputWriter = io.Pipe()
	f.outputReader, f.output = io.Pipe()
	return f
}

func newDoubleFakeConn() (*FakeConn, *FakeConn) {
	var result [2]*FakeConn
	result[0] = newFakeConn()
	result[1] = new(FakeConn)
	result[1].input = result[0].outputReader
	result[1].output = result[0].inputWriter
	return result[0], result[1]
}

func (f *FakeConn) Read(b []byte) (int, error) {
	return f.input.Read(b)
}
func (f *FakeConn) Write(b []byte) (int, error) {
	return f.output.Write(b)
}

func (f *FakeConn) Close() error {
	return nil // NOOP
}

var _ = Describe("Connection", func() {

	var (
		conn     *connection
		fakeConn *FakeConn
	)

	Context("When it's brand new", func() {
		BeforeEach(func() {
			conn = new(connection)
			fakeConn = newFakeConn()
			conn.conn = fakeConn
		})

		It("sendVersionCheck will work", func() {
			c := make(chan string)
			go func() {
				err := conn.sendVersionCheck()
				Expect(err).ToNot(HaveOccurred())
				c <- "Done"
			}()

			m, err := msg.ReadMSG(fakeConn.outputReader)
			Expect(err).ToNot(HaveOccurred())
			Expect(m.IsSetup()).To(BeTrue())
			Expect(m.Payload()).To(Equal([]byte{'4'}))
			Expect(<-c).To(Equal("Done"))
		})
	})

	Context("With two connected connections", func() {
		var offerPairing = func(offer *msg.Msg, name string) {
			It("SendOffer and ReceiveOffer will work against one another for"+name, func() {
				var connections [2]*connection
				connections[0] = new(connection)
				connections[1] = new(connection)
				connections[0].conn, connections[1].conn = newDoubleFakeConn()
				c := make(chan string)
				go func() {
					err := connections[0].receiveOffer()
					Expect(err).ToNot(HaveOccurred())
					c <- "Done"
				}()
				go func() {
					err := connections[1].sendOffer(offer)
					Expect(err).ToNot(HaveOccurred())
					c <- "Done"
				}()

				Expect(<-c).To(Equal("Done"))
				Expect(<-c).To(Equal("Done"))
			})
		}
		offerPairing(msg.NewFileRequestOffer(10, "someid"), "StreamConnection")
		offerPairing(msg.NewSecondaryOffer("someid", "otherid", 11111), "Random Secondary Connection")
	})
})
