package bot

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/iamtakingiteasy/metabot/model"
	"github.com/jmoiron/sqlx"
)

type EventHandlerFunc func(ctx *Context, raw interface{}) error

var GlobalEventHandlers []EventHandlerFunc
var GlobalRegistrars []RegistrarFunc

func AddGlobalEventHandler(handler EventHandlerFunc) {
	GlobalEventHandlers = append(GlobalEventHandlers, handler)
}

func AddGlobalRegistrars(registrar RegistrarFunc) {
	GlobalRegistrars = append(GlobalRegistrars, registrar)
}

type ConfigFile struct {
	Filename string
	Data     Config
}

type Context struct {
	Configs       map[string]*model.Config
	ConfigFile    *ConfigFile
	Session       *discordgo.Session
	Commands      Parser
	Events        chan interface{}
	EventHandlers []EventHandlerFunc
	Db            *sqlx.DB
}

func (ctx *Context) Database() *sqlx.DB {
	return ctx.Db
}

func (ctx *Context) DiscordSession() *discordgo.Session {
	return ctx.Session
}

func (ctx *Context) Connect() error {
	if ctx.ConfigFile.Data.Token == "" {
		return errors.New("bot token is empty, check " + ctx.ConfigFile.Filename)
	}

	if ctx.ConfigFile.Data.Database.Host == "" || ctx.ConfigFile.Data.Database.Port == 0 {
		return errors.New("database host is empty or port is 0, check " + ctx.ConfigFile.Filename)
	}

	dataSourceName := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		ctx.ConfigFile.Data.Database.Host,
		ctx.ConfigFile.Data.Database.Port,
		ctx.ConfigFile.Data.Database.Name,
		ctx.ConfigFile.Data.Database.User,
		ctx.ConfigFile.Data.Database.Password,
	)

	var err error

	log.Println("Connecting to database...")
	ctx.Db, err = sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return err
	}

	err = model.Migrate(ctx.Db)
	if err != nil {
		return err
	}

	configs, err := model.LoadConfigs(ctx)
	if err != nil {
		return err
	}

	for _, c := range configs {
		ctx.Configs[c.GuildDiscordId] = c
	}

	session, err := discordgo.New("Bot " + ctx.ConfigFile.Data.Token)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("couldn't connect")
	}
	ctx.Session = session

	ctx.Session.AddHandler(func(s *discordgo.Session, raw interface{}) {
		ctx.Dispatch(raw)
	})

	err = ctx.Session.Open()
	if err != nil {
		return err
	}

	go ctx.Process()
	return nil
}

func (ctx *Context) Stop() {
	_ = ctx.Session.Close()
}

func (ctx *Context) Dispatch(event interface{}) {
	ctx.Events <- event
}

func (ctx *Context) Process() {
	for raw := range ctx.Events {
		for _, h := range ctx.EventHandlers {
			err := h(ctx, raw)
			if err != nil {
				log.Println(reflect.TypeOf(raw).String(), "warn", err)
			}
		}
	}
}

func NewContext(filename string) *Context {
	ctx := &Context{
		Configs: make(map[string]*model.Config),
		ConfigFile: &ConfigFile{
			Filename: filename,
		},
		Events:        make(chan interface{}),
		EventHandlers: GlobalEventHandlers,
	}
	ctx.Commands = NewParser(ctx)
	for _, r := range GlobalRegistrars {
		err := ctx.Commands.RegisterCommand(r)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return ctx
}
