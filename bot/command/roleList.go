package command

import (
	"strings"

	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
)

type RoleList struct {
	parent Command
}

func (t *RoleList) Name() string {
	return "list"
}

func (t *RoleList) Description() string {
	return "List all available roles"
}

func (t *RoleList) Usage() string {
	return ""
}

func (t *RoleList) SubCommands() map[string]Command {
	return make(map[string]Command)
}

func (t *RoleList) Parent() Command {
	return t.parent
}

func (t *RoleList) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *RoleList) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	var desc string

	rolesConfig, err := c.Cache.GetRoles(c.Guild.ID)
	log.Debugln(rolesConfig)
	if err != nil {
		log.Noteln("Unable to get roles", err)
		rolesConfig = []cache.CacheRole{}
	}

	roles := []string{}
	rolesInsert := make(map[string]bool)
	for _, roleConf := range rolesConfig {
		if roleConf.SelfAssign {
			// Only add a role-name once
			if !rolesInsert[roleConf.Name] {
				roles = append(roles, roleConf.Name)
				rolesInsert[roleConf.Name] = true
			}
		}
	}
	if len(roles) > 0 {
		desc = "Available roles: " + strings.Join(roles, ", ")
		c.Send(desc)
	}
}
