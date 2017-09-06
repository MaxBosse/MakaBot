package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type RoleDel struct{}

func (t *RoleDel) Name() string {
	return "del"
}

func (t *RoleDel) Description() string {
	return "Delete a role to your user"
}

func (t *RoleDel) Usage() string {
	return "[role]"
}

func (t *RoleDel) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *RoleDel) Message(c *Context) {
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
		err = c.Session.GuildMemberRoleRemove(c.Guild.ID, c.Message.Author.ID, roleConfig.RoleID)
		if err != nil {
			_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Error removing role `"+role+"`")
			if err != nil {
				log.Errorln(err)
			}
			return
		}

		_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Removed role `"+role+"`")
		if err != nil {
			log.Errorln(err)
		}
		return
	}

	_, err = c.Session.ChannelMessageSend(c.Channel.ID, "Not allowed to remove role `"+role+"`")
	if err != nil {
		log.Errorln(err)
	}
	return
}
