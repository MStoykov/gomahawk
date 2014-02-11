package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var commandRe = regexp.MustCompile(`"command"\s*:\s*"([^"]+)"`)

type CommandAllocator func() Command

type CommandParser struct {
	registered map[string]CommandAllocator
}

func NewCommandParser() *CommandParser {
	c := new(CommandParser)

	c.registered = make(map[string]CommandAllocator)

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

func (c *CommandParser) Register(commandName string, allocator CommandAllocator) {
	c.registered[commandName] = allocator
}

func (c *CommandParser) ParseCommand(m *Msg) (command Command, err error) {
	m.Uncompress()

	b := commandRe.FindSubmatch(m.Bytes())

	if len(b) != 2 {
		return nil, errors.New("Given Message is not a command")
	}

	commandName := string(b[1])
	if allocator := c.registered[commandName]; allocator != nil {
		command := allocator()
		err := json.Unmarshal(m.Payload(), command)
		if err != nil {
			return nil, err
		}
		return command, nil
	}

	return nil, fmt.Errorf("Not registered Command %s ", commandName)
}
