package model

import (
	"bytes"
	"image/png"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/iamtakingiteasy/metabot/model/tmpl"

	"github.com/bwmarrin/discordgo"

	"github.com/iamtakingiteasy/metabot/bot"
)

type Guild struct {
	Id        uint64      `db:"guilds_id"`
	DiscordId string      `db:"guilds_discord_id"`
	Name      string      `db:"guilds_name"`
	Image     []byte      `db:"guilds_image"`
	Splash    []byte      `db:"guilds_splash"`
	Created   time.Time   `db:"guilds_created"`
	Deleted   pq.NullTime `db:"guilds_deleted"`
}

var (
	queryGuildsLastAll = tmpl.SelectTemplate("guilds", "_last")

	queryGuildsLastByGuildDiscordId = tmpl.SelectTemplate("guilds", "_last",
		"discord_id",
	)
	queryGuildsRevisionsByGuildDiscordId = tmpl.SelectTemplate("guilds", "",
		"discord_id",
	)
	insertGuildsRevision = tmpl.InsertTemplate("guilds", "",
		"discord_id",
		"name",
		"image",
		"splash",
	)
	deleteGuildsByGuildDiscordId = tmpl.DeleteTemplate("guilds", "",
		"discord_id",
	)
)

func QueryGuildsLastAll(ctx bot.Context) ([]*Guild, error) {
	var guilds []*Guild
	rows, err := ctx.Database().NamedQuery(queryGuildsLastAll, &Guild{})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		g := &Guild{}
		err := rows.StructScan(g)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, g)
	}
	return guilds, nil
}

func QueryGuildsLastByDiscordId(ctx bot.Context, discordId string) (*Guild, error) {
	g := &Guild{DiscordId: discordId}
	rows, err := ctx.Database().NamedQuery(queryGuildsLastByGuildDiscordId, g)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(g)
		if err != nil {
			return nil, err
		}
		return g, nil
	}
	return nil, ErrNoRows
}

func QueryGuildsRevisionsByDiscordId(ctx bot.Context, discordId string) ([]*Guild, error) {
	var guilds []*Guild
	rows, err := ctx.Database().NamedQuery(queryGuildsRevisionsByGuildDiscordId, &Guild{DiscordId: discordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		g := &Guild{}
		err := rows.StructScan(g)
		if err != nil {
			return nil, err
		}
		guilds = append(guilds, g)
	}
	return guilds, nil
}

func InsertGuildsRevision(ctx bot.Context, guild *discordgo.Guild) error {
	var icondata, splashdata bytes.Buffer

	if img, err := ctx.DiscordSession().GuildIcon(guild.ID); err == nil {
		err := png.Encode(&icondata, img)
		if err != nil {
			log.Println(err)
		}
	}

	if img, err := ctx.DiscordSession().GuildSplash(guild.ID); err == nil {
		err := png.Encode(&splashdata, img)
		if err != nil {
			log.Println(err)
		}
	}

	rows, err := ctx.Database().NamedQuery(insertGuildsRevision, &Guild{
		DiscordId: guild.ID,
		Name:      guild.Name,
		Image:     icondata.Bytes(),
		Splash:    splashdata.Bytes(),
	})
	if err != nil {
		return err
	}
	for rows.Next() {
		for _, r := range guild.Roles {
			err := InsertRolesRevision(ctx, guild.ID, r)
			if err != nil {
				return err
			}
		}
		for _, m := range guild.Members {
			err := InsertMembersRevision(ctx, m)
			if err != nil {
				return err
			}
		}
		for _, c := range guild.Channels {
			err := InsertChannelsRevision(ctx, c)
			if err != nil {
				return err
			}
		}
		for _, s := range guild.VoiceStates {
			err := InsertGuildsVoiceStatusRevision(ctx, s)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func DeleteGuildsByDiscordId(ctx bot.Context, discordId string) error {
	_, err := ctx.Database().NamedQuery(deleteGuildsByGuildDiscordId, &Guild{DiscordId: discordId})
	return err
}
