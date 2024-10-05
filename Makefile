

run:
	go run ./cmd/api/ .

createdb:
	docker exec -it postgres12.2 createdb --username=nightbot --owner=nightbot todo_db

dropdb:
	docker exec -it postgres12.2 dropdb todo_db

psql:
	docker exec -it postgres12.2 psql -U nightbot -d todo_db

migrate_cli:
	@read -p "Enter migration name: " name; \
	migrate create -seq -ext=.sql -dir=./migrations $$name

docker_up:
	docker-compose up