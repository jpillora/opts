package opts

import (
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

var ansiSequence = regexp.MustCompile(`\x1b\[[0-9;]*m`)

type highlightConfig struct {
	File string   `opts:"help=file to load"`
	Run  struct{} `opts:"mode=cmd, help=run the program"`
}

func highlightHelpText(t *testing.T, width int, isTTY bool) string {
	t.Helper()
	stubTerminal(t, width, isTTY)
	t.Setenv("TERM", "xterm-256color")
	t.Setenv("NO_COLOR", "")
	p, _ := New(&highlightConfig{}).Name("myapp").Version("v1.2.3").
		Author("Jaime").Repo("https://example.com/myapp").
		ParseArgsError([]string{"/bin/myapp", "--help"})
	return p.Help()
}

func TestHelpHighlightsTerminal(t *testing.T) {
	out := highlightHelpText(t, 80, true)
	want := []string{
		ansiBold + "Usage:" + ansiReset,
		ansiBoldCyan + "myapp" + ansiReset,
		ansiCyan + "[options]" + ansiReset,
		ansiCyan + "<command>" + ansiReset,
		ansiBold + "Options:" + ansiReset,
		ansiCyan + "--file, -f" + ansiReset,
		ansiBold + "Commands:" + ansiReset,
		ansiCyan + "run" + ansiReset,
		ansiBold + "Author:" + ansiReset,
		ansiCyan + "Jaime" + ansiReset,
		ansiBold + "Version:" + ansiReset,
		ansiCyan + "v1.2.3" + ansiReset,
		ansiBold + "Read more:" + ansiReset,
		ansiCyan + "https://example.com/myapp" + ansiReset,
	}
	for _, s := range want {
		if !strings.Contains(out, s) {
			t.Errorf("highlighted help does not contain %q:\n%q", s, out)
		}
	}
}

func TestHelpHighlightsErrors(t *testing.T) {
	stubTerminal(t, 80, true)
	t.Setenv("TERM", "xterm")
	t.Setenv("NO_COLOR", "")
	o, err := New(&struct{}{}).Name("myapp").ParseArgsError(
		[]string{"/bin/myapp", "--not-a-real-option"},
	)
	if err == nil {
		t.Fatal("expected invalid option error")
	}
	out := o.Help()
	if !strings.Contains(out, ansiBoldRed+"Error:"+ansiReset) {
		t.Fatalf("error heading is not highlighted:\n%q", out)
	}
	if !strings.Contains(out, ansiRed+"unknown flag: --not-a-real-option"+ansiReset) {
		t.Fatalf("error message is not highlighted:\n%q", out)
	}
}

func TestHelpDoesNotHighlightRedirectedOutput(t *testing.T) {
	out := highlightHelpText(t, 80, false)
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("redirected help contains ANSI escapes:\n%q", out)
	}
}

func TestHelpHighlightOptOut(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{name: "NO_COLOR", key: "NO_COLOR", value: "1"},
		{name: "dumb terminal", key: "TERM", value: "dumb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubTerminal(t, 80, true)
			t.Setenv("TERM", "xterm")
			t.Setenv("NO_COLOR", "")
			t.Setenv(tt.key, tt.value)
			p, _ := New(&highlightConfig{}).Name("myapp").
				ParseArgsError([]string{"/bin/myapp", "--help"})
			if out := p.Help(); strings.Contains(out, "\x1b[") {
				t.Fatalf("opted-out help contains ANSI escapes:\n%q", out)
			}
		})
	}
}

func TestHelpHighlightDoesNotChangeLayout(t *testing.T) {
	plain := highlightHelpText(t, 40, false)
	styled := highlightHelpText(t, 40, true)
	if stripped := ansiSequence.ReplaceAllString(styled, ""); stripped != plain {
		t.Fatalf("styling changed help layout:\nplain:   %q\nstyled:  %q\nstripped: %q", plain, styled, stripped)
	}
}

func TestHighlightedHelpFitsTerminal(t *testing.T) {
	width := 0
	prev := termInfo
	termInfo = func() terminalInfo {
		return terminalInfo{width: width, isTTY: true}
	}
	t.Cleanup(func() { termInfo = prev })
	t.Setenv("TERM", "xterm")
	t.Setenv("NO_COLOR", "")

	for width = 20; width <= 120; width++ {
		p, _ := New(&widthConfig{}).Name("myapp").
			ParseArgsError([]string{"/bin/myapp", "--help"})
		plain := ansiSequence.ReplaceAllString(p.Help(), "")
		for i, line := range strings.Split(plain, "\n") {
			//The usage line is deliberately kept as one semantic unit.
			if strings.HasPrefix(strings.TrimSpace(line), "Usage:") {
				continue
			}
			if columns := utf8.RuneCountInString(line); columns > width {
				t.Fatalf("terminal width %d: line %d is %d columns:\n%s",
					width, i+1, columns, line)
			}
		}
	}
}

func TestHelpHighlightRemovesTrailingAlignmentSpaces(t *testing.T) {
	type config struct {
		A               bool `opts:"short=-"`
		SomethingLonger bool `opts:"short=-, help=has help"`
	}
	stubTerminal(t, 80, true)
	t.Setenv("TERM", "xterm")
	t.Setenv("NO_COLOR", "")
	p, _ := New(&config{}).Name("myapp").
		ParseArgsError([]string{"/bin/myapp", "--help"})
	out := p.Help()
	if !strings.Contains(out, ansiCyan+"--a"+ansiReset+"\n") {
		t.Fatalf("option without help retains hidden trailing spaces:\n%q", out)
	}
}
