#!/usr/bin/env python3
"""Validate only the event's newly introduced commit range and optional PR title."""
import argparse
import json
import re
import subprocess
import sys
from version import increment


def validate_revision(value):
    if not re.fullmatch(r"[0-9a-f]{40}", value):
        raise ValueError("revision must be a full lowercase Git SHA-1")
    return value


def check(base, head, title=None):
    validate_revision(head)
    if base is not None:
        validate_revision(base)
    # Git is used by development/release tooling only, never by factoryctl.
    revisions = [f"{base}..{head}"] if base and base != "0" * 40 else ["-1", head]
    output = subprocess.check_output(["git", "log", "--format=%B%x00", *revisions], text=True)
    messages = [message.strip() for message in output.split("\0") if message.strip()]
    for message in messages:
        increment(message)
    if title is not None:
        increment(title)
    return {"commits_checked": len(messages), "title_checked": title is not None, "base": base, "head": head}


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base")
    parser.add_argument("--head", required=True)
    parser.add_argument("--title")
    args = parser.parse_args()
    try:
        print(json.dumps(check(args.base, args.head, args.title), sort_keys=True))
    except (ValueError, subprocess.CalledProcessError) as error:
        parser.exit(1, f"conventional commit guard: {error}\n")
