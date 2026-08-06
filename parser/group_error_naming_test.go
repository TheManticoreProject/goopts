package parser

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// runGroupErrorSubprocess re-executes the named test in a subprocess with the
// GOOPTS_GROUP_NAMING_SUBPROCESS marker set, and returns the "[!] " error lines it printed.
// ParseFrom prints its error messages and calls os.Exit(1), so they have to be collected from a
// separate process.
func runGroupErrorSubprocess(t *testing.T, testName string) []string {
	t.Helper()

	cmd := exec.Command(os.Args[0], "-test.run="+testName)
	cmd.Env = append(os.Environ(), "GOOPTS_GROUP_NAMING_SUBPROCESS=1")
	out, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("expected the parser to report an unsatisfied group and exit non-zero")
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
		t.Fatalf("registering the argument group failed")
	}

	lines := []string{}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "[!] ") {
			lines = append(lines, strings.TrimPrefix(line, "[!] "))
		}
	}

	return lines
}

// TestGroupErrorMessageIsPrefixedWithGroupName verifies that a constraint error names the group it
// came from, so that several groups produce distinguishable lines.
func TestGroupErrorMessageIsPrefixedWithGroupName(t *testing.T) {
	if os.Getenv("GOOPTS_GROUP_NAMING_SUBPROCESS") == "1" {
		var a, b string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)

		grp, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Source")
		if err != nil {
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&a, "-f", "--from-file", "", false, "h"); err != nil {
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&b, "-s", "--from-stdin", "", false, "h"); err != nil {
			os.Exit(2)
		}

		ap.ParsingState.SetRawArguments([]string{"test"})
		ap.ParseFrom(1, &ap.ParsingState)
		os.Exit(0)
	}

	lines := runGroupErrorSubprocess(t, "TestGroupErrorMessageIsPrefixedWithGroupName")

	expected := "Source: at least one of the arguments \"--from-file\", \"--from-stdin\" needs to be set."
	if len(lines) != 1 || lines[0] != expected {
		t.Fatalf("expected exactly one error line %q, got %v", expected, lines)
	}
}

// TestGroupErrorMessageWithoutGroupNameIsCapitalized verifies that a constraint group with no name
// keeps the unprefixed wording, with the message capitalized because it starts the line.
func TestGroupErrorMessageWithoutGroupNameIsCapitalized(t *testing.T) {
	if os.Getenv("GOOPTS_GROUP_NAMING_SUBPROCESS") == "1" {
		var a, b string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)

		grp, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("")
		if err != nil {
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&a, "-f", "--from-file", "", false, "h"); err != nil {
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&b, "-s", "--from-stdin", "", false, "h"); err != nil {
			os.Exit(2)
		}

		ap.ParsingState.SetRawArguments([]string{"test"})
		ap.ParseFrom(1, &ap.ParsingState)
		os.Exit(0)
	}

	lines := runGroupErrorSubprocess(t, "TestGroupErrorMessageWithoutGroupNameIsCapitalized")

	expected := "At least one of the arguments \"--from-file\", \"--from-stdin\" needs to be set."
	if len(lines) != 1 || lines[0] != expected {
		t.Fatalf("expected exactly one error line %q, got %v", expected, lines)
	}
}

// TestFormatGroupErrorMessage covers the message builder directly, including the empty message
// guard, which the ParseFrom call sites never hit because they always pass a literal.
func TestFormatGroupErrorMessage(t *testing.T) {
	tests := []struct {
		groupName string
		message   string
		expected  string
	}{
		{"Source", "at least one of the arguments \"--a\" needs to be set.", "Source: at least one of the arguments \"--a\" needs to be set."},
		{"", "at least one of the arguments \"--a\" needs to be set.", "At least one of the arguments \"--a\" needs to be set."},
		{"", "", ""},
		{"Proxy", "", "Proxy: "},
	}

	for _, test := range tests {
		if got := formatGroupErrorMessage(test.groupName, test.message); got != test.expected {
			t.Fatalf("formatGroupErrorMessage(%q, %q) = %q, expected %q", test.groupName, test.message, got, test.expected)
		}
	}
}
