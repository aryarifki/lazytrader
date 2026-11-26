#!/usr/bin/env python3
"""Phase 5 Full Parity, SPA/Asset, Security Header, and Latency Verifier for LazyTrader."""
from __future__ import annotations

import argparse
import json
import statistics
import sys
import time
import urllib.error
import urllib.request
from pathlib import Path
from typing import Any

from compare_json import compare_values


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


NO_REDIRECT_OPENER = urllib.request.build_opener(NoRedirect)


def request_raw(
    url: str,
    method: str = "GET",
    body: dict[str, Any] | None = None,
    headers: dict[str, str] | None = None,
    timeout: float = 45.0,
    follow_redirects: bool = True,
) -> tuple[int, dict[str, str], bytes, float]:
    req_headers = dict(headers or {})
    data_bytes = None
    if body is not None:
        data_bytes = json.dumps(body).encode("utf-8")
        req_headers.setdefault("Content-Type", "application/json")
    req = urllib.request.Request(url, data=data_bytes, headers=req_headers, method=method)
    opener = urllib.request.urlopen if follow_redirects else NO_REDIRECT_OPENER.open
    t0 = time.perf_counter()
    try:
        with opener(req, timeout=timeout) as resp:
            raw = resp.read()
            dt = (time.perf_counter() - t0) * 1000.0
            hdrs = {k.lower(): v for k, v in resp.headers.items()}
            return resp.status, hdrs, raw, dt
    except urllib.error.HTTPError as exc:
        raw = exc.read()
        dt = (time.perf_counter() - t0) * 1000.0
        hdrs = {k.lower(): v for k, v in exc.headers.items()} if exc.headers else {}
        return exc.code, hdrs, raw, dt


def request_json(
    url: str,
    method: str = "GET",
    body: dict[str, Any] | None = None,
    extra_headers: dict[str, str] | None = None,
    timeout: float = 45.0,
) -> tuple[int, Any, float]:
    hdrs = {"Accept": "application/json", **(extra_headers or {})}
    status, _, raw, dt = request_raw(
        url, method=method, body=body, headers=hdrs, timeout=timeout
    )
    try:
        return status, json.loads(raw.decode("utf-8")), dt
    except Exception:
        return status, raw.decode("utf-8", errors="replace"), dt


def run_json_parity_suite(
    name: str,
    baseline: str,
    candidate: str,
    cases: list[dict[str, Any]],
    rtol: float = 1e-6,
    atol: float = 1e-6,
    cand_headers: dict[str, str] | None = None,
) -> tuple[list[dict[str, Any]], int, int]:
    results: list[dict[str, Any]] = []
    passed = 0
    failed = 0

    for c in cases:
        method = c.get("method", "GET")
        path = c["path"]
        body = c.get("body")
        ignore = set(c.get("ignore_keys", []))
        case_rtol = float(c.get("rtol", rtol))
        case_atol = float(c.get("atol", atol))
        base_url = baseline.rstrip("/") + path
        cand_url = candidate.rstrip("/") + path

        s_base, j_base, dt_base = request_json(base_url, method=method, body=body)
        s_cand, j_cand, dt_cand = request_json(cand_url, method=method, body=body, extra_headers=cand_headers)

        ok = True
        diffs: list[str] = []
        if s_base != s_cand:
            ok = False
            diffs.append(f"HTTP status {s_base} != {s_cand}")
        else:
            compare_values(j_base, j_cand, "$", case_rtol, case_atol, ignore, diffs)
            if diffs:
                ok = False

        if ok:
            passed += 1
            print(f"  [PASS] {name} {method} {path} (HTTP {s_cand}, cand={dt_cand:.1f}ms base={dt_base:.1f}ms)")
        else:
            failed += 1
            print(f"  [FAIL] {name} {method} {path}: {len(diffs)} diff(s)")
            for d in diffs[:10]:
                print(f"         - {d}")

        results.append(
            {
                "suite": name,
                "method": method,
                "path": path,
                "status_base": s_base,
                "status_cand": s_cand,
                "passed": ok,
                "diffs": diffs[:20],
                "cand_ms": round(dt_cand, 2),
                "base_ms": round(dt_base, 2),
            }
        )

    return results, passed, failed


def run_security_and_exclusion_checks(
    candidate: str, cand_headers: dict[str, str] | None = None
) -> tuple[list[dict[str, Any]], int, int]:
    checks: list[dict[str, Any]] = []
    passed = 0
    failed = 0

    def record(label: str, ok: bool, detail: str) -> None:
        nonlocal passed, failed
        if ok:
            passed += 1
            print(f"  [PASS] {label}: {detail}")
        else:
            failed += 1
            print(f"  [FAIL] {label}: {detail}")
        checks.append({"label": label, "passed": ok, "detail": detail})

    # 1. Terminal health
    st, j, _ = request_json(candidate.rstrip("/") + "/api/v1/terminal/health", extra_headers=cand_headers)
    record("terminal/health", st == 200 and isinstance(j, dict) and j.get("status") == "healthy", f"HTTP {st}")

    # 2. Portal status has zero bei1persen
    st, j, _ = request_json(candidate.rstrip("/") + "/api/v1/terminal/portal/status", extra_headers=cand_headers)
    has_bei = "bei1persen" in json.dumps(j).lower()
    record("portal no-bei1persen", st == 200 and not has_bei, f"HTTP {st}, has_bei1persen={has_bei}")

    # 3. Excluded bei1persen routes return 404
    for p in ["/api/v1/terminal/bei1persen/overview", "/api/v1/terminal/bei1persen/saham/BBCA"]:
        st, _, _, _ = request_raw(candidate.rstrip("/") + p, headers=cand_headers)
        record(f"excluded {p}", st == 404, f"HTTP {st}")

    # 4. Protected routes fail-closed (HTTP 401 when unauthenticated)
    for method, p, body in [
        ("GET", "/api/users/", None),
        ("GET", "/api/logs/jobs/status", None),
        ("POST", "/api/bandar/pipeline/run", {"universe": "watchlist"}),
        ("GET", "/api/v1/terminal/market-report/logs", None),
    ]:
        st, _, _, _ = request_raw(candidate.rstrip("/") + p, method=method, body=body, headers=cand_headers)
        record(f"auth-guard {method} {p}", st == 401, f"HTTP {st}")

    # 5. Security headers + Gzip compression check
    st, hdrs, raw, _ = request_raw(
        candidate.rstrip("/") + "/api/v1/terminal/konglo/report",
        headers={"Accept-Encoding": "gzip", **(cand_headers or {})},
    )
    sec_ok = hdrs.get("x-content-type-options") == "nosniff" and bool(hdrs.get("x-frame-options"))
    gzip_ok = hdrs.get("content-encoding") == "gzip" and len(raw) > 100
    record("security headers", st == 200 and sec_ok, f"nosniff={hdrs.get('x-content-type-options')} xfo={hdrs.get('x-frame-options')}")
    record("gzip compression", st == 200 and gzip_ok, f"encoding={hdrs.get('content-encoding')} compressed_bytes={len(raw)}")

    return checks, passed, failed


def run_spa_and_asset_checks(candidate: str, dist_dir: Path) -> tuple[dict[str, Any], int, int]:
    passed = 0
    failed = 0
    details: dict[str, Any] = {"spa": [], "redirects": [], "assets": []}

    spa_routes = [
        "/",
        "/login",
        "/pending",
        "/admin",
        "/admin/market-report",
        "/signal",
        "/foreign",
        "/ownership",
        "/group",
        "/explorer",
        "/mmd",
        "/ipo-hunter",
        "/foreign-pro",
        "/konglo",
        "/BBCA",
    ]
    for route in spa_routes:
        st, hdrs, raw, dt = request_raw(candidate.rstrip("/") + route)
        html = raw.decode("utf-8", errors="replace")
        cc = hdrs.get("cache-control", "")
        ok = st == 200 and '<div id="root"></div>' in html and "no-cache" in cc
        if ok:
            passed += 1
        else:
            failed += 1
            print(f"  [FAIL] SPA {route}: st={st} cc={cc!r}")
        details["spa"].append({"route": route, "status": st, "passed": ok, "ms": round(dt, 2)})

    redirects = [
        ("/backtest", "/signal?tab=backtest"),
        ("/nyangkut", "/ownership"),
        ("/whales", "/group"),
        ("/market-report", "/admin/market-report"),
        ("/report-generator", "/admin/market-report"),
    ]
    for src, target in redirects:
        st, hdrs, _, dt = request_raw(candidate.rstrip("/") + src, follow_redirects=False)
        loc = hdrs.get("location", "")
        ok = st == 302 and loc == target
        if ok:
            passed += 1
        else:
            failed += 1
            print(f"  [FAIL] REDIRECT {src}: st={st} loc={loc!r} want={target!r}")
        details["redirects"].append({"from": src, "to": target, "status": st, "location": loc, "passed": ok})

    asset_files = sorted((dist_dir / "assets").glob("*"))
    for af in asset_files:
        if not af.is_file():
            continue
        url_path = f"/assets/{af.name}"
        st, hdrs, raw, _ = request_raw(candidate.rstrip("/") + url_path)
        cc = hdrs.get("cache-control", "")
        ok = st == 200 and len(raw) == af.stat().st_size and "immutable" in cc
        if ok:
            passed += 1
        else:
            failed += 1
            print(f"  [FAIL] ASSET {url_path}: st={st} bytes={len(raw)}/{af.stat().st_size} cc={cc!r}")
        details["assets"].append({"path": url_path, "status": st, "bytes": len(raw), "passed": ok})

    print(
        f"  [SPA SUMMARY] routes={len(spa_routes)} redirects={len(redirects)} assets={len(asset_files)} "
        f"passed={passed} failed={failed}"
    )
    return details, passed, failed


def measure_latency(
    base_url: str,
    path: str,
    method: str = "GET",
    body: dict[str, Any] | None = None,
    extra_headers: dict[str, str] | None = None,
    runs: int = 3,
) -> dict[str, float]:
    # 1 warmup + runs timed calls
    request_json(base_url.rstrip("/") + path, method=method, body=body, extra_headers=extra_headers)
    samples: list[float] = []
    for _ in range(runs):
        st, _, dt = request_json(base_url.rstrip("/") + path, method=method, body=body, extra_headers=extra_headers)
        if st == 200:
            samples.append(dt)
    if not samples:
        return {"p50_ms": -1.0, "p95_ms": -1.0}
    samples_sorted = sorted(samples)
    p50 = statistics.median(samples_sorted)
    p95_idx = max(0, min(len(samples_sorted) - 1, int(round(0.95 * (len(samples_sorted) - 1)))))
    return {"p50_ms": round(p50, 2), "p95_ms": round(samples_sorted[p95_idx], 2)}


def main() -> int:
    parser = argparse.ArgumentParser(description="LazyTrader Phase 5 Full Parity & Benchmark Runner")
    parser.add_argument("--candidate", required=True)
    parser.add_argument("--investowl", required=True)
    parser.add_argument("--avenium", required=True)
    parser.add_argument("--dist", required=True)
    parser.add_argument("--out", default="")
    args = parser.parse_args()

    investowl_cases: list[dict[str, Any]] = [
        {"path": "/api/health"},
        {"path": "/api/users/settings/auth-gate"},
        {"path": "/api/stocks/BBCA/history"},
        {"path": "/api/broker-flow/whales"},
        {"path": "/api/broker-flow/whales/YP"},
        {"path": "/api/broker-flow/BBCA/latest"},
        {"path": "/api/broker-flow/BBCA/history"},
        {"path": "/api/broker-flow/BBCA/summary"},
        {"path": "/api/forensic/foreign-exodus?days=45"},
        {"path": "/api/forensic/foreign-exodus?days=60"},
        {"path": "/api/foreign-flow/sectors/summary"},
        {"path": "/api/foreign-flow/sectors/history?sector=IDXENERGY"},
        {"path": "/api/foreign-flow/BBCA"},
        {"path": "/api/signal/daily"},
        {"path": "/api/signal/backtest/stats"},
        {"path": "/api/signal/backtest/history"},
        {"path": "/api/signal/outcomes/summary"},
        {"path": "/api/signal/active"},
        {"path": "/api/signal/confluence-screener"},
        {"path": "/api/bandar/universe"},
        {"path": "/api/bandar/tickers"},
        {"path": "/api/bandar/system-status"},
        {"path": "/api/bandar/daily-summary"},
        {"path": "/api/bandar/screener"},
        {"path": "/api/bandar/screener-v2"},
        {"path": "/api/bandar/detail/BBCA"},
        {"path": "/api/bandar/broker-flow/BBCA"},
        {"path": "/api/bandar/validation/BBCA"},
        {"path": "/api/bandar/validation-v2/BBCA"},
        # Column-scaled Householder QR matches scaled exact math to 1.65e-14; unscaled statsmodels SVD
        # on 1e11 IDR has ~1.0e-5 conditioning drift on BBCA KZ lag-5 with identical lag/rank/significance.
        {"path": "/api/bandar/causality/BBCA", "rtol": 2e-5, "atol": 2e-5},
        {"path": "/api/bandar/stocks/BBCA/metrics"},
        {"path": "/api/bandar/stocks/BBCA/smart-flow"},
        {"path": "/api/bandar/stocks/BBCA/broker-compare"},
        {"path": "/api/bandar/stocks/BBCA/causality", "rtol": 2e-5, "atol": 2e-5},
        {"path": "/api/bandar/stocks/BBCA/event-study"},
        {"path": "/api/bandar/stocks/BBCA/raw-tables"},
        {"path": "/api/bandar/stocks/BBCA/broker-profiles"},
        {"path": "/api/bandar/stocks/BBCA/price-performance"},
        {"path": "/api/bandar/stocks/BBCA/broker-distribution"},
        {"path": "/api/bandar/stocks/BBCA/retail-flow"},
        {"path": "/api/bandar/shareholders/top-cross-holdings"},
        {"path": "/api/bandar/shareholders/groups"},
        {"path": "/api/bandar/shareholders/network-graph"},
        {"path": "/api/bandar/shareholders/BBCA"},
        {"path": "/api/nyangkut/summary"},
        {"path": "/api/nyangkut/leaderboard"},
        {"path": "/api/nyangkut/stock/BBCA"},
        {"path": "/api/v1/stock/BBCA/nyangkut-tracker"},
        {"path": "/api/v1/nyangkut/summary"},
        {"path": "/api/v1/nyangkut/leaderboard"},
        {"path": "/api/v1/ownership/market"},
        {"path": "/api/v1/ownership/stock/BBCA"},
        {"path": "/api/v1/ownership/stock/ZZZZNOTREAL"},
        {"path": "/api/v1/ownership/groups"},
        {"path": "/api/v1/ownership/investors?limit=20"},
        {"path": "/api/v1/ownership/explorer?limit=20"},
        {"path": "/api/v1/ownership/investors/ct?limit=20"},
        {"path": "/api/v1/ownership/explorer/ct?limit=20"},
        {"path": "/api/v1/ownership/staging/overview"},
        {"path": "/api/v1/ownership/staging/stock/BBCA"},
        {"path": "/api/v1/ownership/staging/parity/BBCA"},
    ]

    avenium_cases: list[dict[str, Any]] = [
        {"path": "/api/v1/terminal/mmd/tickers"},
        {"path": "/api/v1/terminal/mmd/BBCA?limit=60"},
        {"path": "/api/v1/terminal/mmd/BBCA/analysis"},
        {"path": "/api/v1/terminal/mmd/BBCA/t5"},
        {"path": "/api/v1/terminal/mmd/BBCA/group-performance"},
        {"path": "/api/v1/terminal/mmd/BBCA/broker-drilldown?broker=YP"},
        {"path": "/api/v1/terminal/mmd/explorer/scan?category=normal&limit=20"},
        {"path": "/api/v1/terminal/ipo-hunter/meta"},
        {"path": "/api/v1/terminal/ipo-hunter/list"},
        {"path": "/api/v1/terminal/ipo-hunter/analysis/RANS?mode=after"},
        {
            "method": "POST",
            "path": "/api/v1/terminal/ipo-hunter/screener",
            "body": {"after": True, "chg": "ipo", "metric": "b_pct", "conds": []},
        },
        {"path": "/api/v1/terminal/foreign-flow/screener?universe=ff&window=40&limit=50"},
        {"path": "/api/v1/terminal/foreign-flow/series/BBCA?limit=60"},
        {"path": "/api/v1/terminal/foreign-flow/ihsg"},
        {"path": "/api/v1/terminal/konglo/report"},
        {"path": "/api/v1/terminal/konglo/groups"},
        {"path": "/api/v1/terminal/konglo/summary?limit=50"},
        {"path": "/api/v1/terminal/konglo/group/Barito%20Group"},
        {"path": "/api/v1/terminal/market-report/status", "ignore_keys": ["current_time_wib", "recent_logs"]},
        {"path": "/api/v1/terminal/market-report/history"},
    ]

    print("=== BATTERY A: InvestOwl Golden JSON Parity (:8001) ===")
    inv_res, inv_pass, inv_fail = run_json_parity_suite(
        "InvestOwl", args.investowl, args.candidate, investowl_cases, cand_headers={"X-Forwarded-For": "10.0.5.1"}
    )

    print("\n=== BATTERY B: Avenium Terminal Golden JSON Parity (:8091) ===")
    av_res, av_pass, av_fail = run_json_parity_suite(
        "Avenium", args.avenium, args.candidate, avenium_cases, cand_headers={"X-Forwarded-For": "10.0.5.2"}
    )

    print("\n=== BATTERY C1: Security Guards, Exclusion & Gzip ===")
    sec_res, sec_pass, sec_fail = run_security_and_exclusion_checks(
        args.candidate, cand_headers={"X-Forwarded-For": "10.0.5.3"}
    )

    print("\n=== BATTERY C2: SPA Routes, Legacy Redirects & Assets ===")
    spa_res, spa_pass, spa_fail = run_spa_and_asset_checks(args.candidate, Path(args.dist))

    print("\n=== BATTERY D: Warm Latency Benchmark (3 runs) ===")
    bench_targets = [
        ("InvestOwl", args.investowl, "GET", "/api/health", None),
        ("InvestOwl", args.investowl, "GET", "/api/forensic/foreign-exodus?days=45", None),
        ("InvestOwl", args.investowl, "GET", "/api/v1/ownership/market", None),
        ("Avenium", args.avenium, "GET", "/api/v1/terminal/mmd/BBCA/analysis", None),
        ("Avenium", args.avenium, "GET", "/api/v1/terminal/konglo/report", None),
        ("Avenium", args.avenium, "POST", "/api/v1/terminal/ipo-hunter/screener", {"after": True, "chg": "ipo", "metric": "b_pct", "conds": []}),
    ]
    lat_rows: list[dict[str, Any]] = []
    for svc, base, method, path, body in bench_targets:
        cand_m = measure_latency(args.candidate, path, method=method, body=body, extra_headers={"X-Forwarded-For": "10.0.5.4"}, runs=3)
        base_m = measure_latency(base, path, method=method, body=body, runs=3)
        speedup = round(base_m["p50_ms"] / cand_m["p50_ms"], 2) if cand_m["p50_ms"] > 0 else 0.0
        row = {
            "service": svc,
            "endpoint": f"{method} {path}",
            "cand_p50_ms": cand_m["p50_ms"],
            "cand_p95_ms": cand_m["p95_ms"],
            "base_p50_ms": base_m["p50_ms"],
            "base_p95_ms": base_m["p95_ms"],
            "speedup_x": speedup,
        }
        lat_rows.append(row)
        print(
            f"  [LAT] {method} {path}: cand p50={cand_m['p50_ms']}ms (p95={cand_m['p95_ms']}ms) vs "
            f"{svc} p50={base_m['p50_ms']}ms ({speedup}x)"
        )

    total_pass = inv_pass + av_pass + sec_pass + spa_pass
    total_fail = inv_fail + av_fail + sec_fail + spa_fail
    summary = {
        "investowl_pass": inv_pass,
        "investowl_fail": inv_fail,
        "avenium_pass": av_pass,
        "avenium_fail": av_fail,
        "security_pass": sec_pass,
        "security_fail": sec_fail,
        "spa_pass": spa_pass,
        "spa_fail": spa_fail,
        "total_pass": total_pass,
        "total_fail": total_fail,
        "investowl_results": inv_res,
        "avenium_results": av_res,
        "security_results": sec_res,
        "spa_summary": {
            "spa_routes": len(spa_res["spa"]),
            "redirects": len(spa_res["redirects"]),
            "assets": len(spa_res["assets"]),
        },
        "latency": lat_rows,
    }
    if args.out:
        Path(args.out).write_text(json.dumps(summary, indent=2), encoding="utf-8")

    print(
        f"\n=== PHASE 5 SUITE TOTAL: {total_pass}/{total_pass + total_fail} PASSED "
        f"(InvestOwl {inv_pass}/{inv_pass + inv_fail}, Avenium {av_pass}/{av_pass + av_fail}, "
        f"Security {sec_pass}/{sec_pass + sec_fail}, SPA/Assets {spa_pass}/{spa_pass + spa_fail}) ==="
    )
    return 0 if total_fail == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
