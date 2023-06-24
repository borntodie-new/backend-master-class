# build a new postgres container
postgresup:
	@echo "build a new postgres container..."
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.0
	@echo "buold a new postgres container success!"
# delete exist container of postgres
postgresdown:
	@echo "stop postgres container..."
	docker stop postgres
	@echo "delete postgres container..."
	docker rm postgres
	@echo "delete postgres container success..."
# create database
createdb:
	@echo "create database..."
	docker exec -it postgres createdb --username=root --owner=root simple_bank
	@echo "create database done..."
# drop database
dropdb:
	@echo "drop database..."
	docker exec -it postgres dropdb simple_bank
	@echo "drop database done..."

# migrate database
migrationup:
	@echo "migrate database for create..."
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up
	@echo "migrate database success..."

migrationdown:
	@echo "migrate database for delete..."
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down
	@echo "migrate database success..."

sqlc:
	sqlc generate

#.PYONY: postgresup createdb dropdb