package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type RoleAdd struct{}

func (t *RoleAdd) Name() string {
	return "add"
}

func (t *RoleAdd) Description() string {
	return "Add a role to your user"
}

func (t *RoleAdd) Usage() string {
	return "[role]"
}

func (t *RoleAdd) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *RoleAdd) Message(c *Context) {
	log.Debugln(c.Invoked + t.Name() + " called")
	var err error

	role := strings.Join(c.Args, " ")
	roleConfig, ok := c.Conf.Roles[role]

	if !ok {
		_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Unknown role `"+role+"`")
		if err != nil {
			log.Errorln(err)
		}
		return
	}

	if selfAssign, ok := roleConfig.Attributes["selfAssign"]; ok && selfAssign == "true" {
		err = c.Session.GuildMemberRoleAdd(c.Guild.ID, c.Message.Author.ID, roleConfig.RoleID)
		if err != nil {
			log.Errorln(err)
			_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Error adding role `"+role+"`")
			if err != nil {
				log.Errorln(err)
			}
			return
		}

		_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Added role `"+role+"`")
		if err != nil {
			log.Errorln(err)
		}
		return
	}

	_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Not allowed to add role `"+role+"`")
	if err != nil {
		log.Errorln(err)
	}
	return
}
