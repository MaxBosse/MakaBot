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
		}
		return cacheChannel, nil
	case CacheRoleKey:
		var cacheRole CacheRole
		err = cache.cacheStmts.getRoleByRoleID.QueryRow(t.RoleID, t.GuildID).Scan(&cacheRole.ID, &cacheRole.SID, &cacheRole.RoleID, &cacheRole.SelfAssign, &cacheRole.Name, &cacheRole.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Role", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheRole.Role, err = cache.session.State.Role(cacheRole.GuildID, cacheRole.RoleID)
		}
		return cacheRole, nil
	case CacheRoleName:
		var cacheRole CacheRole
		err = cache.cacheStmts.getRoleByRoleName.QueryRow(t.RoleName, t.GuildID).Scan(&cacheRole.ID, &cacheRole.SID, &cacheRole.RoleID, &cacheRole.SelfAssign, &cacheRole.Name, &cacheRole.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Role", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheRole.Role, err = cache.session.State.Role(cacheRole.GuildID, cacheRole.RoleID)
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
			}

			roles = append(roles, cacheRole)
		}
		return roles, err
	case CacheMemberGuildKey:
		var cacheMember CacheMember

		err = cache.cacheStmts.getMembers.QueryRow(t.UserID, t.GuildID).Scan(&cacheMember.ID, &cacheMember.SID, &cacheMember.UserID, &cacheMember.Username, &cacheMember.Discriminator, &cacheMember.Avatar, &cacheMember.Nick, &cacheMember.JoinedAt, &cacheMember.GuildID)
		if err != nil {
			log.Debugln("Unable to autoload Member", err)
			return nil, errors.New("unable to autoload")
		}

		if cache.session != nil {
			cacheMember.Member, err = cache.session.State.Member(cacheMember.GuildID, cacheMember.UserID)
		}
		return cacheMember, nil
	case CacheMemberKey:
		rows, err := cache.cacheStmts.getMembersByUserID.Query(t.UserID)
		defer rows.Close()
		if err != nil {
			log.Debugln("Unable to autoload Members", err)
			return nil, errors.New("unable to autoload")
		}

		members := []CacheMember{}
		for rows.Next() {
			var cacheMember CacheMember

			err = rows.Scan(&cacheMember.ID, &cacheMember.SID, &cacheMember.UserID, &cacheMember.Username, &cacheMember.Discriminator, &cacheMember.Avatar, &cacheMember.Nick, &cacheMember.JoinedAt, &cacheMember.GuildID)
			if err != nil {
				log.Debugln("Unable to autoload Member", err)
				return nil, errors.New("unable to autoload")
			}

			if cache.session != nil {
				cacheMember.Member, err = cache.session.State.Member(cacheMember.GuildID, cacheMember.UserID)
			}

			members = append(members, cacheMember)
		}
		return members, err
	case CacheMembers:
		rows, err := cache.cacheStmts.getMembers.Query(t.GuildID)
		defer rows.Close()
		if err != nil {
			log.Debugln("Unable to autoload Members", err)
			return nil, errors.New("unable to autoload")
		}

		members := []CacheMember{}
		for rows.Next() {
			var cacheMember CacheMember

			err = rows.Scan(&cacheMember.ID, &cacheMember.SID, &cacheMember.UserID, &cacheMember.Username, &cacheMember.Discriminator, &cacheMember.Avatar, &cacheMember.Nick, &cacheMember.JoinedAt, &cacheMember.GuildID)
			if err != nil {
				log.Debugln("Unable to autoload Member", err)
				return nil, errors.New("unable to autoload")
			}

			if cache.session != nil {
				cacheMember.Member, err = cache.session.State.Member(cacheMember.GuildID, cacheMember.UserID)
			}

			members = append(members, cacheMember)
		}
		return members, err
	default:
		return nil, errors.New("unable to autoload")

	}
}
