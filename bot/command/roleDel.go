package command

import (
	"strings"

	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

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

func (t *RoleDel) Event(c *Context, event *discordgo.Event) {
	return
}

func (t *RoleDel) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	var err error

	role := strings.Join(c.Args, " ")
	roleConfig, err := c.Cache.GetRoleByName(c.Guild.ID, role)
	if err != nil {
		roleConfig = cache.CacheRole{}
	}

	if roleConfig.SelfAssign {
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
