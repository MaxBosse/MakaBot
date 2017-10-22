package command

import (
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type MusicSkip struct {
	parent Command
}

func (t *MusicSkip) Name() string {
	return "skip"
}

func (t *MusicSkip) Description() string {
	return "Skips a song"
}

func (t *MusicSkip) Usage() string {
	return ""
}

func (t *MusicSkip) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *MusicSkip) Parent() Command {
	return t.parent
}

func (t *MusicSkip) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *MusicSkip) Event(c *Context, event *discordgo.Event) {
	return
}

func (t *MusicSkip) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range c.Guild.VoiceStates {
		if vs.UserID == c.Message.Author.ID {

			music := t.parent.(*Music)
			music.MusicRuntimes[vs.ChannelID].Skip <- true

			c.Send("Skipped the song.")
			return
		}
	}

	c.Send("Error playing music")
	return
}
