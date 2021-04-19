package models

import (
	"sync"
	"time"
)

// Pomodoro defines a single state of a pomodoro sessions
// Usuage of channel for cancel handling
type Pomodoro struct {
	workDuration time.Duration
	onWorkEnd    TaskCallback
	notifyInfo   NotifyInfo
	cancelChan   chan struct{}
	cancel       sync.Once
}

// NotifyInfo defines notification message for users
type NotifyInfo struct {
	titleID string
	userID  User
}

type TaskCallback func(info NotifyInfo, finished bool)
