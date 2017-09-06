package structs

import "time"

type DiscordServer struct {
	GuildID           string
	Nickname          string
	Prefix            string
	AutoDeleteSeconds int
	BotChannels       []string
	Roles             map[string]*DiscordRole
	Attributes        map[string]string
	metricsTickers    time.Ticker
}

type DiscordRole struct {
	RoleID     string
	RoleName   string
	Attributes map[string]string
}
