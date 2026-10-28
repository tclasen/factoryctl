#!/usr/bin/env python3
"""Reject generated review artifacts and recordings in the index or commit history."""
import argparse
import subprocess
import sys

VIDEO_SUFFIXES = (b'.mp4', b'.webm', b'.mov', b'.mkv', b'.avi', b'.m4v')

def git(*args):
    return subprocess.check_output(['git', *args])

def forbidden(name):
    parts = name.lower().split(b'/')
    return b'artifacts' in parts or name.lower().endswith(VIDEO_SUFFIXES)

def check(history=None):
    if history:
        revisions = git('rev-list', history).decode().splitlines()
        inventories = ((rev, git('ls-tree', '-rz', '--name-only', rev)) for rev in revisions)
    else:
        inventories = [('index', git('ls-files', '-z'))]
    violations = set()
    for revision, data in inventories:
        for name in data.split(b'\0'):
            if name and forbidden(name):
                violations.add((revision, name.decode(errors='backslashreplace')))
    for revision, name in sorted(violations):
        print(f'Forbidden repository content: {revision}: {name}', file=sys.stderr)
    if violations:
        print('Keep artifacts and videos local; do not commit or push them.', file=sys.stderr)
    return bool(violations)

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    group = parser.add_mutually_exclusive_group()
    group.add_argument('--staged', action='store_true', help='Check the entire Git index')
    group.add_argument('--history', metavar='REV', help='Check every commit reachable from REV')
    args = parser.parse_args()
    sys.exit(check(args.history))
