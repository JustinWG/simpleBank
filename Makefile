postgres:
	docker run --name simpleBankPostgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

createdb:
	docker exec -it simpleBankPostgres createdb --username=root --owner=root simpleBankPostgres

dropdb:
	docker exec -it simpleBankPostgres dropdb simpleBankPostgres

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simpleBankPostgres?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createdb dropdb postgres migrateup migratedown sqlc