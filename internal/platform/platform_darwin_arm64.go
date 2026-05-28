//go:build darwin && arm64

package platform

// OS and Arch identify the target platform at compile time.
const (
	OS   = "darwin"
	Arch = "arm64"
)
