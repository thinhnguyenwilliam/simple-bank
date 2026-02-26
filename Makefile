# simple-bank\Makefile
.PHONY: migrate_up migrate_down migrate_down1 migrate_force sqlc migrate-create migrate_up1 \
		test test-sqlc server air mock test-api test-util test-token test-only-fuction \
		docker-build docker-run go-build clean compose-up compose-down dbdocs \
		dbdocs-set-password gen evans proto proto-lint proto-gen compose-up-redis \
		proto-2 laptop-server laptop-client

DB_URL=postgres://admin:admin123@localhost:5432/simplebank?sslmode=disable

MIGRATE_PATH=db/migrations
IMAGE_NAME=simple-bank
ENV_FILE=app.env

APP_NAME=simple-bank
BUILD_DIR=bin
PASSWORD=secret123

OPENSSL=openssl
DAYS=3650
CERT_DIR := cert

CA_KEY := $(CERT_DIR)/ca-key.pem
CA_CERT := $(CERT_DIR)/ca.pem
SERVER_KEY := $(CERT_DIR)/server-key.pem
SERVER_CSR := $(CERT_DIR)/server-req.pem
SERVER_CERT := $(CERT_DIR)/server-cert.pem
SERVER_EXT  := $(CERT_DIR)/server.ext


laptop-server:
	go run ./cmd/server

laptop-client:
	go run ./cmd/client

proto-update:
	cd proto && buf dep update

proto: proto-lint proto-gen
proto-2: proto-update proto-gen

proto-lint:
	cd proto && buf lint

proto-gen:
	cd proto && buf generate

evans:
	evans -r repl -p 9091

compose-up-redis:
	docker compose -f docker-compose.redis.yml up -d

gen:
	rm -f pb/*.go
	protoc --proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

dbdocs:
	dbdocs build doc/db.dbml

dbdocs-set-password:
	dbdocs password -s $(PASSWORD) -p simple_bank


compose-up:
	docker compose up -d --build

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


