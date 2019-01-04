package impl

import (
	"github.com/iamtakingiteasy/metabot/model"

	"github.com/bwmarrin/discordgo"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.MessageCreate:
			return model.InsertMessagesRevision(ctx, evt.Message)
		case *discordgo.MessageUpdate:
			return model.InsertMessagesRevision(ctx, evt.Message)
		case *discordgo.MessageDelete:
			return model.DeleteMessagesByGuildChannelMessageDiscordId(ctx, evt.Message.GuildID, evt.Message.ChannelID, evt.Message.ID)
		}
		return nil
	})
}
