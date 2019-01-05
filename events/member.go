package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
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
