# Nuisance

Nuisance is a tiny Windows Pomodoro and website blocker built with Gio.
It blocks distracting websites, refuses to leave your screen, and forces you to focus until the timer ends.
You’re supposed to hate it.

## What it does

- Pomodoro timer (Work / Break / Pause / Idle) with GUI controls (Start / Pause / Reset)
- Blocks websites (configurable list + custom sites) while in Work mode
- Plays sounds from a `sounds/` folder (work/break/button/complete)
- Pomodoro background image from `image/background.png` or `.jpg`
- Custom timing

## Project layout

- `main.go` — app entry, state wiring, event loop
- `ui/` — UI components (`ui.go`) and layout
- `pomodoro/` — timer logic (`pomodoro.go`)
- `sound/` — Windows sound wrapper (`sound.go`) using WinMM (plays files or falls back to system sounds)
- `httpblock/` — website blocker helper
- `window/` — OS window helpers (always-on-top, etc.)
- `sounds/` (runtime) — place your `work_alarm.mp3`, `break_alarm.mp3`, `button.mp3`, `complete.mp3`
- `image/` (runtime) — place `background.png` or `background.jpg` for the Pomodoro background
- `SETUP.md` — quick setup for assets

## Configuration & assets

Create the following folders next to the executable (or in your project root while running locally):

```
nuisance.exe
sounds/
  ├── work_alarm.mp3
  ├── break_alarm.mp3
  ├── button.mp3
  └── complete.mp3
image/
  └── background.png   (or background.jpg)
```

- If sound files are missing the app falls back to Windows system notification sounds.
- The background image is optional; when present it is scaled with "cover" behavior and clipped to the Pomodoro content pane (so it won't overlap the tabs).

See `SETUP.md` for a concise setup checklist.

## Build & run (Windows)

1. Install required Go toolchain.
2. Build:

```powershell
go build -ldflags "-H 'windowsgui'" -o nuisance.exe .
```

3. (Optional) Embed an icon into the executable with `rsrc`:

```powershell
go install github.com/akavel/rsrc@latest
rsrc -arch amd64 -ico assets/icon.ico -o rsrc.syso
go build -ldflags "-H 'windowsgui'" -o nuisance.exe .
```

Place `sounds/` and `image/` next to `nuisance.exe` before running.

## Quick tips & troubleshooting

- Run as administrator to block websites. It access /etc/hosts, so you it needs admin to block websites.
- After closing /etc/hosts switch back to normal 