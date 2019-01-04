package model

import (
	"time"

	"github.com/lib/pq"

	"github.com/iamtakingiteasy/metabot/model/tmpl"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
)

type Channel struct {
	Id              uint64      `db:"channels_id"`
	GuildDiscordId  string      `db:"channels_guild_discord_id"`
	DiscordId       string      `db:"channels_discord_id"`
	Name            string      `db:"channels_name"`
	Type            int         `db:"channels_type"`
	Bitrate         int         `db:"channels_bitrate"`
	ParentDiscordId string      `db:"channels_parent_discord_id"`
	Position        int         `db:"channels_position"`
	Nsfw            bool        `db:"channels_nsfw"`
	Topic           string      `db:"channels_topic"`
	UserLimit       int         `db:"channels_user_limit"`
	Created         time.Time   `db:"channels_created"`
	Deleted         pq.NullTime `db:"channels_deleted"`
}

var (
	queryChannelsLastByGuildDiscordId = tmpl.SelectTemplate("channels", "_last",
		"guild_discord_id",
	)
	queryChannelsLastByGuildChannelDiscordId = tmpl.SelectTemplate("channels", "_last",
		"guild_discord_id",
		"discord_id",
	)
	queryChannelsRevisionsByGuildChannelDiscordId = tmpl.SelectTemplate("channels", "",
		"guild_discord_id",
		"discord_id",
	)
	insertChannelsRevision = tmpl.InsertTemplate("channels", "",
		"guild_discord_id",
		"discord_id",
		"name",
		"type",
		"bitrate",
		"parent_discord_id",
		"position",
		"nsfw",
		"topic",
		"user_limit",
	)
	deleteChannelsByGuildChannelDiscordId = tmpl.DeleteTemplate("channels", "",
		"guild_discord_id",
		"discord_id",
	)
)

func QueryChannelsLastByGuildDiscordId(ctx bot.Context, guildDiscordId string) ([]*Channel, error) {
	var channels []*Channel
	rows, err := ctx.Database().NamedQuery(queryChannelsLastByGuildDiscordId, &Channel{GuildDiscordId: guildDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &Channel{}
		err := rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	return channels, nil
}

func QueryChannelsLastByGuildChannelDiscordId(ctx bot.Context, guildDiscordId, channelDiscordId string) (*Channel, error) {
	c := &Channel{GuildDiscordId: guildDiscordId, DiscordId: channelDiscordId}
	rows, err := ctx.Database().NamedQuery(queryChannelsLastByGuildChannelDiscordId, c)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, ErrNoRows
}

func QueryChannelsRevisionsByGuildChannelDiscordId(ctx bot.Context, guildDiscordId, channelDiscordId string) ([]*Channel, error) {
	var channels []*Channel
	rows, err := ctx.Database().NamedQuery(queryChannelsRevisionsByGuildChannelDiscordId, &Channel{GuildDiscordId: guildDiscordId, DiscordId: channelDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &Channel{}
		err := rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		channels = append(channels, c)
	}
	return channels, nil
}

func InsertChannelsRevision(ctx bot.Context, channel *discordgo.Channel) error {
	_, err := ctx.Database().NamedQuery(insertChannelsRevision, &Channel{
		DiscordId:       channel.ID,
		GuildDiscordId:  channel.GuildID,
		Name:            channel.Name,
		Type:            int(channel.Type),
		Bitrate:         channel.Bitrate,
		ParentDiscordId: channel.ParentID,
		Position:        channel.Position,
		Nsfw:            channel.NSFW,
		Topic:           channel.Topic,
		UserLimit:       channel.UserLimit,
	})
	if err != nil {
		return err
	}
	return nil
}

func DeleteChannelsByGuildChannelDiscordId(ctx bot.Context, guildDiscordId, channelDiscordId string) error {
	_, err := ctx.Database().NamedQuery(deleteChannelsByGuildChannelDiscordId, &Channel{GuildDiscordId: guildDiscordId, DiscordId: channelDiscordId})
	return err
}
