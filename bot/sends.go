package bot

import "github.com/bwmarrin/discordgo"

func (ctx *Context) Prefix(guildId string) string {
	return ctx.Configs[guildId].Prefix
}

func (ctx *Context) Send(channelId, content string) error {
	_, err := ctx.Session.ChannelMessageSend(channelId, content)
	return err
}

func (ctx *Context) SendInfo(channelId string, embed *discordgo.MessageEmbed) error {
	ch, err := ctx.Session.Channel(channelId)
	if err != nil {
		return err
	}
	embed.Color = ctx.Configs[ch.GuildID].ColorInfo
	_, err = ctx.Session.ChannelMessageSendEmbed(channelId, embed)
	return err
}

func (ctx *Context) SendWarn(channelId string, embed *discordgo.MessageEmbed) error {
	ch, err := ctx.Session.Channel(channelId)
	if err != nil {
		return err
	}
	embed.Color = ctx.Configs[ch.GuildID].ColorWarn
	_, err = ctx.Session.ChannelMessageSendEmbed(channelId, embed)
	return err
}

func (ctx *Context) SendError(channelId string, embed *discordgo.MessageEmbed) error {
	ch, err := ctx.Session.Channel(channelId)
	if err != nil {
		return err
	}
	embed.Color = ctx.Configs[ch.GuildID].ColorError
	_, err = ctx.Session.ChannelMessageSendEmbed(channelId, embed)
	return err
}
