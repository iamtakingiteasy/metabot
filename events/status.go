package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	bot.AddGlobalEventHandler(func(ctx *bot.Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.VoiceStateUpdate:
			return model.InsertGuildsVoiceStatusRevision(ctx, evt.VoiceState)
		}
		return nil
	})
}
