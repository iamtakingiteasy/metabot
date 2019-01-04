package impl

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.ChannelCreate:
			return model.InsertChannelsRevision(ctx, evt.Channel)
		case *discordgo.ChannelUpdate:
			return model.InsertChannelsRevision(ctx, evt.Channel)
		case *discordgo.ChannelDelete:
			return model.DeleteChannelsByGuildChannelDiscordId(ctx, evt.Channel.GuildID, evt.Channel.ID)
		}
		return nil
	})
}
