# LazyTrader — Progress Tracker

> **Terakhir Diperbarui:** 2026-11-26  
> **Repo:** `git@github.com:aryarifki/lazytrader.git`

---

## Ringkasan Fase

| Fase | Deskripsi | Status | Catatan / Bukti |
| :---: | :--- | :---: | :--- |
| **0** | Pre-flight Commit & Push (`/root/rev` `1e7cf57`, `/root/investowl` `f1ed6ad`, `/root/SMtracker` `b2744d2`) | ✅ **DONE** | Semua repo clean (`0/0` vs `origin/main`) |
| **1** | Inisialisasi Repo `/root/lazytrader`, struktur modular, & dokumentasi kanonial `docs/` | 🔄 **IN PROGRESS** | Commit #1 |
| **2** | Core Platform Go (`config`, `database`, `cache`, `flags`, `auth`, `middleware`, `spa`) + `health`/`users`/`portal` + Vite React 19 Shell + Parity Harness | ⏳ **NEXT** | Target Commit #2 |
| **3** | Migrasi Backend Domain (11 domain InvestOwl + 5 domain Avenium) + Golden JSON Parity | ⬜ **TODO** | — |
| **4** | Migrasi Frontend Features (seluruh halaman InvestOwl + 5 modul Avenium dipecah modular & MD3 restyle) | ⬜ **TODO** | — |
| **5** | Verifikasi Paritas Total (`45/45` Canary, unit/e2e test, benchmark RAM & disk) | ⬜ **TODO** | — |
| **6** | Deployment & Cutover Cloudflared *(menunggu persetujuan owner)* | ⬜ **TODO** | — |

---

## Tracker Domain Backend (Go)

| Domain | Sumber | Status | Parity Check |
| :--- | :--- | :---: | :---: |
| `core/*` (`config`, `database`, `cache`, `flags`, `auth`, `middleware`, `spa`) | InvestOwl + Avenium | ⬜ TODO | Unit tests |
| `investowl/health` | InvestOwl + Avenium | ⬜ TODO | `compare_json.py` |
| `investowl/users` (9 endpoints + auth-gate) | InvestOwl `routers/auth.py` | ⬜ TODO | `compare_json.py` |
| `terminal/portal` (tanpa `bei1persen`) | Avenium `routers/portal.py` | ⬜ TODO | Unit test no-bei1persen |
| `investowl/stocks`, `broker`, `system` | InvestOwl | ⬜ TODO | — |
| `investowl/bandarmology`, `signal`, `foreign_flow`, `forensic` | InvestOwl | ⬜ TODO | — |
| `investowl/ownership`, `shareholders`, `nyangkut` | InvestOwl | ⬜ TODO | — |
| `terminal/mmd`, `foreign_flow_pro`, `konglo`, `ipo_hunter`, `market_report` | Avenium | ⬜ TODO | — |

## Tracker Modul Frontend (Vite + React 19)

| Modul / Halaman | Sumber | Status | Keterangan |
| :--- | :--- | :---: | :--- |
| Shell Layout (`ClientLayout`, `Sidebar`, `BottomNav`, `TerminalDropdown`, `ThemeToggle`, `RefreshButton`, `LogoutButton`) | InvestOwl + Dropdown Terminal | ⬜ TODO | 100% MD3 InvestOwl |
| `features/auth/` (`LoginPage`, `PendingPage`) & `features/admin/` (`AdminPage`) | InvestOwl | ⬜ TODO | 100% identik |
| `features/home/`, `ticker/`, `signal/`, `foreign/`, `ownership/`, `group/`, `explorer/` | InvestOwl | ⬜ TODO | Fase 4 |
| `features/mmd/`, `ipo-hunter/`, `foreign-pro/`, `konglo/`, `market-report/` | Avenium (pecah modular + MD3) | ⬜ TODO | Fase 4 |
