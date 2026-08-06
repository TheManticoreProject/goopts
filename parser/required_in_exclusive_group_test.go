package parser

import (
	"os"
	"os/exec"
	"testing"
)

// TestRequiredInRequiredMutuallyExclusiveGroupNotUnconditionallyRequired verifies that a
// member of a required mutually exclusive group that was registered with required: true is
// not additionally enforced as unconditionally required. Before the fix, populateMaps put it
// in ap.requiredArguments, which reported it as missing even when the group was satisfied by
// another member, and produced a second error line naming the same argument.
func TestRequiredInRequiredMutuallyExclusiveGroupNotUnconditionallyRequired(t *testing.T) {
	var mode, modeAlt string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)

	grp, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Mode")
	if err != nil {
		t.Fatalf("NewRequiredMutuallyExclusiveArgumentGroup failed: %v", err)
	}
	if err := grp.NewStringArgument(&mode, "-m", "--mode", "", true, "mode"); err != nil {
		t.Fatalf("group.NewStringArgument(--mode) failed: %v", err)
	}
	if err := grp.NewStringArgument(&modeAlt, "-M", "--mode-alt", "", false, "mode alt"); err != nil {
		t.Fatalf("group.NewStringArgument(--mode-alt) failed: %v", err)
	}

	ap.populateMaps(&ap.ParsingState)

	if got := len(ap.requiredArguments); got != 0 {
		t.Fatalf("expected no unconditionally required arguments for a mutually exclusive group, got %d", got)
	}
}

// TestRequiredInNotRequiredMutuallyExclusiveGroupNotUnconditionallyRequired verifies the same
// for the not required mutually exclusive group type, where enforcing a member's required flag
// would force that member and make every other member of the group unusable.
func TestRequiredInNotRequiredMutuallyExclusiveGroupNotUnconditionallyRequired(t *testing.T) {
	var a, b string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)

	grp, err := ap.NewNotRequiredMutuallyExclusiveArgumentGroup("Output")
	if err != nil {
		t.Fatalf("NewNotRequiredMutuallyExclusiveArgumentGroup failed: %v", err)
	}
	if err := grp.NewStringArgument(&a, "-o", "--out-file", "", true, "out file"); err != nil {
		t.Fatalf("group.NewStringArgument(--out-file) failed: %v", err)
	}
	if err := grp.NewStringArgument(&b, "-j", "--out-json", "", false, "out json"); err != nil {
		t.Fatalf("group.NewStringArgument(--out-json) failed: %v", err)
	}

	ap.populateMaps(&ap.ParsingState)

	if got := len(ap.requiredArguments); got != 0 {
		t.Fatalf("expected no unconditionally required arguments for a mutually exclusive group, got %d", got)
	}
}

// TestRequiredInDependentGroupStaysRequired verifies that the fix is limited to mutually
// exclusive groups: a required argument in a dependent group is coherent with the group rule,
// which forces every member once one is set, so it stays unconditionally required.
func TestRequiredInDependentGroupStaysRequired(t *testing.T) {
	var host, port string
	ap := NewParser("test")
	ap.SetOptShowBannerOnRun(false)

	grp, err := ap.NewDependentArgumentGroup("Proxy")
	if err != nil {
		t.Fatalf("NewDependentArgumentGroup failed: %v", err)
	}
	if err := grp.NewStringArgument(&host, "-x", "--proxy-host", "", true, "proxy host"); err != nil {
		t.Fatalf("group.NewStringArgument(--proxy-host) failed: %v", err)
	}
	if err := grp.NewStringArgument(&port, "-X", "--proxy-port", "", false, "proxy port"); err != nil {
		t.Fatalf("group.NewStringArgument(--proxy-port) failed: %v", err)
	}

	ap.populateMaps(&ap.ParsingState)

	if got := len(ap.requiredArguments); got != 1 {
		t.Fatalf("expected 1 unconditionally required argument in a dependent group, got %d", got)
	}
}

// TestExclusiveGroupSatisfiedByOtherMemberParses verifies the user visible consequence of the
// fix: supplying the other member of the group is enough for the parse to succeed. Before the
// fix, ParseFrom recorded "Missing required argument \"--mode\"" and exited 1, so --mode-alt
// could never be used.
//
// ParseFrom calls os.Exit(1) when it records error messages, so the parse runs in a subprocess
// and the exit code is inspected from the parent.
func TestExclusiveGroupSatisfiedByOtherMemberParses(t *testing.T) {
	if os.Getenv("GOOPTS_EXCLUSIVE_GROUP_SUBPROCESS") == "1" {
		var mode, modeAlt string
		ap := NewParser("test")
		ap.SetOptShowBannerOnRun(false)

		grp, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("Mode")
		if err != nil {
			// Registration failure is a different problem; exit 2 to distinguish.
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&mode, "-m", "--mode", "", true, "mode"); err != nil {
			os.Exit(2)
		}
		if err := grp.NewStringArgument(&modeAlt, "-M", "--mode-alt", "", false, "mode alt"); err != nil {
			os.Exit(2)
		}

		ap.ParsingState.SetRawArguments([]string{"--mode-alt", "X"})
		ap.ParseFrom(0, &ap.ParsingState)

		if modeAlt != "X" {
			os.Exit(3)
		}
		os.Exit(0)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExclusiveGroupSatisfiedByOtherMemberParses")
	cmd.Env = append(os.Environ(), "GOOPTS_EXCLUSIVE_GROUP_SUBPROCESS=1")
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
		t.Fatalf("expected the parse to succeed when the group is satisfied by --mode-alt, but the parser reported errors and exited 1")
	case 3:
		t.Fatalf("expected --mode-alt to be set to \"X\"")
	default:
		t.Fatalf("unexpected exit code %d", exitErr.ExitCode())
	}
}
