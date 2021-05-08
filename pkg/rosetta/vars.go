package rosetta

import "errors"

var (
	ErrRateLimited    = errors.New("rate limited")
	EmbedColorDefault = 0x6A5ACD
	EmbedColorError   = 0xE53935
)
