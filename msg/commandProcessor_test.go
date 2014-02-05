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
		processor       *CommandProcessor
		addfilesJSON, _ = ioutil.ReadFile("./fixtures/addfiles.json")
		addfilesMSG     = NewMsg(bytes.NewBuffer(addfilesJSON), JSON|DBOP)
	)

	BeforeEach(func() {
		processor = NewCommandProcessor()
	})

	Describe("Processor", func() {

		It("parses AddFiles", func() {
			command, err := processor.ParseCommand(addfilesMSG)
			Expect(err).ToNot(HaveOccurred())
			addFiles, ok := command.(*AddFiles)
			Expect(ok).To(BeTrue())

			id := addFiles.Files[0].Id
			Expect(id).To(BeNumerically("==", 1))

		})
	})

})
