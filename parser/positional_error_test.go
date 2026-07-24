package parser

import (
	"os"
	"os/exec"
	"testing"
)

// TestPositionalIntInvalidValueReportsError verifies that an invalid value for an
// integer positional argument is surfaced as a parse error and causes ParseFrom to
// exit with a non-zero status, instead of being silently ignored.
//
// ParseFrom calls os.Exit(1) when it records error messages, so the error path is
// exercised in a subprocess and the exit code is inspected from the parent.
func TestPositionalIntInvalidValueReportsError(t *testing.T) {
	if os.Getenv("GOOPTS_POSITIONAL_ERROR_SUBPROCESS") == "1" {
		var count int
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		if err := ap.NewIntPositionalArgument(&count, "count", "the count"); err != nil {
			// Registration failure is a different problem; exit 2 to distinguish.
			os.Exit(2)
		}
		ps := &ParsingState{}
		ps.SetRawArguments([]string{"abc"})
		// "abc" cannot be parsed as an integer: ParseFrom must record an error and os.Exit(1).
		ap.ParseFrom(0, ps)
		// Reaching here means the error was swallowed (the bug): signal failure.
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestPositionalIntInvalidValueReportsError")
	cmd.Env = append(os.Environ(), "GOOPTS_POSITIONAL_ERROR_SUBPROCESS=1")
	err := cmd.Run()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected the parser to exit non-zero on an invalid integer positional, but it exited 0 (error was silently ignored)")
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("expected exit code 1 from a positional parse error, got %d", exitErr.ExitCode())
	}
}
