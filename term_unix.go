//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package opts

import (
	"syscall"
	"unsafe"
)

type winsize struct {
	rows, cols, xpixels, ypixels uint16
}

//terminalWidth returns the number of columns of the terminal attached
//to fd, using the TIOCGWINSZ ioctl. Returns zero when fd is not a terminal.
func terminalWidth(fd uintptr) int {
	ws := winsize{}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)),
	)
	if errno != 0 {
		return 0
	}
	return int(ws.cols)
}
