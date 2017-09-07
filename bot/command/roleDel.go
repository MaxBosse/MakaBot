package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

type RoleDel struct {
	parent Command
}

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

func (t *RoleDel) Parent() Command {
	return t.parent
}

func (t *RoleDel) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *RoleDel) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	var err error

	role := strings.Join(c.Args, " ")
	roleConfig, ok := c.Conf.Roles[role]

	if !ok {
		c.Send("Unknown role `" + role + "`")
		return
	}

	if selfAssign, ok := roleConfig.Attributes["selfAssign"]; ok && selfAssign == "true" {
		err = c.Session.GuildMemberRoleRemove(c.Guild.ID, c.Message.Author.ID, roleConfig.RoleID)
		if err != nil {
			c.Send("Error removing role `" + role + "`")
			return
		}

		c.Send("Removed role `" + role + "`")
		return
	}

	c.Send("Not allowed to remove role `" + role + "`")
	return
}
