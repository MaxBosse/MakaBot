package command

import (
	"strings"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
	"github.com/rylio/ytdl"
)

type MusicPlay struct {
	parent Command
}

func (t *MusicPlay) Name() string {
	return "add"
}

func (t *MusicPlay) Description() string {
	return "Plays a song"
}

func (t *MusicPlay) Usage() string {
	return "[song]"
}

func (t *MusicPlay) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *MusicPlay) Parent() Command {
	return t.parent
}

func (t *MusicPlay) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *MusicPlay) Event(c *Context, event *discordgo.Event) {
	return
}

func (t *MusicPlay) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	var err error
	videoText := strings.Join(c.Args, " ")

	if !strings.HasPrefix(videoText, "https://www.youtube.com/watch?v") {
		c.Send("Please make sure the URL starts with `https://www.youtube.com/watch?v`")
		return
	}

	vid, err := ytdl.GetVideoInfo(videoText)
	if err != nil {
		log.Errorln(err)
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range c.Guild.VoiceStates {
		if vs.UserID == c.Message.Author.ID {
			musicChan := &musicChan{
				c.Guild.ID,
				vs.ChannelID,
				vid,
			}

			music := t.parent.(*Music)
			music.MusicGuilds[c.Guild.ID] <- musicChan
			c.Send("Added " + vid.Title + "to music queue.")
			return
		}
	}

	c.Send("Error playing music")
	return
}
