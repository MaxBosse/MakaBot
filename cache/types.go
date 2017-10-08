package cache

import (
	"github.com/bwmarrin/discordgo"
)

type CacheServerKey struct {
	GuildID string
}

type CacheServer struct {
	ID       string
	GuildID  string
	Enabled  bool
	Nickname string
	Prefix   string
	Name     string
	Guild    *discordgo.Guild
}

type CacheRoleKey struct {
	GuildID string
	RoleID  string
}

type CacheRoleName struct {
	GuildID  string
	RoleName string
}

type CacheRoles struct {
	GuildID string
}

type CacheRole struct {
	ID         string
	SID        string
	GuildID    string
	RoleID     string
	SelfAssign bool
	Name       string
	Role       *discordgo.Role
}

type CacheChannelKey struct {
	ChannelID string
}

type CacheChannelName struct {
	GuildID     string
	ChannelName string
}

type CacheChannel struct {
	ID         string
	SID        string
	GuildID    string
	ChannelID  string
	Listen     bool
	Name       string
	AutoDelete uint
	CType      int
	Channel    *discordgo.Channel
}

type CacheMembers struct {
	GuildID string
}

type CacheMemberKey struct {
	UserID string
}

type CacheMemberGuildKey struct {
	GuildID string
	UserID  string
}

type CacheMember struct {
	ID            string
	SID           string
	GuildID       string
	UserID        string
	Username      string
	Discriminator string
	Avatar        string
	Nick          string
	Member        *discordgo.Member
}
