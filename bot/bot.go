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
	regexUserID    *regexp.Regexp
	discordServers map[string]*structs.DiscordServer
	mem            runtime.MemStats
	tickers        map[string]*time.Ticker
}

func NewMakaBot(metrics *utils.InfluxDB, discordServers []*structs.DiscordServer, mem runtime.MemStats, discordToken string) *MakaBot {
	var err error

	bot := new(MakaBot)
	bot.iDB = metrics
	bot.mem = mem
	bot.discordServers = make(map[string]*structs.DiscordServer)
	bot.tickers = make(map[string]*time.Ticker)

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
	bot.tickers["globalTicker"] = time.NewTicker(time.Second * 10)
	go func() {
		for range bot.tickers["globalTicker"].C {
			bot.CollectGlobalMetrics()
		}
	}()

	bot.dg, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Errorln("Error creating Discord session:", err)
		return nil
	}

	if log.LogFlag <= log.DebugFlag {
		bot.dg.Debug = true
	}

	bot.dg.AddHandler(bot.ready)
	bot.dg.AddHandler(bot.messageCreate)
	bot.dg.AddHandler(bot.guildCreate)
	bot.dg.AddHandler(bot.memberAdd)
	bot.dg.AddHandler(bot.guildMembersChunk)
	bot.dg.AddHandler(bot.memberRemove)
	bot.dg.AddHandler(bot.memberUpdate)

	bot.dg.AddHandler(bot.event)

	err = bot.dg.Open()
	if err != nil {
		log.Errorln("Error opening Discord session:", err)
		return nil
	}

	return bot
}
