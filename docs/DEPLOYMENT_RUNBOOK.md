# LazyTrader — Production Deployment & Cutover Runbook

> **Status Cutover:** ✅ **LIVE di `https://terminal.avenia.my.id`** (2026-11-26)  
> **Service Systemd:** `/etc/systemd/system/lazytrader.service` (template di [deploy/lazytrader.service](file:///root/lazytrader/deploy/lazytrader.service))  
> **Legacy Fallback (`ct.avenia.my.id`):** `avenium-backend.service` (`127.0.0.1:8091`) + `avenium-frontend.service` (`127.0.0.1:3002`)

---

## 1. Arsitektur Produksi & Routing Cloudflared

| Subdomain | Service Target | Unit Systemd | Status |
| :--- | :--- | :--- | :---: |
| **`terminal.avenia.my.id`** | `http://localhost:8099` (LazyTrader Single Binary Go API + Vite SPA) | `lazytrader.service` | ✅ Live HTTPS (`200 OK`) |
| **`ct.avenia.my.id`** (`/api/.*` & `/`) | `http://localhost:8091` & `http://localhost:3002` (Legacy Avenium) | `avenium-backend` + `avenium-frontend` | 🟡 Origin + Tunnel siap; tunggu DNS `CNAME ct` |
| `investowl.avenia.my.id` | `http://localhost:8001` & `http://localhost:3001` (Legacy InvestOwl) | `investowl-backend` + `investowl-frontend` | ✅ Live HTTPS (`200 OK`, tidak tersentuh) |

- **Config Tunnel:** `/etc/cloudflared/config.yml` (Tunnel ID `774c1dad-f104-457d-81cd-08c878bfbd4f`).
- **Backup Pra-Cutover:** `/etc/cloudflared/config.yml.bak.20261126-phase6`.
- **Environment File:** `/root/lazytrader/.env` (permission `600`, gitignored — berisi `APP_HOST=127.0.0.1`, `APP_PORT=8099`, `STATIC_DIST_DIR=/root/lazytrader/frontend/dist`, `APP_BASE_URL=https://terminal.avenia.my.id`, `REDIS_KEY_PREFIX=lt_`, serta kredensial DB/Auth/Admin).

---

## 2. Dua Langkah Manual Eksternal (Dashboard Owner)

Karena token Cloudflare API di server tidak memiliki akses DNS zona `avenia.my.id` (dan Google Cloud Console dikelola langsung oleh owner), lakukan dua langkah berikut di dashboard web:

1. **Cloudflare Dashboard (`avenia.my.id` → DNS → Records):**
   - Tambahkan record **CNAME** *(Proxied / awan oranye)*:
     - **Name:** `ct`
     - **Target:** `774c1dad-f104-457d-81cd-08c878bfbd4f.cfargotunnel.com`
   - Verifikasi setelah propagation: `curl -I https://ct.avenia.my.id/` → `200 OK`.
2. **Google Cloud Console (APIs & Services → Credentials → OAuth 2.0 Client ID):**
   - Pastikan **Authorized JavaScript origins** mencantumkan `https://terminal.avenia.my.id`.
   - Pastikan **Authorized redirect URIs** mencantumkan:
     `https://terminal.avenia.my.id/api/auth/callback/google`

---

## 3. Operasi Harian & Rebuild Aman Container 3,6 GB

```bash
# Cek status service & memori
systemctl status lazytrader.service
/usr/local/bin/mem-guard status

# Rebuild frontend + backend aman
cd /root/lazytrader/frontend && TMPDIR=/var/tmp NODE_OPTIONS=--max-old-space-size=1024 /usr/local/bin/mem-guard build -- npm run build
cd /root/lazytrader/backend && PATH="/usr/local/go/bin:$PATH" GOCACHE=/var/tmp/go-build go build -trimpath -ldflags="-s -w" -o bin/lazytrader ./cmd/server
systemctl restart lazytrader.service
/usr/local/bin/mem-guard clean
```

---

## 4. Prosedur Rollback Instan

Jika sewaktu-waktu perlu mengembalikan `terminal.avenia.my.id` ke stack Avenium lama (`:8091` / `:3002`):

```bash
cp -a /etc/cloudflared/config.yml.bak.20261126-phase6 /etc/cloudflared/config.yml
cloudflared tunnel --config /etc/cloudflared/config.yml ingress validate
systemctl restart cloudflared.service
```
