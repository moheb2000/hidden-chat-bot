include .envrc

.PHONY: run/bot
run/bot:
	go run ./cmd/bot

.PHONY: run/bot/fa
run/bot/fa:
	go run ./cmd/bot -locale=fa-ir

.PHONY: build/bot
build/bot:
	@echo 'Building cmd/bot'
	go build -o=./bin/bot ./cmd/bot

.PHONY: migration/up
migration/up:
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${HIDDEN_CHAT_DB_DSN} up

.PHONY: migration/down
migration/down:
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${HIDDEN_CHAT_DB_DSN} down

.PHONY: migration/new
migration/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -ext=.sql -dir=./migrations ${name}
