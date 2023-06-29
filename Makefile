# build a new postgres container
postgres_up:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:14.0
# delete exist container of postgres
postgres_down:
	docker stop postgres
	docker rm postgres
# create database
create_db:
	docker exec -it postgres createdb --username=root --owner=root simple_bank
# drop database
drop_db:
	docker exec -it postgres dropdb simple_bank

# migrate database
migration_up:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migration_down:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

# migrate database
migration_up1:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migration_down1:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1
# migrate database
migration_up2:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 2

migration_down2:
	migrate -path db/migration -database "postgres://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 2

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb --build_flags=--mod=mod -destination db/mock/store.go github.com/borntodie-new/backend-master-class/db/sqlc Store

server:
	go run main.go

protoc_gen:
	rm -f pb/*.proto
	rm -f doc/swagger/*
	protoc --proto_path=./proto --go_out=./pb --go_opt=paths=source_relative \
	--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger \
	--openapiv2_opt=allow_merge=true,merge_file_name=simpla_bank \
	./proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

redis_up:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

redis_down:
	docker stop redis
	docker rm redis

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)
#.PYONY: postgresup createdb dropdb