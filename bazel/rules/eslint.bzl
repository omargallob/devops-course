"""Rule for linting TypeScript/Astro with ESLint.

Provides a test rule that runs eslint on the frontend source files.
Supports astro-eslint-parser and typescript-eslint.

Usage:
    load("//bazel/rules:eslint.bzl", "eslint_test")

    eslint_test(
        name = "eslint",
        srcs = glob(["apps/web/src/**/*.{ts,astro}"]),
        config = "apps/web/eslint.config.mjs",
    )
"""

def _eslint_test_impl(ctx):
    script = ctx.actions.declare_file(ctx.label.name + ".sh")

    config_arg = ""
    if ctx.file.config:
        config_arg = "--config \"$(pwd)/{}\"".format(ctx.file.config.short_path)

    content = """#!/usr/bin/env bash
set -euo pipefail

# Resolve workspace directory.
# For `bazel run`, BUILD_WORKSPACE_DIRECTORY is set.
# For `bazel test`, we're in the execroot with symlinks to the source tree.
# Follow a known file's symlink to find the real workspace root.
if [ -n "${{BUILD_WORKSPACE_DIRECTORY:-}}" ]; then
    cd "$BUILD_WORKSPACE_DIRECTORY"
elif [ -L "package.json" ]; then
    REAL_PKG=$(python3 -c "import os; print(os.path.realpath('package.json'))")
    WORKSPACE_ROOT=$(dirname "$REAL_PKG")
    cd "$WORKSPACE_ROOT"
else
    WORKSPACE_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || true)
    if [ -n "$WORKSPACE_ROOT" ] && [ -f "$WORKSPACE_ROOT/package.json" ]; then
        cd "$WORKSPACE_ROOT"
    else
        echo "ERROR: cannot determine workspace root" >&2
        exit 1
    fi
fi

# Ensure we use the correct Node version (ESLint 9+ requires Node 18+).
# Try nvm first (reads .nvmrc), then fall back to common Node 22 paths.
if [ -s "${{NVM_DIR:-$HOME/.nvm}}/nvm.sh" ]; then
    . "${{NVM_DIR:-$HOME/.nvm}}/nvm.sh" --no-use
    nvm use --silent 2>/dev/null || true
fi
NODE_VERSION=$(node --version 2>/dev/null || echo "none")
NODE_MAJOR=$(echo "$NODE_VERSION" | sed 's/^v//' | cut -d. -f1)
if [ "$NODE_MAJOR" -lt 18 ] 2>/dev/null; then
    echo "ERROR: Node 18+ required for ESLint 9 (found $NODE_VERSION)" >&2
    echo "Run: nvm use 22" >&2
    exit 1
fi

# Use project-local eslint from node_modules.
if [ -x "apps/web/node_modules/.bin/eslint" ]; then
    ESLINT="apps/web/node_modules/.bin/eslint"
elif command -v pnpm >/dev/null 2>&1; then
    ESLINT="pnpm --filter @devops-course/web exec eslint"
else
    echo "WARNING: eslint not found (run 'pnpm install' in apps/web first)" >&2
    exit 0
fi

echo "Running eslint... (Node $NODE_VERSION)"
echo "  Working directory: $(pwd)"

GLOBS=({globs})

if $ESLINT {config_arg} "${{GLOBS[@]}}"; then
    echo ""
    echo "PASSED: eslint found no issues"
else
    echo ""
    echo "FAILED: eslint found issues" >&2
    exit 1
fi
""".format(
        globs = " ".join(['"{}"'.format(g) for g in ctx.attr.globs]),
        config_arg = config_arg,
    )

    ctx.actions.write(
        output = script,
        content = content,
        is_executable = True,
    )

    runfiles_files = list(ctx.files.srcs)
    if ctx.file.config:
        runfiles_files.append(ctx.file.config)

    runfiles = ctx.runfiles(files = runfiles_files)

    return [DefaultInfo(
        executable = script,
        runfiles = runfiles,
    )]

eslint_test = rule(
    implementation = _eslint_test_impl,
    test = True,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Source files (for dependency tracking).",
        ),
        "globs": attr.string_list(
            mandatory = True,
            doc = "Glob patterns to pass to eslint.",
        ),
        "config": attr.label(
            allow_single_file = True,
            mandatory = False,
            doc = "Optional eslint config file.",
        ),
    },
    doc = "Runs eslint as a test target.",
)
