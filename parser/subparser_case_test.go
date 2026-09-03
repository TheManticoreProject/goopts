package parser

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runParseWithArgs(ap *ArgumentsParser, args []string) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = args
	ap.Parse()
}

func TestSubparserCaseInsensitiveDispatch(t *testing.T) {
	ap := NewParser("t")
	var mode string
	ap.SetupSubParsing("mode", &mode, true)

	var flag bool
	sub := ap.AddSubParser("GroupA", "…")
	if err := sub.NewBoolArgument(&flag, "", "--flag", false, "…"); err != nil {
		t.Fatalf("NewBoolArgument failed: %v", err)
	}

	runParseWithArgs(ap, []string{"prog", "groupa", "--flag"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors, got %v", ap.ParsingState.ErrorMessages)
	}
	// The name stored is the registered one, not the one that was typed on the command line.
	if mode != "GroupA" {
		t.Fatalf("expected mode=GroupA, got %q", mode)
	}
	if !flag {
		t.Fatalf("expected flag=true, got false")
	}
}

func TestSubparserCaseSensitiveDispatchUnchanged(t *testing.T) {
	ap := NewParser("t")
	var mode string
	ap.SetupSubParsing("mode", &mode, false)

	var flag bool
	sub := ap.AddSubParser("GroupA", "…")
	if err := sub.NewBoolArgument(&flag, "", "--flag", false, "…"); err != nil {
		t.Fatalf("NewBoolArgument failed: %v", err)
	}

	runParseWithArgs(ap, []string{"prog", "GroupA", "--flag"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors in case-sensitive mode with exact name, got %v", ap.ParsingState.ErrorMessages)
	}
	if mode != "GroupA" {
		t.Fatalf("expected mode=GroupA, got %q", mode)
	}
}

// captureStdout runs fn and returns everything it wrote to os.Stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	origStdout := os.Stdout
	os.Stdout = writer

	fn()

	os.Stdout = origStdout
	if err := writer.Close(); err != nil {
		t.Fatalf("closing the pipe failed: %v", err)
	}
	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("reading the pipe failed: %v", err)
	}

	return string(output)
}

// TestSubparserNamesAreRegisteredWithTheirOwnCase verifies that enabling case-insensitive matching
// does not rewrite the name a subparser is registered under. Before the fix AddSubParser lowercased
// the map key, which is what leaked into the usage output and into the value bound by
// SetupSubParsing.
func TestSubparserNamesAreRegisteredWithTheirOwnCase(t *testing.T) {
	ap := NewParser("t")
	var mode string
	ap.SetupSubParsing("mode", &mode, true)

	sub := ap.AddSubParser("GroupA", "…")

	if got := ap.SubParsers.Parsers["GroupA"]; got != sub {
		t.Fatalf("expected the subparser to be registered under \"GroupA\", got %v", ap.SubParsers.Parsers)
	}
	if _, exists := ap.SubParsers.Parsers["groupa"]; exists {
		t.Fatalf("expected no lowercased key, got %v", ap.SubParsers.Parsers)
	}
}

// TestSubparserUsageKeepsRegisteredCase verifies that the usage output lists the subparsers with
// the name they were registered with, even when matching is case-insensitive.
func TestSubparserUsageKeepsRegisteredCase(t *testing.T) {
	ap := NewParser("t")
	var mode string
	ap.SetupSubParsing("mode", &mode, true)
	ap.AddSubParser("GroupA", "the A group")
	ap.AddSubParser("GroupZ", "the Z group")

	parsingState := &ParsingState{RawArguments: []string{"prog"}}
	usage := captureStdout(t, func() { ap.UsageFrom(1, parsingState) })

	if !strings.Contains(usage, "<GroupA|GroupZ>") {
		t.Fatalf("expected the usage line to list the registered names, got:\n%s", usage)
	}
	if strings.Contains(usage, "groupa") || strings.Contains(usage, "groupz") {
		t.Fatalf("expected no lowercased names in the usage, got:\n%s", usage)
	}
}

// TestSetSubparserNameCaseSensitiveAfterRegistration verifies that the setting can be changed after
// the subparsers are registered. Before the fix the names were normalized when they were
// registered, so a call made after AddSubParser left the keys in their original case and no
// case-insensitive lookup could ever match them.
func TestSetSubparserNameCaseSensitiveAfterRegistration(t *testing.T) {
	ap := NewParser("t")
	var mode string
	ap.SetupSubParsing("mode", &mode, false)

	var flag bool
	sub := ap.AddSubParser("GroupA", "…")
	if err := sub.NewBoolArgument(&flag, "", "--flag", false, "…"); err != nil {
		t.Fatalf("NewBoolArgument failed: %v", err)
	}

	// Registered first, configured second
	ap.SetSubparserNameCaseSensitive(false)

	runParseWithArgs(ap, []string{"prog", "GROUPA", "--flag"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors, got %v", ap.ParsingState.ErrorMessages)
	}
	if mode != "GroupA" {
		t.Fatalf("expected mode=GroupA, got %q", mode)
	}
	if !flag {
		t.Fatalf("expected flag=true, got false")
	}
}

// TestSubparserNameIsCaseSensitiveByDefault verifies the default, and that the setter and its
// getter agree on the polarity of the setting.
func TestSubparserNameIsCaseSensitiveByDefault(t *testing.T) {
	ap := NewParser("t")

	if !ap.SubparserNameIsCaseSensitive() {
		t.Fatalf("expected subparser names to be case-sensitive by default")
	}

	ap.SetSubparserNameCaseSensitive(false)
	if ap.SubparserNameIsCaseSensitive() {
		t.Fatalf("expected subparser names to be case-insensitive after SetSubparserNameCaseSensitive(false)")
	}

	ap.SetSubparserNameCaseSensitive(true)
	if !ap.SubparserNameIsCaseSensitive() {
		t.Fatalf("expected subparser names to be case-sensitive after SetSubparserNameCaseSensitive(true)")
	}
}

// TestSubparserCaseSensitivityIsInherited verifies that the setting reaches the whole subparser
// tree: the subparsers registered after the call inherit it, and the ones registered before it are
// updated by it.
func TestSubparserCaseSensitivityIsInherited(t *testing.T) {
	ap := NewParser("t")
	var mode, nestedMode string
	ap.SetupSubParsing("mode", &mode, false)

	registeredBefore := ap.AddSubParser("GroupA", "…")
	registeredBefore.SetupSubParsing("nested_mode", &nestedMode, false)
	registeredBefore.AddSubParser("GroupAB", "…")

	ap.SetSubparserNameCaseSensitive(false)

	registeredAfter := ap.AddSubParser("GroupZ", "…")
	if registeredAfter.SubparserNameIsCaseSensitive() {
		t.Fatalf("expected a subparser registered after the call to inherit case-insensitive matching")
	}
	if registeredBefore.SubparserNameIsCaseSensitive() {
		t.Fatalf("expected a subparser registered before the call to be updated by it")
	}

	var flag bool
	nested := registeredBefore.SubParsers.GetSubParser("GroupAB")
	if nested == nil {
		t.Fatalf("expected to find the nested subparser")
	}
	if err := nested.NewBoolArgument(&flag, "", "--flag", false, "…"); err != nil {
		t.Fatalf("NewBoolArgument failed: %v", err)
	}

	runParseWithArgs(ap, []string{"prog", "groupa", "groupab", "--flag"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors, got %v", ap.ParsingState.ErrorMessages)
	}
	if mode != "GroupA" {
		t.Fatalf("expected mode=GroupA, got %q", mode)
	}
	if nestedMode != "GroupAB" {
		t.Fatalf("expected nested_mode=GroupAB, got %q", nestedMode)
	}
	if !flag {
		t.Fatalf("expected flag=true, got false")
	}
}

// TestSubparserCaseSensitivityCanBeOverriddenPerSubparser verifies that a subparser can keep
// case-sensitive names while its parent matches without regard to case.
//
// ParseFrom calls os.Exit on error, so the parse runs in a subprocess.
func TestSubparserCaseSensitivityCanBeOverriddenPerSubparser(t *testing.T) {
	if os.Getenv("GOOPTS_SUBPARSER_CASE_OVERRIDE_SUBPROCESS") == "1" {
		var mode, nestedMode string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		ap.SetupSubParsing("mode", &mode, false)

		sub := ap.AddSubParser("GroupA", "…")
		sub.SetupSubParsing("nested_mode", &nestedMode, false)
		sub.AddSubParser("GroupAB", "…")

		ap.SetSubparserNameCaseSensitive(false)
		sub.SetSubparserNameCaseSensitive(true)

		ap.ParsingState.SetRawArguments([]string{"test", "groupa", "groupab"})
		ap.ParseFrom(1, &ap.ParsingState)

		// The case-sensitive subparser rejects "groupab", so the parse exits before this point
		os.Exit(4)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSubparserCaseSensitivityCanBeOverriddenPerSubparser")
	cmd.Env = append(os.Environ(), "GOOPTS_SUBPARSER_CASE_OVERRIDE_SUBPROCESS=1")
	output, err := cmd.CombinedOutput()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected the parse to report an error and exit 1, got %v with output:\n%s", err, output)
	}
	if exitErr.ExitCode() == 4 {
		t.Fatalf("expected the case-sensitive subparser to reject \"groupab\", output:\n%s", output)
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("unexpected exit code %d, output:\n%s", exitErr.ExitCode(), output)
	}
	if !strings.Contains(string(output), "No subparser with name \"groupab\" was found.") {
		t.Fatalf("expected the nested subparser to report the unknown name, output:\n%s", output)
	}
	// The usage carries the path the parse took, so it shows the case-insensitive parent dispatched
	if !strings.Contains(string(output), "Usage: test groupa <GroupAB>") {
		t.Fatalf("expected the case-insensitive parent to have dispatched to \"GroupA\", output:\n%s", output)
	}
}

// TestGetSubParserHonoursCaseSensitivity verifies that GetSubParser matches names the same way the
// dispatch does, and reports a miss with a nil parser.
func TestGetSubParserHonoursCaseSensitivity(t *testing.T) {
	ap := NewParser("t")
	sub := ap.AddSubParser("GroupA", "…")

	if got := ap.SubParsers.GetSubParser("GroupA"); got != sub {
		t.Fatalf("expected an exact name to match in case-sensitive mode")
	}
	if got := ap.SubParsers.GetSubParser("groupa"); got != nil {
		t.Fatalf("expected a differently cased name not to match in case-sensitive mode, got %v", got)
	}

	ap.SetSubparserNameCaseSensitive(false)

	if got := ap.SubParsers.GetSubParser("groupa"); got != sub {
		t.Fatalf("expected a differently cased name to match in case-insensitive mode")
	}
	if got := ap.SubParsers.GetSubParser("groupb"); got != nil {
		t.Fatalf("expected an unknown name not to match, got %v", got)
	}
}

// TestSubparserCaseInsensitiveMatchIsDeterministic verifies the ambiguous case: when two subparsers
// differ only by case, the exact name still reaches its own subparser, and a name that matches both
// resolves to the same one on every lookup instead of following map iteration order.
func TestSubparserCaseInsensitiveMatchIsDeterministic(t *testing.T) {
	ap := NewParser("t")
	ap.SetSubparserNameCaseSensitive(false)

	lower := ap.AddSubParser("scan", "…")
	upper := ap.AddSubParser("Scan", "…")

	if got := ap.SubParsers.GetSubParser("scan"); got != lower {
		t.Fatalf("expected an exact name to win over a case-insensitive match")
	}
	if got := ap.SubParsers.GetSubParser("Scan"); got != upper {
		t.Fatalf("expected an exact name to win over a case-insensitive match")
	}

	// "SCAN" matches both, and must resolve to the same one every time
	first := ap.SubParsers.GetSubParser("SCAN")
	if first == nil {
		t.Fatalf("expected an ambiguous name to match one of the subparsers")
	}
	for i := 0; i < 32; i++ {
		if got := ap.SubParsers.GetSubParser("SCAN"); got != first {
			t.Fatalf("expected an ambiguous name to resolve to the same subparser on every lookup")
		}
	}
}

// TestSetupSubParsingKeepsRegisteredSubparsers verifies that calling SetupSubParsing after the
// subparsers are registered binds the name value without discarding them.
func TestSetupSubParsingKeepsRegisteredSubparsers(t *testing.T) {
	ap := NewParser("t")
	sub := ap.AddSubParser("GroupA", "…")

	var mode string
	ap.SetupSubParsing("mode", &mode, true)

	if got := len(ap.SubParsers.Parsers); got != 1 {
		t.Fatalf("expected the registered subparser to be kept, got %d subparsers", got)
	}
	if ap.SubParsers.Parsers["GroupA"] != sub {
		t.Fatalf("expected the registered subparser to be kept under its own name")
	}

	runParseWithArgs(ap, []string{"prog", "groupa"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors, got %v", ap.ParsingState.ErrorMessages)
	}
	if mode != "GroupA" {
		t.Fatalf("expected mode=GroupA, got %q", mode)
	}
}

// TestSubparserUnknownNameErrorEchoesTheGivenName verifies that the error names what was typed on
// the command line, not a normalized form of it.
//
// ParseFrom calls os.Exit on error, so the parse runs in a subprocess.
func TestSubparserUnknownNameErrorEchoesTheGivenName(t *testing.T) {
	if os.Getenv("GOOPTS_SUBPARSER_UNKNOWN_NAME_SUBPROCESS") == "1" {
		var mode string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		ap.SetupSubParsing("mode", &mode, true)
		ap.AddSubParser("GroupA", "…")

		ap.ParsingState.SetRawArguments([]string{"test", "GroupB"})
		ap.ParseFrom(1, &ap.ParsingState)

		// "GroupB" matches no subparser, so the parse exits before this point
		os.Exit(4)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestSubparserUnknownNameErrorEchoesTheGivenName")
	cmd.Env = append(os.Environ(), "GOOPTS_SUBPARSER_UNKNOWN_NAME_SUBPROCESS=1")
	output, err := cmd.CombinedOutput()

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected the parse to report an error and exit 1, got %v with output:\n%s", err, output)
	}
	if exitErr.ExitCode() == 4 {
		t.Fatalf("expected the parse to reject the unknown name \"GroupB\", output:\n%s", output)
	}
	if exitErr.ExitCode() != 1 {
		t.Fatalf("unexpected exit code %d, output:\n%s", exitErr.ExitCode(), output)
	}
	if !strings.Contains(string(output), "No subparser with name \"GroupB\" was found.") {
		t.Fatalf("expected the error to quote the given name, output:\n%s", output)
	}
}
