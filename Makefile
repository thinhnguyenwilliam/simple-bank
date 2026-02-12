# simple-bank\Makefile
.PHONY: migrate_up migrate_down migrate_down1 migrate_force sqlc migrate-create migrate_up1 \
		test test-sqlc server air mock test-api test-util test-token test-only-fuction \
		docker-build docker-run go-build clean compose-up compose-down dbdocs \
		dbdocs-set-password gen

DB_URL=postgres://admin:admin123@192.168.1.8:5432/simplebank?sslmode=disable

MIGRATE_PATH=db/migrations
IMAGE_NAME=simple-bank
ENV_FILE=app.env

APP_NAME=simple-bank
BUILD_DIR=bin
PASSWORD=secret123

gen:
	del /Q pb\*.go
	protoc --proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

dbdocs:
	dbdocs build doc/db.dbml

dbdocs-set-password:
	dbdocs password -s $(PASSWORD) -p simple_bank


compose-up:
	docker compose up --build

compose-down:
	docker compose down

compose-down-v:
	docker compose down -v


go-build:
	go build -o $(BUILD_DIR)/$(APP_NAME).exe .

clean:
	if exist $(BUILD_DIR)\$(APP_NAME).exe del $(BUILD_DIR)\$(APP_NAME).exe


# ========================
# Docker
# ========================
docker-build:
	docker build -t $(IMAGE_NAME) .

docker-run:
	docker run --rm \
		--name simple-bank \
		-p 8080:8080 \
		--env-file $(ENV_FILE) \
		$(IMAGE_NAME)


docker-run-bg:
	docker run -d -p 8080:8080 \
		--env-file $(ENV_FILE) \
		--name $(IMAGE_NAME) \
		$(IMAGE_NAME)

docker-stop:
	docker stop $(IMAGE_NAME) || true
	docker rm $(IMAGE_NAME) || true

## Create new migration
## how to use: make migrate-create name=add_users
# make migrate-create name=add_sessions
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

test-token:
	go test -v -count=1 -cover ./token

test-only-fuction:
	go test -run TestJWTMakerCreateAndVerifyToken ./token

server:
	go run .

air:
	air

mock:
	mockgen \
		-package mockdb \
		-destination db/mock/store.go \
		github.com/thinhcompany/simple-bank/db/sqlc Store


