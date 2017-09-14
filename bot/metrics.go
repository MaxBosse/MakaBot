package bot

import (
	"runtime"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

// CollectGlobalMetrics collects global metrics about the bot and environment
// And sends them to influxdb
func (bot *MakaBot) CollectGlobalMetrics() {
	runtime.ReadMemStats(&bot.mem)
	tags := map[string]string{"metric": "server_metrics", "server": "global", "serverID": "-1"}
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

func (bot *MakaBot) CollectGlobalGuildMetrics(s *discordgo.Session, g *discordgo.Guild) {
	roles := make(map[string]int)
	rolesStruct := make(map[string]*discordgo.Role)
	online := make(map[string]int)

	guilds, _ := s.UserGuilds(0, "", "")

	for _, uG := range guilds {
		g, _ := s.State.Guild(uG.ID)

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

		for _, presence := range g.Presences {
			online[string(presence.Status)]++
		}
	}

	tags := map[string]string{"metric": "total_members", "server": "global", "serverID": "-1"}
	fields := map[string]interface{}{
		"totalMembers": g.MemberCount,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}

	for roleName := range roles {
		tags := map[string]string{"metric": "role_members", "server": "global", "serverID": "-1", "roleName": roleName}
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
		tags := map[string]string{"metric": "status_members", "server": "global", "serverID": "-1", "status": status}
		fields := map[string]interface{}{
			"onlineMembers": online[status],
		}

		err := bot.iDB.AddMetric("discord_metrics", tags, fields)
		if err != nil {
			log.Errorln("Error adding Metric:", err)
		}
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

	tags := map[string]string{"metric": "total_members", "server": g.Name, "serverID": g.ID}
	fields := map[string]interface{}{
		"totalMembers": g.MemberCount,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}

	for roleName := range roles {
		tags := map[string]string{"metric": "role_members", "server": g.Name, "serverID": g.ID, "roleName": roleName}
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
		tags := map[string]string{"metric": "status_members", "server": g.Name, "serverID": g.ID, "status": status}
		fields := map[string]interface{}{
			"onlineMembers": online[status],
		}

		err := bot.iDB.AddMetric("discord_metrics", tags, fields)
		if err != nil {
			log.Errorln("Error adding Metric:", err)
		}
	}
}

func (bot *MakaBot) CollectGenericGlobalEventMetric(event string) {
	tags := map[string]string{"event": event, "server": "global", "serverID": "-1"}
	fields := map[string]interface{}{
		"value": 1,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}

func (bot *MakaBot) CollectGenericGuildEventMetric(s *discordgo.Session, g *discordgo.Guild, event string) {
	tags := map[string]string{"event": event, "server": g.Name, "serverID": g.ID}
	fields := map[string]interface{}{
		"value": 1,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}
