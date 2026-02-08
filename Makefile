.PHONY: migrate_up migrate_down migrate_down1 migrate_force sqlc migrate-create migrate_up1 \
		test test-sqlc server air mock test-api test-util

DB_URL=postgres://admin:admin123@192.168.1.8:5432/simplebank?sslmode=disable
# DB_URL=postgres://root:root123@localhost:5432/testdb?sslmode=disable

MIGRATE_PATH=db/migrations

## Create new migration
## how to use: make migrate-create name=add_users
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

migrate_up:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" up

migrate_up1:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" up 1

migrate_down:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" down

migrate_down1:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" down 1

migrate_force:
	migrate -verbose -path "$(MIGRATE_PATH)" -database "$(DB_URL)" force 1

sqlc:
	sqlc generate

test:
	go test -v -count=1 -cover ./...

test-sqlc:
	go test -v -count=1 -cover ./db/sqlc

test-api:
	go test -v -count=1 -cover ./api

test-util:
	go test ./util -v


server:
	go run .

air:
	air

mock:
	mockgen \
		-package mockdb \
		-destination db/mock/store.go \
		simple-bank/db/sqlc Store


