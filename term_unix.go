//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package opts

import (
	"syscall"
	"unsafe"
)

type winsize struct {
	rows, cols, xpixels, ypixels uint16
}

// terminalSize returns the number of columns of the terminal attached to fd,
// whether fd is a terminal, and whether it accepts ANSI sequences, using the
// TIOCGWINSZ ioctl.
func terminalSize(fd uintptr) (int, bool, bool) {
	ws := winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return 0, false, false
	}
	return int(ws.cols), true, true
}
