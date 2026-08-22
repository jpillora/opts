//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !windows

package opts

//terminalWidth is unsupported on this platform, help text falls back
//to the default line width.
func terminalWidth(fd uintptr) int {
	return 0
}
