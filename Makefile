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


# migrate database
migrationup1:
	@echo "migrate database for create..."
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1
	@echo "migrate database success..."

migrationdown1:
	@echo "migrate database for delete..."
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
	@echo "migrate database success..."

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb --build_flags=--mod=mod -destination db/mock/store.go github.com/borntodie-new/backend-master-class/db/sqlc Store

server:
	go run main.go

proto_gen:
	rm -f pb/*.proto
	protoc --proto_path=./proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative ./proto/*.proto

#.PYONY: postgresup createdb dropdb