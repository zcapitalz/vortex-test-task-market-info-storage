.PHONY: up-clickhouse

up-clickhouse:
	docker compose --file=./deploy/docker-compose.yaml up clickhouse

.PHONY: up-clickhouse

up-postgres:
	docker compose --file=./deploy/docker-compose.yaml up postgres

.PHONY: migrate-up-clickhouse

migrate-clickhouse-up:
	migrate -database 'clickhouse://localhost:9000?username=default&password=1312&x-multi-statement=true' -path ./deploy/migrations up

.PHONY: generate-swagger-docs

generate-swagger-docs:
	swag init --generalInfo=./internal/app/app.go --parseInternal --parseDependency --output=./api/v1