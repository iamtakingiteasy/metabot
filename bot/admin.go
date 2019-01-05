package bot

import (
	"errors"

	"github.com/bwmarrin/discordgo"

	"github.com/iamtakingiteasy/metabot/model"
)

func (ctx *Context) OkAdminMember(channelId string, member *discordgo.Member, rules []*model.ConfigAdmin) error {
	for _, rule := range rules {
		switch rule.Type {
		case model.ConfigAdminTypePermissionMask:
			perms, err := ctx.Session.UserChannelPermissions(member.User.ID, channelId)
			if err != nil {
				return err
			}
			if perms&rule.PermissionMask == rule.PermissionMask {
				return nil
			}
		case model.ConfigAdminTypeRoleDiscordId:
			for _, r := range member.Roles {
				if r == rule.RoleDiscordId {
					return nil
				}
			}
		case model.ConfigAdminTypeUserDiscordId:
			if member.User.ID == rule.UserDiscordId {
				return nil
			}
		}
	}
	return errors.New("requires admin privileges")
}

func (ctx *Context) OkAdmin(guildId, channelId, userId string, rules []*model.ConfigAdmin) error {
	member, err := ctx.Session.GuildMember(guildId, userId)
	if err != nil {
		return err
	}
	return ctx.OkAdminMember(channelId, member, rules)
}
