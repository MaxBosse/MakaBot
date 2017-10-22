package command

import (
	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

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

func (t *Role) Event(c *Context, event *discordgo.Event) {
	return
}

func (t *Role) Message(c *Context) {
	log.Debugln(t.Name() + " called")
	if handleSubCommands(c, t) {
		return
	}

	serverConf, err := c.Cache.GetServer(c.Guild.ID)
	if err != nil {
		serverConf = cache.CacheServer{}
	}

	c.Send("Please use `" + serverConf.Prefix + t.Name() + " " + t.Usage() + "`\n" +
		"Use `" + serverConf.Prefix + "help " + t.Name() + "` for more info!")

}
