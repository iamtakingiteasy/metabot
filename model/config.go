package model

import (
	"time"

	"github.com/iamtakingiteasy/metabot/bot"
)

const (
	ConfigAdminTypeUnknown = iota
	ConfigAdminTypePermissionMask
	ConfigAdminTypeRoleDiscordId
	ConfigAdminTypeUserDiscordId
)

const (
	ConfigRestrictRuleEventUnknown = iota
	ConfigRestrictRuleEventGuildJoin
	ConfigRestrictRuleEventVoiceJoin
	ConfigRestrictRuleEventMessageSend
)
const (
	ConfigRestrictRuleActionUnknown = iota
	ConfigRestrictRuleActionAddRole
	ConfigRestrictRuleActionMute
	ConfigRestrictRuleActionKick
	ConfigRestrictRuleActionBan
)

type ConfigRestrictRule struct {
	Id                           uint64        `db:"config_restrict_rules_id"`
	GuildDiscordId               string        `db:"config_restrict_rules_guild_discord_id"`
	Event                        int           `db:"config_restrict_rules_event"`
	SelectorUserDiscordId        string        `db:"config_restrict_rules_selector_user_discord_id"`
	SelectorChannelDiscordId     string        `db:"config_restrict_rules_selector_channel_discord_id"`
	SelectorUserPattern          string        `db:"config_restrict_rules_selector_user_pattern"`
	SelectorChannelPattern       string        `db:"config_restrict_rules_selector_channel_pattern"`
	Action                       int           `db:"config_restrict_rules_action"`
	ActionNotify                 bool          `db:"config_restrict_rules_action_notify"`
	ActionNotifyChannelDiscordId string        `db:"config_restrict_rules_action_notify_channel_discord_id"`
	ActionTime                   time.Duration `db:"config_restrict_rules_action_time"`
	ActionRoleDiscordId          string        `db:"config_restrict_rules_action_role_discord_id"`
}

type ConfigAdmin struct {
	Id             uint64 `db:"config_admins_id"`
	GuildDiscordId string `db:"config_admins_guild_discord_id"`
	Type           int    `db:"config_admins_type"`
	PermissionMask int    `db:"config_admins_permission_mask"`
	RoleDiscordId  string `db:"config_admins_role_discord_id"`
	UserDiscordId  string `db:"config_admins_user_discord_id"`
}

type Config struct {
	GuildDiscordId   string        `db:"config_guild_discord_id"`
	Prefix           string        `db:"config_prefix"`
	ColorError       int           `db:"config_color_error"`
	ColorWarn        int           `db:"config_color_warn"`
	ColorInfo        int           `db:"config_color_info"`
	AutoremoveActive bool          `db:"config_autoremove_active"`
	AutoremoveTime   time.Duration `db:"config_autoremove_time"`
	RestrictActive   bool          `db:"config_restrict_active"`
	Admins           []*ConfigAdmin
	RestrictRules    []*ConfigRestrictRule
}

const (
	loadConfig = `
select c.* from config c
`
	saveConfig = `
insert into config (
  config_guild_discord_id,
  config_prefix,
  config_color_error,
  config_color_warn,
  config_color_info,
  config_autoremove_active,
  config_autoremove_time,
  config_restrict_active
) values (
  :config_guild_discord_id,
  :config_prefix,
  :config_color_error,
  :config_color_warn,
  :config_color_info,
  :config_autoremove_active,
  :config_autoremove_time,
  :config_restrict_active
) on conflict (config_guild_discord_id) do update set 
  config_prefix = :config_prefix,
  config_color_error = :config_color_error,
  config_color_warn = :config_color_warn,
  config_color_info = :config_color_info,
  config_autoremove_active = :config_autoremove_active,
  config_autoremove_time = :config_autoremove_time,
  config_restrict_active = :config_restrict_active
where config.config_guild_discord_id = :config_guild_discord_id
`
	listConfigAdminsByGuildDiscordId = `
select ca.* from config_admins ca where ca.config_admins_guild_discord_id = :config_admins_guild_discord_id order by ca.config_admins_id asc
`
	insertConfigAdmins = `
insert into config_admins (
  config_admins_guild_discord_id,
  config_admins_type,
  config_admins_permission_mask,
  config_admins_role_discord_id,
  config_admins_user_discord_id
) select
    :config_admins_guild_discord_id,
    :config_admins_type,
    :config_admins_permission_mask,
    :config_admins_role_discord_id,
    :config_admins_user_discord_id
  where (
    select count(*) from config_admins where
      config_admins_guild_discord_id = :config_admins_guild_discord_id and
      config_admins_type = :config_admins_type and
      config_admins_permission_mask = :config_admins_permission_mask and
      config_admins_role_discord_id = :config_admins_role_discord_id and
      config_admins_user_discord_id = :config_admins_user_discord_id
  ) = 0
returning *
`
	updateConfigAdminsById = `
update config_admins set
  config_admins_type = :config_admins_type,
  config_admins_permission_mask = :config_admins_permission_mask,
  config_admins_role_discord_id = :config_admins_role_discord_id,
  config_admins_user_discord_id = :config_admins_user_discord_id
where config_admins_id = :config_admins_id
`
	deleteConfigAdminsById = `
delete from config_admins where config_admins_id = :config_admins_id
`
	listConfigRestrictRulesByGuildDiscordId = `
select crr.* from config_restrict_rules crr where crr.config_restrict_rules_guild_discord_id = :config_restrict_rules_guild_discord_id order by crr.config_restrict_rules_id asc
`
	insertConfigRestrictRules = `
insert into config_restrict_rules (
  config_restrict_rules_guild_discord_id,
  config_restrict_rules_event,
  config_restrict_rules_selector_user_discord_id,
  config_restrict_rules_selector_channel_discord_id,
  config_restrict_rules_selector_user_pattern,
  config_restrict_rules_selector_channel_pattern,
  config_restrict_rules_action,
  config_restrict_rules_action_notify,
  config_restrict_rules_action_notify_channel_discord_id,
  config_restrict_rules_action_time,
  config_restrict_rules_action_role_discord_id
) select
    :config_restrict_rules_guild_discord_id,
    :config_restrict_rules_event,
    :config_restrict_rules_selector_user_discord_id,
    :config_restrict_rules_selector_channel_discord_id,
    :config_restrict_rules_selector_user_pattern,
    :config_restrict_rules_selector_channel_pattern,
    :config_restrict_rules_action,
    :config_restrict_rules_action_notify,
    :config_restrict_rules_action_notify_channel_discord_id,
    :config_restrict_rules_action_time,
    :config_restrict_rules_action_role_discord_id
  where (
    select count(*) from config_restrict_rules where
      config_restrict_rules_guild_discord_id = :config_restrict_rules_guild_discord_id and 
      config_restrict_rules_event = :config_restrict_rules_event and
      config_restrict_rules_selector_user_discord_id = :config_restrict_rules_selector_user_discord_id and
      config_restrict_rules_selector_channel_discord_id = :config_restrict_rules_selector_channel_discord_id and
      config_restrict_rules_selector_user_pattern = :config_restrict_rules_selector_user_pattern and
      config_restrict_rules_selector_channel_pattern = :config_restrict_rules_selector_channel_pattern and
      config_restrict_rules_action = :config_restrict_rules_action and
      config_restrict_rules_action_notify = :config_restrict_rules_action_notify and
      config_restrict_rules_action_notify_channel_discord_id = :config_restrict_rules_action_notify_channel_discord_id and
      config_restrict_rules_action_time = :config_restrict_rules_action_time and
      config_restrict_rules_action_role_discord_id = config_restrict_rules_action_role_discord_id
  ) = 0
returning *
`
	updateConfigRestrictRulesById = `
update config_restrict_rules set 
  config_restrict_rules_event = :config_restrict_rules_event,
  config_restrict_rules_selector_user_discord_id = :config_restrict_rules_selector_user_discord_id,
  config_restrict_rules_selector_channel_discord_id = :config_restrict_rules_selector_channel_discord_id,
  config_restrict_rules_selector_user_pattern = :config_restrict_rules_selector_user_pattern,
  config_restrict_rules_selector_channel_pattern = :config_restrict_rules_selector_channel_pattern,
  config_restrict_rules_action = :config_restrict_rules_action,
  config_restrict_rules_action_notify = :config_restrict_rules_action_notify,
  config_restrict_rules_action_notify_channel_discord_id = :config_restrict_rules_action_notify_channel_discord_id,
  config_restrict_rules_action_time = :config_restrict_rules_action_time,
  config_restrict_rules_action_role_discord_id = config_restrict_rules_action_role_discord_id
where config_restrict_rules_id = :config_restrict_rules_id
`
	deleteConfigRestrictRulesById = `
delete from config_restrict_rules where config_restrict_rules_id = :config_restrict_rules_id
`
)

func LoadConfigs(ctx bot.Context) ([]*Config, error) {
	var configs []*Config
	rows, err := ctx.Database().NamedQuery(loadConfig, &Config{})
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		c := &Config{}
		err := rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		err = c.LoadAdmins(ctx)
		if err != nil {
			return nil, err
		}
		err = c.LoadRestrictRules(ctx)
		if err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, nil
}

func (conf *Config) Save(ctx bot.Context) error {
	_, err := ctx.Database().NamedQuery(saveConfig, conf)
	return err
}

func (conf *Config) LoadAdmins(ctx bot.Context) error {
	var admins []*ConfigAdmin
	rows, err := ctx.Database().NamedQuery(listConfigAdminsByGuildDiscordId, &ConfigAdmin{GuildDiscordId: conf.GuildDiscordId})
	if err != nil {
		return err
	}
	for rows.Next() {
		a := &ConfigAdmin{}
		err := rows.StructScan(a)
		if err != nil {
			return err
		}
		admins = append(admins, a)
	}
	conf.Admins = admins
	return nil
}

func (conf *Config) AddAdmin(ctx bot.Context, admin *ConfigAdmin) error {
	admin.GuildDiscordId = conf.GuildDiscordId
	rows, err := ctx.Database().NamedQuery(insertConfigAdmins, admin)
	if err != nil {
		return err
	}
	for rows.Next() {
		a := &ConfigAdmin{}
		err := rows.StructScan(a)
		if err != nil {
			return err
		}
		conf.Admins = append(conf.Admins, a)
	}
	return nil
}

func (admin *ConfigAdmin) Save(ctx bot.Context) error {
	_, err := ctx.Database().NamedQuery(updateConfigAdminsById, admin)
	return err
}

func (conf *Config) DeleteAdminById(ctx bot.Context, adminId uint64) error {
	_, err := ctx.Database().NamedQuery(deleteConfigAdminsById, &ConfigAdmin{Id: adminId})
	if err != nil {
		return err
	}
	for i, a := range conf.Admins {
		if a.Id == adminId {
			conf.Admins = append(conf.Admins[:i], conf.Admins[i+1:]...)
			break
		}
	}
	return nil
}

func (conf *Config) LoadRestrictRules(ctx bot.Context) error {
	var rules []*ConfigRestrictRule
	rows, err := ctx.Database().NamedQuery(listConfigRestrictRulesByGuildDiscordId, &ConfigRestrictRule{GuildDiscordId: conf.GuildDiscordId})
	if err != nil {
		return err
	}
	for rows.Next() {
		r := &ConfigRestrictRule{}
		err := rows.StructScan(r)
		if err != nil {
			return err
		}
		rules = append(rules, r)
	}
	conf.RestrictRules = rules
	return nil
}

func (conf *Config) AddRestrictRule(ctx bot.Context, rule *ConfigRestrictRule) error {
	rule.GuildDiscordId = conf.GuildDiscordId
	rows, err := ctx.Database().NamedQuery(insertConfigRestrictRules, rule)
	if err != nil {
		return err
	}
	for rows.Next() {
		r := &ConfigRestrictRule{}
		err := rows.StructScan(r)
		if err != nil {
			return err
		}
		conf.RestrictRules = append(conf.RestrictRules, r)
	}
	return nil
}

func (rule *ConfigRestrictRule) Save(ctx bot.Context) error {
	_, err := ctx.Database().NamedQuery(updateConfigRestrictRulesById, rule)
	return err
}

func (conf *Config) DeleteRestrictRuleById(ctx bot.Context, ruleId uint64) error {
	_, err := ctx.Database().NamedQuery(deleteConfigRestrictRulesById, &ConfigRestrictRule{Id: ruleId})
	if err != nil {
		return err
	}
	for i, a := range conf.RestrictRules {
		if a.Id == ruleId {
			conf.RestrictRules = append(conf.RestrictRules[:i], conf.RestrictRules[i+1:]...)
			break
		}
	}
	return nil
}
