// Package pomodoro provides functionality for pomodoro techniques.
// Improvement from v0:
//   - ability to check current status of pomodoro
//   - setup a group of pomodoro and break timestamp ie: 25 work - 5 break - 25 work - 15 break etc
//   - ability to pause during a pomodoro
// (plus all the current features)
package pomodoro

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
