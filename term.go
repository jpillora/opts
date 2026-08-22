package opts

import (
	"os"
	"strconv"
)

type terminalInfo struct {
	width int
	isTTY bool
}

// termInfo returns the width of an attached output terminal, and whether an
// output terminal was found at all. It is a variable so tests can simulate
// both terminal dimensions and redirected output.
var termInfo = detectTermInfo

// detectTermInfo inspects the standard output streams for terminal dimensions,
// falling back to the conventional COLUMNS environment variable for width.
// COLUMNS alone does not make redirected output a terminal.
func detectTermInfo() terminalInfo {
	//help text is written to stderr, though stdout is also checked
	//since Help() may be printed there by the program itself
	isTTY := false
	for _, f := range []*os.File{os.Stderr, os.Stdout} {
		if f == nil {
			continue
		}
		w, ok := terminalSize(f.Fd())
		if !ok {
			continue
		}
		isTTY = true
		if w > 0 {
			return terminalInfo{width: w, isTTY: true}
		}
	}
	//COLUMNS is the conventional override, and the only option
	//left when all of the standard streams have been redirected
	if s := os.Getenv("COLUMNS"); s != "" {
		if w, err := strconv.Atoi(s); err == nil && w > 0 {
			return terminalInfo{width: w, isTTY: isTTY}
		}
	}
	return terminalInfo{isTTY: isTTY}
}
