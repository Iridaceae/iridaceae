package pomodoro

import "time"

// TaskRunner defines a runner instance for our tasks.
type TaskRunner struct {
	count    int
	taskID   int
	taskMsg  string
	numPoms  int
	original time.Duration // This is our first task duration. This will be use in case it is paused.
	state    State
	store    *Store
	started  time.Time
	pause    chan bool
	toggle   chan bool
	notifier Notifier
	duration time.Duration
}

func NewMockTaskRunner(task *Task, store *Store, notifier Notifier) (*TaskRunner, error) {
	return &TaskRunner{
		taskID:   task.ID,
		taskMsg:  task.Message,
		numPoms:  task.NumPoms,
		original: task.Duration,
		state:    State(0),
		store:    store,
		pause:    make(chan bool),
		toggle:   make(chan bool),
		notifier: notifier,
		duration: task.Duration,
	}, nil
}

func NewTaskRunner(task *Task) (*TaskRunner, error) {
	return &TaskRunner{
		taskID:   task.ID,
		taskMsg:  task.Message,
		numPoms:  task.NumPoms,
		original: task.Duration,
		state:    State(0),
		pause:    make(chan bool),
		toggle:   make(chan bool),
		duration: 0,
	}, nil
}

func (t *TaskRunner) Start() {
	go func() {
		err := t.run()
		if err != nil {
			panic(err)
		}
	}()
}

func (t *TaskRunner) run() error {
	for t.count < t.numPoms {
		// create a new poms for us to track start and end time of given session
		pom := &Pomodoro{}
		// start this pom
		pom.Start = time.Now()
		// change our state to running
		t.SetState(RUNNING)
		// start a new timer
		timer := time.NewTimer(t.duration)
		// record our start time
		t.started = pom.Start
	loop:
		select {
		case <-timer.C:
			// When we receive sigterm, or breaking signal from user.
			t.SetState(CANCELED)
			t.count++
		case <-t.toggle:
			// catch any toggle when we are not expecting them
			goto loop
		case <-t.pause:
			timer.Stop()
			// record remaining time of current pomodoro
			remains := t.TimeRemaining()
			// change state to pause
			t.SetState(PAUSED)
			// wait for user signal (input or event listener)
			<-t.pause
			// Resume the timer with previous time
			// and resetting our current ones.
			timer.Reset(remains)
			// Change our settings to current
			t.started = time.Now()
			t.duration = remains
			// restore our state to running
			t.SetState(RUNNING)
			goto loop
		}
		pom.End = time.Now()
		// err := t.store -> here we will store our pomodoro with given t.taskID, and poms
		if t.count == t.numPoms {
			break
		}
		err := t.notifier.Notify("Pomodoro", ":timer: Time to take a break :timer:")
		if err != nil {
			return err
		}
		// reset duration in case it was paused.
		t.duration = t.original
		// User concludes the break
		<-t.toggle
	}
	err := t.notifier.Notify("Pomodoro", ":blush: Pomodoro session has completed :blush:")
	if err != nil {
		return err
	}
	t.SetState(COMPLETED)
	return nil
}

func (t *TaskRunner) TimeRemaining() time.Duration {
	return (t.duration - time.Since(t.started)).Truncate(time.Second)
}

func (t *TaskRunner) SetState(state State) {
	t.state = state
}

func (t *TaskRunner) Toogle() {
	t.toggle <- true
}

func (t *TaskRunner) Pause() {
	t.pause <- true
}

func (t *TaskRunner) Status() *Status {
	return &Status{
		State:     t.state,
		Remaining: t.TimeRemaining(),
		Count:     t.count,
		NumPoms:   t.numPoms,
	}
}
