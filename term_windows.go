//go:build windows

package opts

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
)

type coord struct {
	x, y int16
}

type smallRect struct {
	left, top, right, bottom int16
}

type consoleScreenBufferInfo struct {
	size              coord
	cursorPosition    coord
	attributes        uint16
	window            smallRect
	maximumWindowSize coord
}

// terminalSize returns the number of columns of the console attached to fd,
// and whether fd is a console at all.
func terminalSize(fd uintptr) (int, bool) {
	info := consoleScreenBufferInfo{}
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return 0, false
	}
	return int(info.window.right-info.window.left) + 1, true
}
