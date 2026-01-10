package sound

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	winmm         = syscall.NewLazyDLL("winmm.dll")
	procPlaySound = winmm.NewProc("PlaySoundW")
)

const (
	SND_FILENAME = 0x00020000
	SND_ASYNC    = 0x0001
	SND_LOOP     = 0x0008
)

type AlarmPlayer struct {
	stopChan chan struct{}
	playing  bool
}

func NewAlarmPlayer() *AlarmPlayer {
	return &AlarmPlayer{
		stopChan: make(chan struct{}),
		playing:  false,
	}
}

func (ap *AlarmPlayer) Stop() {
	if ap.playing {
		ap.playing = false
		close(ap.stopChan)
		ap.stopChan = make(chan struct{})
		StopSound()
	}
}

func (ap *AlarmPlayer) PlayRepeating(soundPath string) {
	if ap.playing {
		return
	}
	ap.playing = true
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		PlaySound(soundPath)

		for {
			select {
			case <-ticker.C:
				PlaySound(soundPath)
			case <-ap.stopChan:
				return
			}
		}
	}()
}

func PlaySound(soundPath string) error {
	if soundPath == "" {
		return nil
	}

	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		soundPath = "C:\\Windows\\Media\\Windows Notify System Generic.wav"
	}

	soundPtr, err := syscall.UTF16PtrFromString(soundPath)
	if err != nil {
		return err
	}

	procPlaySound.Call(
		uintptr(unsafe.Pointer(soundPtr)),
		0,
		SND_ASYNC|SND_FILENAME,
	)
	return nil
}

func StopSound() {
	procPlaySound.Call(0, 0, 0)
}

func GetSoundPath(filename string) string {
	exePath, err := os.Executable()
	if err != nil {
		return filename
	}
	exeDir := filepath.Dir(exePath)
	soundPath := filepath.Join(exeDir, "sounds", filename)
	return soundPath
}

func PlayWorkStart() {
	PlaySound(GetSoundPath("work_alarm.mp3"))
}

func PlayBreakStart() {
	PlaySound(GetSoundPath("break_alarm.mp3"))
}

func PlayComplete() {
	PlaySound(GetSoundPath("complete.mp3"))
}

func PlayButton() {
	PlaySound(GetSoundPath("button.mp3"))
}
