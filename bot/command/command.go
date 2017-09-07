package command

import (
	"time"

	"github.com/MaxBosse/MakaBot/bot/structs"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands map[string]Command
)

type Context struct {
	Session *discordgo.Session
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	Message *discordgo.Message
	RawText string
	Args    []string
	Conf    *structs.DiscordServer
}

func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	m, err := c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)

	if c.Conf.AutoDeleteSeconds != 0 {
		go waitandDelete(c, m)
	}

	return m, err
}

func (c *Context) Send(message string) (*discordgo.Message, error) {
	m, err := c.Session.ChannelMessageSend(c.Message.ChannelID, message)

	if c.Conf.AutoDeleteSeconds != 0 {
		go waitandDelete(c, m)
	}

	return m, err
}

type Command interface {
	Message(*Context)
	Description() string
	Usage() string
	Name() string
	SubCommands() map[string]Command
	Parent() Command
	SetParent(Command)
}

func init() {
	Commands = make(map[string]Command)
}

func Register(command Command) {
	Commands[command.Name()] = command
}

func waitandDelete(c *Context, m *discordgo.Message) {
	time.Sleep(time.Second * time.Duration(c.Conf.AutoDeleteSeconds))
	c.Session.ChannelMessageDelete(m.ChannelID, m.ID)
	c.Session.ChannelMessageDelete(c.Message.ChannelID, c.Message.ID)
}
