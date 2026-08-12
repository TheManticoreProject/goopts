package parser

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// runBoundsSubprocess re-executes the named test in a subprocess with the
// GOOPTS_BOUNDS_SUBPROCESS marker set to mode, and returns its exit code. A panic in the
// subprocess surfaces as exit code 2, which is what distinguishes the defect from the expected
// exit code 1 of a reported parse error.
func runBoundsSubprocess(t *testing.T, testName, mode string) int {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "GOOPTS_BOUNDS_SUBPROCESS="+mode)
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

// parserWithRequiredArgument builds a parser with a single required argument, used by the
// subprocess halves of the tests below. Registration is a different concern from the bounds
// checks under test, so a failure there exits 3 rather than looking like a reported parse error.
func parserWithRequiredArgument() *ArgumentsParser {
	var name string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)
	if err := ap.NewStringArgument(&name, "-n", "--name", "", true, "name"); err != nil {
		os.Exit(3)
	}

	return ap
}

// TestParseFromEmptyRawArgumentsDoesNotPanic verifies that empty raw arguments are reported through
// the normal error path instead of panicking, both for index 0 and for index 1, the index Parse()
// uses. Before the fix, index 0 panicked in UsageFrom on RawArguments[0] and index 1 panicked in
// ParseFrom on the RawArguments[index:] slice.
func TestParseFromEmptyRawArgumentsDoesNotPanic(t *testing.T) {
	switch os.Getenv("GOOPTS_BOUNDS_SUBPROCESS") {
	case "index0", "index1":
		index := 0
		if os.Getenv("GOOPTS_BOUNDS_SUBPROCESS") == "index1" {
			index = 1
		}
		ap := parserWithRequiredArgument()
		ps := &ParsingState{}
		ps.SetRawArguments([]string{})
		// The required argument is missing, so this must print usage and exit 1, not panic.
		ap.ParseFrom(index, ps)
		os.Exit(0)
	}

	for _, mode := range []string{"index0", "index1"} {
		switch code := runBoundsSubprocess(t, "TestParseFromEmptyRawArgumentsDoesNotPanic", mode); code {
		case 1:
			// Expected: the missing required argument was reported.
		case 2:
			t.Fatalf("%s: ParseFrom panicked on empty raw arguments", mode)
		case 3:
			t.Fatalf("%s: registering the required argument failed", mode)
		default:
			t.Fatalf("%s: expected exit code 1 from a missing required argument, got %d", mode, code)
		}
	}
}

// TestUsageFromIndexPastEndDoesNotPanic verifies that an index beyond the end of the raw arguments
// only limits the subparser prefix of the usage line. Before the fix, UsageFrom read
// RawArguments[k+1] for every k up to index-2 and panicked.
func TestUsageFromIndexPastEndDoesNotPanic(t *testing.T) {
	if os.Getenv("GOOPTS_BOUNDS_SUBPROCESS") == "usage" {
		ap := parserWithRequiredArgument()
		ps := &ParsingState{}
		ps.SetRawArguments([]string{"prog"})
		ap.UsageFrom(5, ps)
		os.Exit(0)
	}

	switch code := runBoundsSubprocess(t, "TestUsageFromIndexPastEndDoesNotPanic", "usage"); code {
	case 0:
		// Expected: usage printed, no panic, no exit.
	case 2:
		t.Fatalf("UsageFrom panicked on an index past the end of the raw arguments")
	case 3:
		t.Fatalf("registering the required argument failed")
	default:
		t.Fatalf("expected UsageFrom to return without exiting, got exit code %d", code)
	}
}

// TestProgramName covers the program name fallbacks used by the usage line.
func TestProgramName(t *testing.T) {
	ps := &ParsingState{}
	ps.SetRawArguments([]string{filepath.Join("usr", "local", "bin", "tool"), "-n", "x"})
	if got := programName(ps); got != "tool" {
		t.Fatalf("expected the base name of the first raw argument, got %q", got)
	}

	// With no raw arguments the running executable is the best available name.
	empty := &ParsingState{}
	expected := filepath.Base(os.Args[0])
	if got := programName(empty); got != expected {
		t.Fatalf("expected the executable name %q, got %q", expected, got)
	}
}
