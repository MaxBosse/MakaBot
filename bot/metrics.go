package bot

import (
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

func (bot *MakaBot) CollectGlobalGuildMetrics(s *discordgo.Session) {
	roles := make(map[string]int)
	rolesStruct := make(map[string]*discordgo.Role)
	online := make(map[string]int)
	totalMembers := 0

	for _, g := range s.State.Guilds {
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

		totalMembers += g.MemberCount
	}

	tags := map[string]string{"metric": "total_members", "server": "global", "serverID": "-1"}
	fields := map[string]interface{}{
		"totalMembers": totalMembers,
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

func (bot *MakaBot) CollectGenericGlobalEventMetric(event *discordgo.Event) {
	tags := map[string]string{"event": event.Type, "server": "global", "serverID": "-1"}
	fields := map[string]interface{}{
		"value": 1,
		"raw":   event.RawData,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}

func (bot *MakaBot) CollectGenericGuildEventMetricByChannelID(s *discordgo.Session, cID string, event *discordgo.Event) {
	c, err := s.State.Channel(cID)
	if err != nil {
		log.Warning("Unable to get Channel", err)
		return
	}

	bot.CollectGenericGuildEventMetricByGuildID(s, c.GuildID, event)
}

func (bot *MakaBot) CollectGenericGuildEventMetricByGuildID(s *discordgo.Session, gID string, event *discordgo.Event) {
	g, err := s.State.Guild(gID)
	if err != nil {
		log.Warning("Unable to get Guild", err)
		return
	}

	bot.CollectGenericGuildEventMetric(s, g, event)
}

func (bot *MakaBot) CollectGenericGuildEventMetric(s *discordgo.Session, g *discordgo.Guild, event *discordgo.Event) {
	tags := map[string]string{"event": event.Type, "server": g.Name, "serverID": g.ID}
	fields := map[string]interface{}{
		"value": 1,
		"raw":   event.RawData,
	}

	err := bot.iDB.AddMetric("discord_metrics", tags, fields)
	if err != nil {
		log.Errorln("Error adding Metric:", err)
	}
}
