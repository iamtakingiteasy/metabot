package impl

import (
	"reflect"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNewParser(t *testing.T) {
	p := NewParser(&Context{})
	defer p.Stop()

	if p == nil {
		t.Error("parser is nil")
	}
}

func TestRawTokenize(t *testing.T) {
	var tokens []string
	var err error

	tokens, err = RawTokenize("a b c")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"a", "b", "c"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("a \" bada boom \" c")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"a", " bada boom ", "c"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("a\\ b\\ c")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"a b c"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("a \"d\\\"d\" c")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"a", "d\"d", "c"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("\\ \\ \\ z")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"   z"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("z\\x")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"z\\x"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("z\\\\x")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"z\\x"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("z\\")
	if err != nil {
		t.Error("non-nil error", err)
	}
	if !reflect.DeepEqual(tokens, []string{"z\\"}) {
		t.Error("invalid tokens")
	}

	tokens, err = RawTokenize("\"z")
	if err == nil {
		t.Error("expected error")
	}
}

func TestParserRegistrar_AddCommand(t *testing.T) {
	parser := NewParser(&Context{})
	defer parser.Stop()

	var err error

	err = parser.RegisterCommand(func(ctx *Context, registrar Registrar) (err error) {
		var foo string
		return registrar.AddCommand(&Descriptor{
			Group: "main",
			Name:  "test",
			Tokens: []*Token{
				{
					Parser:    TypeLiteral,
					Parameter: "test",
					Value:     &foo,
				},
			},
			Handler: func(msg *discordgo.Message) error {
				return nil
			},
		})
	})
	if err != nil {
		t.Error("expected err to be nil", err)
	}

	err = parser.RegisterCommand(func(ctx *Context, registrar Registrar) (err error) {
		var foo string
		return registrar.AddCommand(&Descriptor{
			Group: "main",
			Name:  "test",
			Tokens: []*Token{
				{
					Parser:    TypeLiteral,
					Parameter: "test",
					Value:     &foo,
					Minimum:   1,
				},
			},
			Handler: func(msg *discordgo.Message) error {
				return nil
			},
		})
	})
	if err == nil {
		t.Error("expected err to be non-nil")
	}
}

func TestParserRegistrar_Commands(t *testing.T) {
	parser := NewParser(&Context{})
	defer parser.Stop()

	var err error

	if len(parser.Commands()) != 0 {
		t.Error("expected commands to be empty")
	}

	err = parser.RegisterCommand(func(ctx *Context, registrar Registrar) (err error) {
		var foo string
		return registrar.AddCommand(&Descriptor{
			Group: "main",
			Name:  "test",
			Tokens: []*Token{
				{
					Parser:    TypeLiteral,
					Parameter: "test",
					Value:     &foo,
				},
			},
			Handler: func(msg *discordgo.Message) error {
				return nil
			},
		})
	})
	if err != nil {
		t.Error("expected err to be nil", err)
	}

	if len(parser.Commands()) != 1 {
		t.Error("expected commands to be singular")
	}
}

func TestParserRegistrar_ProcessCommand(t *testing.T) {
	parser := NewParser(&Context{})
	defer parser.Stop()

	var err error

	var run int

	err = parser.RegisterCommand(func(ctx *Context, registrar Registrar) (err error) {
		var foo string
		var strs []string
		var nums []int
		var bools []bool
		return registrar.AddCommand(&Descriptor{
			Group: "main",
			Name:  "test",
			Tokens: []*Token{
				{
					Parser:    TypeLiteral,
					Parameter: "test",
					Value:     &foo,
				},
				{
					Parser:  TypeString,
					Value:   &strs,
					Minimum: 2,
				},
				{
					Parser:  TypeNumber,
					Value:   &nums,
					Minimum: 1,
					Maximum: 3,
				},
				{
					Parser:  TypeBool,
					Value:   &bools,
					Minimum: 0,
					Maximum: -1,
				},
			},
			Handler: func(msg *discordgo.Message) error {
				switch run {
				case 0:
					if foo != "test" {
						t.Error("expected foo to be test, got", foo)
					}
					if !reflect.DeepEqual(strs, []string{"ab", "ba"}) {
						t.Error("expected strs to be [ab ba], got", strs)
					}
					if !reflect.DeepEqual(nums, []int{3, 2}) {
						t.Error("expected nums to be [3 2], got", nums)
					}
					if !reflect.DeepEqual(bools, []bool{true, false, true}) {
						t.Error("expected bools to be [true false true], got", bools)
					}
				case 1:
					if foo != "test" {
						t.Error("expected foo to be test, got", foo)
					}
					if !reflect.DeepEqual(strs, []string{"xc", "zz"}) {
						t.Error("expected strs to be [xc zz], got", strs)
					}
					if !reflect.DeepEqual(nums, []int{5}) {
						t.Error("expected nums to be [5], got", nums)
					}
					if len(bools) > 0 {
						t.Error("expected bools to be [], got", bools)
					}
				}
				return nil
			},
		})
	})
	if err != nil {
		t.Error("expected err to be nil", err)
	}

	run = 0
	err = parser.ProcessCommand(&discordgo.Message{Content: "test ab ba 3 2 true false true"})
	if err != nil {
		t.Error("expected err to be nil", err)
	}
	run = 1
	err = parser.ProcessCommand(&discordgo.Message{Content: "test xc zz 5"})
	if err != nil {
		t.Error("expected err to be nil", err)
	}
	run = 2
	err = parser.ProcessCommand(&discordgo.Message{Content: "test ab ba ff ss true false true"})
	if err == nil {
		t.Error("expected err to be non-nil")
	}
	run = 3
	err = parser.ProcessCommand(&discordgo.Message{Content: "test ab"})
	if err == nil {
		t.Error("expected err to be non-nil")
	}
	run = 4
	err = parser.ProcessCommand(&discordgo.Message{Content: "test \"fff"})
	if err == nil {
		t.Error("expected err to be non-nil")
	}
	run = 5
	err = parser.ProcessCommand(&discordgo.Message{Content: "baka"})
	if err == nil {
		t.Error("expected err to be non-nil")
	}
}
