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
if [ -n "${{BUILD_WORKSPACE_DIRECTORY:-}}" ]; then
    cd "$BUILD_WORKSPACE_DIRECTORY"
else
    WORKSPACE_ROOT=$(git rev-parse --show-toplevel 2>/dev/null || true)
    if [ -n "$WORKSPACE_ROOT" ] && [ -f "$WORKSPACE_ROOT/package.json" ]; then
        cd "$WORKSPACE_ROOT"
    else
        echo "ERROR: cannot determine workspace root" >&2
        exit 1
    fi
fi

if ! command -v npx >/dev/null 2>&1; then
    echo "WARNING: npx not found in PATH, skipping eslint" >&2
    exit 0
fi

echo "Running eslint..."
echo "  Working directory: $(pwd)"

GLOBS=({globs})

if npx eslint {config_arg} "${{GLOBS[@]}}"; then
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
