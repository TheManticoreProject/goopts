package parser

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestRequiredArgumentsOrderIsStable verifies that populateMaps collects required arguments in a
// stable order. Before the fix it ranged over the ap.Groups map, so the order of the names joined
// into the "Missing required arguments" message changed between runs.
//
// populateMaps is called repeatedly: Go randomizes map iteration order on every range statement,
// so a map based implementation would produce a wrong order within a few iterations.
func TestRequiredArgumentsOrderIsStable(t *testing.T) {
	var def, aOne, aTwo, bOne string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)

	if err := ap.NewStringArgument(&def, "-d", "--default-group-arg", "", true, "h"); err != nil {
		t.Fatalf("NewStringArgument failed: %v", err)
	}
	// Registered in reverse alphabetical order to make sure the ordering comes from sorting
	// the group names and not from the registration sequence.
	groupB, err := ap.NewArgumentGroup("Bravo")
	if err != nil {
		t.Fatalf("NewArgumentGroup(Bravo) failed: %v", err)
	}
	if err := groupB.NewStringArgument(&bOne, "-b", "--bravo-one", "", true, "h"); err != nil {
		t.Fatalf("group.NewStringArgument failed: %v", err)
	}
	groupA, err := ap.NewArgumentGroup("Alpha")
	if err != nil {
		t.Fatalf("NewArgumentGroup(Alpha) failed: %v", err)
	}
	if err := groupA.NewStringArgument(&aOne, "-1", "--alpha-one", "", true, "h"); err != nil {
		t.Fatalf("group.NewStringArgument failed: %v", err)
	}
	if err := groupA.NewStringArgument(&aTwo, "-2", "--alpha-two", "", true, "h"); err != nil {
		t.Fatalf("group.NewStringArgument failed: %v", err)
	}

	// The default group is keyed by "", which sorts before any named group.
	expected := []string{"--default-group-arg", "--alpha-one", "--alpha-two", "--bravo-one"}

	for i := 0; i < 50; i++ {
		ap.populateMaps(&ap.ParsingState)

		got := make([]string, 0, len(ap.requiredArguments))
		for _, arg := range ap.requiredArguments {
			got = append(got, arg.GetLongName())
		}
		if strings.Join(got, " ") != strings.Join(expected, " ") {
			t.Fatalf("iteration %d: expected required arguments in order %v, got %v", i, expected, got)
		}
	}
}

// TestGroupErrorMessageOrderIsStable verifies that the group constraint errors are reported in
// alphabetical group order rather than in map iteration order. The groups below are named "Source",
// "Output" and "Mode", so their messages must come out as Mode, Output, Source, identified here by
// the arguments each message names.
//
// ParseFrom prints the messages and calls os.Exit(1), so the parse runs in a subprocess and its
// output is inspected from the parent.
func TestGroupErrorMessageOrderIsStable(t *testing.T) {
	if os.Getenv("GOOPTS_ERROR_ORDER_SUBPROCESS") == "1" {
		var a1, a2, b1, b2, c1, c2 string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)

		source, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Source")
		if err != nil {
			os.Exit(2)
		}
		if err := source.NewStringArgument(&a1, "-f", "--from-file", "", false, "h"); err != nil {
			os.Exit(2)
		}
		if err := source.NewStringArgument(&a2, "-s", "--from-stdin", "", false, "h"); err != nil {
			os.Exit(2)
		}

		output, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Output")
		if err != nil {
			os.Exit(2)
		}
		if err := output.NewStringArgument(&b1, "-o", "--out-file", "", false, "h"); err != nil {
			os.Exit(2)
		}
		if err := output.NewStringArgument(&b2, "-j", "--out-json", "", false, "h"); err != nil {
			os.Exit(2)
		}

		mode, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Mode")
		if err != nil {
			os.Exit(2)
		}
		if err := mode.NewStringArgument(&c1, "-m", "--mode", "", false, "h"); err != nil {
			os.Exit(2)
		}
		if err := mode.NewStringArgument(&c2, "-M", "--mode-alt", "", false, "h"); err != nil {
			os.Exit(2)
		}

		// Mirrors Parse(): raw arguments start with the program name and parsing starts at 1.
		ap.ParsingState.SetRawArguments([]string{"test"})
		ap.ParseFrom(1, &ap.ParsingState)
		os.Exit(0)
	}

	// Each subprocess gets a fresh map iteration order, so a few runs are enough to catch
	// an implementation that depends on it.
	for run := 0; run < 5; run++ {
		cmd := exec.Command(os.Args[0], "-test.run=TestGroupErrorMessageOrderIsStable")
		cmd.Env = append(os.Environ(), "GOOPTS_ERROR_ORDER_SUBPROCESS=1")
		out, err := cmd.CombinedOutput()
		if err == nil {
			t.Fatalf("run %d: expected the parser to report the unsatisfied groups and exit non-zero", run)
		}
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			t.Fatalf("run %d: registering the argument groups failed", run)
		}

		got := []string{}
		for _, line := range strings.Split(string(out), "\n") {
			if !strings.HasPrefix(line, "[!] ") {
				continue
			}
			switch {
			case strings.Contains(line, "--mode"):
				got = append(got, "Mode")
			case strings.Contains(line, "--out-file"):
				got = append(got, "Output")
			case strings.Contains(line, "--from-file"):
				got = append(got, "Source")
			default:
				t.Fatalf("run %d: unexpected error line: %s", run, line)
			}
		}

		expected := "Mode Output Source"
		if strings.Join(got, " ") != expected {
			t.Fatalf("run %d: expected group errors in order %q, got %q", run, expected, strings.Join(got, " "))
		}
	}
}
