package impl

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
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
