package msg_test

import (
	"io"
	"io/ioutil"
	"os"
)

func getReaderToFixtureSafe(name string) (io.Reader, error) {
	return os.Open("./fixtures/" + name)
}

func getReaderToFixture(name string) io.Reader {
	reader, err := getReaderToFixtureSafe(name)
	if err != nil {
		panic(err)
	}
	return reader
}

func getFixtureSafe(name string) ([]byte, error) {
	return ioutil.ReadFile("./fixtures/" + name)
}

func getFixture(name string) []byte {
	b, err := getFixtureSafe(name)
	if err != nil {
		panic(err)
	}
	return b
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
