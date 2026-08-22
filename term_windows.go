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

//terminalWidth returns the number of columns of the console attached to fd.
//Returns zero when fd is not a console.
func terminalWidth(fd uintptr) int {
	info := consoleScreenBufferInfo{}
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(&info)))
	if ret == 0 {
		return 0
	}
	return int(info.window.right-info.window.left) + 1
}
