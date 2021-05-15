// Package pomodoro provides functionality for pomodoro techniques.
// Improvement from v0:
//   - ability to check current status of pomodoro
//   - setup a group of pomodoro and break timestamp ie: 25 work - 5 break - 25 work - 15 break etc
//   - ability to pause during a pomodoro
// (plus all the current features)
package pomodoro

import "time"

// State tracks current state of users pomodoro.
type State int

const (
	RUNNING State = iota + 1
	CANCELED
	COMPLETED
	PAUSED
)

func (s State) String() string {
	switch s {
	case RUNNING:
		return "RUNNING"
	case CANCELED:
		return "CANCELED"
	case COMPLETED:
		return "COMPLETED"
	case PAUSED:
		return "PAUSED"
	}
	return ""
}

// Task will describe our activity.
type Task struct {
	ID       int           `json:"id"`
	Message  string        `json:"message"`
	Duration time.Duration `json:"duration"`
	Tags     []string      `json:"tags"`      // Tags defines user tag with given task.
	CmplPoms []*Pomodoro   `json:"cmpl_poms"` // CmplPoms allows us to check our completed poms during the longevity of the task.
	NumPoms  int           `json:"num_poms"`  // NumPom will give us number of pomodoro running in given task.
}

// TaskByID is a sortable array of tasks.
// Implements Go's sorting interface.
type TaskByID []*Task

func (t TaskByID) Len() int             { return len(t) }
func (t TaskByID) Swap(t1, t2 int)      { t[t1], t[t2] = t[t2], t[t1] }
func (t TaskByID) Less(t1, t2 int) bool { return t[t1].ID < t[t2].ID }

// After returns an array of tasks that is started after given start time.
func After(start time.Time, tasks []*Task) []*Task {
	filtered := make([]*Task, 0)
	for _, t := range tasks {
		if len(t.CmplPoms) > 0 {
			if start.Before(t.CmplPoms[0].Start) {
				filtered = append(filtered, t)
			}
		}
	}
	return filtered
}

// Pomodoro is a unit of time defining our start and end timestamp.
type Pomodoro struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration returns the runtime of our pom.
func (p Pomodoro) Duration() time.Duration {
	return p.End.Sub(p.Start)
}

// Status acts as a middleware for us to track the state of our current pomodoro.
type Status struct {
	State     State         `json:"state"`
	Remaining time.Duration `json:"remaining"`
	Count     int           `json:"count"`
	NumPoms   int           `json:"num_poms"` // NumPoms defines number of pomodoros currently active.
}

// Notifier will act as our middleware to push notification to discord.
type Notifier interface {
	Notify(string, string) error
}

// NoopNotifier is our mock implementation that doesn't do anything.
type NoopNotifier struct{}

func (n NoopNotifier) Notify(string, string) error { return nil }
