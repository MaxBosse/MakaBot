package cache

import (
	"errors"

	"github.com/MaxBosse/MakaBot/log"

	_ "github.com/go-sql-driver/mysql"
)

func (cache *Cache) loader(key interface{}) (interface{}, error) {
	var err error

	switch t := key.(type) {
	case CacheServerKey:
		var cacheServer CacheServer
		err = cache.cacheStmts.getServersByGuildID.QueryRow(t.GuildID).Scan(&cacheServer.ID, &cacheServer.GuildID, &cacheServer.Enabled, &cacheServer.Nickname, &cacheServer.Prefix)
		if err != nil {
			log.Debugln("Unable to autoload Guild", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheServer.Guild, err = cache.session.State.Guild(cacheServer.GuildID)
			if err != nil {
				log.Warningln("Unable to get Guild", err)
			}
		}
		return cacheServer, nil
	case CacheChannelKey:
		var cacheChannel CacheChannel
		err = cache.cacheStmts.getChannelByChannelID.QueryRow(t.ChannelID).Scan(&cacheChannel.ID, &cacheChannel.SID, &cacheChannel.ChannelID, &cacheChannel.Listen, &cacheChannel.Name, &cacheChannel.AutoDelete, &cacheChannel.CType, &cacheChannel.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Channel", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheChannel.Channel, err = cache.session.State.Channel(cacheChannel.ChannelID)
			if err != nil {
				log.Warningln("Unable to get Channel", err)
			}
		}
		return cacheChannel, nil
	case CacheChannelName:
		var cacheChannel CacheChannel

		err = cache.cacheStmts.getChannelByChannelName.QueryRow(t.ChannelName, t.GuildID).Scan(&cacheChannel.ID, &cacheChannel.SID, &cacheChannel.ChannelID, &cacheChannel.Listen, &cacheChannel.Name, &cacheChannel.AutoDelete, &cacheChannel.CType, &cacheChannel.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Channel", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheChannel.Channel, err = cache.session.State.Channel(cacheChannel.ChannelID)
			if err != nil {
				log.Warningln("Unable to get Channel", err)
			}
		}
		return cacheChannel, nil
	case CacheRoleKey:
		var cacheRole CacheRole
		err = cache.cacheStmts.getChannelByRoleID.QueryRow(t.RoleID, t.GuildID).Scan(&cacheRole.ID, &cacheRole.SID, &cacheRole.RoleID, &cacheRole.SelfAssign, &cacheRole.Name, &cacheRole.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Role", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheRole.Role, err = cache.session.State.Role(cacheRole.GuildID, cacheRole.RoleID)
			if err != nil {
				log.Warningln("Unable to get Role", err)
			}
		}
		return cacheRole, nil
	case CacheRoleName:
		var cacheRole CacheRole
		err = cache.cacheStmts.getChannelByRoleName.QueryRow(t.RoleName, t.GuildID).Scan(&cacheRole.ID, &cacheRole.SID, &cacheRole.RoleID, &cacheRole.SelfAssign, &cacheRole.Name, &cacheRole.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Role", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheRole.Role, err = cache.session.State.Role(cacheRole.GuildID, cacheRole.RoleID)
			if err != nil {
				log.Warningln("Unable to get Role", err)
			}
		}
		return cacheRole, nil
	case CacheRoles:
		rows, err := cache.cacheStmts.getRoles.Query(t.GuildID)
		defer rows.Close()
		if err != nil {
			log.Debugln("Unable to autoload Roles", err)
			return nil, errors.New("unable to autoload")
		}

		roles := []CacheRole{}
		for rows.Next() {
			var cacheRole CacheRole

			err = rows.Scan(&cacheRole.ID, &cacheRole.SID, &cacheRole.RoleID, &cacheRole.SelfAssign, &cacheRole.Name, &cacheRole.GuildID)
			if err != nil {
				log.Debugln("Unable to autoload Role", err)
				return nil, errors.New("unable to autoload")
			}

			if cache.session != nil {
				cacheRole.Role, err = cache.session.State.Role(cacheRole.GuildID, cacheRole.RoleID)
				if err != nil {
					log.Warningln("Unable to get Role", err)
				}
			}

			roles = append(roles, cacheRole)
		}
		return roles, nil
	default:
		return nil, errors.New("unable to autoload")

	}
}
