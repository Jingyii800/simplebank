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

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./..

server:
	go run main.go

db_docs:
	dbdocs build doc/db.dbml

db_sql:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Jingyii800/simplebank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/Jingyii800/simplebank/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.json
	protoc --proto_path=proto \
	   --experimental_allow_proto3_optional \
       --go_out=pb \
       --go_opt=paths=source_relative \
       --go-grpc_out=pb \
       --go-grpc_opt=paths=source_relative \
       --grpc-gateway_out=pb \
       --grpc-gateway_opt=paths=source_relative \
       --openapiv2_out=doc/swagger \
       --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
       proto/*.proto
	   statik -src=./doc/swagger -dest=./doc -f

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: createdb, dropdb, postgres, migrateup, migratedown, sqlc, test,server, db_docs, db_sql, mock, proto, evans, redis, new_migration
# .PHONY is to specify it is a command in the Makefile instead of a file