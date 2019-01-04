package impl

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.GuildMemberAdd:
			return model.InsertMembersRevision(ctx, evt.Member)
		case *discordgo.GuildMemberUpdate:
			return model.InsertMembersRevision(ctx, evt.Member)
		case *discordgo.GuildMemberRemove:
			return model.DeleteMembersByGuildUserDiscordId(ctx, evt.Member.GuildID, evt.User.ID)
		}
		return nil
	})
}
