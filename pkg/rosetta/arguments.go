// Package rosetta is iridaceae's internal structured command handler, including a token-bucket rate limiter and middleware handlers.
package rosetta

import (
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

// Codeblock represents a discord codeblock.
type Codeblock struct {
	Language string
	Content  string
}

// Argument extends string.
type Argument string

// String returns raw string value of given argument.
func (a Argument) String() string {
	return string(a)
}

// AsBool returns given argument as boolean.
// Since we are using strconv.ParseBool, it will accept
// 1, t, T, TRUE, true, True, 0, f, F, FALSE, False, false.
func (a Argument) AsBool() (bool, error) {
	return strconv.ParseBool(a.String())
}

// AsInt returns given argument as int.
func (a Argument) AsInt() (int, error) {
	return strconv.Atoi(a.String())
}

// AsInt64 parses given argument into int64.
func (a Argument) AsInt64() (int64, error) {
	return strconv.ParseInt(a.String(), 10, 64)
}

// AsUserMentionID returns id of mentioned user or an empty string if there isn't one.
func (a Argument) AsUserMentionID() string {
	if matches := UserMentionRegex.MatchString(a.String()); !matches {
		return ""
	}
	return UserMentionRegex.FindStringSubmatch(a.String())[1]
}

// AsRoleMentionID returns id of mentioned role or an empty string if there isn't one.
func (a Argument) AsRoleMentionID() string {
	if matches := RoleMentionRegex.MatchString(a.String()); !matches {
		return ""
	}
	return RoleMentionRegex.FindStringSubmatch(a.String())[1]
}

// AsChannelMentionID returns id of mentioned channel or an empty string if there isn't one.
func (a Argument) AsChannelMentionID() string {
	if matches := ChannelMentionRegex.MatchString(a.String()); !matches {
		return ""
	}
	return ChannelMentionRegex.FindStringSubmatch(a.String())[1]
}

// AsDuration parses given argument into a duration.
func (a Argument) AsDuration() (time.Duration, error) {
	return time.ParseDuration(a.String())
}

// Arguments wraps around Arguments.
type Arguments struct {
	raw  string
	args []Argument
}

// FromArguments create a new arguments from given list.
func FromArguments(args []Argument) *Arguments {
	return &Arguments{"", args}
}

// ParseArguments parses raw input message into several arguments.
func ParseArguments(msg string) *Arguments {
	raw := ArgumentsRegex.FindAllString(msg, -1)
	args := make([]Argument, len(raw))

	for idx, r := range raw {
		r = trimPreSuffix(r, "\"")
		args[idx] = Argument(r)
	}

	return &Arguments{
		raw:  msg,
		args: args,
	}
}

// Raw returns raw string of the arguments.
func (a Arguments) Raw() string {
	return a.raw
}

// Args returns list of args.
func (a Arguments) Args() []Argument {
	return a.args
}

// AsSingle returns a singleton of arguments with raw content without args.
func (a Arguments) AsSingle() *Arguments {
	return &Arguments{raw: a.raw}
}

// Len returns length of given arguments.
func (a Arguments) Len() int {
	return len(a.args)
}

// Get returns the nth arguments.
func (a *Arguments) Get(n int) Argument {
	if n < 0 || a.Len() <= n {
		return ""
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
	for _, arg := range a.args {
		raw += arg.String() + " "
	}
	a.raw = strings.TrimSpace(raw)
}

// IndexOf returns the index of a argument in arguments.
func (a Arguments) IndexOf(arg string) int {
	for i, v := range a.args {
		if arg == v.String() {
			return i
		}
	}
	return -1
}

// AsCodeblock parses given arguments as codeblock.
func (a Arguments) AsCodeblock() *Codeblock {
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
