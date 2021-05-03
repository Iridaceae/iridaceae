package jog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// ArgumentsRegex defines regex arguments should match.
	ArgumentsRegex = regexp.MustCompile("(\"[^\"]+\"|[^\\s]+)")
	// UserMentionRegex defines regex user mention should match.
	UserMentionRegex = regexp.MustCompile(`<@!?(\d+)>`)
	// RoleMentionRegex defines regex role mention should match.
	RoleMentionRegex = regexp.MustCompile(`<@&(\d+)>`)
	// ChannelMentionRegex defines regex channel mention should match.
	ChannelMentionRegex = regexp.MustCompile(`<#(\d+)>`)
	// CodeblockRegex defines regex for codeblock to match.
	CodeblockRegex = regexp.MustCompile("(?s)\\n*```(?:([\\w.\\-]*)\\n)?(.*)```")
	// InlineCodeRegex defines regex for inline code to match.
	InlineCodeRegex = regexp.MustCompile("(?s)\\n*`(.*)`")
	// ProgrammingLanguage defines valid language for codeblock.
	ProgrammingLanguage = []string{
		"go",
		"golang",
		"docker",
		"dockerfile",
		"python",
		"java",
		"c",
		"js",
		"jsx",
		"ts",
		"tsx",
		"lua",
		"makefile",
		"json",
	}
)

// Arguments represents arguments that may be used in a command context.
type Arguments struct {
	raw  string
	args []*Argument
}

// Argument represents a single argument.
type Argument struct {
	raw string
}

// Codeblock represents a discord codeblock.
type Codeblock struct {
	Language string
	Content  string
}

// ParseArguments parses raw input message into several arguments.
func ParseArguments(msg string) *Arguments {
	// define raw args from msg
	raw := ArgumentsRegex.FindAllString(msg, -1)
	args := make([]*Argument, len(raw))

	for idx, r := range raw {
		r = trimPreSuffix(r, "\"")
		args[idx] = &Argument{raw: r}
	}

	return &Arguments{
		raw:  msg,
		args: args,
	}
}

// Raw returns raw string of the arguments.
func (a *Arguments) Raw() string {
	return a.raw
}

// AsSingle returns a singleton of arguments with raw content.
func (a *Arguments) AsSingle() *Arguments {
	return &Arguments{raw: a.raw}
}

// Len returns length of given arguments.
func (a *Arguments) Len() int {
	return len(a.args)
}

// Get returns the nth arguments.
func (a *Arguments) Get(n int) *Argument {
	if a.Len() <= n {
		return &Argument{raw: ""}
	}
	return a.args[n]
}

// Remove removes the nth arguments.
func (a *Arguments) Remove(n int) {
	if a.Len() <= n {
		return
	}

	a.args = append(a.args[:n], a.args[n+1:]...)
	// sets new raw string
	raw := ""
	for _, args := range a.args {
		raw += args.raw + " "
	}
	a.raw = strings.TrimSpace(raw)
}

// AsCodeblock parses given arguments as codeblock.
func (a *Arguments) AsCodeblock() *Codeblock {
	raw := a.Raw()

	// check if it is a normal codeblock.
	matches := CodeblockRegex.MatchString(raw)
	if !matches {
		// check if this is inline code.
		matches = InlineCodeRegex.MatchString(raw)
		if matches {
			sub := InlineCodeRegex.FindStringSubmatch(raw)
			return &Codeblock{
				Language: "",
				Content:  sub[1],
			}
		}
		return nil
	}

	// define contents and language
	sub := CodeblockRegex.FindStringSubmatch(raw)
	fmt.Println(sub)
	language := ""
	content := sub[1] + sub[2]
	if sub[1] != "" && arrayContains(ProgrammingLanguage, sub[1], false) {
		language = sub[1]
		content = sub[2]
	}
	return &Codeblock{
		Language: language,
		Content:  content,
	}
}

// Raw returns raw string value of given argument.
func (a *Argument) Raw() string {
	return a.raw
}

// AsBool returns given argument as boolean.
func (a *Argument) AsBool() (bool, error) {
	return strconv.ParseBool(a.raw)
}

// AsInt returns given argument as int.
func (a *Argument) AsInt() (int, error) {
	return strconv.Atoi(a.raw)
}

// AsInt64 parses given argument into int64.
func (a *Argument) AsInt64() (int64, error) {
	return strconv.ParseInt(a.raw, 10, 64)
}

// AsUserMentionID returns id of mentioned user or an empty string if there isn't one.
func (a *Argument) AsUserMentionID() string {
	if matches := UserMentionRegex.MatchString(a.raw); !matches {
		return ""
	}
	return UserMentionRegex.FindStringSubmatch(a.raw)[1]
}

// AsRoleMentionID returns id of mentioned role or an empty string if there isn't one.
func (a *Argument) AsRoleMentionID() string {
	if matches := RoleMentionRegex.MatchString(a.raw); !matches {
		return ""
	}
	return RoleMentionRegex.FindStringSubmatch(a.raw)[1]
}

// AsChannelMentionID returns id of mentioned channel or an empty string if there isn't one.
func (a *Argument) AsChannelMentionID() string {
	if matches := ChannelMentionRegex.MatchString(a.raw); !matches {
		return ""
	}
	return ChannelMentionRegex.FindStringSubmatch(a.raw)[1]
}

// AsDuration parses given argument into a duration.
func (a *Argument) AsDuration() (time.Duration, error) {
	return time.ParseDuration(a.raw)
}
