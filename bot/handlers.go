package bot

import (
	"strings"
	"time"

	"github.com/MaxBosse/MakaBot/bot/command"
	"github.com/MaxBosse/MakaBot/bot/utils"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

func (bot *MakaBot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Debugln("Bot ready")

	bot.CollectGlobalGuildMetrics(s)
	guildTicker := time.NewTicker(time.Second * 10)
	go func() {
		for range guildTicker.C {
			bot.CollectGlobalGuildMetrics(s)
		}
	}()
}

func (bot *MakaBot) guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	log.Debugln("Joining Guild " + event.Guild.Name)

	s.RequestGuildMembers(event.Guild.ID, "", 0)

	guild, _ := s.Guild(event.Guild.ID)
	bot.CollectGuildMetrics(s, guild)
	guildTicker := time.NewTicker(time.Second * 10)
	go func() {
		for range guildTicker.C {
			bot.CollectGuildMetrics(s, guild)
		}
	}()
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

	if _, ok := bot.discordServers[g.ID]; !ok {
		return
	}

	if !bot.discordServers[g.ID].Enabled {
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
		if cmd, ok := command.Commands[split[0]]; ok {
			cmd.Message(context)
			return
		}

		_, err = s.ChannelMessageSend(c.ID, "Unknown function.")
		if err != nil {
			log.Errorln(err)
		}
		return
	}
}

func (bot *MakaBot) guildMembersChunk(s *discordgo.Session, c *discordgo.GuildMembersChunk) {
	for _, g := range s.State.Guilds {
		if g.ID == c.GuildID {
			newm := append(g.Members, c.Members...)
			utils.RemoveDuplicateMembers(&newm)
			g.Members = newm
			break
		}
	}
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

// Only used for metric-collection!
func (bot *MakaBot) event(s *discordgo.Session, event *discordgo.Event) {
	log.Noteln("Event " + event.Type + " called.")
	bot.CollectGenericGlobalEventMetric(event.Type)

	switch t := event.Struct.(type) {
	case *discordgo.Ready:
		for _, g := range t.Guilds {
			bot.CollectGenericGuildEventMetric(s, g, event.Type)
		}
	case *discordgo.GuildCreate:
		bot.CollectGenericGuildEventMetric(s, t.Guild, event.Type)
	case *discordgo.GuildUpdate:
		bot.CollectGenericGuildEventMetric(s, t.Guild, event.Type)
	case *discordgo.VoiceServerUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.VoiceStateUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildDelete:
		bot.CollectGenericGuildEventMetric(s, t.Guild, event.Type)
	case *discordgo.GuildMemberAdd:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildMemberUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildMemberRemove:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildRoleCreate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildRoleUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildRoleDelete:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.GuildEmojisUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.ChannelCreate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.ChannelUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.ChannelDelete:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.MessageCreate:
		c, err := s.State.Channel(t.ChannelID)
		if err != nil {
			log.Warning("Unable to get Channel", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.MessageUpdate:
		c, err := s.State.Channel(t.ChannelID)
		if err != nil {
			log.Warning("Unable to get Channel", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.MessageDelete:
		c, err := s.State.Channel(t.ChannelID)
		if err != nil {
			log.Warning("Unable to get Channel", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.MessageDeleteBulk:
		c, err := s.State.Channel(t.ChannelID)
		if err != nil {
			log.Warning("Unable to get Channel", err)
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	case *discordgo.PresenceUpdate:
		g, err := s.State.Guild(t.GuildID)
		if err != nil {
			log.Warning("Unable to get Guild", err)
			return
		}
		bot.CollectGenericGuildEventMetric(s, g, event.Type)
	}
}
