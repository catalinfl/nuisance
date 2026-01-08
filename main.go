package main

import (
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetWindowPos = user32.NewProc("SetWindowPos")
	procIsWindow     = user32.NewProc("IsWindow")
	procFindWindowW  = user32.NewProc("FindWindowW")
)

const (
	HWND_TOPMOST   = ^uintptr(0)
	HWND_NOTOPMOST = ^uintptr(1)
	SWP_NOMOVE     = 0x0002
	SWP_NOSIZE     = 0x0001
	SWP_SHOWWINDOW = 0x0040
)

func setAlwaysOnTop(hwnd uintptr, enable bool) error {
	if hwnd == 0 {
		return nil
	}
	// Ensure this syscall runs on an OS thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Verify window handle is still valid
	is, _, _ := procIsWindow.Call(hwnd)
	if is == 0 {
		return nil
	}

	var pos uintptr
	if enable {
		pos = HWND_TOPMOST
	} else {
		pos = HWND_NOTOPMOST
	}
	r1, _, err := procSetWindowPos.Call(hwnd, pos, 0, 0, 0, 0, SWP_NOMOVE|SWP_NOSIZE|SWP_SHOWWINDOW)
	if r1 == 0 {
		// err is syscall.Errno with GetLastError; can be 0 in some cases.
		if err != syscall.Errno(0) {
			return err
		}
		return os.ErrInvalid
	}
	return nil
}

func findWindowByTitle(title string) uintptr {
	name, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(name)))
	return hwnd
}

func main() {
	go runApp()
	app.Main()
}

func runApp() {
	w := new(app.Window)
	w.Option(
		app.Size(unit.Dp(400), unit.Dp(300)),
		app.Title("Always On Top"),
	)

	var hwnd atomic.Uintptr
	alwaysOnTop := true

	// Get HWND after short delay (avoid GetForegroundWindow; it can return the wrong window)
	go func() {
		time.Sleep(300 * time.Millisecond)
		for i := 0; i < 50; i++ {
			h := findWindowByTitle("Always On Top")
			if h != 0 {
				hwnd.Store(h)
				_ = setAlwaysOnTop(h, true)
				log.Println("Window set to always on top")
				w.Invalidate()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		log.Println("[WARN] Could not find window handle (FindWindowW)")
	}()

	th := material.NewTheme()
	var toggleBtn widget.Clickable
	var tab1Btn, tab2Btn, tab3Btn widget.Clickable
	currentTab := 0
	var ops op.Ops

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			os.Exit(0)

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
						_ = setAlwaysOnTop(handle, enable)
						w.Invalidate()
					}(h, newState)
					log.Printf("Always On Top: %v\n", alwaysOnTop)
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
