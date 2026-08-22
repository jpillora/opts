//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd && !windows

package opts

// terminalSize is unsupported on this platform, help text falls back to
// the default line width and is never highlighted.
func terminalSize(fd uintptr) (int, bool) {
	return 0, false
}
