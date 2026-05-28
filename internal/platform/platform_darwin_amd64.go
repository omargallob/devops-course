//go:build darwin && amd64

package platform

// OS and Arch identify the target platform at compile time.
const (
	OS   = "darwin"
	Arch = "amd64"
)
