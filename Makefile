run:
	go run ./cmd/api/*.go

postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -p 5432:5432 -d postgres:12.2-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root todo

dropdb:
	docker exec -it postgres12 dropdb todo

psql:
	docker exec -it postgres12 psql -U root -d todo

migrate_cli:
	@read -p "Enter migration name: " name; \
	migrate create -seq -ext=.sql -dir=./migrations $$name

migrateup:
	migrate -path ./migrations -database "postgres://root:root@localhost:5432/todo?sslmode=disable" -verbose up

migratedown:
	migrate -path ./migrations -database "postgres://root:root@localhost:5432/todo?sslmode=disable" -verbose down

docker_up:
	docker-compose up

.PHONY: run createdb dropdb psql migrate_cli docker_up postgres migrateup migratedown