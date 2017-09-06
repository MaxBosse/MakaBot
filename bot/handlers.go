package bot

import (
	"strings"

	"github.com/MaxBosse/MakaBot/bot/command"
	"github.com/MaxBosse/MakaBot/bot/utils"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

func (bot *MakaBot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Debugln("Bot ready")
}

func (bot *MakaBot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
}

func (bot *MakaBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// check if the message starts with our prefix
	if strings.HasPrefix(m.Content, bot.discordServers[g.ID].Prefix) {
		log.Notef("[%s.%s]: %s > %s", g.Name, c.Name, m.Author.Username, m.Content)

		if !utils.StringInSlice(c.Name, bot.discordServers[g.ID].BotChannels) {
			log.Debugln("Channel not in whitelist", c.Name, bot.discordServers[g.ID].BotChannels)
			return
		}

		context := new(command.Context)
		context.Session = s
		context.Guild = g
		context.Channel = c
		context.Message = m.Message
		context.Conf = bot.discordServers[g.ID]

		// Remove the prefix for the raw message
		context.RawText = m.Content[len(bot.discordServers[g.ID].Prefix):]

		split := strings.Split(context.RawText, " ")
		if len(split) > 1 {
			context.Args = split[1:]
		}

		// Execute
		command.Commands[split[0]].Message(context)
		return
	}
}

func (bot *MakaBot) guildMembersChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
}

func (bot *MakaBot) memberAdd(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	log.Noteln("User", event.User.Username, "joined.")
}

func (bot *MakaBot) memberRemove(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	log.Noteln("User", event.User.Username, "removed.")
}

func (bot *MakaBot) memberUpdate(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	log.Noteln("User", event.User.Username, "updated.")
}
