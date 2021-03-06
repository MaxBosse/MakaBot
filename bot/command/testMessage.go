package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type TestMessage struct {
	parent Command
}

func (t *TestMessage) Name() string {
	return "message"
}

func (t *TestMessage) Description() string {
	return "Make the bot say a specific message"
}

func (t *TestMessage) Usage() string {
	return "[text]"
}

func (t *TestMessage) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *TestMessage) Parent() Command {
	return t.parent
}

func (t *TestMessage) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *TestMessage) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	c.Send(strings.Join(c.Args, " "))
}
