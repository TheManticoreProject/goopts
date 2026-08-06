package parser

import (
	"os"
	"os/exec"
	"testing"
)

// runParserSubprocess re-executes the given test in a subprocess with the
// GOOPTS_UNKNOWN_ARGUMENT_SUBPROCESS environment marker set, and returns the exit
// code of that subprocess. ParseFrom calls os.Exit when it records error messages,
// so the error paths have to be exercised out of process.
func runParserSubprocess(t *testing.T, testName string) int {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "GOOPTS_UNKNOWN_ARGUMENT_SUBPROCESS=1")
	err := cmd.Run()

	if err == nil {
		return 0
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("running %s in a subprocess failed: %s", testName, err)
	}
	return exitErr.ExitCode()
}

// TestUnknownArgumentReportsError verifies that flags matching no registered short or
// long name are surfaced as parse errors and cause ParseFrom to exit with a non-zero
// status, instead of being silently ignored.
func TestUnknownArgumentReportsError(t *testing.T) {
	if os.Getenv("GOOPTS_UNKNOWN_ARGUMENT_SUBPROCESS") == "1" {
		var verbose bool
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		if err := ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "verbose"); err != nil {
			// Registration failure is a different problem; exit 2 to distinguish.
			os.Exit(2)
		}
		ap.ParsingState.SetRawArguments([]string{"--typo-flag", "-x"})
		// Neither "--typo-flag" nor "-x" is registered: ParseFrom must record errors and os.Exit(1).
		ap.ParseFrom(0, &ap.ParsingState)
		// Reaching here means the unknown flags were swallowed (the bug): signal failure.
		os.Exit(0)
	}

	if code := runParserSubprocess(t, "TestUnknownArgumentReportsError"); code != 1 {
		if code == 0 {
			t.Fatalf("expected the parser to exit non-zero on unknown flags, but it exited 0 (unknown flags were silently ignored)")
		}
		t.Fatalf("expected exit code 1 from unknown flags, got %d", code)
	}
}

// TestKnownArgumentsWithFlagLikeValuesAreAccepted verifies that a value which looks like
// a flag, such as the negative integer in "--port -1", is recognized as the value of the
// preceding argument and not reported as an unknown flag.
func TestKnownArgumentsWithFlagLikeValuesAreAccepted(t *testing.T) {
	if os.Getenv("GOOPTS_UNKNOWN_ARGUMENT_SUBPROCESS") == "1" {
		var port int
		var name string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		if err := ap.NewIntArgument(&port, "-p", "--port", 0, false, "port"); err != nil {
			os.Exit(2)
		}
		if err := ap.NewStringArgument(&name, "-n", "--name", "", false, "name"); err != nil {
			os.Exit(2)
		}
		ap.ParsingState.SetRawArguments([]string{"--port", "-1", "-n", "-dash-value"})
		// Both values start with "-" but are operands of registered flags: no error expected.
		ap.ParseFrom(0, &ap.ParsingState)
		if port != -1 {
			os.Exit(3)
		}
		if name != "-dash-value" {
			os.Exit(4)
		}
		os.Exit(0)
	}

	switch code := runParserSubprocess(t, "TestKnownArgumentsWithFlagLikeValuesAreAccepted"); code {
	case 0:
	case 1:
		t.Fatalf("expected flag-like values of registered arguments to parse cleanly, but the parser exited 1 (they were reported as unknown arguments)")
	case 3:
		t.Fatalf("expected --port to be set to -1")
	case 4:
		t.Fatalf("expected -n to be set to \"-dash-value\"")
	default:
		t.Fatalf("unexpected exit code %d", code)
	}
}
