package msg

import (
	"regexp"
)

type CommandFilter struct {
}

var re = regexp.MustCompile(`"command"\s*:\s*"([^"]+)"`)

// Reads msg.Msg from c and if it's Command calls the appropriate method on FetchOpsMethod
// if Non Command msg is met it gets send to others
// FetchOpsMethod and others are closed when the c is closed
// this functions blocks and returns any error encountered. On error FetchOpsMethod and otehrs may not be closed.
func FilterCommands(c <-chan *Msg, others chan<- *Msg, f FetchOpsMethod) error {
	defer close(others)
	defer f.Close()

	for m := range c {
		if m.IsCompressed() {
			m.Uncompress()
		}

		b := re.FindSubmatch(m.Bytes())

		if len(b) != 2 {
			others <- m
			continue
		}
		command := string(b[1])

		switch command {
		case "addfiles":
			addFiles, err := NewAddFiles(m)
			if err != nil {
				return err
			}

			f.SendCommand(addFiles)
		case "deletefiles":
			deleteFiles, err := NewDeleteFiles(m)
			if err != nil {
				return err
			}

			f.SendCommand(deleteFiles)
		case "socialaction":
			socialAction, err := NewSocialAction(m)
			if err != nil {
				return err
			}

			f.SendCommand(socialAction)

		case "logplayback":
			logPlayback, err := NewLogPlayBack(m)
			if err != nil {
				return err
			}

			f.SendCommand(logPlayback)

		default:
			//log.Println(m)
		}
		if !m.IsFragment() {
			break
		}
	}

	return nil
}
