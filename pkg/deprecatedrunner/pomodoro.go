// Package pomodoro defines our pomodoro logics.
package deprecatedrunner

import (
	"sync"
	"time"

	"github.com/Iridaceae/iridaceae/internal/database"
)

// Pomodoro defines a single state of a pomodoro sessions.
// Usage of channel to handle cancel signal.
type Pomodoro struct {
	workDuration time.Duration
	onWorkEnd    TaskCallback // NOTE: uses callback as middleware for new command handler.
	notifyInfo   NotifyInfo   // We can just pass straight from discordgo.Event
	cancelChan   chan struct{}
	cancel       sync.Once
}

// NotifyInfo defines notification message for users.
type NotifyInfo struct {
	TitleID string
	User    *database.User
}

// TaskCallback receives NotifyInfo and a boolean to define whether the task is completed or not.
type TaskCallback func(info NotifyInfo, finished bool)

// UserPomodoroMap is a map-like structure to init a single pomodoro in a channel. It has goroutine safe ops.
type UserPomodoroMap struct {
	mutex     sync.Mutex
	userToPom map[string]*Pomodoro
}

// NewUserPomodoroMap creates a ChannelPomMap and prepares it to be used.
func NewUserPomodoroMap() UserPomodoroMap {
	return UserPomodoroMap{userToPom: make(map[string]*Pomodoro)}
}

// NewPomodoro create a new pomodoro and start it using time.NewTimer. onWorkEnd will be called after the goroutine.
func NewPomodoro(workDuration time.Duration, onWorkEnd TaskCallback, notify NotifyInfo) *Pomodoro {
	pom := &Pomodoro{
		workDuration: workDuration,
		onWorkEnd:    onWorkEnd,
		notifyInfo:   notify,
		cancelChan:   make(chan struct{}),
		cancel:       sync.Once{},
	}

	go pom.start()
	return pom
}

// Cancel is used to cancel the current state of the goroutine. sync.Once to prevent panic.
func (pom *Pomodoro) Cancel() {
	pom.cancel.Do(func() {
		close(pom.cancelChan)
	})
}

func (pom *Pomodoro) start() {
	workTimer := time.NewTimer(pom.workDuration)

	select {
	case <-workTimer.C:
		go pom.onWorkEnd(pom.notifyInfo, true)
	case <-pom.cancelChan:
		go pom.onWorkEnd(pom.notifyInfo, false)
	}
}

// CreateIfEmpty will create a new Pomodoro for given user according to their discordID if user has none.
// The pomodoro will then be removed from the mapping once completed or canceled.
func (u *UserPomodoroMap) CreateIfEmpty(duration time.Duration, onWorkEnd TaskCallback, notify NotifyInfo) bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	wasCreated := false
	if _, exists := u.userToPom[notify.User.DiscordID]; !exists {
		doneInMap := func(notify NotifyInfo, completed bool) {
			// only called when it is done then we can use the mutex
			// cancellation won't trigger onWorkEnd since start is already done at this point
			u.RemoveIfExists(notify.User.DiscordID)
			onWorkEnd(notify, completed)
		}

		u.userToPom[notify.User.DiscordID] = NewPomodoro(duration, doneInMap, notify)
		wasCreated = true
	}
	return wasCreated
}

// RemoveIfExists will remove a Pomodoro from given channel i one already exists.
func (u *UserPomodoroMap) RemoveIfExists(discordID string) bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	wasRemoved := false
	if p, exists := u.userToPom[discordID]; exists {
		delete(u.userToPom, discordID)
		p.Cancel()
		wasRemoved = true
	}

	return wasRemoved
}

// Count counts the number of current Pomodoro being tracked.
func (u *UserPomodoroMap) Count() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	return len(u.userToPom)
}
