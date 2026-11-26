# LazyTrader

**LazyTrader** adalah platform penggabungan (*fusion*) **InvestOwl** dan **Avenium Terminal** menjadi **1 aplikasi tunggal yang ultra-ringan** dengan mempertahankan **100% UI & UX InvestOwl**, **100% fitur Avenium (kecuali BEI 1%)**, serta **zero perubahan rumus matematika, logika bisnis, maupun pipeline ML (`SMtracker`)**.

## Tech Stack Terkonfirmasi

- **Backend & Static Server:** Go (`1.27+`, `chi`, `pgx/v5`, `go-redis/v9`, `golang-jwt/jwt/v5`) — 1 single static binary ARM64/AMD64 (~12–25 MB RSS).
- **Frontend:** Vite + React 19 + TypeScript + Tailwind CSS v4 + ECharts + Recharts + SWR + Zustand (Static SPA disajikan langsung oleh binary Go, 0 proses Node.js di produksi).
- **Database & Cache:** PostgreSQL 17 (`bandarmology`) + Redis (`db0` InvestOwl, `db1` Terminal, safe `SCAN + UNLINK`).
- **Batch ETL & ML:** `/root/SMtracker` tetap di Python (dipicu on-demand, tidak diubah).

## Dokumentasi Kanonial

- [docs/FUSION_PLAN.md](./docs/FUSION_PLAN.md) — Rencana master, matriks fitur, dan aturan non-regresi.
- [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) — Arsitektur modular anti-monolith & peta rute.
- [docs/PARITY_CONTRACT.md](./docs/PARITY_CONTRACT.md) — Kontrak paritas numerik/JSON & Canary Exodus.
- [docs/PROGRESS.md](./docs/PROGRESS.md) — Tracker progres per fase dan domain.

## Quick Build & Test

```bash
cp .env.example .env
make test-backend
make build
```
