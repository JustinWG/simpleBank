postgres:
	docker run --name simpleBankPostgres --network simplebankNetwork -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createdb:
	docker exec -it simpleBankPostgres createdb --username=root --owner=root simpleBankPostgres

dropdb:
	docker exec -it simpleBankPostgres dropdb simpleBankPostgres

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/JustinWG/simpleBank/db/sqlc Store

.PHONY: createdb dropdb postgres migrateup migratedown sqlc mock migratedown1 migrateup1