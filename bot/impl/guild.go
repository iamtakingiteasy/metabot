package impl

import (
	"time"

	"github.com/iamtakingiteasy/metabot/model"

	"github.com/bwmarrin/discordgo"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.GuildCreate:
			err := model.InsertGuildsRevision(ctx, evt.Guild)
			if err != nil {
				return err
			}
			if _, ok := ctx.Configs[evt.Guild.ID]; !ok {
				conf := &model.Config{
					GuildDiscordId:   evt.Guild.ID,
					Prefix:           "!",
					ColorError:       0xff0000,
					ColorWarn:        0xffff00,
					ColorInfo:        0x0000ff,
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
