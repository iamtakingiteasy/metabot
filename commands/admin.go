package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
	"github.com/iamtakingiteasy/metabot/model"
)

func init() {
	bot.AddGlobalRegistrars(func(ctx *bot.Context, registrar bot.Registrar) error {
		var err error
		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Lists current admin rules",
			Tokens: []bot.Token{
				bot.TypeLiteral("admin"),
				bot.TypeLiteral("rules"),
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				rules := ctx.Configs[msg.GuildID].Admins
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, rules)
				if err != nil {
					return err
				}
				embed := &discordgo.MessageEmbed{}
				for _, rule := range rules {
					typestr := "unknown"
					sb := &strings.Builder{}

					switch rule.Type {
					case model.ConfigAdminTypePermissionMask:
						typestr = "Permission mask"
						all, parts := ctx.ExplainPermissions(rule.PermissionMask)
						sb.WriteString("Mask: `")
						sb.WriteString(strconv.FormatUint(uint64(rule.PermissionMask), 10))
						sb.WriteString("`\nExplanation:\n")
						for j, agg := range all {
							sb.WriteString("-- ")
							sb.WriteString(agg.Description)
							sb.WriteString(" (")
							for k, p := range agg.Parts {
								sb.WriteString(bot.KnownPermissions[p])
								sb.WriteString("[`")
								sb.WriteString(strconv.FormatUint(uint64(p), 10))
								sb.WriteString("`]")
								if k+1 < len(agg.Parts) {
									sb.WriteString(", ")
								}
							}
							sb.WriteString(")")
							if j+1 < len(all) {
								sb.WriteString("\n")
							}
						}
						if len(all) > 0 && len(parts) > 0 {
							sb.WriteString("\n")
						}
						for k, p := range parts {
							sb.WriteString("-- ")
							sb.WriteString(bot.KnownPermissions[p])
							sb.WriteString("[`")
							sb.WriteString(strconv.FormatUint(uint64(p), 10))
							sb.WriteString("`]")
							if k+1 < len(parts) {
								sb.WriteString("\n")
							}
						}
					case model.ConfigAdminTypeRoleDiscordId:
						typestr = "Role"
						role, err := ctx.Session.State.Role(msg.GuildID, rule.RoleDiscordId)
						sb.WriteString("`")
						if err != nil {
							sb.WriteString("*Unknown role* &")
							sb.WriteString(rule.RoleDiscordId)
						} else {
							sb.WriteString(role.Name)
							sb.WriteString("` - `")
							sb.WriteString(role.ID)
						}
						sb.WriteString("`")
					case model.ConfigAdminTypeUserDiscordId:
						typestr = "User"
						member, err := ctx.Session.GuildMember(msg.GuildID, rule.UserDiscordId)
						sb.WriteString("`")
						if err != nil {
							sb.WriteString("*Unknown user* !")
							sb.WriteString(rule.UserDiscordId)
						} else {
							sb.WriteString(member.User.Username)
							if member.Nick != "" {
								sb.WriteString("` aka `")
								sb.WriteString(member.Nick)
							}
							sb.WriteString("` - `")
							sb.WriteString(member.User.ID)
						}
						sb.WriteString("`")
					}
					num := strconv.FormatUint(uint64(rule.Id), 10)
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  "#" + num + " " + typestr,
						Value: sb.String(),
					})
				}
				return ctx.SendInfo(msg.ChannelID, embed)
			},
		})
		if err != nil {
			return err
		}

		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Lists current admin users",
			Tokens: []bot.Token{
				bot.TypeLiteral("admin"),
				bot.TypeLiteral("list"),
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				rules := ctx.Configs[msg.GuildID].Admins
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, rules)
				if err != nil {
					return err
				}
				embed := &discordgo.MessageEmbed{}
				var members []*discordgo.Member

				after := ""
				for len(members)%1000 == 0 {
					ms, err := ctx.Session.GuildMembers(msg.GuildID, after, 1000)
					if err != nil {
						return err
					}
					if len(ms) == 0 {
						break
					}
					for _, m := range ms {
						if ctx.OkAdminMember(msg.ChannelID, m, rules) == nil {
							members = append(members, m)
						}
					}
					after = ms[len(ms)-1].User.ID
				}

				mapped := make(map[*discordgo.Member][]*model.ConfigAdmin)
				for _, member := range members {
					for _, rule := range rules {
						if ctx.OkAdminMember(msg.ChannelID, member, []*model.ConfigAdmin{rule}) == nil {
							mapped[member] = append(mapped[member], rule)
						}
					}
				}

				for member, rs := range mapped {
					name := &strings.Builder{}
					name.WriteString(member.User.Username)
					if member.Nick != "" {
						name.WriteString(" aka ")
						name.WriteString(member.Nick)
					}
					name.WriteString(" - ")
					name.WriteString(member.User.ID)

					value := &strings.Builder{}
					value.WriteString("due to:\n")
					for i, r := range rs {
						value.WriteString("-- #")
						num := strconv.FormatUint(uint64(r.Id), 10)
						value.WriteString(num)
						value.WriteString(" ")
						switch r.Type {
						case model.ConfigAdminTypePermissionMask:
							value.WriteString("Permission mask `")
							perms, err := ctx.Session.UserChannelPermissions(member.User.ID, msg.ChannelID)
							if err != nil {
								return err
							}
							value.WriteString(strconv.FormatUint(uint64(perms), 10))
							value.WriteString("` is matching required `")
							value.WriteString(strconv.FormatUint(uint64(r.PermissionMask), 10))
							value.WriteString("`")
						case model.ConfigAdminTypeRoleDiscordId:
							value.WriteString("User included in role `")
							role, err := ctx.Session.State.Role(msg.GuildID, r.RoleDiscordId)
							if err != nil {
								return err
							}
							value.WriteString(role.Name)
							value.WriteString("` - `")
							value.WriteString(role.ID)
							value.WriteString("`")
						case model.ConfigAdminTypeUserDiscordId:
							value.WriteString("User included directly")
						}
						if i+1 < len(rs) {
							value.WriteString("\n")
						}
					}
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  name.String(),
						Value: value.String(),
					})
				}
				return ctx.SendInfo(msg.ChannelID, embed)
			},
		})
		if err != nil {
			return err
		}

		errAddValid := errors.New("must be either permission integer, user mention/id/name or role mention/id/name")
		var permission int
		var userId string
		var roleId string
		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Adds new admin rule\n`:ruletype` -- one of `permission` (integer), `role` (mention/id/name) or `user` (mention/id/name)",
			Tokens: []bot.Token{
				bot.TypeLiteral("admin"),
				bot.TypeLiteral("rule"),
				bot.TypeLiteral("add"),
				&bot.TypeOneOf{
					Name:  ":ruletype",
					Error: errAddValid,
					Variants: []bot.Token{
						&bot.TypeNumber{
							Value: &permission,
						},
						&bot.TypeUser{
							Value: &userId,
						},
						&bot.TypeRole{
							Value: &roleId,
						},
					},
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				guild := ctx.Configs[msg.GuildID]
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, guild.Admins)
				if err != nil {
					return err
				}

				if permission != 0 {
					return guild.AddAdmin(ctx, &model.ConfigAdmin{
						GuildDiscordId: msg.GuildID,
						Type:           model.ConfigAdminTypePermissionMask,
						PermissionMask: permission,
					})
				}
				if roleId != "" {
					return guild.AddAdmin(ctx, &model.ConfigAdmin{
						GuildDiscordId: msg.GuildID,
						Type:           model.ConfigAdminTypeRoleDiscordId,
						RoleDiscordId:  roleId,
					})
				}
				if userId != "" {
					return guild.AddAdmin(ctx, &model.ConfigAdmin{
						GuildDiscordId: msg.GuildID,
						Type:           model.ConfigAdminTypeUserDiscordId,
						UserDiscordId:  userId,
					})
				}
				return errAddValid
			},
		})
		if err != nil {
			return err
		}

		errRemoveValid := errors.New("must be valid rule id, optionally prefixed by hash mark (#)")
		var rulenoPat string
		var ruleno uint64
		err = registrar.AddCommand(&bot.Descriptor{
			Group:       "Admin",
			Description: "Removes rule by id",
			Tokens: []bot.Token{
				bot.TypeLiteral("admin"),
				bot.TypeLiteral("rule"),
				bot.TypeLiteral("remove"),
				&bot.TypeOneOf{
					Name:  ":ruleno",
					Error: errRemoveValid,
					Variants: []bot.Token{
						&bot.TypeNumber{
							Value: &ruleno,
						},
						&bot.TypeString{
							Pattern: "^#[0-9]+$",
							Value:   &rulenoPat,
						},
					},
				},
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				guild := ctx.Configs[msg.GuildID]
				err := ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, guild.Admins)
				if err != nil {
					return err
				}

				if ruleno == 0 && rulenoPat != "" {
					ruleno, err = strconv.ParseUint(rulenoPat[1:], 10, 64)
					if err != nil {
						return err
					}
				}

				if ruleno == 0 {
					return errRemoveValid
				}

				for i, rule := range guild.Admins {
					if rule.Id == ruleno {
						fmt.Println(guild.Admins)
						var newrules []*model.ConfigAdmin
						for n, r := range guild.Admins {
							if n != i {
								newrules = append(newrules, r)
							}
						}
						if ctx.OkAdmin(msg.GuildID, msg.ChannelID, msg.Author.ID, newrules) != nil {
							return errors.New("removing this rule would remove last rule making you admin in first place")
						}
						return guild.DeleteAdminById(ctx, ruleno)
					}
				}
				return errors.New("no rule with such id")
			},
		})
		if err != nil {
			return err
		}
		return nil
	})
}
