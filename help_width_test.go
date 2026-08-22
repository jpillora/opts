package opts

import (
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

// TestMain ensures the help text tests below do not depend on the
// dimensions of the terminal they happen to be run in.
func TestMain(m *testing.M) {
	termWidth = func() int { return 0 }
	os.Exit(m.Run())
}

// stubTermWidth simulates a terminal of the given width for the duration
// of the test. A width of zero simulates no terminal at all.
func stubTermWidth(t *testing.T, w int) {
	prev := termWidth
	termWidth = func() int { return w }
	t.Cleanup(func() { termWidth = prev })
}

// widthConfig has help text long enough to wrap in a narrow terminal
type widthConfig struct {
	Listen  string   `opts:"help=the address to listen on. it may be a host or a port or both"`
	Retries int      `opts:"help=number of times to retry before giving up"`
	Serve   struct{} `opts:"mode=cmd, help=serve the current directory over http until interrupted"`
	Backup  struct{} `opts:"mode=cmd, help=back everything up"`
}

func widthHelp(t *testing.T, term int, with func(Opts) Opts) string {
	t.Helper()
	stubTermWidth(t, term)
	o := New(&widthConfig{}).Name("myapp")
	if with != nil {
		o = with(o)
	}
	p, _ := o.ParseArgsError([]string{"/bin/prog", "--help"})
	return p.Help()
}

// a wide terminal is capped at the default line width
func TestHelpAutoWide(t *testing.T) {
	check(t, widthHelp(t, 200, nil), `
  Usage: myapp [options] <command>

  Options:
  --listen, -l   the address to listen on. it may be a host or a port or both
  --retries, -r  number of times to retry before giving up
  --help, -h     display help

  Commands:
  · backup  back everything up
  · serve   serve the current directory over http until interrupted

`)
}

// no terminal (piped output) falls back to the default line width
func TestHelpAutoNoTerminal(t *testing.T) {
	check(t, widthHelp(t, 0, nil), widthHelp(t, 200, nil))
}

// a medium terminal wraps both columns to fit
func TestHelpAutoMedium(t *testing.T) {
	check(t, widthHelp(t, 60, nil), `
  Usage: myapp [options] <command>

  Options:
  --listen, -l   the address to listen on. it may be a host
                 or a port or both
  --retries, -r  number of times to retry before giving up
  --help, -h     display help

  Commands:
  · backup  back everything up
  · serve   serve the current directory over http until
            interrupted

`)
}

// a narrow terminal stacks the help text under each name. note the two
// lists are measured independently, the short command names still leave
// enough room for a second column
func TestHelpAutoNarrow(t *testing.T) {
	check(t, widthHelp(t, 36, nil), `
  Usage: myapp [options] <command>

  Options:
  --listen, -l
    the address to listen on. it may
    be a host or a port or both
  --retries, -r
    number of times to retry before
    giving up
  --help, -h
    display help

  Commands:
  · backup  back everything up
  · serve   serve the current
            directory over http
            until interrupted

`)
}

// an explicit line width is always respected, terminal or not
func TestHelpExplicitWidth(t *testing.T) {
	with := func(o Opts) Opts { return o.SetLineWidth(50) }
	expect := `
  Usage: myapp [options] <command>

  Options:
  --listen, -l   the address to listen on. it may be
                 a host or a port or both
  --retries, -r  number of times to retry before
                 giving up
  --help, -h     display help

  Commands:
  · backup  back everything up
  · serve   serve the current directory over http
            until interrupted

`
	//same output in a wide terminal, a narrow terminal and no terminal
	check(t, widthHelp(t, 200, with), expect)
	check(t, widthHelp(t, 30, with), expect)
	check(t, widthHelp(t, 0, with), expect)
}

// an explicit line width set on the root is inherited by subcommands
func TestHelpExplicitWidthInherited(t *testing.T) {
	stubTermWidth(t, 200)
	type sub struct {
		Path string `opts:"help=where the backup should be written to. defaults to the working directory"`
	}
	p, _ := New(&struct{}{}).Name("myapp").SetLineWidth(50).
		AddCommand(New(&sub{}).Name("backup")).
		ParseArgsError([]string{"/bin/prog", "backup"})
	check(t, p.Selected().Help(), `
  Usage: myapp backup [options]

  Options:
  --path, -p  where the backup should be written to.
              defaults to the working directory
  --help, -h  display help

`)
}

// help text must never be wider than the terminal rendering it
func TestHelpFitsTerminal(t *testing.T) {
	type fits struct {
		Listen  string   `opts:"help=the address to listen on. it may be a host or a port or both"`
		Verbose bool     `opts:"help=log everything that happens while the program is running"`
		Run     struct{} `opts:"mode=cmd, help=run the thing until it is done or something goes wrong"`
		Ps      struct{} `opts:"mode=cmd, help=list what is currently running"`
	}
	for term := 20; term <= 120; term++ {
		stubTermWidth(t, term)
		p, _ := New(&fits{}).Name("myapp").Version("1.2.3").
			Summary("a program which exists only to be rendered at many different widths").
			ParseArgsError([]string{"/bin/prog", "--help"})
		for i, line := range strings.Split(p.Help(), "\n") {
			//the usage line is a single unit, it is never wrapped
			if strings.HasPrefix(strings.TrimSpace(line), "Usage:") {
				continue
			}
			//the "·" command bullet is multi-byte, count columns
			if w := utf8.RuneCountInString(line); w > term {
				t.Fatalf("terminal width %d: line %d is %d chars:\n%s",
					term, i+1, w, line)
			}
		}
	}
}

// a name too long for the terminal cannot be wrapped, though its help
// text still is
func TestHelpNameWiderThanTerminal(t *testing.T) {
	type long struct {
		SomeVeryLongFlagName bool `opts:"help=a flag with a name longer than the terminal is wide"`
	}
	stubTermWidth(t, 24)
	p, _ := New(&long{}).Name("myapp").ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, p.Help(), `
  Usage: myapp [options]

  Options:
  --some-very-long-flag-name, -s
    a flag with a name
    longer than the
    terminal is wide
  --help, -h
    display help

`)
}

// disabling padAll gives back the columns it was reserving
func TestHelpAutoNoPadAll(t *testing.T) {
	stubTermWidth(t, 40)
	p, _ := New(&widthConfig{}).Name("myapp").DisablePadAll().
		ParseArgsError([]string{"/bin/prog", "--help"})
	check(t, p.Help(), `Usage: myapp [options] <command>

Options:
--listen, -l   the address to listen on.
               it may be a host or a
               port or both
--retries, -r  number of times to retry
               before giving up
--help, -h     display help

Commands:
· backup  back everything up
· serve   serve the current directory
          over http until interrupted
`)
}
