-include .env
.SILENT:
CURRENT_DIR=$(shell pwd)
DB_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

swag-init:
	swag init -g api/api.go -o api/docs
	
run:
	go run cmd/main.go

up:
	docker-compose up -d

down:
	docker-compose down 

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-up1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migrate-down1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

proto-gen:
	rm -rf genproto
	./scripts/gen-proto.sh ${CURRENT_DIR}

test:
	go test -v -cover ./...

pull-sub-module:
	git submodule update --init --recursive

update-sub-module:
	git submodule update --remote --merge 

lint:
	golangci-lint run

cache:
	go clean -testcache
