package main

import (
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
	"github.com/catalinfl/nuisance/ui"
	"github.com/catalinfl/nuisance/window"
)

func main() {
	var handler window.WindowHandler = window.WindowHandler{}
	go runApp(&handler)
	app.Main()
}

func runApp(winHandler *window.WindowHandler) {
	w := new(app.Window)
	w.Option(
		app.Size(unit.Dp(400), unit.Dp(300)),
		app.Title("Always On Top"),
	)

	b := httpblock.Blocker{
		Token: "nuisance",
		Sites: []string{"www.facebook.com", "www.youtube.com"},
	}

	var hwnd atomic.Uintptr

	var cleanupOnce sync.Once
	cleanup := func() {
		_ = b.RemoveBlockEntries()
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

	b.RemoveBlockEntries()
	b.AddBlockEntries()

	th := material.NewTheme()
	btns := ui.NewButtons()
	state := &ui.AppState{
		CurrentTab:  0,
		AlwaysOnTop: true,
	}
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
			}
			if btns.Tab2.Clicked(gtx) {
				state.CurrentTab = 1
			}
			if btns.Tab3.Clicked(gtx) {
				state.CurrentTab = 2
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

			// Render UI
			ui.Layout(gtx, th, btns, state)

			e.Frame(gtx.Ops)
		}
	}
}
