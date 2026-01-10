package pomodoro

import (
	"time"
)

type Mode int

const (
	WorkMode Mode = iota
	BreakMode
	PauseMode
	IdleMode
)

type PomodoroTimer struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	Mode          Mode
	Remaining     time.Duration
	quit          chan struct{}
	Updates       chan time.Duration
	previousMode  Mode
}

func NewPomodoroTimer(workMinutes, breakMinutes int) *PomodoroTimer {
	return &PomodoroTimer{
		WorkDuration:  time.Duration(workMinutes) * time.Minute,
		BreakDuration: time.Duration(breakMinutes) * time.Minute,
		Mode:          IdleMode,
		quit:          make(chan struct{}),
		Updates:       make(chan time.Duration, 10),
	}
}

func (pt *PomodoroTimer) Start() {
	if pt.Mode != IdleMode {
		return
	}

	pt.Mode = WorkMode
	pt.Remaining = pt.WorkDuration
	pt.quit = make(chan struct{})
	go pt.run()
}

func (pt *PomodoroTimer) UpdateDurations(workMinutes, breakMinutes int) {
	pt.WorkDuration = time.Duration(workMinutes) * time.Minute
	pt.BreakDuration = time.Duration(breakMinutes) * time.Minute
}

func (pt *PomodoroTimer) Stop() {
	if pt.Mode == IdleMode {
		return
	}

	pt.Mode = IdleMode
	pt.Remaining = 0

	select {
	case <-pt.quit:
	default:
		close(pt.quit)
	}
}

func (pt *PomodoroTimer) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			pt.Remaining -= 1 * time.Second
			select {
			case pt.Updates <- pt.Remaining:
			default:
			}

			if pt.Remaining <= 0 {
				pt.switchMode()
			}

		case <-pt.quit:
			return
		}
	}
}

func (pt *PomodoroTimer) Pause() {
	if pt.Mode != WorkMode && pt.Mode != BreakMode {
		return
	}

	pt.previousMode = pt.Mode
	pt.Mode = PauseMode

	select {
	case <-pt.quit:
	default:
		close(pt.quit)
	}
}

func (pt *PomodoroTimer) Resume() {
	if pt.Mode != PauseMode {
		return
	}

	pt.Mode = pt.previousMode
	pt.quit = make(chan struct{})
	go pt.run()
}

func (pt *PomodoroTimer) switchMode() {
	if pt.Mode == WorkMode {
		pt.Mode = BreakMode
		pt.Remaining = pt.BreakDuration
	} else {
		pt.Stop()
	}
}

func (pt *PomodoroTimer) Shutdown() {
	if pt.Mode != IdleMode {
		pt.Mode = IdleMode
		select {
		case <-pt.quit:
		default:
			close(pt.quit)
		}
	}
	close(pt.Updates)
}
