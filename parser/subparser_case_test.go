package parser

import (
	"os"
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

	runParseWithArgs(ap, []string{"prog", "GroupA", "--flag"})

	if len(ap.ParsingState.ErrorMessages) != 0 {
		t.Fatalf("expected no errors, got %v", ap.ParsingState.ErrorMessages)
	}
	if mode != "groupa" {
		t.Fatalf("expected mode=groupa, got %q", mode)
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
