package model

import (
	"time"

	"github.com/iamtakingiteasy/metabot/api"

	"github.com/iamtakingiteasy/metabot/model/tmpl"
	"github.com/lib/pq"

	"github.com/bwmarrin/discordgo"
)

type Member struct {
	Id             uint64      `db:"members_id"`
	GuildDiscordId string      `db:"members_guild_discord_id"`
	UserDiscordId  string      `db:"members_user_discord_id"`
	Nick           string      `db:"members_nick"`
	JoinedAt       time.Time   `db:"members_joined_at"`
	Created        time.Time   `db:"members_created"`
	Deleted        pq.NullTime `db:"members_deleted"`
}

var (
	queryMembersLastByGuildDiscordId = tmpl.SelectTemplate("members", "_last",
		"guild_discord_id",
	)
	queryMembersLastByGuildUserDiscordId = tmpl.SelectTemplate("members", "_last",
		"guild_discord_id",
		"user_discord_id",
	)
	queryMembersRevisionsByGuildUserDiscordId = tmpl.SelectTemplate("members", "",
		"guild_discord_id",
		"user_discord_id",
	)
	insertMembersRevision = tmpl.InsertTemplate("members", "",
		"guild_discord_id",
		"user_discord_id",
		"nick",
		"joined_at",
	)
	deleteMembersByGuildUserDiscordId = tmpl.DeleteTemplate("members", "",
		"guild_discord_id",
		"user_discord_id",
	)
)

func QueryMembersLastByGuildDiscordId(ctx api.Context, guildDiscordId string) ([]*Member, error) {
	var members []*Member
	rows, err := ctx.Database().NamedQuery(queryMembersLastByGuildDiscordId, &Member{GuildDiscordId: guildDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := &Member{}
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func QueryMembersLastByGuildUserDiscordId(ctx api.Context, guildDiscordId, userDiscordId string) (*Member, error) {
	m := &Member{GuildDiscordId: guildDiscordId, UserDiscordId: userDiscordId}
	rows, err := ctx.Database().NamedQuery(queryMembersLastByGuildUserDiscordId, m)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, ErrNoRows
}

func QueryMembersRevisionsByGuildUserDiscordId(ctx api.Context, guildDiscordId, userDiscordId string) ([]*Member, error) {
	var members []*Member
	rows, err := ctx.Database().NamedQuery(queryMembersRevisionsByGuildUserDiscordId, &Member{GuildDiscordId: guildDiscordId, UserDiscordId: userDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := &Member{}
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func InsertMembersRevision(ctx api.Context, member *discordgo.Member) error {
	joined, err := member.JoinedAt.Parse()
	if err != nil {
		return err
	}
	rows, err := ctx.Database().NamedQuery(insertMembersRevision, &Member{
		UserDiscordId:  member.User.ID,
		GuildDiscordId: member.GuildID,
		Nick:           member.Nick,
		JoinedAt:       joined,
	})
	if err != nil {
		return err
	}
	for rows.Next() {
		user, err := ctx.DiscordSession().User(member.User.ID)
		if err != nil {
			return err
		}
		err = InsertUsersRevision(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteMembersByGuildUserDiscordId(ctx api.Context, guildDiscordId, userDiscordId string) error {
	_, err := ctx.Database().NamedQuery(deleteMembersByGuildUserDiscordId, &Member{GuildDiscordId: guildDiscordId, UserDiscordId: userDiscordId})
	return err
}
