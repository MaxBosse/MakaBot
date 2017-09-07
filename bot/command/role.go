package command

import "github.com/MaxBosse/MakaBot/log"

type Role struct {
	subCommands map[string]Command
}

func init() {
	cmd := new(Role)
	cmd.subCommands = make(map[string]Command)
	cmd.subCommands["add"] = new(RoleAdd)
	cmd.subCommands["del"] = new(RoleDel)
	cmd.subCommands["list"] = new(RoleList)

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

func (t *Role) Message(c *Context) {
	log.Debugln(c.Invoked + t.Name() + " called")

	// Handle sub-commands
	if len(c.Args) != 0 {
		cmd, ok := t.subCommands[c.Args[0]]
		if ok {
			c.Invoked += t.Name() + " "
			c.Args = c.Args[1:]
			cmd.Message(c)

			return
		}
	}

	c.Send("Please use `" + c.Conf.Prefix + t.Name() + " " + t.Usage() + "`\n" +
		"Use `" + c.Conf.Prefix + "help " + t.Name() + "` for more info!")

}
