#!/usr/bin/env python3
"""Golden JSON Parity Comparator for LazyTrader vs Baseline (InvestOwl :8001 / Avenium :8091).

Uses Python stdlib only. Exits 0 on exact parity (within float tolerance), 1 on mismatch.
"""
from __future__ import annotations

import argparse
import json
import math
import sys
import urllib.error
import urllib.request
from typing import Any


def fetch_json(url: str, timeout: float = 10.0) -> tuple[int, Any]:
    req = urllib.request.Request(url, headers={"Accept": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            return resp.status, json.loads(resp.read().decode("utf-8"))
    except urllib.error.HTTPError as exc:
        body = exc.read().decode("utf-8", errors="replace")
        try:
            return exc.code, json.loads(body)
        except Exception:
            return exc.code, body


def compare_values(a: Any, b: Any, path: str, rtol: float, atol: float, ignore: set[str], diffs: list[str]) -> None:
    if type(a) is not type(b):
        if isinstance(a, (int, float)) and isinstance(b, (int, float)) and not isinstance(a, bool) and not isinstance(b, bool):
            pass
        else:
            diffs.append(f"{path}: type mismatch {type(a).__name__} ({a!r}) != {type(b).__name__} ({b!r})")
            return

    if isinstance(a, dict):
        keys_a = set(a.keys()) - ignore
        keys_b = set(b.keys()) - ignore
        for k in sorted(keys_a - keys_b):
            diffs.append(f"{path}.{k}: missing in candidate")
        for k in sorted(keys_b - keys_a):
            diffs.append(f"{path}.{k}: extra in candidate")
        for k in sorted(keys_a & keys_b):
            compare_values(a[k], b[k], f"{path}.{k}", rtol, atol, ignore, diffs)
    elif isinstance(a, list):
        if len(a) != len(b):
            diffs.append(f"{path}: list length {len(a)} != {len(b)}")
            return
        for idx, (va, vb) in enumerate(zip(a, b)):
            compare_values(va, vb, f"{path}[{idx}]", rtol, atol, ignore, diffs)
    elif isinstance(a, float) or isinstance(b, float):
        fa, fb = float(a), float(b)
        if not math.isclose(fa, fb, rel_tol=rtol, abs_tol=atol):
            diffs.append(f"{path}: float {fa} != {fb} (rtol={rtol}, atol={atol})")
    elif a != b:
        diffs.append(f"{path}: {a!r} != {b!r}")


def main() -> int:
    parser = argparse.ArgumentParser(description="Compare baseline vs candidate JSON API responses")
    parser.add_argument("--baseline", required=True, help="Baseline base URL, e.g. http://127.0.0.1:8001")
    parser.add_argument("--candidate", required=True, help="Candidate base URL, e.g. http://127.0.0.1:8099")
    parser.add_argument("--path", required=True, help="API path, e.g. /api/health")
    parser.add_argument("--rtol", type=float, default=1e-6)
    parser.add_argument("--atol", type=float, default=1e-6)
    parser.add_argument("--ignore-keys", default="", help="Comma-separated keys to ignore")
    args = parser.parse_args()

    ignore = {k.strip() for k in args.ignore_keys.split(",") if k.strip()}
    u_base = args.baseline.rstrip("/") + args.path
    u_cand = args.candidate.rstrip("/") + args.path

    s_base, j_base = fetch_json(u_base)
    s_cand, j_cand = fetch_json(u_cand)

    if s_base != s_cand:
        print(f"[FAIL] {args.path}: HTTP status {s_base} != {s_cand}")
        return 1

    diffs: list[str] = []
    compare_values(j_base, j_cand, "$", args.rtol, args.atol, ignore, diffs)
    if diffs:
        print(f"[FAIL] {args.path}: {len(diffs)} diff(s):")
        for d in diffs[:20]:
            print(f"  - {d}")
        return 1

    print(f"[PASS] {args.path}: 100% parity (HTTP {s_cand})")
    return 0


if __name__ == "__main__":
    sys.exit(main())
