package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
)

type Context interface {
	DiscordSession() *discordgo.Session
	Database() *sqlx.DB
}
