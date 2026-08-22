package opts

import (
	"os"
	"strconv"
)

//termWidth returns the width of the attached terminal in columns, or zero
//when it cannot be detected (for example, when the help text is being piped
//into another program). It is a variable so tests can simulate terminals.
var termWidth = detectTermWidth

//detectTermWidth inspects the standard streams for terminal dimensions,
//falling back to the conventional COLUMNS environment variable.
func detectTermWidth() int {
	//help text is written to stderr, though stdout is also checked
	//since Help() may be printed there by the program itself
	for _, f := range []*os.File{os.Stderr, os.Stdout} {
		if f == nil {
			continue
		}
		if w := terminalWidth(f.Fd()); w > 0 {
			return w
		}
	}
	//COLUMNS is the conventional override, and the only option
	//left when all of the standard streams have been redirected
	if s := os.Getenv("COLUMNS"); s != "" {
		if w, err := strconv.Atoi(s); err == nil && w > 0 {
			return w
		}
	}
	return 0
}
