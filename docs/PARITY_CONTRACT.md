# LazyTrader — Kontrak Paritas & Verifikasi Non-Regresi

## 1. Aturan Zero Mutation
1. Tidak ada rumus matematika, bobot skor, logika statistik, atau algoritma trading dari `/root/investowl`, `/root/rev`, maupun `/root/SMtracker` yang boleh diubah.
2. Seluruh kolom JSON, nama field, tipe data, satuan (Rupiah, lot, persen, rasio), serta kode status HTTP dan struktur error (`{"detail": "..."}`) harus identik dengan backend sumber (`:8001` dan `:8091`).
3. **Canary Makro IHSG (`hybrid_p5_p15`):** Wajib lulus `45/45` API parity dan `0/1.034` mismatch terhadap `docs/STATE_SNAPSHOT.json` sebelum cutover.

## 2. Protokol Golden JSON Diff (`tools/parity/compare_json.py`)
Setiap endpoint yang selesai dimigrasi di Fase 2–3 wajib diverifikasi menggunakan:
```bash
python3 tools/parity/compare_json.py \
  --baseline http://127.0.0.1:8001 \
  --candidate http://127.0.0.1:8099 \
  --path /api/<endpoint> \
  --rtol 1e-6 --atol 1e-6
```
- Tipe data berbeda (mis. `int` vs `string`, `null` vs `0`) → **GAGAL**.
- Key hilang atau berlebih → **GAGAL**.
- Selisih float melebihi toleransi `1e-6` → **GAGAL**.

## 3. Penanganan Endpoint Statistik Khusus
- `/api/bandar/validation/{ticker}` (Student-t p-value): Diimplementasikan di Go dengan verifikasi golden test terhadap `scipy.stats.t.sf`.
- `/api/bandar/causality/{ticker}` (ADF & Granger Causality OLS F-test): Diimplementasikan dengan golden test ketat terhadap `statsmodels`, atau dieksekusi via helper Python on-demand jika diperlukan demi jaminan paritas 100% tanpa menahan RAM idle.
