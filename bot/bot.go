package bot

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/MaxBosse/MakaBot/cache"
	"github.com/MaxBosse/MakaBot/log"
	"github.com/MaxBosse/MakaBot/utils"
	"github.com/bwmarrin/discordgo"

	_ "github.com/go-sql-driver/mysql"
)

type MakaBot struct {
	dg          *discordgo.Session
	iDB         *utils.InfluxDB
	db          *sql.DB
	regexUserID *regexp.Regexp
	cache       *cache.Cache
	tickers     map[string]*time.Ticker
}

func NewMakaBot(discordToken string, metrics *utils.InfluxDB, db *sql.DB, cache *cache.Cache) *MakaBot {
	var err error

	bot := new(MakaBot)
	bot.iDB = metrics
	bot.tickers = make(map[string]*time.Ticker)
	bot.db = db
	bot.cache = cache

	bot.dg, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Errorln("Error creating Discord session:", err)
		return nil
	}

	/*if log.LogFlag <= log.DebugFlag {
		bot.dg.Debug = true
	}*/

	bot.dg.AddHandler(bot.ready)
	bot.dg.AddHandler(bot.messageCreate)
	bot.dg.AddHandler(bot.guildCreate)
	bot.dg.AddHandler(bot.memberAdd)
	bot.dg.AddHandler(bot.guildMembersChunk)
	bot.dg.AddHandler(bot.memberRemove)
	bot.dg.AddHandler(bot.memberUpdate)
	bot.dg.AddHandler(bot.roleUpdate)
	bot.dg.AddHandler(bot.roleDelete)
	bot.dg.AddHandler(bot.channelUpdate)
	bot.dg.AddHandler(bot.channelDelete)
	bot.dg.AddHandler(bot.guildUpdate)

	bot.dg.AddHandler(bot.event)

	err = bot.dg.Open()
	if err != nil {
		log.Errorln("Error opening Discord session:", err)
		return nil
	}

	return bot
}
