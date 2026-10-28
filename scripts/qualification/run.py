#!/usr/bin/env python3
"""Local-only, append-only performance/packaging evidence; never publishes.

Run from the repository root. Each output destination must be new and beneath
artifacts/completion-qualification/. Failures retain command logs and manifests.
Non-host binaries are inspected as data; only the current Mac executes tooling.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import platform
import subprocess
import sys
import time


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def source_manifest(root):
    paths = [root / "go.mod", root / "go.sum"]
    for directory in ("cmd", "internal", "scripts/qualification", "scripts/release"):
        paths.extend(p for p in (root / directory).rglob("*") if p.is_file() and "__pycache__" not in p.parts)
    return {str(p.relative_to(root)): digest(p) for p in sorted(paths) if p.is_file()}


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", required=True, help="new versioned attempt directory under artifacts/completion-qualification/")
    parser.add_argument("--skip-scale", action="store_true")
    parser.add_argument("--skip-builds", action="store_true")
    args = parser.parse_args()
    root = Path.cwd().resolve()
    output = Path(args.output).resolve()
    allowed = root / "artifacts/completion-qualification"
    if not output.is_relative_to(allowed) or output == allowed or output.exists():
        parser.error("output must be a NEW descendant of artifacts/completion-qualification/; existing attempts are preserved")
    if platform.system() != "Darwin":
        parser.error("only current Mac execution is qualified; this runner does not authorize foreign execution")
    output.mkdir(parents=True)
    env = dict(os.environ, CGO_ENABLED="0", GOCACHE="/private/tmp/factoryctl-go-cache", GOTMPDIR="/private/tmp")
    before = source_manifest(root)
    report = {"started_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()), "source_before": before, "source_before_sha256": hashlib.sha256(json.dumps(before, sort_keys=True).encode()).hexdigest(), "host": platform.platform(), "skip_scale": args.skip_scale, "skip_builds": args.skip_builds, "commands": [], "limits": ["Mac execution only", "foreign targets built and inspected, never executed", "no release/tag/publication", "Darwin and Windows retain platform system interfaces; not fully static"]}

    def write():
        (output / "report.json").write_text(json.dumps(report, indent=2, sort_keys=True) + "\n")

    def run(name, command, timeout=120, extra_env=None):
        started = time.monotonic()
        result = {"name": name, "argv": command, "timeout_seconds": timeout}
        report["commands"].append(result)
        write()
        try:
            proc = subprocess.run(command, env=dict(env, **(extra_env or {})), capture_output=True, timeout=timeout)
            (output / f"{name}.stdout").write_bytes(proc.stdout)
            (output / f"{name}.stderr").write_bytes(proc.stderr)
            result.update(exit_code=proc.returncode, elapsed_seconds=time.monotonic()-started)
            if proc.returncode:
                raise RuntimeError(f"{name} exited {proc.returncode}; preserved logs in {output}")
            return proc.stdout
        except subprocess.TimeoutExpired as error:
            (output / f"{name}.stdout").write_bytes(error.stdout or b"")
            (output / f"{name}.stderr").write_bytes(error.stderr or b"")
            result.update(timed_out=True, elapsed_seconds=time.monotonic()-started)
            raise
        finally:
            write()

    try:
        report["go_version"] = run("go-version", ["go", "version"]).decode().strip()
        report["hardware"] = run("hardware", ["sysctl", "-n", "machdep.cpu.brand_string", "hw.memsize", "hw.logicalcpu"]).decode().splitlines()
        if not args.skip_scale:
            runner = output / "scale-runner"
            run("build-scale", ["go", "build", "-o", str(runner), "./scripts/qualification/scale"])
            # Discovery correctly honors the owning repository's artifacts/
            # ignore rule. A sibling-repository fixture must therefore live
            # outside this checkout, not inside its ignored evidence tree.
            corpus_path = Path("/private/tmp") / ("factoryctl-qualification-" + hashlib.sha256(str(output).encode()).hexdigest()[:16])
            report["retained_corpus"] = str(corpus_path)
            data = run("scale", [str(runner), "--corpus", str(corpus_path)], timeout=600)
            report["scale"] = json.loads(data)
            report["scale_runner_sha256"] = digest(runner)
        if not args.skip_builds:
            binaries = []
            for target_os in ("darwin", "linux", "windows"):
                for arch in ("arm64", "amd64"):
                    name = f"factoryctl-{target_os}-{arch}" + (".exe" if target_os == "windows" else "")
                    binary = output / name
                    run(name, ["go", "build", "-trimpath", "-o", str(binary), "./cmd/factoryctl"], extra_env={"GOOS": target_os, "GOARCH": arch})
                    binaries.append(str(binary))
            inspector = output / "inspect-runner"
            run("build-inspector", ["go", "build", "-o", str(inspector), "./scripts/qualification/inspect"])
            report["packaging"] = json.loads(run("inspect", [str(inspector)] + binaries))
        report["source_after"] = source_manifest(root)
        report["source_stable_during_run"] = before == report["source_after"]
        report["completed"] = True
    except (RuntimeError, subprocess.TimeoutExpired, OSError, ValueError) as error:
        report["error"] = str(error)
        report["completed"] = False
        write()
        print(str(error), file=sys.stderr)
        return 1
    finally:
        write()
    print(json.dumps({"report": str(output / "report.json"), "completed": True, "source_stable_during_run": report["source_stable_during_run"]}))
    return 0


if __name__ == "__main__":
    sys.exit(main())
