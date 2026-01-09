package window

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type WindowHandler struct{}

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

func (h *WindowHandler) FindWindowByTitle(title string) uintptr {
	name, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(name)))
	return hwnd
}

func (h *WindowHandler) SetAlwaysOnTop(hwnd uintptr, enable bool) error {
	if hwnd == 0 {
		return nil
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

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
		if err != syscall.Errno(0) {
			return err
		}
		return os.ErrInvalid
	}
	return nil
}
