package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type TestMessage struct{}

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

func (t *TestMessage) Message(c *Context) {
	log.Debugln(c.Invoked + t.Name() + " called")
	var err error

	_, err = c.Session.ChannelMessageSend(c.Channel.ID, strings.Join(c.Args, " "))
	if err != nil {
		log.Errorln(err)
	}
}
