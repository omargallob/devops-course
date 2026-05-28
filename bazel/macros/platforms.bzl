"""Macro for declaring cross-compilation platforms.

Usage:
    load("//bazel/macros:platforms.bzl", "declare_platforms")

    declare_platforms(
        platforms = [
            ("linux_amd64", "linux", "x86_64"),
            ("linux_arm64", "linux", "aarch64"),
            ("darwin_amd64", "darwin", "x86_64"),
            ("darwin_arm64", "darwin", "aarch64"),
        ],
    )
"""

_OS_CONSTRAINT = {
    "linux": "@platforms//os:linux",
    "darwin": "@platforms//os:macos",
    "macos": "@platforms//os:macos",
}

_CPU_CONSTRAINT = {
    "x86_64": "@platforms//cpu:x86_64",
    "amd64": "@platforms//cpu:x86_64",
    "aarch64": "@platforms//cpu:aarch64",
    "arm64": "@platforms//cpu:aarch64",
}

def declare_platforms(platforms, visibility = None):
    """Declare platform targets from a list of tuples.

    Args:
        platforms: List of (name, os, cpu) tuples.
        visibility: Optional visibility for all generated platforms.
    """
    for (name, os, cpu) in platforms:
        os_constraint = _OS_CONSTRAINT.get(os)
        if not os_constraint:
            fail("Unknown OS '{}' for platform '{}'. Known: {}".format(
                os,
                name,
                ", ".join(_OS_CONSTRAINT.keys()),
            ))

        cpu_constraint = _CPU_CONSTRAINT.get(cpu)
        if not cpu_constraint:
            fail("Unknown CPU '{}' for platform '{}'. Known: {}".format(
                cpu,
                name,
                ", ".join(_CPU_CONSTRAINT.keys()),
            ))

        kwargs = {
            "name": name,
            "constraint_values": [os_constraint, cpu_constraint],
        }
        if visibility:
            kwargs["visibility"] = visibility

        native.platform(**kwargs)
