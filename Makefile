postgres:
	docker run --name postgres16 -p 5433:5433 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16.1-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres16 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:password@localhost:5433/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:password@localhost:5433/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb, dropdb, postgres, migrateup, migratedown, sqlc
# .PHONY is to specify it is a command in the Makefile instead of a file