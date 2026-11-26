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
| **3** | Migrasi Backend Domain (11 domain InvestOwl + 5 domain Avenium) + Golden JSON Parity | 🔄 **IN PROGRESS** | `8/11` domain InvestOwl + `portal` DONE (Commit `915a292` + Commit 2); 24/24 read parity PASS (`:8001`); Canary `0/1032` mismatch; `VmRSS` **27,4 MB** |
| **4** | Migrasi Frontend Features (seluruh halaman InvestOwl + 5 modul Avenium dipecah modular & MD3 restyle) | ⬜ **TODO** | — |
| **5** | Verifikasi Paritas Total (`45/45` Canary, unit/e2e test, benchmark RAM & disk) | ⬜ **TODO** | — |
| **6** | Deployment & Cutover Cloudflared *(menunggu persetujuan owner)* | ⬜ **TODO** | — |

---

## Tracker Domain Backend (Go)

| Domain | Sumber | Status | Parity Check |
| :--- | :--- | :---: | :---: |
| `core/*` (`config`, `database`, `cache`, `flags`, `auth`, `middleware`, `spa`) | InvestOwl + Avenium | ✅ DONE | 5/5 unit tests PASS |
| `investowl/health` | InvestOwl + Avenium | ✅ DONE | `compare_json.py` PASS |
| `investowl/users` (9 endpoints + auth-gate) | InvestOwl `routers/auth.py` | ✅ DONE | `auth-gate` parity PASS |
| `terminal/portal` (tanpa `bei1persen`) | Avenium `routers/portal.py` | ✅ DONE | Unit test + curl no-bei1persen PASS |
| `investowl/stocks`, `broker`, `system` (16 endpoints + `shared/*`) | InvestOwl | ✅ DONE | Commit `915a292` — 100% JSON parity (6/6 endpoints) + `/api/logs/*` auth & dry-run smoke PASS; `system_test.go` PASS |
| `investowl/forensic`, `foreign_flow`, `signal` (13 endpoints + `shared/{stats,causality,conviction}`) | InvestOwl | ✅ DONE | 100% JSON parity (18/18 endpoints vs `:8001` incl. `confluence-screener`, HTTP 404); Canary `0/1032` mismatch (`absolute` & `hybrid_p5_p15`); `VmRSS` 28.032 kB |
| `investowl/bandarmology`, `ownership`, `shareholders`, `nyangkut` | InvestOwl | ⏳ NEXT | Fase 3 Wave 3C–3D |
| `terminal/mmd`, `foreign_flow_pro`, `konglo`, `ipo_hunter`, `market_report` | Avenium | ⬜ TODO | Fase 3 Wave 3E |

## Tracker Modul Frontend (Vite + React 19)

| Modul / Halaman | Sumber | Status | Keterangan |
| :--- | :--- | :---: | :--- |
| Shell Layout (`ClientLayout`, `Sidebar`, `BottomNav`, `TerminalDropdown`, `ThemeToggle`, `RefreshButton`, `LogoutButton`) | InvestOwl + Dropdown Terminal | ✅ DONE | 100% MD3 InvestOwl + `vite build` 4.71s |
| `features/auth/` (`LoginPage`, `PendingPage`) & `features/admin/` (`AdminPage`) | InvestOwl | ✅ DONE | `tsc --noEmit` PASS |
| `features/home/`, `ticker/`, `signal/`, `foreign/`, `ownership/`, `group/`, `explorer/` | InvestOwl | ⬜ TODO | Fase 4 |
| `features/mmd/`, `ipo-hunter/`, `foreign-pro/`, `konglo/`, `market-report/` | Avenium (pecah modular + MD3) | ⬜ TODO | Fase 4 |
