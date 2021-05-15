package pomodoro

import "fmt"

func OutputStatus(status Status) string {
	state := ""
	if status.State >= RUNNING {
		state = string(status.State.String()[0])
	}
	if status.State == RUNNING {
		return fmt.Sprintf("%s [%d/%d] %s", state, status.Count, status.NumPoms, status.Remaining)
	}
	return fmt.Sprintf("%s [%d/%d] -", state, status.Count, status.NumPoms)
}
