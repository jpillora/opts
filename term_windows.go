//go:build windows

package opts

import (
	"syscall"
	"unsafe"

	"github.com/mattn/go-isatty"
)

var (
	kernel32                       = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
	procGetConsoleMode             = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode             = kernel32.NewProc("SetConsoleMode")

	isCygwinTerminal           = isatty.IsCygwinTerminal
	getConsoleScreenBufferInfo = func(fd uintptr, info *consoleScreenBufferInfo) bool {
		ret, _, _ := procGetConsoleScreenBufferInfo.Call(fd, uintptr(unsafe.Pointer(info)))
		return ret != 0
	}
	getConsoleMode = func(fd uintptr, mode *uint32) bool {
		ret, _, _ := procGetConsoleMode.Call(fd, uintptr(unsafe.Pointer(mode)))
		return ret != 0
	}
	setConsoleMode = func(fd uintptr, mode uint32) bool {
		ret, _, _ := procSetConsoleMode.Call(fd, uintptr(mode))
		return ret != 0
	}
)

const (
	enableProcessedOutput           = 0x0001
	enableVirtualTerminalProcessing = 0x0004
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

// terminalSize returns the number of columns of the terminal attached to fd,
// whether fd is a terminal, and whether it accepts ANSI sequences. Windows
// consoles require virtual-terminal processing to be explicitly enabled.
// Cygwin and MSYS2 terminals use pipes rather than Windows console handles,
// but their pseudo-terminals accept ANSI sequences directly.
func terminalSize(fd uintptr) (int, bool, bool) {
	if isCygwinTerminal(fd) {
		return 0, true, true
	}

	var mode uint32
	if !getConsoleMode(fd, &mode) {
		return 0, false, false
	}

	info := consoleScreenBufferInfo{}
	width := 0
	if getConsoleScreenBufferInfo(fd, &info) {
		width = int(info.window.right-info.window.left) + 1
	}
	return width, true, enableVirtualTerminal(fd, mode)
}

func enableVirtualTerminal(fd uintptr, mode uint32) bool {
	want := mode | enableProcessedOutput | enableVirtualTerminalProcessing
	return want == mode || setConsoleMode(fd, want)
}
