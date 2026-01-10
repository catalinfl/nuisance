package sound

import (
	"syscall"
	"unsafe"
)

var (
	winmm         = syscall.NewLazyDLL("winmm.dll")
	procPlaySound = winmm.NewProc("PlaySoundW")
)

const (
	SND_FILENAME = 0x00020000
	SND_ASYNC    = 0x0001
)

func PlaySystemSound(soundName string) error {
	var soundPtr *uint16
	if soundName != "" {
		var err error
		soundPtr, err = syscall.UTF16PtrFromString(soundName)
		if err != nil {
			return err
		}
	}

	procPlaySound.Call(
		uintptr(unsafe.Pointer(soundPtr)),
		0,
		SND_ASYNC|SND_FILENAME,
	)
	return nil
}

func PlayWorkStart() {
	PlaySystemSound("C:\\Windows\\Media\\Windows Notify System Generic.wav")
}

func PlayBreakStart() {
	PlaySystemSound("C:\\Windows\\Media\\Windows Notify Calendar.wav")
}

func PlayComplete() {
	PlaySystemSound("C:\\Windows\\Media\\Windows Notify Messaging.wav")
}
