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

type User struct {
	Id            uint64      `db:"users_id"`
	DiscordId     string      `db:"users_discord_id"`
	Name          string      `db:"users_name"`
	Discriminator string      `db:"users_discriminator"`
	Bot           bool        `db:"users_bot"`
	Email         string      `db:"users_email"`
	Avatar        []byte      `db:"users_avatar"`
	Locale        string      `db:"users_locale"`
	Verified      bool        `db:"users_verified"`
	Created       time.Time   `db:"users_created"`
	Deleted       pq.NullTime `db:"users_deleted"`
}

var (
	queryUsersLastAll = tmpl.SelectTemplate("users", "_last")

	queryUsersLastByUserDiscordId = tmpl.SelectTemplate("users", "_last",
		"discord_id",
	)

	queryUsersRevisionsByUserDiscordId = tmpl.SelectTemplate("users", "",
		"discord_id",
	)
	insertUsersRevision = tmpl.InsertTemplate("users", "",
		"discord_id",
		"name",
		"discriminator",
		"bot",
		"email",
		"avatar",
		"locale",
		"verified",
	)
	deleteUsersByUserDiscordId = tmpl.DeleteTemplate("users", "",
		"discord_id",
	)
)

func QueryUsersLastAll(ctx bot.Context) ([]*User, error) {
	var users []*User
	rows, err := ctx.Database().NamedQuery(queryUsersLastAll, &User{})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		u := &User{}
		err := rows.StructScan(u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func QueryUsersLastByUserDiscordId(ctx bot.Context, discordId string) (*User, error) {
	u := &User{DiscordId: discordId}
	rows, err := ctx.Database().NamedQuery(queryUsersLastByUserDiscordId, u)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(u)
		if err != nil {
			return nil, err
		}
		return u, nil
	}
	return nil, ErrNoRows
}

func QueryUsersRevisionsByUserDiscordId(ctx bot.Context, discordId string) ([]*User, error) {
	var users []*User
	rows, err := ctx.Database().NamedQuery(queryUsersRevisionsByUserDiscordId, &User{DiscordId: discordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		u := &User{}
		err := rows.StructScan(u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func InsertUsersRevision(ctx bot.Context, user *discordgo.User) error {
	var avatardata bytes.Buffer

	if img, err := ctx.DiscordSession().UserAvatar(user.ID); err == nil {
		err := png.Encode(&avatardata, img)
		if err != nil {
			log.Println(err)
		}
	}

	_, err := ctx.Database().NamedQuery(insertUsersRevision, &User{
		DiscordId:     user.ID,
		Name:          user.Username,
		Discriminator: user.Discriminator,
		Bot:           user.Bot,
		Email:         user.Email,
		Avatar:        avatardata.Bytes(),
		Locale:        user.Locale,
		Verified:      user.Verified,
	})
	return err
}

func DeleteUsersByUserDiscordId(ctx bot.Context, discordId string) error {
	_, err := ctx.Database().NamedQuery(deleteUsersByUserDiscordId, &User{DiscordId: discordId})
	return err
}
