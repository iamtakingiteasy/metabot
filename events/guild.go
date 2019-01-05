package events

import (
	"time"

	"github.com/iamtakingiteasy/metabot/bot"

	"github.com/iamtakingiteasy/metabot/model"

	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.GuildCreate:
			err := model.InsertGuildsRevision(ctx, evt.Guild)
			if err != nil {
				return err
			}
			if _, ok := ctx.Configs[evt.Guild.ID]; !ok {
				conf := &model.Config{
					GuildDiscordId:   evt.Guild.ID,
					Prefix:           "changeme!",
					ColorError:       0xffaaaa,
					ColorWarn:        0xffffaa,
					ColorInfo:        0xaaaaff,
					AutoremoveActive: true,
					AutoremoveTime:   1 * time.Minute,
					RestrictActive:   false,
				}
				err := conf.Save(ctx)
				if err != nil {
					return err
				}
				err = conf.AddAdmin(ctx, &model.ConfigAdmin{
					Type:           model.ConfigAdminTypePermissionMask,
					PermissionMask: discordgo.PermissionAdministrator,
				})
				if err != nil {
					return err
				}
				ctx.Configs[evt.Guild.ID] = conf
			}
			return nil
		case *discordgo.GuildUpdate:
			return model.InsertGuildsRevision(ctx, evt.Guild)
		case *discordgo.GuildDelete:
			return model.DeleteGuildsByDiscordId(ctx, evt.Guild.ID)
		}
		return nil
	})
}
