package commands

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
)

func init() {
	bot.AddGlobalRegistrars(func(ctx *bot.Context, registrar bot.Registrar) error {
		var prefix string
		var err error
		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Sets prefix to bot commands",
			Tokens: []bot.Token{
				bot.TypeLiteral("set"),
				bot.TypeLiteral("prefix"),
				&bot.TypeString{
					Value: &prefix,
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, ctx.Configs[msg.GuildID].Admins)
				if err != nil {
					return err
				}
				once := false
				for _, r := range ([]rune)(prefix) {
					if !unicode.IsPrint(r) || r == ' ' {
						return errors.New("should only be printable characters")
					}
					once = true
				}
				if !once {
					return errors.New("should be non-empty string")
				}
				cfg := ctx.Configs[msg.GuildID]
				cfg.Prefix = prefix
				err = cfg.Save(ctx)
				if err != nil {
					return err
				}
				return ctx.Send(msg.ChannelID, fmt.Sprintf("Try writing `%shelp`", prefix))
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
