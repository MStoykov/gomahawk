package msg_test

import (
	"encoding/json"
	"io/ioutil"

	. "github.com/MStoykov/gomahawk/msg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getFixtureSafe(name string) ([]byte, error) {
	return ioutil.ReadFile("./fixtures/" + name)
}

func getJSONFixture(name string) []byte {
	b, err := getJSONFixtureSafe(name)
	if err != nil {
		panic(err)
	}
	return b
}

func getJSONFixtureSafe(name string) ([]byte, error) {
	return getFixtureSafe(name + ".json")
}

var _ = Describe("CommandParser", func() {
	var (
		msgWrapper = func(b []byte) *Msg {
			return NewMsg(b, JSON|DBOP)
		}
		processor *CommandParser
		command   Command
		err       error
	)

	BeforeEach(func() {
		processor = NewCommandParser()
	})

	var CommonAssertions = func(jsonB []byte) {
		BeforeEach(func() {
			command, err = processor.ParseCommand(msgWrapper(jsonB))
		})

		It("parses it", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be marshalled to an equal json", func() {
			j, err := json.Marshal(command)
			Expect(err).ToNot(HaveOccurred())
			Expect(j).To(MatchJSON(jsonB))
		})

	}

	var testMatrix = []struct {
		name  string
		jsonB []byte
	}{
		{"AddFiles", getJSONFixture("addfiles")},
		{"DeleteFiles", getJSONFixture("deletefiles")},
		{"LogPlayback for start", getJSONFixture("logplayback")},
		{"LogPlayback for end", getJSONFixture("logplayback2")},
		{"CreatePlaylist", getJSONFixture("createplaylist")},
		{"DeletePlaylist", getJSONFixture("deleteplaylist")},
		{"RenamePlaylist", getJSONFixture("renameplaylist")},
		{"SetPlaylistRevision for new playlist", getJSONFixture("setplaylistrevision")},
		{"SetPlaylistRevision for not new playlist", getJSONFixture("setplaylistrevision2")},
		{"SetCollectionAttributes", getJSONFixture("setcollectionattributes")},
		{"SocialAction love", getJSONFixture("socialaction")},
		{"SocialAction unlove", getJSONFixture("socialaction2")},
	}

	for _, test := range testMatrix {
		Context("Given"+test.name, func() {
			CommonAssertions(test.jsonB)
		})
	}
})
