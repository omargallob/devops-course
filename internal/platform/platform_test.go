package platform

import (
	"runtime"
	"testing"
)

func TestSupportedPlatform(t *testing.T) {
	// This test compiles only on supported platforms due to build constraints.
	// If it compiles and runs, the platform is supported.
	supported := map[string][]string{
		"linux":  {"amd64", "arm64"},
		"darwin": {"amd64", "arm64"},
	}

	archs, ok := supported[runtime.GOOS]
	if !ok {
		t.Fatalf("unsupported OS: %s", runtime.GOOS)
	}

	found := false
	for _, a := range archs {
		if a == runtime.GOARCH {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("unsupported arch %s for OS %s", runtime.GOARCH, runtime.GOOS)
	}

	if OS != runtime.GOOS {
		t.Errorf("platform.OS = %q, want %q", OS, runtime.GOOS)
	}
	if Arch != runtime.GOARCH {
		t.Errorf("platform.Arch = %q, want %q", Arch, runtime.GOARCH)
	}
}
