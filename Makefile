.PHONY: dev dev-db dev-api dev-web migrate migrate-down sqlc test

dev:
	overmind start

dev-db:
	docker compose up -d

dev-api: dev-db
	cd api && air

dev-web:
	cd web && npm run dev

migrate:
	cd api && go run main.go migrate up

migrate-down:
	cd api && go run main.go migrate down

sqlc:
	cd api && sqlc generate

test:
	(cd api && go test ./...) && (cd web && npm run build)
