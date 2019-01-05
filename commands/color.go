package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
)

func init() {
	bot.AddGlobalRegistrars(func(ctx *bot.Context, registrar bot.Registrar) error {
		var infoColor, warnColor, errorColor int
		var err error
		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Sets color for info messages",
			Tokens: []bot.Token{
				bot.TypeLiteral("set"),
				bot.TypeLiteral("color"),
				bot.TypeLiteral("info"),
				&bot.TypeColor{
					Value: &infoColor,
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, ctx.Configs[msg.GuildID].Admins)
				if err != nil {
					return err
				}
				cfg := ctx.Configs[msg.GuildID]
				cfg.ColorInfo = infoColor
				return cfg.Save(ctx)
			},
		})
		if err != nil {
			return err
		}

		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Sets color for warning messages",
			Tokens: []bot.Token{
				bot.TypeLiteral("set"),
				bot.TypeLiteral("color"),
				bot.TypeLiteral("warn"),
				&bot.TypeColor{
					Value: &warnColor,
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, ctx.Configs[msg.GuildID].Admins)
				if err != nil {
					return err
				}
				cfg := ctx.Configs[msg.GuildID]
				cfg.ColorWarn = warnColor
				return cfg.Save(ctx)
			},
		})
		if err != nil {
			return err
		}

		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Sets color for error messages",
			Tokens: []bot.Token{
				bot.TypeLiteral("set"),
				bot.TypeLiteral("color"),
				bot.TypeLiteral("error"),
				&bot.TypeColor{
					Value: &errorColor,
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, ctx.Configs[msg.GuildID].Admins)
				if err != nil {
					return err
				}
				cfg := ctx.Configs[msg.GuildID]
				cfg.ColorError = errorColor
				return cfg.Save(ctx)
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
