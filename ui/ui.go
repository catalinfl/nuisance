package ui

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AppState struct {
	CurrentTab  int
	AlwaysOnTop bool
	HwndValid   bool
}

type Buttons struct {
	Tab1   *widget.Clickable
	Tab2   *widget.Clickable
	Tab3   *widget.Clickable
	Toggle *widget.Clickable
}

func NewButtons() *Buttons {
	return &Buttons{
		Tab1:   new(widget.Clickable),
		Tab2:   new(widget.Clickable),
		Tab3:   new(widget.Clickable),
		Toggle: new(widget.Clickable),
	}
}

// tabs with buttons
func TabBar(gtx layout.Context, th *material.Theme, btns *Buttons, currentTab int) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, btns.Tab1, "Pomodoro")
			if currentTab == 0 {
				btn.Background = th.Palette.ContrastBg
			}
			return btn.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, btns.Tab2, "Clock")
			if currentTab == 1 {
				btn.Background = th.Palette.ContrastBg
			}
			return btn.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, btns.Tab3, "Settings")
			if currentTab == 2 {
				btn.Background = th.Palette.ContrastBg
			}
			return btn.Layout(gtx)
		}),
	)
}

func TabContent(gtx layout.Context, th *material.Theme, currentTab int) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		var text string
		switch currentTab {
		case 0:
			text = "Pomodoro Timer\n\nFocus on your work!\nThe window stays always on top."
		case 1:
			text = "Clock\n\nTrack your time\nwhile working on other tasks."
		case 2:
			text = "Settings\n\nConfigure your preferences\nand blocked sites."
		}
		return material.H5(th, text).Layout(gtx)
	})
}

func ToggleButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, alwaysOnTop bool) layout.Dimensions {
	text := "Deactivate Always On Top"
	if !alwaysOnTop {
		text = "Activate Always On Top"
	}
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Button(th, btn, text).Layout(gtx)
	})
}

func Layout(gtx layout.Context, th *material.Theme, btns *Buttons, state *AppState) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Tab buttons
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return TabBar(gtx, th, btns, state.CurrentTab)
		}),

		// Content area
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return TabContent(gtx, th, state.CurrentTab)
		}),
	)
}
