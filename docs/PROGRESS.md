# LazyTrader — Progress Tracker

> **Terakhir Diperbarui:** 2026-11-26  
> **Repo:** `git@github.com:aryarifki/lazytrader.git`

---

## Ringkasan Fase

| Fase | Deskripsi | Status | Catatan / Bukti |
| :---: | :--- | :---: | :--- |
| **0** | Pre-flight Commit & Push (`/root/rev` `1e7cf57`, `/root/investowl` `f1ed6ad`, `/root/SMtracker` `b2744d2`) | ✅ **DONE** | Semua repo clean (`0/0` vs `origin/main`) |
| **1** | Inisialisasi Repo `/root/lazytrader`, struktur modular, & dokumentasi kanonial `docs/` | ✅ **DONE** | Commit `f1b0cfa` |
| **2** | Core Platform Go (`config`, `database`, `cache`, `flags`, `auth`, `middleware`, `spa`) + `health`/`users`/`portal` + Vite React 19 Shell + Parity Harness | ✅ **DONE** | 7/7 Go unit tests PASS; `/api/health` & `/api/users/settings/auth-gate` 100% parity vs `:8001`; `VmRSS` terukur **16,5 MB** |
| **3** | Migrasi Backend Domain (11 domain InvestOwl + 5 domain Avenium) + Golden JSON Parity | ✅ **DONE** | **`11/11` router InvestOwl + `5/5` domain Avenium (`26` endpoint) + `health` + `portal` 100% DONE!** (`915a292`..`3a3d7d9` + Wave 3E-3); **108/108** InvestOwl + **85/85** Terminal parity PASS (`:8001` / `:8091`, zero fails); Canary `0/1032` mismatch; `VmRSS` Wave 3E-3 **31.540 kB** (`~30,8 MB`, `VmHWM` 33.200 kB) vs legacy `avenium-backend` ~357 MB |
| **4** | Migrasi Frontend Features (seluruh halaman InvestOwl + 5 modul Avenium dipecah modular & MD3 restyle) | ✅ **DONE (Wave 4A–4D Complete)** | **100% DONE (16/16 rute aktif, zero `FeaturePlaceholder`)!** Wave 4A (`shared`, `home`, `ticker`, `signal`, `App.tsx`): 31 file, `vite build` 10.30s, `VmRSS` 48.944 kB; Wave 4B (`research`, `foreign`, `ownership`, `group`, `explorer`): 110 file (12.438 LOC), `vite build` 26.11s, SPA 4/4 + 13/13 API PASS, `VmRSS` 34.628 kB; Wave 4C (`terminal.ts`, `mmd`, `ipo-hunter`): 33 file (6.568 LOC), `vite build` 21.38s, SPA 2/2 + JS 3/3 + API 11/11 PASS, `VmRSS` 23.756 kB; **Wave 4D** (`foreign-pro`, `konglo`, `market-report` + perluas `terminal.ts` + hapus `FeaturePlaceholder`): 42 file fitur baru (8.416 LOC, maks 335 LOC/file), `tsc --noEmit` + `GREP_CLEAN`, `vite build` 22.18s (`mem-guard`), SPA 3/3 + JS 4/4 + API 9/9 HTTP 200 + auth guard 401 (`VmRSS` 33.020 kB / ~32,2 MB, `VmHWM` 35.236 kB) |
| **5** | Verifikasi Paritas Total (`45/45` Canary, unit/e2e test, benchmark RAM & disk) | ✅ **DONE** | **142/142 (`100%`) Parity Suite PASS** (`61/61` InvestOwl `:8001`, `20/20` Avenium `:8091`, `10/10` Security/No-BEI1%/Gzip, `51/51` SPA/Redirects/Assets); **19/19 paket Go test PASS**; Canary `0/1032` both-mode replay + **`45/45` Live API**; `VmRSS` post-load **`32,8 MiB`** (`33.636 kB`) vs 4-service legacy **`~1.019 MiB`** (**-96,8% RAM**, **-98,7% disk** `19,8 MiB` vs `1,48 GiB`), latency **2,45x–193,01x** lebih cepat → [PHASE5_VERIFICATION_REPORT.md](file:///root/lazytrader/docs/PHASE5_VERIFICATION_REPORT.md) |
| **6** | Deployment & Cutover Cloudflared *(menunggu persetujuan owner)* | ✅ **DONE** | `lazytrader.service` (`127.0.0.1:8099`) `active`, idle `VmRSS` **`13.944 kB`** (`~13,6 MiB`, `VmHWM` `21.452 kB`); `https://terminal.avenia.my.id` live (`7/7` SPA, `/api/health`, `portal` no-`bei1persen`, `bei1persen` `404`, `exodus?days=45`, `konglo`, `market-report/logs` `401`, Google OAuth `302`); `ct.avenia.my.id` dikonfigurasi ke `:8091`/`:3002` + CORS `https://ct.avenia.my.id` aktif (tunggu CNAME `ct` di Cloudflare Dashboard) → [DEPLOYMENT_RUNBOOK.md](file:///root/lazytrader/docs/DEPLOYMENT_RUNBOOK.md) |

---

## Tracker Domain Backend (Go)

| Domain | Sumber | Status | Parity Check |
| :--- | :--- | :---: | :---: |
| `core/*` (`config`, `database`, `cache`, `flags`, `auth`, `middleware`, `spa`) | InvestOwl + Avenium | ✅ DONE | 5/5 unit tests PASS |
| `investowl/health` | InvestOwl + Avenium | ✅ DONE | `compare_json.py` PASS |
| `investowl/users` (9 endpoints + auth-gate) | InvestOwl `routers/auth.py` | ✅ DONE | `auth-gate` parity PASS |
| `terminal/portal` (tanpa `bei1persen`) | Avenium `routers/portal.py` | ✅ DONE | Unit test + curl no-bei1persen PASS |
| `investowl/stocks`, `broker`, `system` (16 endpoints + `shared/*`) | InvestOwl | ✅ DONE | Commit `915a292` — 100% JSON parity (6/6 endpoints) + `/api/logs/*` auth & dry-run smoke PASS; `system_test.go` PASS |
| `investowl/forensic`, `foreign_flow`, `signal` (13 endpoints + `shared/{stats,causality,conviction}`) | InvestOwl | ✅ DONE | Commit `2e9e471` — 100% JSON parity (18/18 endpoints vs `:8001` incl. `confluence-screener`, HTTP 404); Canary `0/1032` mismatch (`absolute` & `hybrid_p5_p15`); `VmRSS` 28.032 kB |
| `investowl/bandarmology` (30 endpoints + `affiliation` + `causality_breakdown`) | InvestOwl | ✅ DONE | Commit `b4124d1` — 100% JSON parity (26/26 read endpoints vs `:8001`, HTTP 401 auth guard on `pipeline/run`); `bandarmology_test.go` + `stats_test.go` PASS; `VmRSS` 30.856 kB |
| `investowl/nyangkut` & `shareholders` (11 handlers / 15 URL + `shared/ct_groups.go`) | InvestOwl | ✅ DONE | Commit `7aa6f5b` — 100% JSON parity (17/17 paths vs `:8001` incl. `JELI`, aliases & 404); `/refresh` SAFE SCAN+UNLINK & `/api/bandar/system-status` 200; `VmRSS` 18.948 kB |
| `investowl/ownership` (14 endpoints + `ct_baseline_repo`, `ct_parity`, `ct_explorer`) | InvestOwl | ✅ DONE | Commit `b4e711d` — 100% JSON parity (41/41 cases vs `:8001` incl. `stock`, `investors`, `market`, `groups`, `explorer`, `investors/ct`, `explorer/ct` 6 varian, `staging/*`, & 404); `ownership_test.go` + mock `sync-catalog` PASS; `VmRSS` 38.496 kB (`~37,6 MB`) |
| `terminal/shared`, `foreign_flow_pro`, `konglo`, `market_report` (15 endpoints) | Avenium | ✅ DONE | Commit `e51e03c` — 100% JSON parity (26/26 cases vs `:8091`: 12 `foreign-flow`, 8 `konglo`, 6 `market-report` + 401/admin unmasked parity); dual InvestOwl Admin Session + `X-Admin-Key`; zero mutation `/root/rev/reports/execution.log`; `VmRSS` 28.436 kB (`~27,8 MB`) |
| `terminal/mmd` (7 endpoints) | Avenium | ✅ DONE | Commit `3a3d7d9` — 100% JSON parity (23/23 cases + 3/3 HTTP 422 checks vs `:8091` across `tickers`, `explorer/scan` normal/advanced/bigPlayer, `analysis`, `t5`, `group-performance`, `broker-drilldown`, `detail` & 404); 7/7 unit tests + full `go test ./...` PASS; `VmRSS` **22.272 kB** (`~21,7 MB`, `VmHWM` 23.732 kB) |
| `terminal/ipo_hunter` (4 endpoints) | Avenium | ✅ DONE | Wave 3E-3 — 100% JSON parity (**36/36** GET/POST cases + **4/4** HTTP 422 checks vs `:8091` across `meta`, `list` semua filter/bucket/sort, `analysis` after/periodic/override/group/404, & `screener` periodic/after/comma-conds); 7/7 unit tests + full `go test ./...` PASS; indexed SQL + streaming screener tanpa preload JSONB 281,7 MB → `VmRSS` **31.540 kB** (`~30,8 MB`, `VmHWM` 33.200 kB) |

## Tracker Modul Frontend (Vite + React 19)

| Modul / Halaman | Sumber | Status | Keterangan |
| :--- | :--- | :---: | :--- |
| Shell Layout (`ClientLayout`, `Sidebar`, `BottomNav`, `TerminalDropdown`, `ThemeToggle`, `RefreshButton`, `LogoutButton`) | InvestOwl + Dropdown Terminal | ✅ DONE | 100% MD3 InvestOwl + `vite build` 4.71s |
| `features/auth/` (`LoginPage`, `PendingPage`) & `features/admin/` (`AdminPage`) | InvestOwl | ✅ DONE | `tsc --noEmit` PASS |
| Shared (`TickerSearch`, `InvestOwlLoader`, `useForeignFlow`) + `features/home/`, `ticker/`, `signal/` (Wave 4A) | InvestOwl (`/`, `/:ticker`, `/signal`, `/backtest`) | ✅ DONE | 31 file modular + routing di `App.tsx`; `tsc --noEmit` + `vite build` 10.30s (`mem-guard`); SPA (`/` `/BBCA` `/signal` `/backtest` + JS bundle) & 12/12 API smoke HTTP 200 PASS; `VmRSS` 48.944 kB (~47,8 MB) |
| `features/foreign/`, `ownership/`, `group/`, `explorer/` + `components/research/` (Wave 4B) | InvestOwl (`/foreign`, `/ownership`, `/group`, `/explorer`) | ✅ DONE | 110 file modular (12.438 LOC) + routing di `App.tsx`; `tsc --noEmit` PASS, zero Next.js/BEI1%, `vite build` 26.11s (`mem-guard`); SPA 4/4 (`/foreign` `/ownership` `/group` `/explorer` + JS bundle) & 13/13 API smoke HTTP 200 PASS; `VmRSS` 34.628 kB (~33,8 MB), `VmHWM` 37.028 kB |
| `lib/api/terminal.ts` + `features/mmd/` + `features/ipo-hunter/` (Wave 4C) | Avenium (`/mmd`, `/ipo-hunter`, MD3 restyle, tanpa `bei1persen`) | ✅ DONE | 33 file modular (6.568 LOC: `terminal.ts` 157, `mmd` 14 file / 3.012 LOC, `ipo-hunter` 18 file / 3.399 LOC) + routing di `App.tsx`; `tsc --noEmit` PASS, zero Next.js/BEI1%, `vite build` 21.38s (`mem-guard`, `MmdPage` 91,29 kB, `IpoHunterPage` 76,20 kB); SPA (`/mmd`, `/ipo-hunter`) & 11/11 API HTTP 200 PASS (`mmd` 6 + `ipo-hunter` 4 + `ownership` 1); `VmRSS` 23.756 kB (~23,2 MB), `VmHWM` 26.896 kB |
| `features/foreign-pro/`, `konglo/`, `market-report/` (Wave 4D) | Avenium (`/foreign-pro`, `/konglo`, `/admin/market-report`, MD3 restyle, tanpa `bei1persen`) | ✅ DONE | 42 file modular (8.416 LOC: `foreign-pro` 23 file / 4.821 LOC maks 335 LOC, `konglo` 9 file / 1.673 LOC maks 267 LOC, `market-report` 10 file / 1.922 LOC maks 311 LOC) + `terminal.ts` (279 LOC) + `App.tsx` wired & `FeaturePlaceholder.tsx` dihapus; `tsc --noEmit` PASS, `GREP_CLEAN`, `vite build` 22.18s (`mem-guard`: `ForeignProPage` 111,54 kB, `KongloPage` 46,54 kB, `MarketReportPage` 58,51 kB); SPA 3/3 + JS 4/4 + API 9/9 HTTP 200 + `market-report/logs` 401 guard PASS; `VmRSS` 33.020 kB (~32,2 MB), `VmHWM` 35.236 kB |
