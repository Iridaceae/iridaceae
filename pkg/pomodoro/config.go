package pomodoro

import (
	"time"

	"github.com/Iridaceae/iridaceae/pkg/configmanager"
)

var (
	DateTimeFmt, _ = configmanager.RegisterOption("pomodoro.datetimefmt", "date time format", time.RFC3339)
	SocketPath, _  = configmanager.RegisterOption("pomodoro.socketpath", "socket path for our pomodoro server", "/var/run/pomodoro.sock")
)
