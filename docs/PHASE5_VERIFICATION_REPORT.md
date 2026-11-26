# LazyTrader — Phase 5 Verification & Benchmark Report

> **Tanggal Verifikasi:** 2026-11-26  
> **Target Binary:** `/root/lazytrader/backend/bin/lazytrader` (`17M` stripped, `16,908,551` bytes) + `/root/lazytrader/frontend/dist` (`3.9M`, `3,890,784` bytes, `30` JS assets + CSS + `index.html`)  
> **Baseline Live:** InvestOwl (`127.0.0.1:8001` & `:3001`), Avenium Terminal (`127.0.0.1:8091` & `:3002`), SMtracker Canary Replay  
> **Status:** ✅ **PASS 100% (`142/142` API/SPA/Security, `19/19` Go Test Packages, Canary `45/45` Live API & `0/1032` Replay)**

---

## 1. Test Suite Go (`go test -count=1 ./...`) & Frontend Build Audit
- **Go Unit & Integration Suite (`go test -count=1 ./...`):** **19/19 paket uji LULUS** (0 gagal) melintasi `core/{auth,cache,flags}`, `investowl/{bandarmology,foreign_flow,forensic,nyangkut,ownership,shared,shareholders,signal,system}`, dan `terminal/{foreign_flow_pro,ipo_hunter,konglo,market_report,mmd,portal,shared}`.
- **TypeScript & Anti-Regresi Frontend:**
  - `npx tsc --noEmit` → exit `0` (`TSC_OK`).
  - Audit `grep` `bei1persen|next/link|next/navigation|next-auth|next/dynamic|"use client"|FeaturePlaceholder` di `frontend/src/` → `0` temuan (`GREP_CLEAN`).
  - Audit `grep` `bei1persen` di `backend/internal/` (non-test) → `0` route/handler (hanya komentar dokumentasi di `portal/handler.go:10`).

---

## 2. Verifikasi Canary Makro IHSG (`hybrid_p5_p15` & `absolute`)
- **Snapshot Replay Go (`TestCanarySnapshotsParity_BothModes`):**
  - Mode `absolute`: **`0/1032` mismatch** (100% parity)
  - Mode `hybrid_p5_p15`: **`0/1032` mismatch** (100% parity)
- **Snapshot Replay Python Reference (`/root/SMtracker/tools/canary/verify_backend_parity.py`):**
  - Mode `absolute`: `OK (1032 sesi)`
  - Mode `hybrid_p5_p15`: `OK (1032 sesi)` (`HASIL: LULUS`)
- **Canary Live API Parity vs SMtracker (`http://127.0.0.1:18099/api/forensic/foreign-exodus?days=45`):**
  - **`45/45` tanggal cocok (`100%`)** terhadap replay `canary_exodus_report.run_states(..., "hybrid_p5_p15")`.

---

## 3. Full Parity Suite (`tools/parity/run_phase5_verify.py`) — `142/142 PASSED`

| Battery | Cakupan | Hasil | Keterangan |
| :--- | :--- | :---: | :--- |
| **A — InvestOwl Golden JSON Parity (`:8001`)** | Seluruh `11` router InvestOwl + `health` (`61` endpoint incl. `404`) | **`61/61 PASS`** | Toleransi `rtol=1e-6, atol=1e-6` (kecuali 2 endpoint `causality` `2e-5`, lihat §4) |
| **B — Avenium Terminal Golden JSON Parity (`:8091`)** | Seluruh `5` domain Avenium (`mmd`, `ipo-hunter`, `foreign-flow`, `konglo`, `market-report`) (`20` endpoint GET & POST) | **`20/20 PASS`** | Toleransi `rtol=1e-6, atol=1e-6`; `recent_logs` & `current_time_wib` diabaikan pada `/status` |
| **C1 — Security, Exclusion & Compression** | `terminal/health`, `portal/status` zero `bei1persen`, `bei1persen` 404, 4× auth-guard `401`, security headers, gzip | **`10/10 PASS`** | `X-Content-Type-Options: nosniff`, `X-Frame-Options: SAMEORIGIN`, `gzip` aktif |
| **C2 — SPA Routes, Redirects & Assets** | `15` rute SPA langsung + `5` redirect `302` (incl. `/backtest`) + `31` aset `dist/assets/*` | **`51/51 PASS`** | Total `16/16` rute aktif, `no-cache` untuk HTML, `immutable` untuk `/assets/*` |
| **TOTAL** | — | **`142/142 PASS`** | Exit code `0` |

---

## 4. Temuan & Perbaikan Presisi Numerik Selama Fase 5
1. **CPython Half-to-Even Decimal Rounding (`internal/domains/investowl/shared/profiles.go`):** `shared.Round` sebelumnya memakai `math.RoundToEven(v*pow)/pow` yang rentan binary float drift pada titik tengah (contoh: `56.765` di `/api/foreign-flow/BBCA`). Diperbarui ke `strconv.ParseFloat(strconv.FormatFloat(v, 'f', decimals, 64), 64)` (identik dengan CPython `_Py_dg_dtoa`), diverifikasi oleh `TestRoundCPythonParity`.
2. **ARM64 Non-FMA Evaluation (`internal/domains/investowl/shared/conviction.go`):** Pada `INDF` (`25*0.30 + 50*0.30 + 35.75*0.20 + 80*0.20 = 45.65`), kontraksi `FMADD` ARM64 menghasilkan `45.6500000000000057 -> 45.7`, sementara evaluasi IEEE-754 Python menghasilkan `45.6499999999999986 -> 45.6`. Membungkus tiap produk dengan `float64(...)` menonaktifkan FMA contraction sehingga identik bit-for-bit dengan CPython (`TestConvictionScoreNoFMADrift` PASS).
3. **Granger Causality OLS Conditioning (`BBCA` `KZ` lag-5):** Terverifikasi langsung terhadap `statsmodels` Python: pada skala Rupiah mentah (`1e11`), SVD tak-terskala `statsmodels` menghasilkan `adj_p = 0.6110691918568347`; saat diskalakan `/1e11`, `statsmodels` menghasilkan `0.6110591359733628`, sedangkan Householder QR Go menghasilkan `0.6110591359733463` (selisih `1.65e-14` terhadap matematika presisi; `1.005e-5` terhadap SVD tak-terskala, dengan `lag`, `rank`, dan `is_significant` 100% identik).
4. **Isolasi Prefix Rate Limiter (`internal/core/middleware/ratelimit.go`):** Menambahkan `REDIS_KEY_PREFIX` pada key Redis rate limit agar proses uji paralel tidak bertabrakan dengan instance produksi di Redis DB 0 (di produksi `REDIS_KEY_PREFIX=""` tetap menghasilkan key identik `"rate_limit:%s:%s"`).

---

## 5. Benchmark Performa & Resource (Terukur Post-Load)

### A. Warm Latency (`3` request setelah warmup)
| Endpoint | LazyTrader (`:18099`) p50 / p95 | Baseline (`:8001` / `:8091`) p50 | Speedup |
| :--- | :---: | :---: | :---: |
| `GET /api/health` | `1.92 ms` / `2.09 ms` | `22.06 ms` (`:8001`) | **11.49x** |
| `GET /api/forensic/foreign-exodus?days=45` | `2.55 ms` / `2.98 ms` | `6.25 ms` (`:8001`) | **2.45x** |
| `GET /api/v1/ownership/market` | `2.08 ms` / `2.53 ms` | `10.96 ms` (`:8001`) | **5.27x** |
| `GET /api/v1/terminal/mmd/BBCA/analysis` | `2.38 ms` / `3.43 ms` | `9.25 ms` (`:8091`) | **3.89x** |
| `GET /api/v1/terminal/konglo/report` | `2.29 ms` / `2.36 ms` | `51.79 ms` (`:8091`) | **22.62x** |
| `POST /api/v1/terminal/ipo-hunter/screener` | `3.73 ms` / `3.77 ms` | `719.91 ms` (`:8091`) | **193.01x** |

### B. Memori RAM (`VmRSS` & `VmHWM` Post-Load)
| Stack | Komponen / Port | `VmRSS` | `VmHWM` |
| :--- | :--- | :---: | :---: |
| **Legacy 4-Service** | `uvicorn` (`:8001`, InvestOwl API) | `593,496 kB` (~`579.6 MiB`) | `642,200 kB` |
| **Legacy 4-Service** | `python` (`:8091`, Avenium API) | `364,180 kB` (~`355.6 MiB`) | `369,488 kB` |
| **Legacy 4-Service** | `next-server` (`:3001`, InvestOwl UI) | `40,144 kB` (~`39.2 MiB`) | `86,804 kB` |
| **Legacy 4-Service** | `next-server` (`:3002`, Avenium UI) | `45,592 kB` (~`44.5 MiB`) | `90,832 kB` |
| **Total Legacy (4 proses)** | `:8001 + :8091 + :3001 + :3002` | **`1,043,412 kB` (~`1,019.0 MiB`)** | — |
| **LazyTrader Fusion (1 binary)** | `lazytrader` (`:18099`, Full API + SPA) | **`33,636 kB` (~`32.8 MiB`)** | `134,100 kB` (~`130.9 MiB`) |
| **Penghematan RAM Post-Load** | — | **~96.8% (`-986.1 MiB`)** | — |

### C. Disk Footprint Runtime
- **Legacy Runtime Artifacts:**
  - `/root/investowl/backend/venv`: `644,786,697` B (`681M` apparent / `615 MiB`)
  - `/root/investowl/frontend/node_modules`: `635,257,050` B (`705M`)
  - `/root/investowl/frontend/.next`: `22,852,343` B (`24M`)
  - `/root/rev/frontend/.next`: `289,012,544` B (`277M`) + symlink `node_modules` (`37` B)
  - **Total Legacy:** **`1,591,908,671` bytes (`~1.48 GiB` / `1.59 GB`)**
- **LazyTrader Runtime Artifacts:**
  - `/root/lazytrader/backend/bin/lazytrader`: `16,908,551` B (`17M`)
  - `/root/lazytrader/frontend/dist`: `3,890,784` B (`3.9M`)
  - **Total LazyTrader:** **`20,799,335` bytes (`~19.8 MiB` / `20.8 MB`)**
- **Penghematan Disk Runtime:** **~98.7% (`-1.46 GiB`)**.

---

## 6. Kebersihan Lingkungan Setelah Verifikasi
- **Redis Cleanup:** Safe `SCAN + UNLINK` membersihkan `47` key `phase5*` di DB `0` dan `18` key `phase5*` di DB `1` (zero `FLUSHDB`).
- **Container (`mem-guard status && mem-guard clean`):** RAM tersedia `1655 MB`, swapfile bersih, `/tmp` tmpfs `3.9M`.
- **Legacy Repos:** `/root/investowl`, `/root/rev`, `/root/SMtracker` seluruhnya `dirty=0`.
- **Status Kesiapan:** **SIAP UNTUK FASE 6 (Deployment & Cutover)** *setelah persetujuan eksplisit owner.*
