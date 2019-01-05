package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.GuildRoleCreate:
			return model.InsertRolesRevision(ctx, evt.GuildID, evt.Role)
		case *discordgo.GuildRoleUpdate:
			return model.InsertRolesRevision(ctx, evt.GuildID, evt.Role)
		case *discordgo.GuildRoleDelete:
			return model.DeleteRolesByGuildRoleDiscordId(ctx, evt.GuildID, evt.RoleID)
		}
		return nil
	})
}
