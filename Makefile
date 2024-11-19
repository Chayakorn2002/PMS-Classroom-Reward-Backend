start:
	go run cmd/app/main.go

gen-sqlc:
	sqlc generate

migrate-up:
	go run cmd/migration/main.go --migrate:up

migrate-down:
	go run cmd/migration/main.go --migrate:down --step=$(step)

migrate-flush:
	go run cmd/migration/main.go --migrate:flush

migrate-make:
	go run cmd/migration/main.go --migrate:make --name=$(name)

docker-migrate-up:
	docker compose exec api go run cmd/migration/main.go --migrate:up

docker-migrate-down:
	docker compose exec api go run cmd/migration/main.go --migrate:down --step=$(step)

docker-migrate-reset:
	docker compose exec api go run cmd/migration/main.go --migrate:reset

docker-migrate-make:
	docker compose exec api go run cmd/migration/main.go --migrate:make --name=$(name)