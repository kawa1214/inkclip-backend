DB_URL=postgresql://root:secret@localhost:5432/bookmark?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres
createdb:
	docker exec -it postgres createdb --username=root --owner=root bookmark
dropdb:
	docker exec -it postgres dropdb bookmark
migrateup:
	migrate -path db/migration -database "$(DB_URL)" --verbose up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)"  --verbose down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

testcoverage:
	go test -coverprofile=.coverage/coverage.out ./...
	go tool cover -html=.coverage/coverage.out -o .coverage/coverage.html

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/inkclip/backend/db/sqlc Store

swag:
	swag init

openswag:
	open http://0.0.0.0:8080/swagger/index.html

openmail:
	open http://localhost:1080/

air:
	air -c .air.toml

gosec:
	gosec -tests ./...

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock gosec testcoverage swag openswag air openmail
