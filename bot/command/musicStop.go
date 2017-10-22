package command

import (
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type MusicStop struct {
	parent Command
}

func (t *MusicStop) Name() string {
	return "stop"
}

func (t *MusicStop) Description() string {
	return "Stops a song"
}

func (t *MusicStop) Usage() string {
	return ""
}

func (t *MusicStop) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *MusicStop) Parent() Command {
	return t.parent
}

func (t *MusicStop) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *MusicStop) Event(c *Context, event *discordgo.Event) {
	return
}

func (t *MusicStop) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range c.Guild.VoiceStates {
		if vs.UserID == c.Message.Author.ID {

			music := t.parent.(*Music)
			music.MusicRuntimes[c.Guild.ID].Stop <- true

			c.Send("Stopped playing.")
			return
		}
	}

	c.Send("Error playing music")
	return
}
