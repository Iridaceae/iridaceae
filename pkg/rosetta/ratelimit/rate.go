// Package ratelimit provides a basic token-bucket limiter for rosetta router.
// This can be used to prevent users from spamming the bot and overload discord API.
package ratelimit

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/Iridaceae/iridaceae/pkg/rosetta"
)

// RateLimiter implements a middleware that manage multiple rate limiters.
// This can also be parsed a custom Manager instance if you want to handle limiters differently.
type RateLimiter struct {
	m Manager
}

// New returns a new instance of Rate Limiter.
func New(m ...Manager) *RateLimiter {
	var man Manager
	if len(m) > 0 && m[0] != nil {
		man = m[0]
	} else {
		man = newInternalManager(10 * time.Minute)
	}
	return &RateLimiter{man}
}

func (r *RateLimiter) Handle(cmd rosetta.Command, ctx rosetta.Context, layer rosetta.MiddlewareLayer) (bool, error) {
	c, ok := cmd.(rosetta.LimitedConfig)
	if !ok {
		return true, nil
	}

	var gid string
	switch {
	case c.IsLimiterGlobal():
		gid = "__global__"
	case ctx.GetChannel().Type == discordgo.ChannelTypeDM || ctx.GetChannel().Type == discordgo.ChannelTypeGroupDM:
		gid = "__dm__"
	default:
		gid = ctx.GetGuild().ID
	}

	limiter := r.m.GetBucket(cmd, ctx.GetUser().ID, gid)
	if k, next := limiter.Take(); !k {
		_, _ = ctx.RespondEmbedError(fmt.Sprintf("You are being rate limited.\nWait %s before using this command again.", next.String()), rosetta.ErrRateLimited)
		return false, rosetta.ErrRateLimited
	}
	return true, nil
}

func (r *RateLimiter) GetLayer() rosetta.MiddlewareLayer {
	return rosetta.LayerBeforeCommand
}
