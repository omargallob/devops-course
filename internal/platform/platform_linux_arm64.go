//go:build linux && arm64

package platform

// OS and Arch identify the target platform at compile time.
const (
	OS   = "linux"
	Arch = "arm64"
)
