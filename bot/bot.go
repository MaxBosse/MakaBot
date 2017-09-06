package bot

import (
	"regexp"
	"runtime"
	"time"

	"github.com/HeroesAwaken/GoAwaken/core"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/bwmarrin/discordgo"
)

type MakaBot struct {
	dg             *discordgo.Session
	iDB            *core.InfluxDB
	batchTicker    *time.Ticker
	regexUserID    *regexp.Regexp
	discordServers map[string]*DiscordServer
	mem            runtime.MemStats
}

type DiscordServer struct {
	GuildID           string
	Nickname          string
	Prefix            string
	AutoDeleteSeconds int
	BotChannels       []string
	Roles             map[string]*DiscordRole
	Attributes        map[string]string
	metricsTickers    time.Ticker
}

type DiscordRole struct {
	RoleID     string
	RoleName   string
	Attributes map[string]string
}

func NewMakaBot(metrics *core.InfluxDB, discordServers []*DiscordServer, mem runtime.MemStats, discordToken string) *MakaBot {
	var err error

	bot := new(MakaBot)
	bot.iDB = metrics
	bot.mem = mem
	bot.discordServers = make(map[string]*DiscordServer)

	// Generate Roles-Maps
	for _, discordServer := range discordServers {
		tmpRoles := make(map[string]*DiscordRole)

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
	bot.batchTicker = time.NewTicker(time.Second * 10)
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
