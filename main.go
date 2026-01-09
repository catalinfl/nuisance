package main

import (
	"sync"
	"sync/atomic"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/catalinfl/nuisance/httpblock"
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
	alwaysOnTop := true

	var cleanupOnce sync.Once
	cleanup := func() {
		_ = b.RemoveBlockEntries()
	}

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-sigCh
	// 	cleanupOnce.Do(cleanup)
	// 	w.Perform(system.ActionClose)
	// }()

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
	var toggleBtn widget.Clickable
	var tab1Btn, tab2Btn, tab3Btn widget.Clickable
	currentTab := 0
	var ops op.Ops

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			cleanupOnce.Do(cleanup)
			return

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if tab1Btn.Clicked(gtx) {
				currentTab = 0
			}
			if tab2Btn.Clicked(gtx) {
				currentTab = 1
			}
			if tab3Btn.Clicked(gtx) {
				currentTab = 2
			}

			if toggleBtn.Clicked(gtx) {
				h := hwnd.Load()
				if h != 0 {
					alwaysOnTop = !alwaysOnTop
					newState := alwaysOnTop
					go func(handle uintptr, enable bool) {
						_ = winHandler.SetAlwaysOnTop(handle, enable)
						w.Invalidate()
					}(h, newState)
				}
			}

			layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{}.Layout(gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &tab1Btn, "Pomodoro")
							if currentTab == 0 {
								btn.Background = th.Palette.ContrastBg
							}
							return btn.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &tab2Btn, "Clock")
							if currentTab == 1 {
								btn.Background = th.Palette.ContrastFg
							}
							return btn.Layout(gtx)
						}),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							btn := material.Button(th, &tab3Btn, "Settings")
							if currentTab == 2 {
								btn.Background = th.Palette.ContrastBg
							}
							return btn.Layout(gtx)
						}),
					)
				}),

				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						var text string
						switch currentTab {
						case 0:
							text = "Tab 1\n\nThe window stays\nalways on top!"
							return material.H5(th, text).Layout(gtx)
						case 1:
							text = "Tab 2\n\nOpen other programs,\nthis one won't disappear."
						case 2:
							text = "Tab 3\n\nSwitch between tabs!"
						}
						return material.H5(th, text).Layout(gtx)
					})
				}),

				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					text := "Deactivate Always On Top"
					if !alwaysOnTop {
						text = "Activate Always On Top"
					}
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Button(th, &toggleBtn, text).Layout(gtx)
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}
