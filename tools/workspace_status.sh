#!/usr/bin/env bash
# workspace_status.sh — provides stamp variables for Bazel builds
set -euo pipefail

echo "STABLE_GIT_COMMIT $(git rev-parse HEAD 2>/dev/null || echo unknown)"
echo "STABLE_GIT_BRANCH $(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"
echo "STABLE_GIT_TAG $(git describe --tags --always --dirty 2>/dev/null || echo unknown)"
echo "BUILD_TIMESTAMP $(date -u +%Y-%m-%dT%H:%M:%SZ)"
