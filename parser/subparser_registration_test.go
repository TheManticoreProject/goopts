package parser

import (
	"os"
	"os/exec"
	"testing"
)

// TestAddSubParserWithoutSetupSubParsing verifies that registering a subparser on a parser that
// never went through SetupSubParsing works. Before the fix it panicked with "assignment to entry
// in nil map", because only SetupSubParsing allocated the map that AddSubParser writes to.
func TestAddSubParserWithoutSetupSubParsing(t *testing.T) {
	ap := NewParser("test")

	sub := ap.AddSubParser("scan", "scan things")

	if sub == nil {
		t.Fatalf("expected AddSubParser to return a parser")
	}
	if got := len(ap.SubParsers.Parsers); got != 1 {
		t.Fatalf("expected 1 registered subparser, got %d", got)
	}
	if ap.SubParsers.Parsers["scan"] != sub {
		t.Fatalf("expected the returned parser to be registered under \"scan\"")
	}
	if !ap.SubParsers.Enabled {
		t.Fatalf("expected registering a subparser to enable subparsing, otherwise the subparser is silently ignored while parsing")
	}
}

// TestAddSubParserNestedWithoutSetupSubParsing verifies the nesting case: the parsers created by
// AddSubParser never go through SetupSubParsing, so before the fix no subparser could ever have a
// subparser of its own.
func TestAddSubParserNestedWithoutSetupSubParsing(t *testing.T) {
	var chosen string
	ap := NewParser("test")
	ap.SetupSubParsing("mode", &chosen, true)

	sub := ap.AddSubParser("scan", "scan things")
	nested := sub.AddSubParser("deep", "deep scan")

	if nested == nil {
		t.Fatalf("expected AddSubParser to return a parser")
	}
	if got := len(sub.SubParsers.Parsers); got != 1 {
		t.Fatalf("expected 1 subparser nested under \"scan\", got %d", got)
	}
	if sub.SubParsers.Parsers["deep"] != nested {
		t.Fatalf("expected the returned parser to be registered under \"deep\"")
	}
	if !sub.SubParsers.Enabled {
		t.Fatalf("expected nesting a subparser to enable subparsing on the parent subparser")
	}
}

// TestParseFromDispatchesWithoutSetupSubParsing verifies that a subparser registered without
// SetupSubParsing is actually dispatched to, and that the missing name pointer — SetupSubParsing is
// what binds one — does not cause a nil pointer dereference in ParseFrom.
//
// ParseFrom calls os.Exit on error, so the parse runs in a subprocess.
func TestParseFromDispatchesWithoutSetupSubParsing(t *testing.T) {
	if os.Getenv("GOOPTS_SUBPARSER_REGISTRATION_SUBPROCESS") == "1" {
		var host string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)

		sub := ap.AddSubParser("scan", "scan things")
		if err := sub.NewStringArgument(&host, "-H", "--host", "", true, "host"); err != nil {
			os.Exit(3)
		}

		ap.ParsingState.SetRawArguments([]string{"test", "scan", "--host", "h1"})
		ap.ParseFrom(1, &ap.ParsingState)

		if host != "h1" {
			os.Exit(4)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestParseFromDispatchesWithoutSetupSubParsing")
	cmd.Env = append(os.Environ(), "GOOPTS_SUBPARSER_REGISTRATION_SUBPROCESS=1")
	err := cmd.Run()

	if err == nil {
		return
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("running the parse in a subprocess failed: %s", err)
	}
	switch exitErr.ExitCode() {
	case 1:
		t.Fatalf("expected the parse to succeed, but the parser reported errors and exited 1")
	case 2:
		t.Fatalf("the parse panicked")
	case 3:
		t.Fatalf("registering the subparser argument failed")
	case 4:
		t.Fatalf("expected --host of the \"scan\" subparser to be set to \"h1\"")
	default:
		t.Fatalf("unexpected exit code %d", exitErr.ExitCode())
	}
}

// TestParseFromDispatchesToNestedSubParser verifies dispatch two levels deep, where the
// intermediate parser was registered by AddSubParser and never had SetupSubParsing called on it.
func TestParseFromDispatchesToNestedSubParser(t *testing.T) {
	if os.Getenv("GOOPTS_SUBPARSER_NESTED_SUBPROCESS") == "1" {
		var chosen, target string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)
		ap.SetupSubParsing("mode", &chosen, true)

		sub := ap.AddSubParser("scan", "scan things")
		nested := sub.AddSubParser("deep", "deep scan")
		if err := nested.NewStringArgument(&target, "-t", "--target", "", true, "target"); err != nil {
			os.Exit(3)
		}

		ap.ParsingState.SetRawArguments([]string{"test", "scan", "deep", "--target", "t1"})
		ap.ParseFrom(1, &ap.ParsingState)

		if chosen != "scan" {
			os.Exit(4)
		}
		if target != "t1" {
			os.Exit(5)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestParseFromDispatchesToNestedSubParser")
	cmd.Env = append(os.Environ(), "GOOPTS_SUBPARSER_NESTED_SUBPROCESS=1")
	err := cmd.Run()

	if err == nil {
		return
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("running the parse in a subprocess failed: %s", err)
	}
	switch exitErr.ExitCode() {
	case 1:
		t.Fatalf("expected the parse to succeed, but the parser reported errors and exited 1")
	case 2:
		t.Fatalf("the parse panicked")
	case 3:
		t.Fatalf("registering the nested subparser argument failed")
	case 4:
		t.Fatalf("expected the root subparser name to be recorded as \"scan\"")
	case 5:
		t.Fatalf("expected --target of the nested \"deep\" subparser to be set to \"t1\"")
	default:
		t.Fatalf("unexpected exit code %d", exitErr.ExitCode())
	}
}
