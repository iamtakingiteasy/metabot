package bot

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

var ErrInvalidPrefix = errors.New("invalid prefix")

type RegistrarFunc func(ctx *Context, registrar Registrar) error
type HandlerFunc func(msg *discordgo.Message, tokens []string) error

type Token interface {
	Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error)
	ClearValue()
	SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error
	Minimum() int
	Maximum() int
	fmt.Stringer
}

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
	Group       string
	Description string
	Tokens      []Token
	Handler     HandlerFunc
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

type TypeLiteral string

func (parser TypeLiteral) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	if token != string(parser) {
		return nil, errors.New(fmt.Sprintf("`%s` not matched expected `%s`", token, parser))
	}
	return token, nil
}

func (parser TypeLiteral) ClearValue() {

}

func (parser TypeLiteral) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return nil
}

func (parser TypeLiteral) Minimum() int {
	return 1
}

func (parser TypeLiteral) Maximum() int {
	return 1
}

func (parser TypeLiteral) String() string {
	return string(parser)
}

type TypeString struct {
	Pattern string
	Value   interface{}
}

func (parser *TypeString) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	if parser.Pattern != "" {
		patt, err := regexp.Compile(parser.Pattern)
		if err != nil {
			return nil, err
		}
		if !patt.MatchString(token) {
			return nil, errors.New("does not match pattern")
		}
	}
	return token, nil
}

func (parser *TypeString) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeString) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeString) Minimum() int {
	return 1
}

func (parser *TypeString) Maximum() int {
	return 1
}

func (parser *TypeString) String() string {
	if parser.Pattern != "" {
		return ":string(" + parser.Pattern + ")"
	}
	return ":string"
}

type TypeNumber struct {
	Value interface{}
}

func (parser *TypeNumber) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	return strconv.ParseFloat(token, 64)
}

func (parser *TypeNumber) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeNumber) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeNumber) Minimum() int {
	return 1
}

func (parser *TypeNumber) Maximum() int {
	return 1
}

func (parser *TypeNumber) String() string {
	return ":number"
}

type TypeBool struct {
	Value interface{}
}

func (parser *TypeBool) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	return strconv.ParseBool(token)
}

func (parser *TypeBool) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeBool) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeBool) Minimum() int {
	return 1
}

func (parser *TypeBool) Maximum() int {
	return 1
}

func (parser *TypeBool) String() string {
	return ":bool"
}

type TypePattern struct {
	Value interface{}
}

func (parser *TypePattern) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	_, e := regexp.Compile(token)
	return token, e
}

func (parser *TypePattern) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypePattern) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypePattern) Minimum() int {
	return 1
}

func (parser *TypePattern) Maximum() int {
	return 1
}

func (parser *TypePattern) String() string {
	return ":pattern"
}

type TypeUser struct {
	Value interface{}
}

func (parser *TypeUser) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	for _, u := range msg.Mentions {
		if token == "<@"+u.ID+">" || token == "<@!"+u.ID+">" {
			return u.ID, nil
		}
	}

	u, err := ctx.DiscordSession().User(token)
	if err == nil {
		return u.ID, nil
	}

	g, err := ctx.DiscordSession().Guild(msg.GuildID)
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

func (parser *TypeUser) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeUser) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeUser) Minimum() int {
	return 1
}

func (parser *TypeUser) Maximum() int {
	return 1
}

func (parser *TypeUser) String() string {
	return ":user"
}

type TypeRole struct {
	Value interface{}
}

func (parser *TypeRole) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	for _, rid := range msg.MentionRoles {
		if token == "<@&"+rid+">" {
			role, err := ctx.DiscordSession().State.Role(msg.GuildID, rid)
			if err != nil {
				return nil, err
			}
			return role.ID, nil
		}
	}

	role, err := ctx.DiscordSession().State.Role(msg.GuildID, token)
	if err == nil {
		return role.ID, nil
	}

	roles, err := ctx.DiscordSession().GuildRoles(msg.GuildID)
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

func (parser *TypeRole) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeRole) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeRole) Minimum() int {
	return 1
}

func (parser *TypeRole) Maximum() int {
	return 1
}

func (parser *TypeRole) String() string {
	return ":role"
}

type TypeColor struct {
	Value interface{}
}

func (parser *TypeColor) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	err := errors.New(fmt.Sprintf("`%s` not a hex clolor", token))
	if len(token) < 3 {
		return nil, err
	}
	if token[0] == '#' {
		token = token[1:]
	}
	bit := uint(20)
	if len(token) != 3 && len(token) != 6 {
		return nil, err
	}
	var color uint
	for _, r := range token {
		var n uint
		if r >= '0' && r <= '9' {
			n = uint(r - '0')
		} else if r >= 'a' && r <= 'f' {
			n = uint(r-'a') + 10
		} else if r >= 'A' && r <= 'F' {
			n = uint(r-'A') + 10
		} else {
			return nil, err
		}
		color |= n << bit
		bit -= 4
		if len(token) == 3 {
			color |= n << bit
			bit -= 4
		}
	}
	return int(color), nil
}

func (parser *TypeColor) ClearValue() {
	ClearValue(parser.Value)
}

func (parser *TypeColor) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	return SetValue(parser.Value, value)
}

func (parser *TypeColor) Minimum() int {
	return 1
}

func (parser *TypeColor) Maximum() int {
	return 1
}

func (parser *TypeColor) String() string {
	return ":color"
}

type TypeOneOfResult struct {
	Variant Token
	Value   interface{}
}

type TypeOneOf struct {
	Name     string
	Error    error
	Variants []Token
}

func (parser *TypeOneOf) Parse(ctx *Context, msg *discordgo.Message, token string) (interface{}, error) {
	for _, t := range parser.Variants {
		v, err := t.Parse(ctx, msg, token)
		if err == nil {
			return &TypeOneOfResult{
				Variant: t,
				Value:   v,
			}, nil
		}
	}
	return nil, parser.Error
}

func (parser *TypeOneOf) ClearValue() {
	for _, v := range parser.Variants {
		v.ClearValue()
	}
}

func (parser *TypeOneOf) SetValue(ctx *Context, msg *discordgo.Message, token string, value interface{}) error {
	res := value.(*TypeOneOfResult)
	return res.Variant.SetValue(ctx, msg, token, res.Value)
}

func (parser *TypeOneOf) Minimum() int {
	return 1
}

func (parser *TypeOneOf) Maximum() int {
	return 1
}

func (parser *TypeOneOf) String() string {
	return parser.Name
}

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

func ClearValue(dest interface{}) {
	if dest == nil {
		return
	}
	v := reflect.Indirect(reflect.ValueOf(dest))
	z := reflect.Zero(v.Type())
	v.Set(z)
}

func SetValue(dest, value interface{}) (ret error) {
	defer func() {
		if err := recover(); err != nil {
			if ferr, ok := err.(error); ok {
				ret = ferr
			} else {
				ret = errors.New(fmt.Sprint(err))
			}
		}
	}()

	if dest == nil {
		return nil
	}
	destvalue := reflect.Indirect(reflect.ValueOf(dest))
	if !destvalue.CanSet() {
		return errors.New("cannot reflective set")
	}

	tv := reflect.ValueOf(value)
	t := destvalue.Type()

	if destvalue.Kind() == reflect.Slice {
		destvalue.Set(reflect.Append(destvalue, tv.Convert(t.Elem())))
	} else {
		destvalue.Set(tv.Convert(t))
	}
	return nil
}

func (pr *ParserRegistrar) ProcessInternal(msg *discordgo.Message) error {
	rawtokens, err := RawTokenize(pr.Context.Prefix(msg.GuildID), msg.Content)
	if err == ErrInvalidPrefix {
		return nil
	}
	if err != nil {
		return err
	}
	if len(rawtokens) == 0 {
		return nil
	}

	for _, v := range pr.Descriptors {
		for _, d := range v {
			for _, t := range d.Tokens {
				t.ClearValue()
			}
			n := 0
			last := -1
		tokenloop:
			for l, t := range d.Tokens {
				for x := 0; x < t.Minimum(); x++ {
					if n >= len(rawtokens) {
						break tokenloop
					}
					res, err := t.Parse(pr.Context, msg, rawtokens[n])
					if err != nil {
						break tokenloop
					}
					err = t.SetValue(pr.Context, msg, rawtokens[n], res)
					if err != nil {
						break tokenloop
					}
					n++
				}

				last = l
				for x := t.Minimum(); x < t.Maximum(); x++ {
					if n >= len(rawtokens) {
						continue tokenloop
					}
					res, err := t.Parse(pr.Context, msg, rawtokens[n])
					if err != nil {
						continue tokenloop
					}
					err = t.SetValue(pr.Context, msg, rawtokens[n], res)
					if err != nil {
						continue tokenloop
					}
					n++
				}
			}
			if last == len(d.Tokens)-1 {
				err := d.Handler(msg, rawtokens)
				if err != nil {
					ferr := pr.Context.Session.MessageReactionAdd(msg.ChannelID, msg.ID, "✖")
					if ferr != nil {
						log.Println(ferr)
					}
					return err
				}
				ferr := pr.Context.Session.MessageReactionAdd(msg.ChannelID, msg.ID, "✔")
				if ferr != nil {
					log.Println(ferr)
				}
				return err
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
	err := <-ret
	if err != nil {
		err := pr.Context.SendError(msg.ChannelID, &discordgo.MessageEmbed{
			Description: err.Error(),
		})
		if err != nil {
			log.Println(err)
		}
	}
	return err
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
	pr.Descriptors[descriptor.Group] = append(pr.Descriptors[descriptor.Group], descriptor)
	return nil
}

func RawTokenize(prefix, raw string) (result []string, err error) {
	const (
		prefx = iota
		space
		value
		quote
	)

	prefrunes := []rune(prefix)
	if len(prefrunes) == 0 {
		return nil, errors.New("invalid prefix configured")
	}

	var quotechar rune
	var token []rune
	runes := ([]rune)(raw)
	escaping := false
	state := space

	poff := 0

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
				if poff > 0 {
					switch r {
					case '\\':
						escaping = true
						state = value
					case '"', '\'':
						quotechar = r
						state = quote
					default:
						token = append(token, r)
						state = value
					}
				} else {
					if prefrunes[0] != r {
						return nil, ErrInvalidPrefix
					}
					poff++
					if poff == len(prefrunes) {
						state = space
					} else {
						state = prefx
					}
				}
			}
		case prefx:
			if prefrunes[poff] != r {
				return nil, ErrInvalidPrefix
			}
			poff++
			if poff == len(prefrunes) {
				state = space
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
			if r == quotechar {
				if escaping {
					token = append(token, quotechar)
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
		if state == quote && token[len(token)-1] != quotechar {
			return nil, errors.New(fmt.Sprintf("mismatched quotes in `%c%s`", quotechar, string(token)))
		}

		result = append(result, string(token))
	}
	return
}
