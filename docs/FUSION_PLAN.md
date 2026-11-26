# LazyTrader — Master Fusion Plan (InvestOwl × Avenium)

> **Repositori:** `git@github.com:aryarifki/lazytrader.git`  
> **Tech Stack Terkonfirmasi:** Go (`1.27+`, 1 Binary API + Static SPA Server) + Vite + React 19 SPA  
> **Navigasi Terkonfirmasi:** Menu Dropdown **"Terminal"** terpisah di Header InvestOwl (`/mmd`, `/ipo-hunter`, `/foreign-pro`, `/konglo`, `/admin/market-report`)  
> **Status:** Fase 1 Inisialisasi Repositori

---

## 1. Kontrak Non-Regresi (Aturan Emas)

1. **100% UI & UX InvestOwl:**
   - Seluruh warna, tipografi, radius, animasi, dan layout menggunakan token Material Design 3 warm-pastel InvestOwl (`--md-sys-color-*`, `--color-positive`, `--color-negative`), `ClientLayout`, `Sidebar`, `BottomNav`, serta `Inter` + `.font-mono-nums`.
   - Seluruh 5 modul dari Avenium di-restyle ke token CSS & komponen UI InvestOwl.
2. **Pengecualian BEI 1%:**
   - Modul `bei1persen` (`bei1persen_engine.py`, `routers/bei1persen.py`, `app/bei1persen/page.tsx`, `db.json`, `aux.json`) **tidak digabungkan** karena InvestOwl sudah memiliki `/ownership`, `/group`, `/explorer`.
   - Di `/api/v1/terminal/portal/status` & `/entitlements`, entri `bei1persen` dihapus.
3. **Zero Logic / Math / ML / SMtracker Mutation:**
   - Dilarang mengubah rumus matematika, bobot scoring, algoritma, threshold, atau pipeline ML di `/root/investowl`, `/root/rev`, dan `/root/SMtracker`.
   - **Canary Makro IHSG** (`hybrid_p5_p15`) wajib mempertahankan **100% paritas** (`45/45` API parity, `0/1.034` mismatch).
   - **Redis:** Hanya `SCAN + UNLINK` dengan prefix terlindungi `rate_limit:` dan `setting:`. Dilarang `FLUSHDB` / `FLUSHALL`.
   - **Pipeline Harian (`run_daily.sh`):** Tetap dijalankan **manual oleh operator** (tidak ada cron 22:00).
4. **Anti-Monolith & Maintainability:**
   - Setiap fungsi/halaman/domain wajib dipisah ke file dan folder modular (target ≤ 250–350 baris per file).

---

## 2. Baseline Repositori Sumber (Verified)

| Repositori | Path Lokal | Remote | Commit HEAD |
| :--- | :--- | :--- | :--- |
| **InvestOwl** | `/root/investowl` | `git@github.com:aryarifki/investowl.git` | `f1ed6ad` |
| **Avenium** | `/root/rev` | `git@github.com:aryarifki/Avenium.git` | `1e7cf57` |
| **SMtracker** | `/root/SMtracker` | `git@github.com:aryarifki/SMtracker.git` | `b2744d2` |

---

## 3. Matriks Penggabungan Fitur

- **InvestOwl Core (100%):**
  - Halaman: `/`, `/[ticker]`, `/signal`, `/backtest`, `/foreign`, `/ownership`, `/group`, `/explorer`, `/admin`, `/login`, `/pending`.
  - Backend (11 domain, 99 endpoint): `auth`/`users`, `bandarmology`, `broker`, `foreign_flow`, `forensic`, `nyangkut`, `ownership`, `shareholders`, `signal`, `stocks`, `system` + `health`.
- **Avenium Terminal Modules (5 Modul + Portal, 28 endpoint):**
  - `mmd` (MM Detector Pro, 7 endpoint) → `/mmd`
  - `ipo-hunter` (IPO Hunter Pro, 4 endpoint) → `/ipo-hunter`
  - `foreign-flow` (Foreign Flow Pro, 3 endpoint) → `/foreign-pro`
  - `konglo` (Konglo Screener, 4 endpoint) → `/konglo`
  - `market-report` (Market Report Generator Admin, 8 endpoint, **wajib admin session**) → `/admin/market-report`
  - `portal` (2 endpoint, tanpa `bei1persen`)
