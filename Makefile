export PATH := /usr/local/go/bin:$(PATH)
export TMPDIR := /var/tmp
export GOCACHE := /var/tmp/go-build

.PHONY: test-backend build-backend build-frontend build dev-backend dev-frontend parity-smoke

test-backend:
	cd backend && go test -v ./...

build-backend:
	cd backend && go build -ldflags="-s -w" -o bin/lazytrader ./cmd/server

build-frontend:
	cd frontend && NODE_OPTIONS=--max-old-space-size=1024 npx tsc --noEmit && NODE_OPTIONS=--max-old-space-size=1024 npm run build

build: build-frontend build-backend

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

parity-smoke:
	python3 tools/parity/compare_json.py --baseline http://127.0.0.1:8001 --candidate http://127.0.0.1:8099 --path /api/health
	python3 tools/parity/compare_json.py --baseline http://127.0.0.1:8001 --candidate http://127.0.0.1:8099 --path /api/users/settings/auth-gate
