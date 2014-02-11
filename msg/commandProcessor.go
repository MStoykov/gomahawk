package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var commandRe = regexp.MustCompile(`"command"\s*:\s*"([^"]+)"`)

type commandParser func(*Msg) (Command, error)

type CommandProcessor struct {
	registered map[string]commandParser
}

type Generator func() Command

func unmarshalFuncFor(generator Generator) commandParser {
	return func(m *Msg) (Command, error) {
		command := generator()
		err := json.Unmarshal(m.Payload(), command)
		if err != nil {
			return nil, err
		}
		return command, nil
	}
}

func NewCommandProcessor() *CommandProcessor {
	c := new(CommandProcessor)

	c.registered = make(map[string]commandParser)

	c.Register("addfiles", func() Command {
		return new(AddFiles)
	})
	c.Register("deletefiles", func() Command {
		return new(DeleteFiles)
	})
	c.Register("logplayback", func() Command {
		return new(LogPlayback)
	})
	c.Register("socialaction", func() Command {
		return new(SocialAction)
	})
	c.Register("deleteplaylist", func() Command {
		return new(DeletePlaylist)
	})
	c.Register("createplaylist", func() Command {
		return new(CreatePlaylist)
	})
	c.Register("renameplaylist", func() Command {
		return new(RenamePlaylist)
	})
	c.Register("setplaylistrevision", func() Command {
		return new(SetPlaylistRevision)
	})
	c.Register("setcollectionattributes", func() Command {
		return new(SetCollectionAttributes)
	})
	return c
}

func (c *CommandProcessor) Register(commandName string, g Generator) {
	c.registered[commandName] = unmarshalFuncFor(g)
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
