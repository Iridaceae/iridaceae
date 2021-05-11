package log_test

import (
	"os"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/log"

	"github.com/rs/zerolog"
)

var hookFunc = log.MapperHook{}

func setup() {
	// UNIX time is a lot faster and smaller than most timestamp.
	zerolog.TimeFieldFormat = ""
	// setup a static time to stdout
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2020, 4, 20, 4, 20, 0o4, 0, time.UTC)
	}
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger().Hook(hookFunc)
	log.L = log.NewZ(zlog)
}

func ExampleNew() {
	log.New()
	setup()
	log.Info().Msg("hello world")
	// Output: {"level":"INFO","time":1587356404,"message":"hello world"}
}

func ExamplePrint() {
	setup()
	log.Print("hello world")
	// Output: {"level":"DEBUG","time":1587356404,"message":"hello world"}
}

func ExamplePrintf() {
	setup()
	log.Printf("hello %s", "world")
	// Output: {"level":"DEBUG","time":1587356404,"message":"hello world"}
}

func ExampleLog() {
	setup()
	log.Log().Msg("hello world")
	// Output: {"time":1587356404,"message":"hello world"}
}

func ExampleTrace() {
	setup()
	log.Trace().Msg("hello world")
	// Output: {"level":"TRACE","time":1587356404,"message":"hello world"}
}

func ExampleDebug() {
	setup()
	log.Debug().Msg("hello world")
	// Output: {"level":"DEBUG","time":1587356404,"message":"hello world"}
}
