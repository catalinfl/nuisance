package ui

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AppState struct {
	CurrentTab     int
	AlwaysOnTop    bool
	HwndValid      bool
	BlockedSites   map[string]bool
	PomodoroTime   string
	PomodoroMode   string
	WorkMinutes    int
	BreakMinutes   int
	CustomWebsites []string
	WebsiteInput   string
}

type Buttons struct {
	Tab1      *widget.Clickable
	Tab2      *widget.Clickable
	Toggle    *widget.Clickable
	PomoPlay  *widget.Clickable
	PomoPause *widget.Clickable
	PomoReset *widget.Clickable
}

type SettingsButtons struct {
	BlockFacebook  *widget.Clickable
	BlockYouTube   *widget.Clickable
	BlockTwitter   *widget.Clickable
	BlockReddit    *widget.Clickable
	BlockInstagram *widget.Clickable
	BlockTikTok    *widget.Clickable
	BlockWhatsapp  *widget.Clickable
	Theme          *widget.Clickable
	WorkInc        *widget.Clickable
	WorkDec        *widget.Clickable
	BreakInc       *widget.Clickable
	BreakDec       *widget.Clickable
	AddWebsite     *widget.Clickable
	WebsiteEditor  *widget.Editor
}

func NewButtons() *Buttons {
	return &Buttons{
		Tab1:      new(widget.Clickable),
		Tab2:      new(widget.Clickable),
		Toggle:    new(widget.Clickable),
		PomoPlay:  new(widget.Clickable),
		PomoPause: new(widget.Clickable),
		PomoReset: new(widget.Clickable),
	}
}

func NewSettingsButtons() *SettingsButtons {
	editor := new(widget.Editor)
	editor.SingleLine = true
	editor.Submit = true
	return &SettingsButtons{
		BlockFacebook:  new(widget.Clickable),
		BlockYouTube:   new(widget.Clickable),
		BlockTwitter:   new(widget.Clickable),
		BlockReddit:    new(widget.Clickable),
		BlockInstagram: new(widget.Clickable),
		BlockTikTok:    new(widget.Clickable),
		BlockWhatsapp:  new(widget.Clickable),
		Theme:          new(widget.Clickable),
		WorkInc:        new(widget.Clickable),
		WorkDec:        new(widget.Clickable),
		BreakInc:       new(widget.Clickable),
		BreakDec:       new(widget.Clickable),
		AddWebsite:     new(widget.Clickable),
		WebsiteEditor:  editor,
	}
}

// tabs with buttons
func TabBar(gtx layout.Context, th *material.Theme, btns *Buttons, currentTab int) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, btns.Tab1, "Pomodoro")
			btn.Inset = layout.UniformInset(unit.Dp(6))
			btn.TextSize = unit.Sp(11)
			if currentTab == 0 {
				btn.Background = th.Palette.ContrastBg
			}
			return btn.Layout(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(th, btns.Tab2, "Settings")
			btn.Inset = layout.UniformInset(unit.Dp(6))
			btn.TextSize = unit.Sp(11)
			if currentTab == 1 {
				btn.Background = th.Palette.ContrastBg
			}
			return btn.Layout(gtx)
		}),
	)
}

func TabContent(gtx layout.Context, th *material.Theme, btns *Buttons, state *AppState, currentTab int) layout.Dimensions {
	if currentTab == 0 {
		return PomodoroContent(gtx, th, btns, state)
	}
	return layout.Dimensions{}
}

func PomodoroContent(gtx layout.Context, th *material.Theme, btns *Buttons, state *AppState) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H4(th, state.PomodoroTime)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(th, state.PomodoroMode)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.PomoPlay, "Start")
						btn.Inset = layout.UniformInset(unit.Dp(6))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.PomoPause, "Pause")
						btn.Inset = layout.UniformInset(unit.Dp(6))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.PomoReset, "Reset")
						btn.Inset = layout.UniformInset(unit.Dp(6))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
				)
			}),
		)
	})
}

func SettingsContent(gtx layout.Context, th *material.Theme, mainBtns *Buttons, btns *SettingsButtons, state *AppState) layout.Dimensions {
	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, "Pomodoro Timer")
				label.TextSize = unit.Sp(14)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(th, "Work (min):")
						label.TextSize = unit.Sp(12)
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.WorkDec, "-")
						btn.Inset = layout.UniformInset(unit.Dp(4))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(th, fmt.Sprintf("%d", state.WorkMinutes))
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.WorkInc, "+")
						btn.Inset = layout.UniformInset(unit.Dp(4))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(th, "Break (min):")
						label.TextSize = unit.Sp(12)
						return label.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.BreakDec, "-")
						btn.Inset = layout.UniformInset(unit.Dp(4))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							label := material.Body1(th, fmt.Sprintf("%d", state.BreakMinutes))
							label.TextSize = unit.Sp(14)
							return label.Layout(gtx)
						})
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.BreakInc, "+")
						btn.Inset = layout.UniformInset(unit.Dp(4))
						btn.TextSize = unit.Sp(12)
						return btn.Layout(gtx)
					}),
				)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, "Block Websites")
				label.TextSize = unit.Sp(14)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockFacebook, "Facebook", state.BlockedSites["facebook"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockYouTube, "YouTube", state.BlockedSites["youtube"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockTwitter, "Twitter/X", state.BlockedSites["twitter"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockReddit, "Reddit", state.BlockedSites["reddit"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockInstagram, "Instagram", state.BlockedSites["instagram"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockTikTok, "TikTok", state.BlockedSites["tiktok"])
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return settingsButton(gtx, th, btns.BlockWhatsapp, "WhatsApp", state.BlockedSites["whatsapp"])
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(th, "Custom Websites:")
				label.TextSize = unit.Sp(12)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						ed := material.Editor(th, btns.WebsiteEditor, "www.example.com")
						ed.TextSize = unit.Sp(12)
						return ed.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(th, btns.AddWebsite, "+ Add")
						btn.Inset = layout.UniformInset(unit.Dp(4))
						btn.TextSize = unit.Sp(11)
						return btn.Layout(gtx)
					}),
				)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					customWebsiteList(th, state)...,
				)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, "Appearance")
				label.TextSize = unit.Sp(14)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, btns.Theme, "🎨 Toggle Theme")
				btn.Inset = layout.UniformInset(unit.Dp(4))
				btn.TextSize = unit.Sp(12)
				return btn.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.H6(th, "Window")
				label.TextSize = unit.Sp(14)
				return label.Layout(gtx)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return ToggleButton(gtx, th, mainBtns.Toggle, state.AlwaysOnTop)
			}),
		)
	})
}

func customWebsiteList(th *material.Theme, state *AppState) []layout.FlexChild {
	var children []layout.FlexChild
	for i, site := range state.CustomWebsites {
		idx := i
		website := site
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(th, "• "+website)
						label.TextSize = unit.Sp(11)
						return label.Layout(gtx)
					}),
				)
			})
		}))
		_ = idx
	}
	return children
}

func settingsButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, name string, isBlocked bool) layout.Dimensions {
	return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Clickable(gtx, btn, func(gtx layout.Context) layout.Dimensions {
			bg := th.Palette.Bg
			if isBlocked {
				bg = th.Palette.ContrastBg
			}

			return layout.Background{}.Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					defer op.Offset(image.Point{}).Push(gtx.Ops).Pop()
					paint.ColorOp{Color: bg}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					return layout.Dimensions{Size: gtx.Constraints.Min}
				},
				func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(6)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
							// Icon and name
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(layout.Spacer{Width: unit.Dp(6)}.Layout),
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										label := material.Body1(th, name)
										label.TextSize = unit.Sp(12)
										return label.Layout(gtx)
									}),
								)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								var statusText string
								if isBlocked {
									statusText = ""
								} else {
									statusText = "✓"
								}
								label := material.Body2(th, statusText)
								label.TextSize = unit.Sp(12)
								return label.Layout(gtx)
							}),
						)
					})
				},
			)
		})
	})
}

func ToggleButton(gtx layout.Context, th *material.Theme, btn *widget.Clickable, alwaysOnTop bool) layout.Dimensions {
	text := "Unpin"
	if !alwaysOnTop {
		text = "Pin on Top"
	}
	return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		b := material.Button(th, btn, text)
		b.Inset = layout.UniformInset(unit.Dp(4))
		b.TextSize = unit.Sp(12)
		return b.Layout(gtx)
	})
}

func Layout(gtx layout.Context, th *material.Theme, btns *Buttons, settingsBtns *SettingsButtons, state *AppState) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Tab buttons
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return TabBar(gtx, th, btns, state.CurrentTab)
		}),

		// Content area
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if state.CurrentTab == 1 {
				return SettingsContent(gtx, th, btns, settingsBtns, state)
			}
			return TabContent(gtx, th, btns, state, state.CurrentTab)
		}),
	)
}
