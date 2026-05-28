// Package platform provides a compile-time check that ensures the server
// binary is only built on supported platforms: Linux (amd64/arm64) and
// macOS (amd64/arm64).
//
// Supported targets:
//   - linux/amd64  (Ubuntu, Arch, etc.)
//   - linux/arm64  (Ubuntu, Arch on ARM)
//   - darwin/amd64 (macOS Intel)
//   - darwin/arm64 (macOS Apple Silicon)
package platform
