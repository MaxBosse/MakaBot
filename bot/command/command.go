package command

import (
	"time"

	"github.com/MaxBosse/MakaBot/cache"
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
	Cache   *cache.Cache
}

func (c *Context) SendEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	m, err := c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed)

	channelConf, err := c.Cache.GetChannel(c.Message.ChannelID)
	if err != nil {
		channelConf = cache.CacheChannel{}
	}
	if channelConf.AutoDelete != 0 {
		go waitandDelete(c, m)
	}

	return m, err
}

func (c *Context) Send(message string) (*discordgo.Message, error) {
	m, err := c.Session.ChannelMessageSend(c.Message.ChannelID, message)

	channelConf, err := c.Cache.GetChannel(c.Message.ChannelID)
	if err != nil {
		channelConf = cache.CacheChannel{}
	}

	if channelConf.AutoDelete != 0 {
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
	Event(*Context, *discordgo.Event)
}

func init() {
	Commands = make(map[string]Command)
}

func Register(command Command) {
	Commands[command.Name()] = command
}

func handleSubCommands(c *Context, command Command) bool {
	// Handle sub-commands
	if len(c.Args) != 0 {
		cmd, ok := command.SubCommands()[c.Args[0]]
		if ok {
			cmd.SetParent(command)
			c.Args = c.Args[1:]
			cmd.Message(c)

			return true
		}
	}
	return false
}

func waitandDelete(c *Context, m *discordgo.Message) {
	channelConf, err := c.Cache.GetChannel(m.ChannelID)
	if err != nil {
		channelConf = cache.CacheChannel{}
	}

	time.Sleep(time.Second * time.Duration(channelConf.AutoDelete))
	c.Session.ChannelMessageDelete(m.ChannelID, m.ID)
	c.Session.ChannelMessageDelete(c.Message.ChannelID, c.Message.ID)
}
