package impl

import (
	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	AddGlobalEventHandler(func(ctx *Context, raw interface{}) error {
		switch evt := raw.(type) {
		case *discordgo.UserUpdate:
			return model.InsertUsersRevision(ctx, evt.User)
		}
		return nil
	})
}
