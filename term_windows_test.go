//go:build windows

package opts

import "testing"

func stubWindowsTerminalAPIs(t *testing.T) {
	t.Helper()
	oldIsCygwinTerminal := isCygwinTerminal
	oldGetConsoleScreenBufferInfo := getConsoleScreenBufferInfo
	oldGetConsoleMode := getConsoleMode
	oldSetConsoleMode := setConsoleMode
	t.Cleanup(func() {
		isCygwinTerminal = oldIsCygwinTerminal
		getConsoleScreenBufferInfo = oldGetConsoleScreenBufferInfo
		getConsoleMode = oldGetConsoleMode
		setConsoleMode = oldSetConsoleMode
	})
}

func TestWindowsTerminalEnablesVirtualTerminalProcessing(t *testing.T) {
	stubWindowsTerminalAPIs(t)
	isCygwinTerminal = func(uintptr) bool { return false }
	getConsoleScreenBufferInfo = func(_ uintptr, info *consoleScreenBufferInfo) bool {
		info.window.left = 4
		info.window.right = 83
		return true
	}
	getConsoleMode = func(_ uintptr, mode *uint32) bool {
		*mode = 0x0002
		return true
	}
	var setMode uint32
	setConsoleMode = func(_ uintptr, mode uint32) bool {
		setMode = mode
		return true
	}

	width, isTTY, supportsANSI := terminalSize(42)
	if width != 80 || !isTTY || !supportsANSI {
		t.Fatalf("terminalSize() = (%d, %t, %t), want (80, true, true)", width, isTTY, supportsANSI)
	}
	wantMode := uint32(0x0002 | enableProcessedOutput | enableVirtualTerminalProcessing)
	if setMode != wantMode {
		t.Fatalf("SetConsoleMode mode = %#x, want %#x", setMode, wantMode)
	}
}

func TestWindowsTerminalWithoutVirtualTerminalSupportStaysPlain(t *testing.T) {
	stubWindowsTerminalAPIs(t)
	isCygwinTerminal = func(uintptr) bool { return false }
	getConsoleScreenBufferInfo = func(_ uintptr, info *consoleScreenBufferInfo) bool {
		info.window.right = 79
		return true
	}
	getConsoleMode = func(_ uintptr, mode *uint32) bool { return true }
	setConsoleMode = func(uintptr, uint32) bool { return false }

	width, isTTY, supportsANSI := terminalSize(42)
	if width != 80 || !isTTY || supportsANSI {
		t.Fatalf("terminalSize() = (%d, %t, %t), want (80, true, false)", width, isTTY, supportsANSI)
	}
}

func TestWindowsTerminalKeepsExistingVirtualTerminalMode(t *testing.T) {
	stubWindowsTerminalAPIs(t)
	isCygwinTerminal = func(uintptr) bool { return false }
	getConsoleScreenBufferInfo = func(uintptr, *consoleScreenBufferInfo) bool { return false }
	getConsoleMode = func(_ uintptr, mode *uint32) bool {
		*mode = enableProcessedOutput | enableVirtualTerminalProcessing
		return true
	}
	setConsoleMode = func(uintptr, uint32) bool {
		t.Fatal("SetConsoleMode called for an already configured terminal")
		return false
	}

	width, isTTY, supportsANSI := terminalSize(42)
	if width != 0 || !isTTY || !supportsANSI {
		t.Fatalf("terminalSize() = (%d, %t, %t), want (0, true, true)", width, isTTY, supportsANSI)
	}
}

func TestWindowsCygwinTerminalAcceptsANSI(t *testing.T) {
	stubWindowsTerminalAPIs(t)
	isCygwinTerminal = func(uintptr) bool { return true }
	getConsoleMode = func(uintptr, *uint32) bool {
		t.Fatal("Windows console API called for a Cygwin terminal")
		return false
	}

	width, isTTY, supportsANSI := terminalSize(42)
	if width != 0 || !isTTY || !supportsANSI {
		t.Fatalf("terminalSize() = (%d, %t, %t), want (0, true, true)", width, isTTY, supportsANSI)
	}
}

func TestWindowsRedirectedOutputIsNotTerminal(t *testing.T) {
	stubWindowsTerminalAPIs(t)
	isCygwinTerminal = func(uintptr) bool { return false }
	getConsoleMode = func(uintptr, *uint32) bool { return false }

	width, isTTY, supportsANSI := terminalSize(42)
	if width != 0 || isTTY || supportsANSI {
		t.Fatalf("terminalSize() = (%d, %t, %t), want (0, false, false)", width, isTTY, supportsANSI)
	}
}
