"""Rule for formatting/linting with Prettier.

Provides a test rule that runs prettier --check on source files.
Fails the test if any files are not formatted.

Usage:
    load("//bazel/rules:prettier.bzl", "prettier_test")

    prettier_test(
        name = "prettier",
        srcs = glob(["**/*.ts", "**/*.astro", "**/*.css", "**/*.json"]),
        config = ".prettierrc.yaml",
    )
"""

def _prettier_test_impl(ctx):
    script = ctx.actions.declare_file(ctx.label.name + ".sh")

    config_arg = ""
    if ctx.file.config:
        config_arg = "--config \"$(pwd)/{}\"".format(ctx.file.config.short_path)

    ignore_arg = ""
    if ctx.file.ignore:
        ignore_arg = "--ignore-path \"$(pwd)/{}\"".format(ctx.file.ignore.short_path)

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
    echo "WARNING: npx not found in PATH, skipping prettier" >&2
    exit 0
fi

echo "Running prettier --check..."
echo "  Working directory: $(pwd)"

GLOBS=({globs})

if npx prettier --check {config_arg} {ignore_arg} "${{GLOBS[@]}}"; then
    echo ""
    echo "PASSED: all files formatted correctly"
else
    echo ""
    echo "FAILED: some files need formatting" >&2
    echo "Run: npx prettier --write ${{GLOBS[*]}}" >&2
    exit 1
fi
""".format(
        globs = " ".join(['"{}"'.format(g) for g in ctx.attr.globs]),
        config_arg = config_arg,
        ignore_arg = ignore_arg,
    )

    ctx.actions.write(
        output = script,
        content = content,
        is_executable = True,
    )

    runfiles_files = list(ctx.files.srcs)
    if ctx.file.config:
        runfiles_files.append(ctx.file.config)
    if ctx.file.ignore:
        runfiles_files.append(ctx.file.ignore)

    runfiles = ctx.runfiles(files = runfiles_files)

    return [DefaultInfo(
        executable = script,
        runfiles = runfiles,
    )]

prettier_test = rule(
    implementation = _prettier_test_impl,
    test = True,
    attrs = {
        "srcs": attr.label_list(
            allow_files = True,
            doc = "Source files (for dependency tracking in runfiles).",
        ),
        "globs": attr.string_list(
            mandatory = True,
            doc = "Glob patterns to pass to prettier (e.g. 'apps/web/src/**/*.{ts,astro,css}').",
        ),
        "config": attr.label(
            allow_single_file = True,
            mandatory = False,
            doc = "Optional prettier config file.",
        ),
        "ignore": attr.label(
            allow_single_file = True,
            mandatory = False,
            doc = "Optional .prettierignore file.",
        ),
    },
    doc = "Runs prettier --check as a test target.",
)
