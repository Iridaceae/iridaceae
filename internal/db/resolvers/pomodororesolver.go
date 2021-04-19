package resolvers

import (
	"sync"
	"time"

	"github.com/aarnphm/iris/internal/db/models"
)

// NewPom create a new pomodoro and start it using time.NewTimer. onWorkEnd will be called after the goroutine
func NewPom(workDuration time.Duration, onWorkEnd models.TaskCallback, notify models.NotifyInfo) *models.Pomodoro {
	pom := &models.Pomodoro{
		workDuration,
		onWorkEnd,
		notify,
		make(chan struct{}),
		sync.Once{},
	}

	go pom.startPom()
	return pom
}

// Cancel is used to cancel the current state of the goroutine. sync.Once to prevent panic
func (pom *models.Pomodoro) Cancel() {
	pom.cancel.Do(func() {
		close(pom.cancelChan)
	})
}

func (pom *models.Pomodoro) startPom() {
	workTimer := time.NewTimer(pom.workDuration)

	select {
	case <-workTimer.C:
		go pom.onWorkEnd(pom.NotifyInfo, true)
	case <-pom.cancelChan:
		go pom.onWorkEnd(pom.NotifyInfo, false)
	}
}
