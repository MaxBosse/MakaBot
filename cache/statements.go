package cache

import (
	"database/sql"

	"github.com/MaxBosse/MakaBot/log"

	_ "github.com/go-sql-driver/mysql"
)

type cacheStmts struct {
	setServer               *sql.Stmt
	setChannel              *sql.Stmt
	setRole                 *sql.Stmt
	getServersByGuildID     *sql.Stmt
	getChannelByChannelID   *sql.Stmt
	getChannelByChannelName *sql.Stmt
	getChannelByRoleID      *sql.Stmt
	getChannelByRoleName    *sql.Stmt
	getRoles                *sql.Stmt
	removeServer            *sql.Stmt
	removeChannel           *sql.Stmt
	removeRole              *sql.Stmt
}

func (cache *Cache) prepareStatements() *cacheStmts {
	var stmts cacheStmts
	var err error

	stmts.removeRole, err = cache.db.Prepare(
		"DELETE " +
			"FROM " +
			"	roles " +
			"WHERE " +
			"	roles.id = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement removeRole.", err.Error())
	}

	stmts.removeChannel, err = cache.db.Prepare(
		"DELETE " +
			"FROM " +
			"	channels " +
			"WHERE " +
			"	channels.id = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement removeChannel.", err.Error())
	}

	stmts.removeServer, err = cache.db.Prepare(
		"DELETE " +
			"FROM " +
			"	servers " +
			"WHERE " +
			"	servers.id = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement removeServer.", err.Error())
	}

	stmts.getRoles, err = cache.db.Prepare(
		"SELECT " +
			"	roles.id, " +
			"	roles.sID, " +
			"	roles.roleID, " +
			"	roles.selfAssign, " +
			"	roles.name, " +
			"	servers.guildID " +
			"FROM " +
			"	roles " +
			"LEFT JOIN servers " +
			"	ON roles.sID = servers.id " +
			"WHERE " +
			"	servers.guildID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getRoles.", err.Error())
	}

	stmts.setServer, err = cache.db.Prepare(
		"INSERT INTO servers " +
			"	(guildID, enabled, nickname, prefix, name) " +
			"VALUES " +
			"	(?, ?, ?, ?, ?) " +
			"ON DUPLICATE KEY UPDATE " +
			"	enabled=VALUES(enabled), " +
			"	nickname=VALUES(nickname), " +
			"	prefix=VALUES(prefix), " +
			"	name=VALUES(name) ")
	if err != nil {
		log.Fatalln("Could not prepare statement setServer.", err.Error())
	}

	stmts.setChannel, err = cache.db.Prepare(
		"INSERT INTO channels " +
			"	(sID, channelID, listen, name, autoDeleteSec, cType) " +
			"VALUES " +
			"	(?, ?, ?, ?, ?, ?) " +
			"ON DUPLICATE KEY UPDATE " +
			"	listen=VALUES(listen), " +
			"	autoDeleteSec=VALUES(autoDeleteSec), " +
			"	cType=VALUES(cType), " +
			"	name=VALUES(name) ")
	if err != nil {
		log.Fatalln("Could not prepare statement setChannel.", err.Error())
	}

	stmts.setRole, err = cache.db.Prepare(
		"INSERT INTO roles " +
			"	(sID, roleID, selfAssign, name) " +
			"VALUES " +
			"	(?, ?, ?, ?) " +
			"ON DUPLICATE KEY UPDATE " +
			"	selfAssign=VALUES(selfAssign), " +
			"	name=VALUES(name) ")
	if err != nil {
		log.Fatalln("Could not prepare statement setRole.", err.Error())
	}

	stmts.getServersByGuildID, err = cache.db.Prepare(
		"SELECT " +
			"	servers.id, " +
			"	servers.guildID, " +
			"	servers.enabled, " +
			"	servers.nickname, " +
			"	servers.prefix " +
			"FROM " +
			"	servers " +
			"WHERE " +
			"	servers.guildID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getServersByGuildID.", err.Error())
	}

	stmts.getChannelByChannelID, err = cache.db.Prepare(
		"SELECT " +
			"	channels.id, " +
			"	channels.sID, " +
			"	channels.channelID, " +
			"	channels.listen, " +
			"	channels.name, " +
			"	channels.autoDeleteSec, " +
			"	channels.cType, " +
			"	servers.guildID " +
			"FROM " +
			"	channels " +
			"LEFT JOIN servers " +
			"	ON channels.sID = servers.id " +
			"WHERE " +
			"	channels.channelID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getChannelByChannelID.", err.Error())
	}

	stmts.getChannelByChannelName, err = cache.db.Prepare(
		"SELECT " +
			"	channels.id, " +
			"	channels.sID, " +
			"	channels.channelID, " +
			"	channels.listen, " +
			"	channels.name, " +
			"	channels.autoDeleteSec, " +
			"	channels.cType, " +
			"	servers.guildID " +
			"FROM " +
			"	channels " +
			"LEFT JOIN servers " +
			"	ON channels.sID = servers.id " +
			"WHERE " +
			"	channels.name = ? " +
			"	AND servers.guildID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getChannelByChannelName.", err.Error())
	}

	stmts.getChannelByRoleID, err = cache.db.Prepare(
		"SELECT " +
			"	roles.id, " +
			"	roles.sID, " +
			"	roles.roleID, " +
			"	roles.selfAssign, " +
			"	roles.name, " +
			"	servers.guildID " +
			"FROM " +
			"	roles " +
			"LEFT JOIN servers " +
			"	ON roles.sID = servers.id " +
			"WHERE " +
			"	roles.roleID = ? " +
			"	AND servers.guildID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getChannelByRoleID.", err.Error())
	}

	stmts.getChannelByRoleName, err = cache.db.Prepare(
		"SELECT " +
			"	roles.id, " +
			"	roles.sID, " +
			"	roles.roleID, " +
			"	roles.selfAssign, " +
			"	roles.name, " +
			"	servers.guildID " +
			"FROM " +
			"	roles " +
			"LEFT JOIN servers " +
			"	ON roles.sID = servers.id " +
			"WHERE " +
			"	roles.name = ? " +
			"	AND servers.guildID = ? ")
	if err != nil {
		log.Fatalln("Could not prepare statement getChannelByRoleName.", err.Error())
	}

	return &stmts
}
