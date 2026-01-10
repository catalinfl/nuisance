package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/catalinfl/nuisance/httpblock"
	"github.com/catalinfl/nuisance/pomodoro"
	"github.com/catalinfl/nuisance/sound"
	"github.com/catalinfl/nuisance/ui"
	"github.com/catalinfl/nuisance/window"
)

func main() {
	var handler window.WindowHandler = window.WindowHandler{}
	go runApp(&handler)
	app.Main()
}

func updateBlocker(b *httpblock.Blocker, state *ui.AppState) {
	var sitesToBlock []string

	if state.BlockedSites["facebook"] {
		sitesToBlock = append(sitesToBlock, "www.facebook.com", "facebook.com")
	}
	if state.BlockedSites["youtube"] {
		sitesToBlock = append(sitesToBlock, "www.youtube.com", "youtube.com", "m.youtube.com")
	}
	if state.BlockedSites["twitter"] {
		sitesToBlock = append(sitesToBlock, "www.twitter.com", "twitter.com", "x.com", "www.x.com")
	}
	if state.BlockedSites["reddit"] {
		sitesToBlock = append(sitesToBlock, "www.reddit.com", "reddit.com")
	}
	if state.BlockedSites["instagram"] {
		sitesToBlock = append(sitesToBlock, "www.instagram.com", "instagram.com")
	}
	if state.BlockedSites["tiktok"] {
		sitesToBlock = append(sitesToBlock, "www.tiktok.com", "tiktok.com")
	}
	if state.BlockedSites["whatsapp"] {
		sitesToBlock = append(sitesToBlock, "www.web.whatsapp.com", "web.whatsapp.com")
	}

	for _, site := range state.CustomWebsites {
		sitesToBlock = append(sitesToBlock, site)
	}

	b.Sites = sitesToBlock
}

func runApp(winHandler *window.WindowHandler) {
	w := new(app.Window)
	w.Option(
		app.Size(unit.Dp(360), unit.Dp(240)),
		app.Title("Always On Top"),
	)

	// sites to block
	defaultSites := []string{
		"www.facebook.com", "facebook.com",
		"www.youtube.com", "youtube.com", "m.youtube.com",
		"www.twitter.com", "twitter.com", "x.com", "www.x.com",
		"www.reddit.com", "reddit.com",
		"www.instagram.com", "instagram.com",
		"www.tiktok.com", "tiktok.com",
		"www.web.whatsapp.com", "web.whatsapp.com",
	}

	b := httpblock.Blocker{
		Token: "nuisance",
		Sites: defaultSites,
	}

	var hwnd atomic.Uintptr

	// Initialize Pomodoro timer
	pomoTimer := pomodoro.NewPomodoroTimer(25, 5)

	var cleanupOnce sync.Once
	cleanup := func() {
		_ = b.RemoveBlockEntries()
		pomoTimer.Shutdown()
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cleanupOnce.Do(cleanup)
		w.Perform(system.ActionClose)
	}()

	// Find window and set always on top
	go func() {
		time.Sleep(300 * time.Millisecond)
		for range 50 {
			h := winHandler.FindWindowByTitle("Always On Top")
			if h != 0 {
				hwnd.Store(h)
				_ = winHandler.SetAlwaysOnTop(h, true)
				w.Invalidate()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Remove any existing blocks on startup
	b.RemoveBlockEntries()

	th := material.NewTheme()
	btns := ui.NewButtons()
	settingsBtns := ui.NewSettingsButtons()
	state := &ui.AppState{
		CurrentTab:     0,
		AlwaysOnTop:    true,
		BlockedSites:   make(map[string]bool),
		PomodoroTime:   "25:00",
		PomodoroMode:   "Ready",
		WorkMinutes:    25,
		BreakMinutes:   5,
		CustomWebsites: []string{},
	}

	state.BlockedSites["facebook"] = true
	state.BlockedSites["youtube"] = true
	state.BlockedSites["twitter"] = true
	state.BlockedSites["reddit"] = true
	state.BlockedSites["instagram"] = true
	state.BlockedSites["tiktok"] = true
	state.BlockedSites["whatsapp"] = true

	var isBlocking atomic.Bool
	isBlocking.Store(false)

	go func() {
		var lastMode pomodoro.Mode = pomodoro.IdleMode

		for remaining := range pomoTimer.Updates {
			minutes := int(remaining.Minutes())
			seconds := int(remaining.Seconds()) % 60
			state.PomodoroTime = fmt.Sprintf("%02d:%02d", minutes, seconds)

			currentMode := pomoTimer.Mode

			// check for mode changes to handle blocking
			if currentMode != lastMode {
				if currentMode == pomodoro.WorkMode && !isBlocking.Load() {
					// Start work - block websites and play sound
					isBlocking.Store(true)
					updateBlocker(&b, state)
					go func() {
						sound.PlayWorkStart()
						_ = b.AddBlockEntries()
					}()
				} else if currentMode == pomodoro.BreakMode && isBlocking.Load() {
					// Break time - unblock and play sound
					isBlocking.Store(false)
					go func() {
						sound.PlayBreakStart()
						_ = b.RemoveBlockEntries()
					}()
				} else if currentMode == pomodoro.IdleMode && isBlocking.Load() {
					// Complete - unblock and play sound
					isBlocking.Store(false)
					go func() {
						sound.PlayComplete()
						_ = b.RemoveBlockEntries()
					}()
				}
				lastMode = currentMode
			}

			switch currentMode {
			case pomodoro.WorkMode:
				state.PomodoroMode = "Work Time"
			case pomodoro.BreakMode:
				state.PomodoroMode = "Break Time"
			case pomodoro.PauseMode:
				state.PomodoroMode = "Paused"
			case pomodoro.IdleMode:
				state.PomodoroMode = "Ready"
			}
			w.Invalidate()
		}
	}()

	var ops op.Ops

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			cleanupOnce.Do(cleanup)
			return

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle tab clicks
			if btns.Tab1.Clicked(gtx) {
				state.CurrentTab = 0
				w.Option(app.Size(unit.Dp(360), unit.Dp(240)))
			}
			if btns.Tab2.Clicked(gtx) {
				state.CurrentTab = 1
				w.Option(app.Size(unit.Dp(360), unit.Dp(700)))
			}

			// Handle toggle click
			if btns.Toggle.Clicked(gtx) {
				h := hwnd.Load()
				if h != 0 {
					state.AlwaysOnTop = !state.AlwaysOnTop
					newState := state.AlwaysOnTop
					go func(handle uintptr, enable bool) {
						_ = winHandler.SetAlwaysOnTop(handle, enable)
						w.Invalidate()
					}(h, newState)
				}
			}

			// Handle Pomodoro button clicks
			if btns.PomoPlay.Clicked(gtx) {
				if pomoTimer.Mode == pomodoro.PauseMode {
					pomoTimer.Resume()
				} else if pomoTimer.Mode == pomodoro.IdleMode {
					// Update blocker with current settings before starting
					updateBlocker(&b, state)
					pomoTimer.Start()
				}
			}
			if btns.PomoPause.Clicked(gtx) {
				pomoTimer.Pause()
			}
			if btns.PomoReset.Clicked(gtx) {
				pomoTimer.Stop()
				// Update blocker with current settings
				updateBlocker(&b, state)
				state.PomodoroTime = fmt.Sprintf("%02d:00", state.WorkMinutes)
				state.PomodoroMode = "Ready"
			}

			// Handle settings duration adjustments
			if settingsBtns.WorkInc.Clicked(gtx) {
				if state.WorkMinutes < 60 {
					state.WorkMinutes += 5
					pomoTimer.UpdateDurations(state.WorkMinutes, state.BreakMinutes)
					if pomoTimer.Mode == pomodoro.IdleMode {
						state.PomodoroTime = fmt.Sprintf("%02d:00", state.WorkMinutes)
					}
				}
			}
			if settingsBtns.WorkDec.Clicked(gtx) {
				if state.WorkMinutes > 5 {
					state.WorkMinutes -= 5
					pomoTimer.UpdateDurations(state.WorkMinutes, state.BreakMinutes)
					if pomoTimer.Mode == pomodoro.IdleMode {
						state.PomodoroTime = fmt.Sprintf("%02d:00", state.WorkMinutes)
					}
				}
			}
			if settingsBtns.BreakInc.Clicked(gtx) {
				if state.BreakMinutes < 30 {
					state.BreakMinutes += 5
					pomoTimer.UpdateDurations(state.WorkMinutes, state.BreakMinutes)
				}
			}
			if settingsBtns.BreakDec.Clicked(gtx) {
				if state.BreakMinutes > 5 {
					state.BreakMinutes -= 5
					pomoTimer.UpdateDurations(state.WorkMinutes, state.BreakMinutes)
				}
			}

			// Handle settings button clicks
			if settingsBtns.BlockFacebook.Clicked(gtx) {
				state.BlockedSites["facebook"] = !state.BlockedSites["facebook"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockYouTube.Clicked(gtx) {
				state.BlockedSites["youtube"] = !state.BlockedSites["youtube"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockTwitter.Clicked(gtx) {
				state.BlockedSites["twitter"] = !state.BlockedSites["twitter"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockReddit.Clicked(gtx) {
				state.BlockedSites["reddit"] = !state.BlockedSites["reddit"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockInstagram.Clicked(gtx) {
				state.BlockedSites["instagram"] = !state.BlockedSites["instagram"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockTikTok.Clicked(gtx) {
				state.BlockedSites["tiktok"] = !state.BlockedSites["tiktok"]
				updateBlocker(&b, state)
			}
			if settingsBtns.BlockWhatsapp.Clicked(gtx) {
				state.BlockedSites["whatsapp"] = !state.BlockedSites["whatsapp"]
				updateBlocker(&b, state)
			}
			if settingsBtns.AddWebsite.Clicked(gtx) {
				website := settingsBtns.WebsiteEditor.Text()
				if website != "" && website != "www.example.com" {
					state.CustomWebsites = append(state.CustomWebsites, website)
					settingsBtns.WebsiteEditor.SetText("")
					updateBlocker(&b, state)
				}
			}
			if settingsBtns.Theme.Clicked(gtx) {
				// TODO: toggle theme
			}

			ui.Layout(gtx, th, btns, settingsBtns, state)

			e.Frame(gtx.Ops)
		}
	}
}
