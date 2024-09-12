include .envrc

# ------------ Run App ---------------
.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${DB_DSN} -port=${APP_PORT} -secret-key=${SECRET_KEY}


# ------------ Migration ---------------
.PHONY: migration/new
migration/new:
	migrate create -seq -ext=.sql -dir ./migrations ${name}

.PHONY: migration/up
migration/up:
	migrate -path ./migrations -database=${DB_DSN} up