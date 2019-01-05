package bot

import (
	"github.com/bwmarrin/discordgo"
)

type AllPermission struct {
	Description string
	Combined    int
	Parts       []int
}

var AllTextPermissions = &AllPermission{
	Description: "All text permissions",
	Combined: discordgo.PermissionReadMessages |
		discordgo.PermissionSendMessages |
		discordgo.PermissionSendTTSMessages |
		discordgo.PermissionManageMessages |
		discordgo.PermissionEmbedLinks |
		discordgo.PermissionAttachFiles |
		discordgo.PermissionReadMessageHistory |
		discordgo.PermissionMentionEveryone |
		discordgo.PermissionUseExternalEmojis,
	Parts: []int{
		discordgo.PermissionReadMessages,
		discordgo.PermissionSendMessages,
		discordgo.PermissionSendTTSMessages,
		discordgo.PermissionManageMessages,
		discordgo.PermissionEmbedLinks,
		discordgo.PermissionAttachFiles,
		discordgo.PermissionReadMessageHistory,
		discordgo.PermissionMentionEveryone,
		discordgo.PermissionUseExternalEmojis,
	},
}
var AllVoicePermissions = &AllPermission{
	Description: "All voice permissions",
	Combined: discordgo.PermissionVoiceConnect |
		discordgo.PermissionVoiceSpeak |
		discordgo.PermissionVoiceMuteMembers |
		discordgo.PermissionVoiceDeafenMembers |
		discordgo.PermissionVoiceMoveMembers |
		discordgo.PermissionVoiceUseVAD,
	Parts: []int{
		discordgo.PermissionVoiceConnect,
		discordgo.PermissionVoiceSpeak,
		discordgo.PermissionVoiceMuteMembers,
		discordgo.PermissionVoiceDeafenMembers,
		discordgo.PermissionVoiceMoveMembers,
		discordgo.PermissionVoiceUseVAD,
	},
}

var AllChannelPermissions = &AllPermission{
	Description: "All channel permissions",
	Combined: AllTextPermissions.Combined | AllVoicePermissions.Combined |
		discordgo.PermissionCreateInstantInvite |
		discordgo.PermissionManageRoles |
		discordgo.PermissionManageChannels |
		discordgo.PermissionAddReactions |
		discordgo.PermissionViewAuditLogs |
		discordgo.PermissionChangeNickname,
	Parts: append(
		append(AllTextPermissions.Parts, AllVoicePermissions.Parts...),
		discordgo.PermissionCreateInstantInvite,
		discordgo.PermissionManageRoles,
		discordgo.PermissionManageChannels,
		discordgo.PermissionAddReactions,
		discordgo.PermissionViewAuditLogs,
		discordgo.PermissionChangeNickname,
	),
}

var AllKnownPermissions = &AllPermission{
	Description: "All known permissions",
	Combined: AllChannelPermissions.Combined |
		discordgo.PermissionKickMembers |
		discordgo.PermissionBanMembers |
		discordgo.PermissionManageServer |
		discordgo.PermissionAdministrator |
		discordgo.PermissionManageWebhooks |
		discordgo.PermissionManageEmojis |
		discordgo.PermissionManageNicknames,
	Parts: append(
		AllChannelPermissions.Parts,
		discordgo.PermissionKickMembers,
		discordgo.PermissionBanMembers,
		discordgo.PermissionManageServer,
		discordgo.PermissionAdministrator,
		discordgo.PermissionManageWebhooks,
		discordgo.PermissionManageEmojis,
		discordgo.PermissionManageNicknames,
	),
}

var AggregatePermissions = []*AllPermission{
	AllKnownPermissions,
	AllChannelPermissions,
	AllVoicePermissions,
	AllTextPermissions,
}

var KnownPermissions = map[int]string{
	discordgo.PermissionViewAuditLogs:       "View audit logs",
	discordgo.PermissionAddReactions:        "Add reactions",
	discordgo.PermissionManageServer:        "Manage server",
	discordgo.PermissionManageChannels:      "Manage channels",
	discordgo.PermissionAdministrator:       "Administrator",
	discordgo.PermissionBanMembers:          "Ban members",
	discordgo.PermissionKickMembers:         "Kick members",
	discordgo.PermissionCreateInstantInvite: "Create instant invite",
	discordgo.PermissionManageEmojis:        "Manage emojis",
	discordgo.PermissionManageWebhooks:      "Mange webhooks",
	discordgo.PermissionManageRoles:         "Manage roles",
	discordgo.PermissionManageNicknames:     "Manage nicknames",
	discordgo.PermissionChangeNickname:      "Change nickname",
	discordgo.PermissionVoiceUseVAD:         "Use voice activity",
	discordgo.PermissionVoiceMoveMembers:    "Move voice members",
	discordgo.PermissionVoiceMuteMembers:    "Mute voice members",
	discordgo.PermissionVoiceDeafenMembers:  "Deafen voice members",
	discordgo.PermissionVoiceSpeak:          "Speak on voice",
	discordgo.PermissionVoiceConnect:        "Connect to voice",
	discordgo.PermissionUseExternalEmojis:   "Use external emojis",
	discordgo.PermissionMentionEveryone:     "Mention everyone",
	discordgo.PermissionReadMessageHistory:  "Read message history",
	discordgo.PermissionAttachFiles:         "Attach files",
	discordgo.PermissionEmbedLinks:          "Embed links",
	discordgo.PermissionManageMessages:      "Manage messages",
	discordgo.PermissionSendTTSMessages:     "Send TTS messages",
	discordgo.PermissionSendMessages:        "Send text messages",
	discordgo.PermissionReadMessages:        "Read text messages",
}

func (ctx *Context) ExplainPermissions(permissions int) (all []*AllPermission, rest []int) {
	for _, agg := range AggregatePermissions {
		if agg.Combined&permissions == agg.Combined {
			all = append(all, agg)
			permissions = permissions &^ agg.Combined
		}
	}

	for {
		found := false
		for p := range KnownPermissions {
			if permissions&p == p {
				rest = append(rest, p)
				permissions = permissions &^ p
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	return
}
