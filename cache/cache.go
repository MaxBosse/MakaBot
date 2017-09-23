package cache

import (
	"database/sql"

	"github.com/MaxBosse/MakaBot/log"
	"github.com/MaxBosse/MakaBot/utils"
	"github.com/bluele/gcache"
	"github.com/bwmarrin/discordgo"

	_ "github.com/go-sql-driver/mysql"
)

type Cache struct {
	gcache     *gcache.Cache
	db         *sql.DB
	iDB        *utils.InfluxDB
	session    *discordgo.Session
	cacheStmts *cacheStmts
}

func NewCache(metrics *utils.InfluxDB, db *sql.DB) *Cache {
	cache := new(Cache)
	cache.iDB = metrics
	cache.db = db
	cache.cacheStmts = cache.prepareStatements()

	tmpCache := gcache.New(500).
		LRU().
		LoaderFunc(cache.loader).
		Build()

	cache.gcache = &tmpCache

	return cache
}

func (cache *Cache) UpdateSession(s *discordgo.Session) {
	cache.session = s
}

func (cache *Cache) GetRoles(gID string) ([]CacheRole, error) {
	value := []CacheRole{}
	key := CacheRoles{
		GuildID: gID,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.([]CacheRole)
	}

	return value, err
}

func (cache *Cache) GetRole(gID, rID string) (CacheRole, error) {
	value := CacheRole{}
	key := CacheRoleKey{
		GuildID: gID,
		RoleID:  rID,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.(CacheRole)
	}

	return value, err
}

func (cache *Cache) GetRoleByName(gID, rName string) (CacheRole, error) {
	value := CacheRole{}
	key := CacheRoleName{
		GuildID:  gID,
		RoleName: rName,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.(CacheRole)
	}

	return value, err
}

func (cache *Cache) GetChannel(cID string) (CacheChannel, error) {
	value := CacheChannel{}
	key := CacheChannelKey{
		ChannelID: cID,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.(CacheChannel)
	}

	return value, err
}

func (cache *Cache) GetChannelByName(gID, cName string) (CacheChannel, error) {
	value := CacheChannel{}
	key := CacheChannelName{
		GuildID:     gID,
		ChannelName: cName,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.(CacheChannel)
	}

	return value, err
}

func (cache *Cache) GetServer(gID string) (CacheServer, error) {
	value := CacheServer{}
	key := CacheServerKey{
		GuildID: gID,
	}
	valueI, err := cache.Get(key)
	if err == nil {
		value = valueI.(CacheServer)
	}

	return value, err
}

func (cache *Cache) Get(key interface{}) (interface{}, error) {
	return (*cache.gcache).Get(key)
}

func (cache *Cache) Set(key interface{}, value interface{}) error {
	var err error

	switch t := value.(type) {
	case CacheServer:
		_, err = cache.cacheStmts.setServer.Exec(t.GuildID, t.Enabled, t.Nickname, t.Prefix, t.Name)
		if err != nil {
			log.Warningln("Unable to set server", err)
		}
	case CacheChannel:
		_, err = cache.cacheStmts.setChannel.Exec(t.SID, t.ChannelID, t.Listen, t.Name, t.AutoDelete, t.CType)
		if err != nil {
			log.Warningln("Unable to set channel", err)
		}
	case CacheRole:
		_, err = cache.cacheStmts.setRole.Exec(t.SID, t.RoleID, t.SelfAssign, t.Name)
		if err != nil {
			log.Warningln("Unable to set channel", err)
		}

	default:
		return (*cache.gcache).Set(key, value)
	}

	return err
}
