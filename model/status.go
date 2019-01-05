package model

import (
	"time"

	"github.com/iamtakingiteasy/metabot/api"

	"github.com/lib/pq"

	"github.com/iamtakingiteasy/metabot/model/tmpl"

	"github.com/bwmarrin/discordgo"
)

type GuildVoiceStatus struct {
	Id               uint64      `db:"guilds_voice_status_id"`
	GuildDiscordId   string      `db:"guilds_voice_status_guild_discord_id"`
	ChannelDiscordId string      `db:"guilds_voice_status_channel_discord_id"`
	UserDiscordId    string      `db:"guilds_voice_status_user_discord_id"`
	Created          time.Time   `db:"guilds_voice_status_created"`
	Deleted          pq.NullTime `db:"guilds_voice_status_deleted"`
}

var (
	insertGuildsVoiceStatusRevision = tmpl.InsertTemplate("guilds_voice_status", "",
		"guild_discord_id",
		"channel_discord_id",
		"user_discord_id",
	)
)

func InsertGuildsVoiceStatusRevision(ctx api.Context, state *discordgo.VoiceState) error {
	_, err := ctx.Database().NamedQuery(insertGuildsVoiceStatusRevision, &GuildVoiceStatus{
		GuildDiscordId:   state.GuildID,
		ChannelDiscordId: state.ChannelID,
		UserDiscordId:    state.UserID,
	})
	return err
}
