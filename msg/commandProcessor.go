package msg

import (
	"errors"
	"fmt"
	"regexp"
)

var commandRe = regexp.MustCompile(`"command"\s*:\s*"([^"]+)"`)

type CommandParser func(*Msg) (Command, error)

type CommandProcessor struct {
	registered map[string]CommandParser
}

func NewCommandProcessor() *CommandProcessor {
	c := new(CommandProcessor)

	c.registered = make(map[string]CommandParser)

	c.registered["addfiles"] = func(m *Msg) (Command, error) { return NewAddFiles(m) }
	c.registered["logplayback"] = func(m *Msg) (Command, error) { return NewLogPlayBack(m) }

	return c
}

func (c *CommandProcessor) ParseCommand(m *Msg) (command Command, err error) {
	if m.IsCompressed() {
		m.Uncompress()
	}

	b := commandRe.FindSubmatch(m.Bytes())

	if len(b) != 2 {
		return nil, errors.New("Given Message is not a command")
	}

	commandName := string(b[1])
	if parseFunction := c.registered[commandName]; parseFunction != nil {
		return parseFunction(m)
	}

	return nil, fmt.Errorf("Not registered Command %s ", commandName)
}
