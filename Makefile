.PHONY: migrate_up migrate_down migrate_down1 migrate_force sqlc test test-sqlc

DB_URL=postgres://admin:admin123@192.168.1.8:5432/simplebank?sslmode=disable
MIGRATE_PATH=C:/GoCode/github.com/thinhcompany/simple-bank/db/migrations

migrate_up:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" up

migrate_down:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" down

migrate_down1:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" down 1

migrate_force:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" force 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

test-sqlc:
	go test -v -cover ./db/sqlc

