package command

import (
	"fmt"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type Help struct {
	parent Command
}

func init() {
	Register(new(Help))
}

func (t *Help) Name() string {
	return "help"
}

func (t *Help) Description() string {
	return "Get the usage of a command"
}

func (t *Help) Usage() string {
	return ""
}

func (t *Help) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *Help) Parent() Command {
	return t.parent
}

func (t *Help) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *Help) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	var desc, invoked string

	if t.parent != nil {
		tmp := t.parent
		for {
			invoked += tmp.Name() + " "
			if tmp.Parent() != nil {
				tmp = tmp.Parent()
			} else {
				break
			}
		}
	}

	if t.parent != nil {
		desc = "Commands:\n"
		desc = fmt.Sprintf("`%s%s` - %s", c.Conf.Prefix+invoked, t.parent.Usage(), t.parent.Description())
		if len(t.parent.SubCommands()) != 0 {
			desc += "\n\nSubcommands:"

			for subCommandName, subCommand := range t.parent.SubCommands() {
				desc += fmt.Sprintf("\n`%s%s %s` - %s", c.Conf.Prefix+invoked, subCommandName, subCommand.Usage(), subCommand.Description())
			}
		}

		t.createEmbedMessage(c, desc)
		return
	}

	if len(Commands) != 0 {
		desc += "Commands:"

		for subCommandName, subCommand := range Commands {
			desc += fmt.Sprintf("\n`%s%s %s` - %s", c.Conf.Prefix+invoked, subCommandName, subCommand.Usage(), subCommand.Description())
		}
	}

	t.createEmbedMessage(c, desc)
}

func (t *Help) createEmbedMessage(c *Context, desc string) {
	embed := discordgo.MessageEmbed{}
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    c.Message.Author.Username,
		IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", c.Message.Author.ID, c.Message.Author.Avatar),
	}
	embed.Description = desc
	embed.Description += "\n\n"
	embed.Description += fmt.Sprintf("[MakaBot](https://github.com/MaxBosse/MakaBot) - use `%s[command] help` for more info!", c.Conf.Prefix)
	c.SendEmbed(&embed)
}
