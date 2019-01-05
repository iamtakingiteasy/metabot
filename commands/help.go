package commands

import (
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/bot"
)

func init() {
	bot.AddGlobalRegistrars(func(ctx *bot.Context, registrar bot.Registrar) error {
		return registrar.AddCommand(&bot.Descriptor{
			Group:       "General",
			Description: "Displays helpful message",
			Tokens: []bot.Token{
				bot.TypeLiteral("help"),
			},
			Handler: func(msg *discordgo.Message, tokens []string) error {
				embed := &discordgo.MessageEmbed{}
				mapped := make(map[string][]*bot.Descriptor)
				for _, cmd := range ctx.Commands.Commands() {
					mapped[cmd.Group] = append(mapped[cmd.Group], cmd)
				}
				var sorted []string
				for k := range mapped {
					sorted = append(sorted, k)
				}
				sort.Strings(sorted)
				for _, g := range sorted {
					name := "**" + g + "**"
					value := &strings.Builder{}
					for z := 0; z < len(g); z++ {
						value.WriteString("\\`")
					}
					embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
						Name:  name,
						Value: value.String(),
					})
					for _, v := range mapped[g] {
						sb := &strings.Builder{}
						sb.WriteString("`")
						sb.WriteString(ctx.Prefix(msg.GuildID))
						for i, t := range v.Tokens {
							sb.WriteString(t.String())
							if i < len(v.Tokens)-1 {
								sb.WriteString(" ")
							}
						}
						sb.WriteString("`")
						embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
							Name:  sb.String(),
							Value: "-- " + v.Description,
						})
					}
				}
				return ctx.SendInfo(msg.ChannelID, embed)
			},
		})
	})
}
