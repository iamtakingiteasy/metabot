package impl

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

type RegistrarFunc func(ctx *Context, registrar Registrar) error
type HandlerFunc func(msg *discordgo.Message) error
type ParserFunc func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error)

type Parser interface {
	ProcessCommand(message *discordgo.Message) error
	RegisterCommand(registrar RegistrarFunc) error
	Commands() []*Descriptor
	Stop()
}

type Registrar interface {
	AddCommand(descriptor *Descriptor) error
}

type Descriptor struct {
	Group   string
	Name    string
	Tokens  []*Token
	Handler HandlerFunc
}

type Token struct {
	Parser    ParserFunc
	Parameter string
	Value     interface{}
	Minimum   int
	Maximum   int
}

type ParserEventAdd struct {
	Handler RegistrarFunc
	Ret     chan error
}

type ParserEventProcess struct {
	Raw *discordgo.Message
	Ret chan error
}

type ParserRegistrar struct {
	Descriptors map[string][]*Descriptor
	Control     chan interface{}
	Context     *Context
}

var (
	TypeLiteral ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		if token != parameter {
			return nil, errors.New(fmt.Sprintf("`%s` not matched expected `%s`", token, parameter))
		}
		return token, nil
	}

	TypeString ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		if parameter != "" {
			patt, err := regexp.Compile(parameter)
			if err != nil {
				return nil, err
			}
			if !patt.MatchString(token) {
				return nil, errors.New("does not match pattern")
			}
		}
		return token, nil
	}

	TypeNumber ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		return strconv.ParseFloat(token, 64)
	}

	TypeBool ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		return strconv.ParseBool(token)
	}

	TypePattern ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		_, e := regexp.Compile(token)
		return token, e
	}

	TypeUser ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		for _, u := range msg.Mentions {
			if token == "<@"+u.ID+">" || token == "<@!"+u.ID+">" {
				return u.ID, nil
			}
		}

		u, err := ctx.Session.User(token)
		if err == nil {
			return u.ID, nil
		}

		g, err := ctx.Session.Guild(msg.GuildID)
		if err != nil {
			return nil, err
		}

		for _, m := range g.Members {
			name := m.Nick
			if name == "" {
				name = m.User.Username
			}
			if name == token {
				return m.User.ID, nil
			}
		}

		return nil, errors.New("user name not found")
	}

	TypeRole ParserFunc = func(ctx *Context, msg *discordgo.Message, parameter, token string) (interface{}, error) {
		for _, rid := range msg.MentionRoles {
			if token == "<@&"+rid+">" {
				role, err := ctx.Session.State.Role(msg.GuildID, rid)
				if err != nil {
					return nil, err
				}
				return role.ID, nil
			}
		}

		role, err := ctx.Session.State.Role(msg.GuildID, token)
		if err == nil {
			return role.ID, nil
		}

		roles, err := ctx.Session.GuildRoles(msg.GuildID)
		if err != nil {
			return nil, err
		}

		for _, r := range roles {
			if r.Name == token {
				return r.ID, nil
			}
		}

		return nil, errors.New("role name not found")
	}
)

func NewParser(context *Context) Parser {
	pr := &ParserRegistrar{
		Descriptors: make(map[string][]*Descriptor),
		Control:     make(chan interface{}),
		Context:     context,
	}

	go pr.Run()

	return pr
}

func (pr *ParserRegistrar) Run() {
	for raw := range pr.Control {
		switch evt := raw.(type) {
		case *ParserEventAdd:
			evt.Ret <- evt.Handler(pr.Context, pr)
		case *ParserEventProcess:
			evt.Ret <- pr.ProcessInternal(evt.Raw)
		}
	}
}

func (pr *ParserRegistrar) Stop() {
	close(pr.Control)
}

func ClearValue(token *Token) {
	v := reflect.Indirect(reflect.ValueOf(token.Value))
	z := reflect.Zero(v.Type())
	v.Set(z)
}

func SetValue(token *Token, value interface{}) {
	v := reflect.Indirect(reflect.ValueOf(token.Value))
	tv := reflect.ValueOf(value)
	t := v.Type()
	if v.Kind() == reflect.Slice {
		v.Set(reflect.Append(v, tv.Convert(t.Elem())))
	} else {
		v.Set(tv.Convert(t))
	}
}

func (pr *ParserRegistrar) ProcessInternal(msg *discordgo.Message) error {
	rawtokens, err := RawTokenize(msg.Content)
	if err != nil {
		return err
	}

	for _, v := range pr.Descriptors {
		for _, d := range v {
			n := 0
			last := 0
		tokenloop:
			for l, t := range d.Tokens {
				var min, max int
				if t.Minimum > 0 {
					min = t.Minimum
				} else if t.Maximum == 0 {
					min = 1
				}

				if t.Maximum < 0 {
					max = math.MaxInt32
				} else if t.Maximum > min {
					max = t.Maximum
				} else {
					max = min
				}

				ClearValue(t)
				for x := 0; x < min; x++ {
					if n >= len(rawtokens) {
						break tokenloop
					}
					res, err := t.Parser(pr.Context, msg, t.Parameter, rawtokens[n])
					if err != nil {
						break tokenloop
					}
					SetValue(t, res)
					n++
				}

				last = l
				for x := min; x < max; x++ {
					if n >= len(rawtokens) {
						continue tokenloop
					}
					res, err := t.Parser(pr.Context, msg, t.Parameter, rawtokens[n])
					if err != nil {
						continue tokenloop
					}
					SetValue(t, res)
					n++
				}
			}
			if last == len(d.Tokens)-1 {
				return d.Handler(msg)
			}
		}
	}

	return errors.New(fmt.Sprintf("no matching command `%s`", msg.Content))
}

func (pr *ParserRegistrar) ProcessCommand(msg *discordgo.Message) error {
	ret := make(chan error)
	defer close(ret)

	pr.Control <- &ParserEventProcess{
		Raw: msg,
		Ret: ret,
	}
	return <-ret
}

func (pr *ParserRegistrar) RegisterCommand(handler RegistrarFunc) error {
	ret := make(chan error)
	defer close(ret)

	pr.Control <- &ParserEventAdd{
		Handler: handler,
		Ret:     ret,
	}
	return <-ret
}

func (pr *ParserRegistrar) Commands() []*Descriptor {
	var descs []*Descriptor
	for _, v := range pr.Descriptors {
		descs = append(descs, v...)
	}
	return descs
}

func (pr *ParserRegistrar) AddCommand(descriptor *Descriptor) error {
	for i, t := range descriptor.Tokens {
		if t.Minimum > 0 || t.Maximum > 0 {
			if reflect.Indirect(reflect.ValueOf(t.Value)).Kind() != reflect.Slice {
				return errors.New(fmt.Sprintf("token %d of %s in group %s is having multiple values, but value is not a slice", i, descriptor.Name, descriptor.Group))
			}
		}
	}
	pr.Descriptors[descriptor.Group] = append(pr.Descriptors[descriptor.Group], descriptor)
	return nil
}

func RawTokenize(raw string) (result []string, err error) {
	const (
		space = iota
		value
		quote
	)

	var token []rune
	runes := ([]rune)(raw)
	escaping := false
	state := space

	escapeadd := func(r rune) {
		if r == '\\' {
			if escaping {
				token = append(token, r)
				escaping = false
			} else {
				escaping = true
			}
		} else {
			if escaping {
				token = append(token, '\\')
				escaping = false
			}
			token = append(token, r)
		}
	}

	tokenadd := func() {
		result = append(result, string(token))
		token = nil
	}

	for _, r := range runes {
		switch state {
		case space:
			if !unicode.IsSpace(r) {
				switch r {
				case '\\':
					escaping = true
					state = value
				case '"':
					state = quote
				default:
					token = append(token, r)
					state = value
				}
			}
		case value:
			if unicode.IsSpace(r) {
				if escaping {
					escaping = false
					token = append(token, r)
				} else {
					state = space
					tokenadd()
				}
			} else {
				escapeadd(r)
			}
		case quote:
			if r == '"' {
				if escaping {
					token = append(token, '"')
				} else {
					state = space
					tokenadd()
				}
				escaping = false
			} else {
				escapeadd(r)
			}
		}
	}

	if escaping {
		token = append(token, '\\')
	}

	if len(token) > 0 {
		if state == quote && token[len(token)-1] != '"' {
			return nil, errors.New(fmt.Sprintf("mismatched quotes in `%s`", string(token)))
		}

		result = append(result, string(token))
	}
	return
}
