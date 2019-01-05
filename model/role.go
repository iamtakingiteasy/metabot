package model

import (
	"time"

	"github.com/iamtakingiteasy/metabot/api"

	"github.com/lib/pq"

	"github.com/iamtakingiteasy/metabot/model/tmpl"

	"github.com/bwmarrin/discordgo"
)

type Role struct {
	Id             uint64      `db:"roles_id"`
	GuildDiscordId string      `db:"roles_guild_discord_id"`
	DiscordId      string      `db:"roles_discord_id"`
	Name           string      `db:"roles_name"`
	Color          int         `db:"roles_color"`
	Position       int         `db:"roles_position"`
	Permissions    int         `db:"roles_permissions"`
	Created        time.Time   `db:"roles_created"`
	Deleted        pq.NullTime `db:"roles_deleted"`
}

var (
	queryRolesLastByGuildDiscordId = tmpl.SelectTemplate("roles", "_last",
		"guild_discord_id",
	)

	queryRolesLastByGuildRoleDiscordId = tmpl.SelectTemplate("roles", "_last",
		"guild_discord_id",
		"discord_id",
	)

	queryRolesRevisionsByGuildRoleDiscordId = tmpl.SelectTemplate("roles", "",
		"guild_discord_id",
		"discord_id",
	)
	insertRolesRevision = tmpl.InsertTemplate("roles", "",
		"guild_discord_id",
		"discord_id",
		"name",
		"color",
		"position",
		"permissions",
	)
	deleteRolesByGuildRoleDiscordId = tmpl.DeleteTemplate("roles", "",
		"guild_discord_id",
		"discord_id",
	)
)

func QueryRolesLastByGuildDiscordId(ctx api.Context, guildDiscordId string) ([]*Role, error) {
	var roles []*Role
	rows, err := ctx.Database().NamedQuery(queryRolesLastByGuildDiscordId, &Role{GuildDiscordId: guildDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := &Role{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func QueryRolesLastByGuildRoleDiscordId(ctx api.Context, guildDiscordId, roleDiscordId string) (*Role, error) {
	r := &Role{GuildDiscordId: guildDiscordId, DiscordId: roleDiscordId}
	rows, err := ctx.Database().NamedQuery(queryRolesLastByGuildRoleDiscordId, r)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, ErrNoRows
}

func QueryRolesRevisionsByGuildRoleDiscordId(ctx api.Context, guildDiscordId, roleDiscordId string) ([]*Role, error) {
	var roles []*Role
	rows, err := ctx.Database().NamedQuery(queryRolesRevisionsByGuildRoleDiscordId, &Role{GuildDiscordId: guildDiscordId, DiscordId: roleDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		r := &Role{}
		err := rows.StructScan(r)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func InsertRolesRevision(ctx api.Context, guildId string, role *discordgo.Role) error {
	_, err := ctx.Database().NamedQuery(insertRolesRevision, &Role{
		DiscordId:      role.ID,
		GuildDiscordId: guildId,
		Name:           role.Name,
		Color:          role.Color,
		Position:       role.Position,
		Permissions:    role.Permissions,
	})
	return err
}

func DeleteRolesByGuildRoleDiscordId(ctx api.Context, guildDiscordId, roleDiscordId string) error {
	_, err := ctx.Database().NamedQuery(deleteRolesByGuildRoleDiscordId, &Role{GuildDiscordId: guildDiscordId, DiscordId: roleDiscordId})
	return err
}
