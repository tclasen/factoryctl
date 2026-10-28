#!/usr/bin/env python3
"""Check conventional messages and calculate the bounded pre-1.0 release candidate.

Input is a JSON array of complete commit messages, newest or oldest first.
This program neither creates tags nor publishes releases. Git collection belongs
to the caller so the exact candidate range can be retained with the evidence.
"""
import argparse
import json
import re
import sys

HEADER = re.compile(r"^(feat|fix|docs|refactor|perf|test|build|ci|chore|style|revert)(?:\([^()\r\n]+\))?(!)?: (\S.*)$")
VERSION = re.compile(r"^0\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$")


def increment(message):
    lines = message.splitlines()
    match = HEADER.fullmatch(lines[0] if lines else "")
    if not match:
        raise ValueError("commit title must use a supported Conventional Commit type")
    breaking = bool(match[2]) or any(re.match(r"^BREAKING(?: CHANGE|-CHANGE): \S", line) for line in lines[1:])
    if breaking or match[1] == "feat":
        return 2
    return 1 if match[1] in ("fix", "perf") else 0


def calculate(messages, previous=None, proposed=None):
    if not isinstance(messages, list) or any(not isinstance(m, str) for m in messages):
        raise ValueError("messages must be a JSON array of complete commit messages")
    rank = max((increment(m) for m in messages), default=0)
    if previous is None:
        result = "0.1.0"
    else:
        match = VERSION.fullmatch(previous)
        if not match:
            raise ValueError("previous release must be a canonical 0.x.y version; major 1+ is unauthorized")
        minor, patch = map(int, match.groups())
        result = f"0.{minor + 1}.0" if rank == 2 else f"0.{minor}.{patch + 1}" if rank == 1 else None
    if proposed is not None:
        if not VERSION.fullmatch(proposed):
            raise ValueError("major 1+ and noncanonical versions are unauthorized")
        if proposed != result:
            raise ValueError("proposed version does not equal calculated version (or no release is due)")
    return {"version": result, "release_due": result is not None, "initial": previous is None, "increment": ["none", "patch", "minor"][rank]}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("messages", help="JSON array file, or - for stdin")
    parser.add_argument("--previous")
    parser.add_argument("--proposed")
    args = parser.parse_args()
    try:
        with (sys.stdin if args.messages == "-" else open(args.messages)) as source:
            result = calculate(json.load(source), args.previous, args.proposed)
        print(json.dumps(result, sort_keys=True))
    except (ValueError, OSError) as error:
        parser.exit(1, f"version guard: {error}\n")


if __name__ == "__main__":
    main()
