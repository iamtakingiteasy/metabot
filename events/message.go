package events

import (
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"

	"github.com/bwmarrin/discordgo"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.MessageCreate:
			err := model.InsertMessagesRevision(ctx, evt.Message)
			if err != nil {
				return err
			}
			return ctx.Commands.ProcessCommand(evt.Message)
		case *discordgo.MessageUpdate:
			return model.InsertMessagesRevision(ctx, evt.Message)
		case *discordgo.MessageDelete:
			return model.DeleteMessagesByGuildChannelMessageDiscordId(ctx, evt.Message.GuildID, evt.Message.ChannelID, evt.Message.ID)
		}
		return nil
	})
}
