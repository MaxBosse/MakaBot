package command

import (
	"fmt"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type Help struct{}

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
	return "[command]"
}

func (t *Help) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *Help) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	var desc string

	if len(c.Args) != 0 {
		command := c.Args[0]

		cmd, ok := Commands[command]
		if !ok {
			t.createEmbedMessage(c, "No command called `"+command+"` found.")
			return
		}

		desc = fmt.Sprintf("`%s%s %s` - %s", c.Conf.Prefix, command, cmd.Usage(), cmd.Description())

		if len(cmd.SubCommands()) != 0 {
			desc += "\n\nSubcommands:"

			for subCommandName, subCommand := range cmd.SubCommands() {
				desc += fmt.Sprintf("\n`%s%s %s` - %s", c.Conf.Prefix, command, subCommandName, subCommand.Description())
			}
		}

		t.createEmbedMessage(c, desc)
		return
	}

	desc = "Commands:"
	desc += fmt.Sprintf(" `%shelp [command]` for more info!", c.Conf.Prefix)
	if len(Commands) != 0 {
		desc += "\n\nSubcommands:"

		for subCommandName, subCommand := range Commands {
			desc += fmt.Sprintf("\n`%s%s` - %s", c.Conf.Prefix, subCommandName, subCommand.Description())
		}
	}

	t.createEmbedMessage(c, desc)
}

func (t *Help) createEmbedMessage(c *Context, desc string) {
	log.Debugln(c.Message.Author.ID)
	log.Debugln(c.Message.Author.Avatar)
	embed := discordgo.MessageEmbed{}
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    c.Message.Author.Username,
		IconURL: fmt.Sprintf("https://discordapp.com/api/users/%s/avatars/%s.jpg", c.Message.Author.ID, c.Message.Author.Avatar),
	}
	embed.Description = desc
	embed.Description += "\n\n"
	embed.Description += "[MakaBot](https://github.com/MaxBosse/MakaBot)"
	c.SendEmbed(&embed)
}
