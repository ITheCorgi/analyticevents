.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune

.PHONY: codegen
codegen: ## Generate server.
	oapi-codegen --config oapi-cfg.yaml ./api/swagger.yaml

.PHONY: build
build: ## Build server.
	GOOS=darwin GOARCH=amd64 go build -o analyticevents cmd/analytic/main.go

.PHONY: migrate-up
migrate_file: build ## make migrate up.
	./getground-party migrate up
