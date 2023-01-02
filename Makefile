DB_URL=postgresql://root:secret@postgres:5432/bookmark?sslmode=disable

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
	mockgen -package mockdb -destination db/mock/store.go github.com/inkclip/backend/db/sqlc Store &&  mockgen -package mockmail -destination mail/mock/client.go github.com/inkclip/backend/mail Client
swag:
	swag init
openswag:
	open http://0.0.0.0:8080/swagger/index.html
openmail:
	open http://localhost:1080/
air:
	air -c .air.toml
gosec:
	gosec -exclude=G101 -tests ./...

.PHONY: dropdb migrateup migratedown sqlc test server mock gosec testcoverage swag openswag air openmail
