package impl

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.VoiceStateUpdate:
			return model.InsertGuildsVoiceStatusRevision(ctx, evt.VoiceState)
		}
		return nil
	})
}
