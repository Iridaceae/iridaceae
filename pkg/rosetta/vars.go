package rosetta

import (
	"errors"
)

const (
	ObjectMapKeyRouter = "rosetta_router"
)

// ErrorType is the type of error happened while parsing in router.
type ErrorType int

const (
	ErrTypeGuildPrefixGetter ErrorType = iota + 1
	ErrTypeGetChannel
	ErrTypeGetGuild
	ErrTypeCommandNotFound
	ErrTypeNotExecutableInDM
	ErrTypeMiddleware
	ErrTypeCommandExec
	ErrTypeDeleteCommandMessage
)

var (
	// ErrRateLimited is thrown when users spam commands =).
	ErrRateLimited = errors.New("rate limited")

	// ErrInvokeDoesNotExists is thrown when given command invoker doesn't exists.
	ErrInvokeDoesNotExists = errors.New("given invoke doesn't exists")

	// ErrGuildPrefixGetter is thrown GuildPrefixGetter failed.
	ErrGuildPrefixGetter = errors.New("error while getting guild prefix")

	// ErrGetChannel is thrown while getting channel.
	ErrGetChannel = errors.New("error while getting channel")

	// ErrGetGuild is thrown while getting guild.
	ErrGetGuild = errors.New("error while getting guild")

	// ErrCommandNotFound is thrown when command is not found.
	ErrCommandNotFound = errors.New("command not found")

	// ErrNotExecutableInDMs is thrown when command is not allowed to run in DMs.
	ErrNotExecutableInDMs = errors.New("command is not executable in DMs")

	// ErrMiddleware is thrown when middleware failed unexpectedly.
	ErrMiddleware = errors.New("middleware error")

	// ErrCommandExec is thrown when command exec failed.
	ErrCommandExec = errors.New("command failed to execute")

	// ErrDeleteCommandMessage is thrown when error occurred when deleting message.
	ErrDeleteCommandMessage = errors.New("failed while deleting command message")

	EmbedColorDefault = 0x6A5ACD
	EmbedColorError   = 0xE53935
)
