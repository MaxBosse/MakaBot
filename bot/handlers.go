package bot

import (
	"strings"
	"time"

	"github.com/MaxBosse/MakaBot/bot/command"
	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

func (bot *MakaBot) ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Debugln("Bot ready")

	bot.cache.UpdateSession(s)

	// Stop ticker if already running to start with new session
	if _, ok := bot.tickers["globalGuildTicker"]; ok {
		bot.tickers["globalGuildTicker"].Stop()
	}

	bot.CollectGlobalGuildMetrics(s)
	bot.tickers["globalGuildTicker"] = time.NewTicker(time.Second * 10)
	go func() {
		for range bot.tickers["globalGuildTicker"].C {
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

	// Check if Server already exists in DB otherwise create new one with default values
	guildKey := cache.CacheServerKey{
		GuildID: event.Guild.ID,
	}
	_, err := bot.cache.Get(guildKey)
	if err != nil {
		guildValue := cache.CacheServer{
			GuildID:  event.Guild.ID,
			Name:     event.Guild.Name,
			Nickname: "MakaBot",
			Prefix:   "!",
		}
		bot.cache.Set(guildKey, guildValue)
	}

	serverConf, err := bot.cache.GetServer(event.Guild.ID)
	if err != nil {
		log.Warningln("Unable to get server", err)
	}

	// Get all roles and put into cache/db if not already
	roles, err := s.GuildRoles(event.Guild.ID)
	if err != nil {
		log.Warningln("Unable to get guild roles", err)
	}
	for _, role := range roles {
		// Check if Server already exists in DB otherwise create new one with default values
		roleKey := cache.CacheRoleKey{
			GuildID: event.Guild.ID,
			RoleID:  role.ID,
		}
		_, err := bot.cache.Get(roleKey)
		if err != nil {
			roleValue := cache.CacheRole{
				SID:     serverConf.ID,
				GuildID: serverConf.GuildID,
				RoleID:  role.ID,
				Name:    role.Name,
			}
			bot.cache.Set(roleKey, roleValue)
		}
	}

	// Get all channels and put into cache/db if not already
	channels, err := s.GuildChannels(event.Guild.ID)
	if err != nil {
		log.Warningln("Unable to get guild channels", err)
	}
	for _, channel := range channels {
		// Check if Server already exists in DB otherwise create new one with default values
		channelKey := cache.CacheChannelKey{
			ChannelID: channel.ID,
		}
		_, err := bot.cache.Get(channelKey)
		if err != nil {
			channelValue := cache.CacheChannel{
				SID:       serverConf.ID,
				GuildID:   serverConf.GuildID,
				ChannelID: channel.ID,
				Name:      channel.Name,
				CType:     int(channel.Type),
			}
			bot.cache.Set(channelKey, channelValue)
		}
	}

	// METRIC COLLECTION
	// Stop ticker if already running to start with new session
	if _, ok := bot.tickers["guildTicker_"+event.Guild.ID]; ok {
		bot.tickers["guildTicker_"+event.Guild.ID].Stop()
	}

	guild, _ := s.Guild(event.Guild.ID)
	bot.CollectGuildMetrics(s, guild)
	bot.tickers["guildTicker_"+event.Guild.ID] = time.NewTicker(time.Second * 10)
	go func() {
		for range bot.tickers["guildTicker_"+event.Guild.ID].C {
			bot.CollectGuildMetrics(s, guild)
		}
	}()
	// END METRIC COLLECTION
}

func (bot *MakaBot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
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

	serverConf, err := bot.cache.GetServer(c.GuildID)
	if err != nil {
		serverConf = cache.CacheServer{}
	}

	channelConf, err := bot.cache.GetChannel(m.ChannelID)
	if err != nil {
		channelConf = cache.CacheChannel{}
	}

	if !serverConf.Enabled || !channelConf.Listen {
		return
	}

	// check if the message starts with our prefix
	if strings.HasPrefix(m.Content, serverConf.Prefix) {
		log.Notef("[%s.%s]: %s > %s", g.Name, c.Name, m.Author.Username, m.Content)

		context := new(command.Context)
		context.Session = s
		context.Guild = g
		context.Channel = c
		context.Message = m.Message
		context.Cache = bot.cache

		// Remove the prefix for the raw message
		context.RawText = m.Content[len(serverConf.Prefix):]

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
			RemoveDuplicateMembers(&newm)
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

func (bot *MakaBot) roleUpdate(s *discordgo.Session, event *discordgo.GuildRoleUpdate) {
	log.Noteln("Role", event.Role.Name, "updated.")
	roleKey := cache.CacheRoleKey{
		GuildID: event.GuildID,
		RoleID:  event.Role.ID,
	}
	roleConfI, err := bot.cache.Get(roleKey)
	if err != nil {
		return
	}
	roleConf := roleConfI.(cache.CacheRole)
	roleConf.Name = event.Role.Name
	bot.cache.Set(roleKey, roleConf)
}

func (bot *MakaBot) roleDelete(s *discordgo.Session, event *discordgo.GuildRoleDelete) {
	log.Noteln("Role", event.RoleID, "removed.")
	bot.cache.DeleteRole(event.GuildID, event.RoleID)
}

func (bot *MakaBot) channelUpdate(s *discordgo.Session, event *discordgo.ChannelUpdate) {
	log.Noteln("Channel", event.Name, "updated.")
	channelKey := cache.CacheChannelKey{
		ChannelID: event.ID,
	}
	channelConfI, err := bot.cache.Get(channelKey)
	if err != nil {
		return
	}
	channelConf := channelConfI.(cache.CacheChannel)
	channelConf.Name = event.Name
	channelConf.CType = int(event.Type)
	bot.cache.Set(channelKey, channelConf)
}

func (bot *MakaBot) channelDelete(s *discordgo.Session, event *discordgo.ChannelDelete) {
	log.Noteln("Channel", event.Name, "removed.")
	bot.cache.DeleteChannel(event.ID)
}

func (bot *MakaBot) guildUpdate(s *discordgo.Session, event *discordgo.GuildUpdate) {
	log.Noteln("Guild", event.Name, "updated.")

	if event.Guild.MemberCount == 0 {
		s.RequestGuildMembers(event.Guild.ID, "", 0)
	}

	guildKey := cache.CacheServerKey{
		GuildID: event.ID,
	}
	guildConfI, err := bot.cache.Get(guildKey)
	if err != nil {
		return
	}
	guildConf := guildConfI.(cache.CacheServer)
	guildConf.Name = event.Name
	bot.cache.Set(guildKey, guildConf)
}

// Only used for metric-collection!
func (bot *MakaBot) event(s *discordgo.Session, event *discordgo.Event) {
	log.Debugln("Event " + event.Type + " called.")
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
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.VoiceStateUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildDelete:
		bot.CollectGenericGuildEventMetric(s, t.Guild, event.Type)
	case *discordgo.GuildMemberAdd:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildMemberUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildMemberRemove:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildRoleCreate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildRoleUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildRoleDelete:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.GuildEmojisUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.ChannelCreate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.ChannelUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.ChannelDelete:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.MessageCreate:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	case *discordgo.MessageUpdate:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	case *discordgo.MessageDelete:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	case *discordgo.MessageDeleteBulk:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	case *discordgo.PresenceUpdate:
		bot.CollectGenericGuildEventMetricByGuildID(s, t.GuildID, event.Type)
	case *discordgo.TypingStart:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	case *discordgo.ChannelPinsUpdate:
		bot.CollectGenericGuildEventMetricByChannelID(s, t.ChannelID, event.Type)
	}

}
