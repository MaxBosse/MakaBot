package bot

import (
	"regexp"
	"runtime"
	"time"

	"github.com/MaxBosse/MakaBot/bot/structs"
	"github.com/MaxBosse/MakaBot/bot/utils"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type MakaBot struct {
	dg             *discordgo.Session
	iDB            *utils.InfluxDB
	batchTicker    *time.Ticker
	regexUserID    *regexp.Regexp
	discordServers map[string]*structs.DiscordServer
	mem            runtime.MemStats
}

func NewMakaBot(metrics *utils.InfluxDB, discordServers []*structs.DiscordServer, mem runtime.MemStats, discordToken string) *MakaBot {
	var err error

	bot := new(MakaBot)
	bot.iDB = metrics
	bot.mem = mem
	bot.discordServers = make(map[string]*structs.DiscordServer)

	// Generate Roles-Maps
	for _, discordServer := range discordServers {
		tmpRoles := make(map[string]*structs.DiscordRole)

		for index, key := range discordServer.Roles {

			// Create a roles[1521512578] == DiscordRole matching
			if index != "" {
				key.RoleID = index
				tmpRoles[index] = key
			} else if key.RoleID != "" {
				tmpRoles[key.RoleID] = key
			} else {
				log.Fatalln("Could not determine RoleID -> Role", index, key)
			}

			// Create a roles["SysAdmin"] == DiscordRole matching
			if key.RoleName != "" {
				tmpRoles[key.RoleName] = key
			}
		}
		discordServer.Roles = tmpRoles

		bot.discordServers[discordServer.GuildID] = discordServer
	}

	// Collect memory statistics
	bot.CollectGlobalMetrics()
	bot.batchTicker = time.NewTicker(time.Second * 1)
	go func() {
		for range bot.batchTicker.C {
			bot.CollectGlobalMetrics()
		}
	}()

	bot.dg, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Errorln("Error creating Discord session:", err)
		return nil
	}

	bot.dg.AddHandler(bot.ready)
	bot.dg.AddHandler(bot.messageCreate)
	bot.dg.AddHandler(bot.guildCreate)
	bot.dg.AddHandler(bot.memberAdd)
	bot.dg.AddHandler(bot.guildMembersChunk)
	bot.dg.AddHandler(bot.memberRemove)
	bot.dg.AddHandler(bot.memberUpdate)

	err = bot.dg.Open()
	if err != nil {
		log.Errorln("Error opening Discord session:", err)
		return nil
	}

	return bot
}

// CollectGlobalMetrics collects global metrics about the bot and environment
// And sends them to influxdb
func (bot *MakaBot) CollectGlobalMetrics() {
	runtime.ReadMemStats(&bot.mem)
	tags := map[string]string{"metric": "server_metrics", "server": "global"}
	fields := map[string]interface{}{
		"memAlloc":      int(bot.mem.Alloc),
		"memTotalAlloc": int(bot.mem.TotalAlloc),
		"memHeapAlloc":  int(bot.mem.HeapAlloc),
		"memHeapSys":    int(bot.mem.HeapSys),
	}

	err := bot.iDB.AddMetric("server_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}

func (bot *MakaBot) CollectGuildMetrics(s *discordgo.Session, g *discordgo.Guild) {
	roles := make(map[string]int)
	rolesStruct := make(map[string]*discordgo.Role)

	for _, member := range g.Members {
		for _, role := range member.Roles {
			_, ok := rolesStruct[role]

			if !ok {
				dRole, err := s.State.Role(g.ID, role)
				if err != nil {
					log.Errorln("Could not get discord role")
					return
				}

				rolesStruct[role] = dRole
			}

			roles[rolesStruct[role].Name]++
		}
	}

	online := make(map[string]int)
	for _, presence := range g.Presences {
		online[string(presence.Status)]++
	}

	tags := map[string]string{"metric": "total_members", "server": g.Name}
	fields := map[string]interface{}{
		"totalMembers": g.MemberCount,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}

	for roleName := range roles {
		tags := map[string]string{"metric": "role_members", "server": g.Name, "roleName": roleName}
		fields := map[string]interface{}{
			"totalMembers": roles[roleName],
			//"onlineMembers": online["roles"][roleName],
		}

		err := bot.iDB.AddMetric("discord_metrics", tags, fields)
		if err != nil {
			log.Errorln("Error adding Metric:", err)
		}
	}

	for status := range online {
		tags := map[string]string{"metric": "status_members", "server": g.Name, "status": status}
		fields := map[string]interface{}{
			"onlineMembers": online[status],
		}

		err := bot.iDB.AddMetric("discord_metrics", tags, fields)
		if err != nil {
			log.Errorln("Error adding Metric:", err)
		}
	}
}
