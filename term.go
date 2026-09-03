package opts

import (
	"os"
	"strconv"
)

type terminalInfo struct {
	width        int
	isTTY        bool
	supportsANSI bool
}

// termInfo returns the width and capabilities of an attached output terminal.
// It is a variable so tests can simulate terminal dimensions, ANSI support,
// and redirected output.
var termInfo = detectTermInfo

// detectTermInfo inspects the standard output streams for terminal dimensions,
// falling back to the conventional COLUMNS environment variable for width.
// COLUMNS alone does not make redirected output a terminal.
func detectTermInfo() terminalInfo {
	//help text is written to stderr, though stdout is also checked
	//since Help() may be printed there by the program itself
	terminal := terminalInfo{}
	for _, f := range []*os.File{os.Stderr, os.Stdout} {
		if f == nil {
			continue
		}
		w, isTTY, supportsANSI := terminalSize(f.Fd())
		if !isTTY {
			continue
		}
		terminal.isTTY = true
		terminal.supportsANSI = terminal.supportsANSI || supportsANSI
		if terminal.width == 0 && w > 0 {
			terminal.width = w
		}
	}
	//COLUMNS is the conventional override, and the only option
	//left when terminal dimensions could not be detected
	if terminal.width == 0 {
		s := os.Getenv("COLUMNS")
		if w, err := strconv.Atoi(s); err == nil && w > 0 {
			terminal.width = w
		}
	}
	return terminal
}
