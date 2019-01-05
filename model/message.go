package model

import (
	"time"

	"github.com/iamtakingiteasy/metabot/api"

	"github.com/lib/pq"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model/tmpl"
)

type EmbedField struct {
	Id      uint64      `db:"embeds_fields_id"`
	EmbedId uint64      `db:"embeds_fields_embed_id"`
	Name    string      `db:"embeds_fields_name"`
	Value   string      `db:"embeds_fields_value"`
	Inline  bool        `db:"embeds_fields_inline"`
	Created time.Time   `db:"embeds_fields_created"`
	Deleted pq.NullTime `db:"embeds_fields_deleted"`
}

type Embed struct {
	Id                 uint64      `db:"embeds_id"`
	MessageId          uint64      `db:"embeds_message_id"`
	Timestamp          string      `db:"embeds_timestamp"`
	Type               string      `db:"embeds_type"`
	Color              int         `db:"embeds_color"`
	Description        string      `db:"embeds_description"`
	Title              string      `db:"embeds_title"`
	Url                string      `db:"embeds_url"`
	AuthorName         string      `db:"embeds_author_name"`
	AuthorIconUrl      string      `db:"embeds_author_icon_url"`
	AuthorProxyIconUrl string      `db:"embeds_author_proxy_icon_url"`
	AuthorUrl          string      `db:"embeds_author_url"`
	ProviderName       string      `db:"embeds_provider_name"`
	ProviderUrl        string      `db:"embeds_provider_url"`
	ThumbnailWidth     int         `db:"embeds_thumbnail_width"`
	ThumbnailHeight    int         `db:"embeds_thumbnail_height"`
	ThumbnailUrl       string      `db:"embeds_thumbnail_url"`
	ThumbnailProxyUrl  string      `db:"embeds_thumbnail_proxy_url"`
	ImageWidth         int         `db:"embeds_image_width"`
	ImageHeight        int         `db:"embeds_image_height"`
	ImageUrl           string      `db:"embeds_image_url"`
	ImageProxyUrl      string      `db:"embeds_image_proxy_url"`
	VideoWidth         int         `db:"embeds_video_width"`
	VideoHeight        int         `db:"embeds_video_height"`
	VideoUrl           string      `db:"embeds_video_url"`
	VideoProxyUrl      string      `db:"embeds_video_proxy_url"`
	FooterText         string      `db:"embeds_footer_text"`
	FooterIconUrl      string      `db:"embeds_footer_icon_url"`
	FooterProxyIconUrl string      `db:"embeds_footer_proxy_icon_url"`
	Created            time.Time   `db:"embeds_created"`
	Deleted            pq.NullTime `db:"embeds_deleted"`
	Fields             []*EmbedField
}

type Message struct {
	Id               uint64      `db:"messages_id"`
	GuildDiscordId   string      `db:"messages_guild_discord_id"`
	ChannelDiscordId string      `db:"messages_channel_discord_id"`
	UserDiscordId    string      `db:"messages_user_discord_id"`
	DiscordId        string      `db:"messages_discord_id"`
	Type             int         `db:"messages_type"`
	WebhookDiscordId string      `db:"messages_webhook_discord_id"`
	Content          string      `db:"messages_content"`
	Created          time.Time   `db:"messages_created"`
	Deleted          pq.NullTime `db:"messages_deleted"`
	Embeds           []*Embed
}

var (
	queryEmbedsFieldsByEmbedId = tmpl.SelectTemplate("embeds_fields", "",
		"embed_id",
	)

	insertEmbedsFieldsRevision = tmpl.InsertTemplate("embeds_fields", "",
		"embed_id",
		"name",
		"value",
		"inline",
	)

	queryEmbedsByMessageId = tmpl.SelectTemplate("embeds", "",
		"messages_id",
	)
	insertEmbedsRevision = tmpl.InsertTemplate("embeds", "",
		"message_id",
		"timestamp",
		"type",
		"color",
		"description",
		"title",
		"url",
		"author_name",
		"author_icon_url",
		"author_proxy_icon_url",
		"author_url",
		"provider_name",
		"provider_url",
		"thumbnail_width",
		"thumbnail_height",
		"thumbnail_url",
		"thumbnail_proxy_url",
		"image_width",
		"image_height",
		"image_url",
		"image_proxy_url",
		"video_width",
		"video_height",
		"video_url",
		"video_proxy_url",
		"footer_text",
		"footer_icon_url",
		"footer_proxy_icon_url",
	)

	queryMessagesLastByGuildChannelDiscordId = tmpl.SelectTemplate("messages", "_last",
		"guild_discord_id",
		"channel_discord_id",
	)
	queryMessagesLastByGuildChannelMessageDiscordId = tmpl.SelectTemplate("messages", "_last",
		"guild_discord_id",
		"channel_discord_id",
		"discord_id",
	)
	queryMessagesRevisionsByGuildChannelMessageDiscordId = tmpl.SelectTemplate("messages", "",
		"guild_discord_id",
		"channel_discord_id",
		"discord_id",
	)
	insertMessagesRevision = tmpl.InsertTemplate("messages", "",
		"guild_discord_id",
		"channel_discord_id",
		"user_discord_id",
		"discord_id",
		"type",
		"webhook_discord_id",
		"content",
	)
	deleteMessagesByGuildChannelMessageDiscordId = tmpl.DeleteTemplate("messages", "",
		"guild_discord_id",
		"channel_discord_id",
		"discord_id",
	)
)

func QueryEmbedsFieldsByEmbed(ctx api.Context, embed *Embed) (*Embed, error) {
	var fields []*EmbedField
	rows, err := ctx.Database().NamedQuery(queryEmbedsFieldsByEmbedId, &EmbedField{EmbedId: embed.Id})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := &EmbedField{}
		err := rows.StructScan(f)
		if err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	embed.Fields = fields
	return embed, nil
}

func QueryEmbedsByMessage(ctx api.Context, message *Message) (*Message, error) {
	var embeds []*Embed
	rows, err := ctx.Database().NamedQuery(queryEmbedsByMessageId, &Embed{MessageId: message.Id})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		e := &Embed{}
		err := rows.StructScan(e)
		if err != nil {
			return nil, err
		}
		e, err = QueryEmbedsFieldsByEmbed(ctx, e)
		if err != nil {
			return nil, err
		}
		embeds = append(embeds, e)
	}
	message.Embeds = embeds
	return message, nil
}

func QueryMessagesLastByGuildChannelDiscordId(ctx api.Context, guildDiscordId, channelDiscordId string) ([]*Message, error) {
	var messages []*Message
	rows, err := ctx.Database().NamedQuery(queryMessagesLastByGuildChannelDiscordId, &Message{GuildDiscordId: guildDiscordId, ChannelDiscordId: channelDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := &Message{}
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		m, err = QueryEmbedsByMessage(ctx, m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func QueryMessagesLastByGuildChannelMessageDiscordId(ctx api.Context, guildDiscordId, channelDiscordId, messageDiscordId string) (*Message, error) {
	m := &Message{GuildDiscordId: guildDiscordId, ChannelDiscordId: channelDiscordId, DiscordId: messageDiscordId}
	rows, err := ctx.Database().NamedQuery(queryMessagesLastByGuildChannelMessageDiscordId, m)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		return QueryEmbedsByMessage(ctx, m)
	}
	return nil, ErrNoRows
}

func QueryMessagesRevisionsByGuildChannelMessageDiscordId(ctx api.Context, guildDiscordId, channelDiscordId, messageDiscordId string) ([]*Message, error) {
	var messages []*Message
	rows, err := ctx.Database().NamedQuery(queryMessagesRevisionsByGuildChannelMessageDiscordId, &Message{GuildDiscordId: guildDiscordId, ChannelDiscordId: channelDiscordId, DiscordId: messageDiscordId})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		m := &Message{}
		err := rows.StructScan(m)
		if err != nil {
			return nil, err
		}
		m, err = QueryEmbedsByMessage(ctx, m)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func InsertMessagesRevision(ctx api.Context, message *discordgo.Message) error {
	user := message.Author
	if user == nil {
		user = ctx.DiscordSession().State.User
	}
	m := &Message{
		GuildDiscordId:   message.GuildID,
		UserDiscordId:    user.ID,
		DiscordId:        message.ID,
		Type:             int(message.Type),
		WebhookDiscordId: message.WebhookID,
		Content:          message.Content,
	}
	rows, err := ctx.Database().NamedQuery(insertMessagesRevision, m)
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.StructScan(m)
		if err != nil {
			return err
		}
		for _, er := range message.Embeds {
			authorName := ""
			authorIconUrl := ""
			authorProxyIconUrl := ""
			authorUrl := ""

			if er.Author != nil {
				authorName = er.Author.Name
				authorIconUrl = er.Author.IconURL
				authorProxyIconUrl = er.Author.ProxyIconURL
				authorUrl = er.Author.URL
			}

			providerName := ""
			providerUrl := ""

			if er.Provider != nil {
				providerName = er.Provider.Name
				providerUrl = er.Provider.URL
			}

			thumbnailWidth := 0
			thumbnailHeight := 0
			thumbnailUrl := ""
			thumbnailProxyUrl := ""

			if er.Thumbnail != nil {
				thumbnailWidth = er.Thumbnail.Width
				thumbnailHeight = er.Thumbnail.Height
				thumbnailUrl = er.Thumbnail.URL
				thumbnailProxyUrl = er.Thumbnail.ProxyURL
			}

			imageWidth := 0
			imageHeight := 0
			imageUrl := ""
			imageProxyUrl := ""

			if er.Image != nil {
				imageWidth = er.Image.Width
				imageHeight = er.Image.Height
				imageUrl = er.Image.URL
				imageProxyUrl = er.Image.ProxyURL
			}

			videoWidth := 0
			videoHeight := 0
			videoUrl := ""
			videoProxyUrl := ""

			if er.Video != nil {
				videoWidth = er.Video.Width
				videoHeight = er.Video.Height
				videoUrl = er.Video.URL
				videoProxyUrl = er.Video.ProxyURL
			}

			footerText := ""
			footerIconUrl := ""
			footerProxyIconUrl := ""

			if er.Footer != nil {
				footerText = er.Footer.Text
				footerIconUrl = er.Footer.IconURL
				footerProxyIconUrl = er.Footer.ProxyIconURL
			}

			e := &Embed{
				MessageId:          m.Id,
				Timestamp:          er.Timestamp,
				Type:               er.Type,
				Color:              er.Color,
				Description:        er.Description,
				Title:              er.Title,
				Url:                er.URL,
				AuthorName:         authorName,
				AuthorIconUrl:      authorIconUrl,
				AuthorProxyIconUrl: authorProxyIconUrl,
				AuthorUrl:          authorUrl,
				ProviderName:       providerName,
				ProviderUrl:        providerUrl,
				ThumbnailWidth:     thumbnailWidth,
				ThumbnailHeight:    thumbnailHeight,
				ThumbnailUrl:       thumbnailUrl,
				ThumbnailProxyUrl:  thumbnailProxyUrl,
				ImageWidth:         imageWidth,
				ImageHeight:        imageHeight,
				ImageUrl:           imageUrl,
				ImageProxyUrl:      imageProxyUrl,
				VideoWidth:         videoWidth,
				VideoHeight:        videoHeight,
				VideoUrl:           videoUrl,
				VideoProxyUrl:      videoProxyUrl,
				FooterText:         footerText,
				FooterIconUrl:      footerIconUrl,
				FooterProxyIconUrl: footerProxyIconUrl,
			}
			rows, err := ctx.Database().NamedQuery(insertEmbedsRevision, e)
			if err != nil {
				return err
			}
			for rows.Next() {
				err := rows.StructScan(e)
				if err != nil {
					return err
				}
				for _, fr := range e.Fields {
					_, err := ctx.Database().NamedQuery(insertEmbedsFieldsRevision, &EmbedField{
						EmbedId: e.Id,
						Name:    fr.Name,
						Value:   fr.Value,
						Inline:  fr.Inline,
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func DeleteMessagesByGuildChannelMessageDiscordId(ctx api.Context, guildDiscordId, channelDiscordId, messageDiscordId string) error {
	_, err := ctx.Database().NamedQuery(deleteMessagesByGuildChannelMessageDiscordId, &Message{GuildDiscordId: guildDiscordId, ChannelDiscordId: channelDiscordId, DiscordId: messageDiscordId})
	return err
}
