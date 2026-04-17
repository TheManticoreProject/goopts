package parser

import (
	"testing"
)

func TestPopulateMaps_DefaultGroupRegisteredOnce(t *testing.T) {
	ap := NewParser("t")

	var a, b int
	if err := ap.NewIntArgument(&a, "a", "aaa", 0, true, "h"); err != nil {
		t.Fatalf("NewIntArgument(a) failed: %v", err)
	}
	if err := ap.NewIntArgument(&b, "", "bbb", 0, false, "h"); err != nil {
		t.Fatalf("NewIntArgument(b) failed: %v", err)
	}

	ap.populateMaps()

	if got := len(ap.allArguments); got != 2 {
		t.Fatalf("expected 2 arguments in allArguments, got %d", got)
	}
	if got := len(ap.requiredArguments); got != 1 {
		t.Fatalf("expected 1 argument in requiredArguments, got %d", got)
	}

	// Calling populateMaps again must not grow the slices.
	ap.populateMaps()
	ap.populateMaps()

	if got := len(ap.allArguments); got != 2 {
		t.Fatalf("expected 2 arguments in allArguments after repeated calls, got %d", got)
	}
	if got := len(ap.requiredArguments); got != 1 {
		t.Fatalf("expected 1 argument in requiredArguments after repeated calls, got %d", got)
	}
}

func TestPopulateMaps_MixedGroupsCountedOnce(t *testing.T) {
	ap := NewParser("t")

	var a int
	if err := ap.NewIntArgument(&a, "", "a", 0, true, "h"); err != nil {
		t.Fatalf("NewIntArgument failed: %v", err)
	}

	grp, err := ap.NewArgumentGroup("G")
	if err != nil {
		t.Fatalf("NewArgumentGroup failed: %v", err)
	}
	var b int
	if err := grp.NewIntArgument(&b, "", "b", 0, true, "h"); err != nil {
		t.Fatalf("group.NewIntArgument failed: %v", err)
	}

	ap.populateMaps()

	if got := len(ap.allArguments); got != 2 {
		t.Fatalf("expected 2 arguments in allArguments, got %d", got)
	}
	if got := len(ap.requiredArguments); got != 2 {
		t.Fatalf("expected 2 required arguments, got %d", got)
	}
}
