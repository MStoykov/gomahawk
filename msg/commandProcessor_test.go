package msg_test

import (
	"bytes"
	"io/ioutil"

	. "github.com/MStoykov/gomahawk/msg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CommandProcessor", func() {
	var (
		msgWrapper = func(b []byte) *Msg {
			return NewMsg(bytes.NewBuffer(b), JSON|DBOP)
		}
		processor           *CommandProcessor
		addfilesJSON, _     = ioutil.ReadFile("./fixtures/addfiles.json")
		deletefilesJSON, _  = ioutil.ReadFile("./fixtures/deletefiles.json")
		logplaybackJSON, _  = ioutil.ReadFile("./fixtures/logplayback.json")
		logplayback2JSON, _ = ioutil.ReadFile("./fixtures/logplayback2.json")
	)

	BeforeEach(func() {
		processor = NewCommandProcessor()
	})

	Describe("Processor", func() {

		It("parses AddFiles", func() {
			command, err := processor.ParseCommand(msgWrapper(addfilesJSON))
			Expect(err).ToNot(HaveOccurred())
			addFiles, ok := command.(*AddFiles)
			Expect(ok).To(BeTrue())

			id := addFiles.Files[0].Id
			Expect(id).To(BeNumerically("==", 1))

		})

		It("parses DeleteFiles", func() {
			command, err := processor.ParseCommand(msgWrapper(deletefilesJSON))
			Expect(err).ToNot(HaveOccurred())
			deleteFiles, ok := command.(*DeleteFiles)
			Expect(ok).To(BeTrue())

			id := deleteFiles.Ids[2]
			Expect(id).To(BeNumerically("==", 353))

		})

		It("parses LogPlayback", func() {
			command, err := processor.ParseCommand(msgWrapper(logplaybackJSON))
			Expect(err).ToNot(HaveOccurred())
			deleteFiles, ok := command.(*LogPlayback)
			Expect(ok).To(BeTrue())

			Expect(deleteFiles.Track).To(Equal("Brass Monkey"))
			Expect(deleteFiles.Action).To(BeNumerically("==", 1))

			command, err = processor.ParseCommand(msgWrapper(logplayback2JSON))
			Expect(err).ToNot(HaveOccurred())
			deleteFiles, ok = command.(*LogPlayback)
			Expect(ok).To(BeTrue())

			Expect(deleteFiles.Track).To(Equal("Brass Monkey"))
			Expect(deleteFiles.Action).To(BeNumerically("==", 2))

		})
	})

})
