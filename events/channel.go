package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
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
