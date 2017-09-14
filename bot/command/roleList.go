package command

import "github.com/MaxBosse/MakaBot/log"
import "strings"

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

	roles := []string{}
	rolesInsert := make(map[string]bool)
	for _, roleConf := range c.Conf.Roles {
		if selfAssign, ok := roleConf.Attributes["selfAssign"]; ok && selfAssign == "true" {

			// Only add a role-name once
			if !rolesInsert[roleConf.RoleName] {
				roles = append(roles, roleConf.RoleName)
				rolesInsert[roleConf.RoleName] = true
			}
		}
	}

	desc = "Available roles: " + strings.Join(roles, ", ")
	c.Send(desc)
}
