package command

import "github.com/MaxBosse/MakaBot/log"

type Role struct {
	subCommands map[string]Command
	parent      Command
}

func init() {
	cmd := new(Role)
	cmd.subCommands = make(map[string]Command)
	cmd.subCommands["add"] = new(RoleAdd)
	cmd.subCommands["del"] = new(RoleDel)
	cmd.subCommands["list"] = new(RoleList)
	cmd.subCommands["help"] = new(Help)

	Register(cmd)
}

func (t *Role) Name() string {
	return "role"
}

func (t *Role) Description() string {
	return "Allows Role-Management"
}

func (t *Role) Usage() string {
	return "[command]"
}

func (t *Role) SubCommands() map[string]Command {
	return t.subCommands
}

func (t *Role) Parent() Command {
	return t.parent
}

func (t *Role) SetParent(cmd Command) {
	t.parent = cmd
}

func (t *Role) Message(c *Context) {
	log.Debugln(t.Name() + " called")

	// Handle sub-commands
	if len(c.Args) != 0 {
		cmd, ok := t.subCommands[c.Args[0]]
		if ok {
			cmd.SetParent(t)
			c.Args = c.Args[1:]
			cmd.Message(c)

			return
		}
	}

	c.Send("Please use `" + c.Conf.Prefix + t.Name() + " " + t.Usage() + "`\n" +
		"Use `" + c.Conf.Prefix + "help " + t.Name() + "` for more info!")

}
