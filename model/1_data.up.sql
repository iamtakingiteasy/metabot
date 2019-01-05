create table config
(
  config_guild_discord_id text primary key,
  config_prefix text,
  config_color_error bigint,
  config_color_warn bigint,
  config_color_info bigint,
  config_autoremove_active boolean,
  config_autoremove_time bigint,
  config_restrict_active boolean
);

create table config_admins
(
  config_admins_id bigserial primary key,
  config_admins_guild_discord_id text,
  config_admins_type bigint,
  config_admins_permission_mask bigint,
  config_admins_role_discord_id text,
  config_admins_user_discord_id text
);

create table config_restrict_rules
(
  config_restrict_rules_id bigserial primary key,
  config_restrict_rules_guild_discord_id text,
  config_restrict_rules_event bigint,
  config_restrict_rules_selector_user_discord_id text,
  config_restrict_rules_selector_channel_discord_id text,
  config_restrict_rules_selector_user_pattern text,
  config_restrict_rules_selector_channel_pattern text,
  config_restrict_rules_action bigint,
  config_restrict_rules_action_notify boolean,
  config_restrict_rules_action_notify_channel_discord_id text,
  config_restrict_rules_action_time bigint,
  config_restrict_rules_action_role_discord_id text
);

create table guilds
(
  guilds_id bigserial primary key,
  guilds_discord_id text,
  guilds_name text,
  guilds_image bytea,
  guilds_splash bytea,
  guilds_created timestamp default now(),
  guilds_deleted timestamp
);

create view guilds_last as
select a.* from guilds a left join guilds b on
      a.guilds_discord_id = b.guilds_discord_id and
      a.guilds_id < b.guilds_id
where b.guilds_id is null;

create table roles
(
  roles_id bigserial primary key,
  roles_guild_discord_id text,
  roles_discord_id text,
  roles_name text,
  roles_color bigint,
  roles_position bigint,
  roles_permissions bigint,
  roles_created timestamp default now(),
  roles_deleted timestamp
);

create view roles_last as
select a.* from roles a left join roles b on
      a.roles_discord_id = b.roles_discord_id and
      a.roles_id < b.roles_id
where b.roles_id is null;

create table users
(
  users_id bigserial primary key,
  users_discord_id text,
  users_name text,
  users_discriminator text,
  users_bot boolean,
  users_email text,
  users_avatar bytea,
  users_locale text,
  users_verified boolean,
  users_created timestamp default now(),
  users_deleted timestamp
);

create view users_last as
select a.* from users a left join users b on
      a.users_discord_id = b.users_discord_id and
      a.users_id < b.users_id
where b.users_id is null;

create table members
(
  members_id bigserial primary key,
  members_guild_discord_id text,
  members_user_discord_id text,
  members_nick text,
  members_joined_at timestamp,
  members_created timestamp default now(),
  members_deleted timestamp
);

create view members_last as
select a.* from members a left join members b on
      a.members_guild_discord_id = b.members_guild_discord_id and
      a.members_user_discord_id = b.members_user_discord_id and
      a.members_id < b.members_id
where b.members_id is null;

create table channels
(
  channels_id bigserial primary key,
  channels_guild_discord_id text,
  channels_discord_id text,
  channels_name text,
  channels_type bigint,
  channels_bitrate bigint,
  channels_parent_discord_id text,
  channels_position bigint,
  channels_nsfw boolean,
  channels_topic text,
  channels_user_limit bigint,
  channels_created timestamp default now(),
  channels_deleted timestamp
);

create view channels_last as
select a.* from channels a left join channels b on
      a.channels_discord_id = b.channels_discord_id and
      a.channels_id < b.channels_id
where b.channels_id is null;

create table guilds_voice_status
(
  guilds_voice_status_id bigserial primary key,
  guilds_voice_status_time timestamp default now(),
  guilds_voice_status_guild_discord_id text,
  guilds_voice_status_channel_discord_id text,
  guilds_voice_status_user_discord_id text,
  guilds_voice_status_created timestamp default now(),
  guilds_voice_status_deleted timestamp
);

create view guilds_voice_status_last as
select a.* from guilds_voice_status a left join guilds_voice_status b on
      a.guilds_voice_status_guild_discord_id = b.guilds_voice_status_guild_discord_id and
      a.guilds_voice_status_user_discord_id = b.guilds_voice_status_user_discord_id and
      a.guilds_voice_status_id < b.guilds_voice_status_id
where b.guilds_voice_status_id is null;

create table messages
(
  messages_id bigserial primary key,
  messages_guild_discord_id text,
  messages_channel_discord_id text,
  messages_user_discord_id text,
  messages_discord_id text,
  messages_type bigint,
  messages_webhook_discord_id text,
  messages_content text,
  messages_created timestamp default now(),
  messages_deleted timestamp
);

create view messages_last as
select a.* from messages a left join messages b on
      a.messages_guild_discord_id = b.messages_guild_discord_id and
      a.messages_channel_discord_id = b.messages_channel_discord_id and
      a.messages_user_discord_id = b.messages_user_discord_id and
      a.messages_discord_id = b.messages_discord_id and
      a.messages_id < b.messages_id
where b.messages_id is null;

create table embeds
(
  embeds_id bigserial primary key,
  embeds_message_id bigint references messages(messages_id),
  embeds_timestamp text,
  embeds_type text,
  embeds_color bigint,
  embeds_description text,
  embeds_title text,
  embeds_url text,
  embeds_author_name text,
  embeds_author_icon_url text,
  embeds_author_proxy_icon_url text,
  embeds_author_url text,
  embeds_provider_name text,
  embeds_provider_url text,
  embeds_thumbnail_width bigint,
  embeds_thumbnail_height bigint,
  embeds_thumbnail_url text,
  embeds_thumbnail_proxy_url text,
  embeds_image_width bigint,
  embeds_image_height bigint,
  embeds_image_url text,
  embeds_image_proxy_url text,
  embeds_video_width bigint,
  embeds_video_height bigint,
  embeds_video_url text,
  embeds_video_proxy_url text,
  embeds_footer_text text,
  embeds_footer_icon_url text,
  embeds_footer_proxy_icon_url text,
  embeds_created timestamp default now(),
  embeds_deleted timestamp
);

create table embeds_fields
(
  embeds_fields_id bigserial primary key,
  embeds_fields_embed_id bigint references embeds(embeds_id),
  embeds_fields_name text,
  embeds_fields_value text,
  embeds_fields_inline boolean,
  embeds_fields_created timestamp default now(),
  embeds_fields_deleted timestamp
);
