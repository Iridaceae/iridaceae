package log_test

import (
	"os"
	"time"

	"github.com/Iridaceae/iridaceae/pkg/sclog"

	"github.com/Iridaceae/iridaceae/pkg/sclog/log"

	"github.com/rs/zerolog"
)

// this is verbatim taken from rs/zerolog/log.
// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init().
func setup() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""
	// In order to always output a static time to stdout for these
	// examples to pass, we need to override zerolog.TimestampFunc
	// and log.Logger globals -- you would not normally need to do this
	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	zlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = log.SetZ(zlog)
}

// TODO: New
func ExampleNew() {
	log.New()
	setup()
	log.Info("hello world")
	// Output: {"level":"INFO","time":1199811905,"message":"hello world"}
}

func ExampleSetZ() {
	setup()
	sclog.Mapper().Set("test", "v1")
	sclog.AddGlobalFields("test")
	log.Info("hello world")
	// Output: {"level":"INFO","time":1199811905,"test":"v1","message":"hello world"}
}

func ExampleTrace() {
	setup()
	log.Trace("hello world")
	// Output: {"level":"TRACE","time":1199811905,"message":"hello world"}
}

func ExampleTracef() {
	setup()
	log.Tracef("hello %s", "world")
	// Output: {"level":"TRACE","time":1199811905,"message":"hello world"}
}

func ExampleDebug() {
	setup()
	log.Debug("hello world")
	// Output: {"level":"DEBUG","time":1199811905,"message":"hello world"}
}

func ExampleDebugf() {
	setup()
	log.Debugf("hello %s", "world")
	// Output: {"level":"DEBUG","time":1199811905,"message":"hello world"}
}
func ExampleInfo() {
	setup()
	log.Info("hello world")
	// Output: {"level":"INFO","time":1199811905,"message":"hello world"}
}

func ExampleInfof() {
	setup()
	log.Infof("hello %s", "world")
	// Output: {"level":"INFO","time":1199811905,"message":"hello world"}
}

func ExampleWarn() {
	setup()
	log.Warn("test")
	// Output: {"level":"WARN","time":1199811905,"message":"test"}
}

func ExampleWarnf() {
	setup()
	log.Warnf("%s", "test")
	// Output: {"level":"WARN","time":1199811905,"message":"test"}
}

func ExampleError() {
	setup()
	log.Error("test error")
	// Output: {"level":"ERROR","time":1199811905,"message":"test error"}
}

func ExampleErrorf() {
	setup()
	log.Errorf("test %s", "error")
	// Output: {"level":"ERROR","time":1199811905,"message":"test error"}
}

// TODO: Panic and Fatal

func ExampleLog() {
	setup()
	log.Log("hello world")
	// Output: {"time":1199811905,"message":"hello world"}
}

func ExampleLogf() {
	setup()
	log.Logf("hello %s", "world")
	// Output: {"time":1199811905,"message":"hello world"}
}
func ExamplePrint() {
	setup()
	log.Print("hello world")
	// Output: {"level":"DEBUG","time":1199811905,"message":"hello world"}
}

func ExamplePrintf() {
	setup()
	log.Printf("hello %s", "world")
	// Output: {"level":"DEBUG","time":1199811905,"message":"hello world"}
}
