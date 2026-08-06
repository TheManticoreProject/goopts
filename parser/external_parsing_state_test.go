package parser

import (
	"testing"
)

// TestParseFromExternalParsingStateNamedArgument verifies that ParseFrom initializes the
// ParsedArguments maps of the ParsingState it was given, instead of the parser's own
// ap.ParsingState. Before the fix, recording the successfully parsed "-n bob" panicked
// with "assignment to entry in nil map" because the written-to maps were never allocated.
func TestParseFromExternalParsingStateNamedArgument(t *testing.T) {
	var name string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)
	if err := ap.NewStringArgument(&name, "-n", "--name", "default", false, "name"); err != nil {
		t.Fatalf("NewStringArgument failed: %v", err)
	}

	ps := &ParsingState{}
	ps.SetRawArguments([]string{"-n", "bob"})
	ap.ParseFrom(0, ps)

	if name != "bob" {
		t.Fatalf("expected --name to be %q, got %q", "bob", name)
	}
	if got := len(ps.GetErrorMessages()); got != 0 {
		t.Fatalf("expected no error messages, got %d: %v", got, ps.GetErrorMessages())
	}
	if _, exists := ps.ParsedArguments.ShortNameToArgument["-n"]; !exists {
		t.Fatalf("expected \"-n\" to be recorded in the parsing state that was passed to ParseFrom")
	}
	if _, exists := ps.ParsedArguments.LongNameToArgument["--name"]; !exists {
		t.Fatalf("expected \"--name\" to be recorded in the parsing state that was passed to ParseFrom")
	}
}

// TestParseFromExternalParsingStatePositionalArgument verifies the same invariant for
// positional arguments, which are recorded through AddPositionalArgument.
func TestParseFromExternalParsingStatePositionalArgument(t *testing.T) {
	var count int
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)
	if err := ap.NewIntPositionalArgument(&count, "count", "the count"); err != nil {
		t.Fatalf("NewIntPositionalArgument failed: %v", err)
	}

	ps := &ParsingState{}
	ps.SetRawArguments([]string{"7"})
	ap.ParseFrom(0, ps)

	if count != 7 {
		t.Fatalf("expected count to be 7, got %d", count)
	}
	if got := len(ps.GetErrorMessages()); got != 0 {
		t.Fatalf("expected no error messages, got %d: %v", got, ps.GetErrorMessages())
	}
	if _, exists := ps.ParsedArguments.PositionalArguments["count"]; !exists {
		t.Fatalf("expected \"count\" to be recorded in the parsing state that was passed to ParseFrom")
	}
}
