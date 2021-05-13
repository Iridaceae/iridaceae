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
	ErrTypeGuildPrefixGetter ErrorType = 1 << iota
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

	// ErrCommandNotFound is thrown when command is not found.
	ErrCommandNotFound = errors.New("command not found")

	// ErrNotExecutableInDMs is thrown when command is not allowed to run in DMs.
	ErrNotExecutableInDMs = errors.New("command is not executable in DMs")

	// ErrInvokeDoesNotExists is thrown when given command invoker doesn't exists.
	ErrInvokeDoesNotExists = errors.New("given invoke doesn't exists")

	EmbedColorDefault = 0x6A5ACD
	EmbedColorError   = 0xE53935
)
