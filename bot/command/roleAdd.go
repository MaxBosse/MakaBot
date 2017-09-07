package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type RoleAdd struct {
	parent Command
}

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

func (t *RoleAdd) Parent() Command {
	return t.parent
}

func (t *RoleAdd) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *RoleAdd) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	var err error

	role := strings.Join(c.Args, " ")
	roleConfig, ok := c.Conf.Roles[role]

	if !ok {
		c.Send("Unknown role `" + role + "`")
		return
	}

	if selfAssign, ok := roleConfig.Attributes["selfAssign"]; ok && selfAssign == "true" {
		err = c.Session.GuildMemberRoleAdd(c.Guild.ID, c.Message.Author.ID, roleConfig.RoleID)
		if err != nil {
			log.Errorln(err)
			c.Send("Error adding role `" + role + "`")
			return
		}

		c.Send("Added role `" + role + "`")
		return
	}

	c.Send("Not allowed to add role `" + role + "`")
	return
}
