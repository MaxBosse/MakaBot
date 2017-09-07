package command

import "github.com/MaxBosse/MakaBot/log"

type Test struct {
	subCommands map[string]Command
}

func init() {
	test := new(Test)
	test.subCommands = make(map[string]Command)
	test.subCommands["message"] = new(TestMessage)

	Register(test)
}

func (t *Test) Name() string {
	return "test"
}

func (t *Test) Description() string {
	return "Test-Description"
}

func (t *Test) Usage() string {
	return ""
}

func (t *Test) SubCommands() map[string]Command {
	return t.subCommands
}

func (t *Test) Message(c *Context) {
	log.Debugln(c.Invoked + t.Name() + " called")

	// Handle sub-commands
	if len(c.Args) != 0 {
		cmd, ok := t.subCommands[c.Args[0]]
		if ok {
			c.Invoked += t.Name() + " "
			c.Args = c.Args[1:]
			cmd.Message(c)

			return
		}
	}

	c.Send("This is a test message.")
}
